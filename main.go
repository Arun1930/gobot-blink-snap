package main

import (
  "log"
  "os/exec"
  "strings"
  "time"

  "gobot.io/x/gobot"
  "gobot.io/x/gobot/drivers/gpio"
  "gobot.io/x/gobot/platforms/beaglebone"
  "gobot.io/x/gobot/platforms/dragonboard"
  "gobot.io/x/gobot/platforms/raspi"
)

func main() {
  out, err := exec.Command("uname", "-r").Output()
  if err != nil {
    log.Fatal(err)
  }
  var adaptor gobot.Adaptor
  var pin string
  kernelRelease := string(out)
  if strings.Contains(kernelRelease, "raspi2") {
    adaptor = raspi.NewAdaptor()
    pin = "7"
  } else if strings.Contains(kernelRelease, "snapdragon") {
    adaptor = dragonboard.NewAdaptor()
    pin = "GPIO_A"
  } else {
    adaptor = beaglebone.NewAdaptor()
    pin = "P8_7"
  }
  digitalWriter, ok := adaptor.(gpio.DigitalWriter)
  if !ok {
    log.Fatal("Invalid adaptor")
  }
  led := gpio.NewLedDriver(digitalWriter, pin)

  work := func() {
    gobot.Every(1*time.Second, func() {
      led.Toggle()
    })
  }

  robot := gobot.NewRobot("snapbot",
    []gobot.Connection{adaptor},
    []gobot.Device{led},
    work,
  )

  robot.Start()
}
