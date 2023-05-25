package main

import (
	"fmt"
	"github.com/Tinkerforge/go-api-bindings/dual_button_bricklet"
	"github.com/Tinkerforge/go-api-bindings/io16_bricklet"
	"github.com/Tinkerforge/go-api-bindings/ipconnection"
	"github.com/Tinkerforge/go-api-bindings/stepper_brick"
	"github.com/multiplay/go-cticker"
	"time"
)

const ADDR string = "localhost:4223"

const IO16UID string = "gqN"
const HOURSTEPUID string = "5Wqtru"
const MINUTESTEPUID string = "6RsvAH"
const BTNUID string = "mz2"

func main() {

	ipcon := ipconnection.New()
	defer ipcon.Close()
	io, _ := io16_bricklet.New(IO16UID, &ipcon) // Create device object.
	hour, _ := stepper_brick.New(HOURSTEPUID, &ipcon)
	minute, _ := stepper_brick.New(MINUTESTEPUID, &ipcon)
	db, _ := dual_button_bricklet.New(BTNUID, &ipcon)

	err := ipcon.Connect(ADDR) // Connect to brickd.
	if err != nil {
		return
	}

	defer ipcon.Disconnect()

	goToCurrentTimeSetTime(&io, &hour, &minute)

	autoTick := true // Show current Time on Clock, will be set to false when going to blue

	db.SetSelectedLEDState(0, 3)
	db.SetSelectedLEDState(1, 3)

	db.RegisterStateChangedCallback(func(buttonL dual_button_bricklet.ButtonState, buttonR dual_button_bricklet.ButtonState, ledL dual_button_bricklet.LEDState, ledR dual_button_bricklet.LEDState) {
		if buttonL == dual_button_bricklet.ButtonStatePressed {
			autoTick = false
			goToBlue(&io, &hour, &minute)
		}
		if buttonR == dual_button_bricklet.ButtonStatePressed {
			autoTick = true
			goToCurrentTimeSetTime(&io, &hour, &minute)
		}
	})

	t := cticker.New(time.Minute, time.Second)

	for tick := range t.C {
		fmt.Println(tick)
		if autoTick {
			goToCurrentTime(&io, &hour, &minute, tick.Hour(), tick.Minute())
		}
	}

}

func goToCurrentTimeSetTime(io *io16_bricklet.IO16Bricklet, hour *stepper_brick.StepperBrick, minute *stepper_brick.StepperBrick) {
	cHour, cMinute, _ := time.Now().Clock()
	goToCurrentTime(io, hour, minute, cHour, cMinute)
}

func goToCurrentTime(io *io16_bricklet.IO16Bricklet, hour *stepper_brick.StepperBrick, minute *stepper_brick.StepperBrick, cHour int, cMinute int) {
	driveToPos(io, hour, lookupHour(cHour), 'a')
	driveToPos(io, minute, lookupMinute(cMinute), 'b')
}

func goToBlue(io *io16_bricklet.IO16Bricklet, hour *stepper_brick.StepperBrick, minute *stepper_brick.StepperBrick) {
	driveToPos(io, hour, lookupHour(99), 'a')
	driveToPos(io, minute, lookupMinute(99), 'b')
}

func driveToPos(io *io16_bricklet.IO16Bricklet, stepper *stepper_brick.StepperBrick, pos byte, port rune) {

	stepper.SetMotorCurrent(800)
	stepper.SetStepMode(8)
	stepper.SetMaxVelocity(2000)
	stepper.SetSpeedRamping(0, 0)

	stepper.Enable()
	stepper.DriveBackward()
	for {
		// Get current value from port A as bitmask.
		valueMaskA, _ := io.GetPort(port)
		if valueMaskA == pos {
			stepper.FullBrake()
			break
		}
	}

	stepper.FullBrake()
	stepper.Disable()

}

func lookupHour(hourSearch int) byte {
	hour := map[int]byte{
		99: 0b11110001,
		0:  0b11011101,
		1:  0b11101111,
		2:  0b11101110,
		3:  0b11001110,
		4:  0b11101101,
		5:  0b11001101,
		6:  0b11001100,
		7:  0b11111010,
		8:  0b11011011,
		9:  0b11111001,
		10: 0b11111000,
		11: 0b11011000,
		12: 0b11101010,
		13: 0b11001011,
		14: 0b11101001,
		15: 0b11101000,
		16: 0b11001000,
		17: 0b11110111,
		18: 0b11010111,
		19: 0b11010110,
		20: 0b11110100,
		21: 0b11010100,
		22: 0b11100111,
		23: 0b11000111,
	}

	if val, ok := hour[hourSearch]; ok {
		return val
	}
	return 0b11110001
}
func lookupMinute(minuteSearch int) byte {
	minutes := map[int]byte{
		0:  0b11001011,
		1:  0b11001010,
		2:  0b11101001,
		3:  0b11101000,
		4:  0b11001001,
		5:  0b11001000,
		6:  0b11110111,
		7:  0b11110110,
		8:  0b11010111,
		9:  0b11010110,
		10: 0b11110101,
		11: 0b11110100,
		12: 0b11010101,
		13: 0b11010100,
		14: 0b11100111,
		15: 0b11100110,
		16: 0b11000111,
		17: 0b11000110,
		18: 0b11100101,
		19: 0b11100100,
		20: 0b11000101,
		21: 0b11000100,
		22: 0b11110011,
		23: 0b11110010,
		24: 0b11010011,
		25: 0b11010010,
		26: 0b11110001,
		27: 0b11110000,
		28: 0b11010001,
		29: 0b11010000,
		30: 0b11100011,
		98: 0b11100010,
		31: 0b11000011,
		32: 0b11000010,
		33: 0b11100001,
		34: 0b11100000,
		35: 0b11000001,
		36: 0b11111110,
		37: 0b11011111,
		38: 0b11011110,
		39: 0b11111101,
		40: 0b11111100,
		41: 0b11011101,
		42: 0b11011100,
		43: 0b11101111,
		44: 0b11101110,
		45: 0b11001111,
		46: 0b11001110,
		47: 0b11101101,
		48: 0b11101100,
		49: 0b11001101,
		50: 0b11001100,
		51: 0b11111011,
		52: 0b11111010,
		53: 0b11011011,
		54: 0b11011010,
		55: 0b11111001,
		56: 0b11111000,
		57: 0b11011001,
		58: 0b11011000,
		59: 0b11101011,
		99: 0b11101010,
	}
	if val, ok := minutes[minuteSearch]; ok {
		return val
	}
	return 0b11101010
}
