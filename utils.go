package main

import (
	"bufio"
	"fmt"
	"github.com/Tinkerforge/go-api-bindings/io16_bricklet"
	"github.com/Tinkerforge/go-api-bindings/ipconnection"
	"github.com/Tinkerforge/go-api-bindings/stepper_brick"
	"log"
	"os"
	"strconv"
	"time"
)

func ShowAllMinutes() {

	ipcon := ipconnection.New()
	defer ipcon.Close()
	io, _ := io16_bricklet.New(IO16UID, &ipcon) // Create device object.
	minute, _ := stepper_brick.New(MINUTESTEPUID, &ipcon)

	err := ipcon.Connect(ADDR) // Connect to brickd.
	if err != nil {
		return
	}

	defer ipcon.Disconnect()

	for i := 0; i < 60; i++ {

		fmt.Printf("Displaying %d (%b)\n", i, lookupMinute(i))

		driveToPos(&io, &minute, lookupMinute(i), 'b')
		time.Sleep(time.Millisecond * 100)
	}

}

func ReadAllPos(io *io16_bricklet.IO16Bricklet, stepper *stepper_brick.StepperBrick, port rune) {
	//setup stepper Stepper
	stepper.SetMotorCurrent(800)
	stepper.SetStepMode(8)
	stepper.SetMaxVelocity(40)
	stepper.SetSpeedRamping(0, 0)

	//for i := 0; i < len(codes); i = i + 1 {

	stepper.Enable()

	f, err := os.Create("file.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for i := 0; i < 600; i = i + 1 {

		stepper.SetSteps(-200)

		// Get current value from port A as bitmask.
		valueMaskA, _ := io.GetPort(port)

		fmt.Printf("%b -> %d\n", valueMaskA, valueMaskA)
		f.WriteString(fmt.Sprintf("%b -> %d\n", valueMaskA, valueMaskA))
		time.Sleep(time.Second * 1)
	}

	//}

	stepper.FullBrake()
	stepper.Disable()

}

func ConfigMinutes(io *io16_bricklet.IO16Bricklet, stepper *stepper_brick.StepperBrick) {

	//setup stepper Stepper
	stepper.SetMotorCurrent(800)
	stepper.SetStepMode(8)
	stepper.SetMaxVelocity(2000)
	stepper.SetSpeedRamping(0, 0)

	//for i := 0; i < len(codes); i = i + 1 {

	codes := []byte{
		0b11000001,
		0b11000010,
		0b11000011,
		0b11000100,
		0b11000101,
		0b11000110,
		0b11000111,
		0b11001000,
		0b11001001,
		0b11001011,
		0b11001100,
		0b11001101,
		0b11001110,
		0b11001111,
		0b11010000,
		0b11010001,
		0b11010010,
		0b11010011,
		0b11010100,
		0b11010101,
		0b11010110,
		0b11010111,
		0b11011000,
		0b11011001,
		0b11011010,
		0b11011011,
		0b11011100,
		0b11011101,
		0b11011110,
		0b11011111,
		0b11100000,
		0b11100001,
		0b11100010,
		0b11100011,
		0b11100100,
		0b11100101,
		0b11100110,
		0b11100111,
		0b11101000,
		0b11101001,
		0b11101010,
		0b11101011,
		0b11101100,
		0b11101101,
		0b11101110,
		0b11101111,
		0b11110000,
		0b11110001,
		0b11110010,
		0b11110011,
		0b11110100,
		0b11110101,
		0b11110110,
		0b11110111,
		0b11111000,
		0b11111001,
		0b11111010,
		0b11111011,
		0b11111100,
		0b11111101,
		0b11111110,
	}

	for i := 0; i < len(codes); i = i + 1 {

		stepper.Enable()
		stepper.DriveBackward()

		fmt.Printf("Looking for %b\n", codes[i])

		for {

			// Get current value from port B as bitmask.
			valueMaskB, _ := io.GetPort('b')

			if valueMaskB == codes[i] {
				stepper.FullBrake()
				stepper.Stop()
				stepper.Disable()
				time.Sleep(time.Second * 2)
				break
			}
		}

	}

	//}

	stepper.FullBrake()
	stepper.Disable()

}

func Ask() {
	fmt.Print("insert y value here: ")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	fmt.Println(input.Text())

	ipcon := ipconnection.New()
	defer ipcon.Close()
	io, _ := io16_bricklet.New(IO16UID, &ipcon) // Create device object.
	hour, _ := stepper_brick.New(HOURSTEPUID, &ipcon)
	//minute, _ := stepper_brick.New(MINUTESTEPUID, &ipcon)

	err := ipcon.Connect(ADDR) // Connect to brickd.
	if err != nil {
		return
	}

	defer ipcon.Disconnect()
	// Don't use device before ipcon is connected.
	newPos, err := strconv.Atoi(input.Text())

	driveToPos(&io, &hour, lookupHour(newPos), 'a')

	Ask()
}
