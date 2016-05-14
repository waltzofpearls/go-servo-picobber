package servo

import (
	"log"
	"math"
	"time"

	"bitbucket.org/gmcbay/i2c"
)

const (
	I2C_BUS  byte = 1
	I2C_ADDR byte = 0x6F

	SUBADR1      byte = 0x02
	SUBADR2      byte = 0x03
	SUBADR3      byte = 0x04
	MODE1        byte = 0x00
	PRESCALE     byte = 0xFE
	LED0_ON_L    byte = 0x06
	LED0_ON_H    byte = 0x07
	LED0_OFF_L   byte = 0x08
	LED0_OFF_H   byte = 0x09
	ALLLED_ON_L  byte = 0xFA
	ALLLED_ON_H  byte = 0xFB
	ALLLED_OFF_L byte = 0xFC
	ALLLED_OFF_H byte = 0xFD
)

type Servo struct {
	i2cbus *i2c.I2CBus
}

func NewServo() (*Servo, error) {
	i2cbus, err := i2c.Bus(I2C_BUS)
	if err != nil {
		return nil, err
	}
	log.Println("Reseting PCA9685")
	if err = i2cbus.WriteByte(I2C_ADDR, MODE1, 0x01); err != nil {
		return nil, err
	}
	return &Servo{i2cbus}, nil
}

// Sets the PWM frequency
func (sv *Servo) SetPwmFreq(freq uint8) error {
	prescaleval := 25000000.0          // 25MHz
	prescaleval = prescaleval / 4096.0 // 12-bit
	prescaleval = prescaleval / float64(freq)
	prescaleval = prescaleval - 1.0
	log.Printf("Setting PWM frequency to %d Hz", freq)
	log.Printf("Estimated pre-scale: %d", prescaleval)
	prescale := math.Floor(prescaleval + 0.5)
	log.Printf("Final pre-scale: %d", prescale)

	list, err := sv.i2cbus.ReadByteBlock(I2C_ADDR, MODE1, 1)
	if err != nil {
		return err
	}
	oldmode := list[0]
	newmode := (oldmode & 0x7F) | 0x10                                   // sleep
	if err = sv.i2cbus.WriteByte(I2C_ADDR, MODE1, newmode); err != nil { // go to sleep
		return err
	}
	if err = sv.i2cbus.WriteByte(I2C_ADDR, PRESCALE, uint8(math.Floor(prescale))); err != nil {
		return err
	}
	if err = sv.i2cbus.WriteByte(I2C_ADDR, MODE1, oldmode); err != nil {
		return err
	}
	time.Sleep(5 * time.Millisecond)
	if err = sv.i2cbus.WriteByte(I2C_ADDR, MODE1, oldmode|0x80); err != nil {
		return err
	}
	return nil
}

// Sets a single PWM channel
func (sv *Servo) SetPwm(channel uint8, on, off int) error {
	if err := sv.i2cbus.WriteByte(I2C_ADDR, LED0_ON_L+4*channel, byte(on&0xFF)); err != nil {
		return err
	}
	if err := sv.i2cbus.WriteByte(I2C_ADDR, LED0_ON_H+4*channel, byte(on>>8)); err != nil {
		return err
	}
	if err := sv.i2cbus.WriteByte(I2C_ADDR, LED0_OFF_L+4*channel, byte(off&0xFF)); err != nil {
		return err
	}
	if err := sv.i2cbus.WriteByte(I2C_ADDR, LED0_OFF_H+4*channel, byte(off>>8)); err != nil {
		return err
	}
	return nil
}
