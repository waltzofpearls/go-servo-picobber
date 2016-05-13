# go-servo-picobber

A golang port for
[Geekroo PWM Servo PiCobber](http://www.robotshop.com/ca/en/pwm-servo-driver-picobber-raspberry-pi-hat.html)
[python driver](https://github.com/geekroo/Geekroo-PiCobber-PWMServo).

The PWM Servo Driver PiCobber Raspberry Pi HAT has the ability to drive
up to 16 servos or PWM outputs with over 12C with only 2 pins. The PWM
controller on-board drives all the 16 channels and it does not require
any additional Raspberry Pi processing overhead.

## Install the driver

```shell
go get github.com/waltzofpearls/go-servo-picobber
```

## Example

```go
package main

import (
    "log"
    "time"

    "github.com/waltzofpearls/go-servo-picobber"
)

func main() {
    servoMin := 150 // Min pulse length out of 4096
    servoMax := 600 // Max pulse length out of 4096

    sv, err := servo.NewServo()
    if err != nil {
        log.Println(err)
    }
    sv.SetPwmFreq(60) // Set frequency to 60 Hz
    for {
        // Change speed of continuous servo on channel O
        sv.SetPwm(0, 0, servoMin)
        time.Sleep(1 * time.Second)
        sv.SetPwm(0, 0, servoMax)
        time.Sleep(1 * time.Second)
    }
}
```
