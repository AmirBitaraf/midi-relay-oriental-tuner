package main

import (
	"fmt"
	"os"
	"time"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	_ "gitlab.com/gomidi/midi/v2/drivers/portmididrv"
)

func writer(out drivers.Out) (chan bool, chan midi.Message) {
	quit := make(chan bool)
	message := make(chan midi.Message)

	go func() {
		out.Open()
		defer out.Close()

		outSender, err := midi.SendTo(out)
		if err != nil {
			fmt.Println("Error Initiating Output Port!")
			return
		}

		for {
			select {
			case msg := <-message:
				outSender(msg)
				time.Sleep(time.Microsecond * 1200)
			case <-quit:
				return
			}
		}

	}()
	return quit, message
}

func relay(in drivers.In, out drivers.Out) chan bool {
	quit := make(chan bool)
	go func() {
		outQuit, outSender := writer(out)

		fmt.Println("Ports:\nIn: " + in.String() + "\nOut: " + out.String())

		for ch := 0; ch < 16; ch++ {
			outSender <- midi.Pitchbend(uint8(ch), 0)
		}

		outSender <- midi.Pitchbend(3, -8192/4)

		stop, err := midi.ListenTo(in, func(msg midi.Message, timestampms int32) {

			var bt []byte
			var ch, key, vel uint8
			var control, value uint8

			switch {
			case msg.GetSysEx(&bt):
				fmt.Printf("got sysex: % X\n", bt)
				outSender <- msg
			case msg.GetNoteStart(&ch, &key, &vel):
				fmt.Printf("starting note %s on channel %v with velocity %v\n", midi.Note(key), ch, vel)

				if midi.Note(key).Name() == "B" {
					outSender <- midi.NoteOn(3, key, vel)
				} else {
					outSender <- midi.NoteOn(0, key, vel)
				}

			case msg.GetNoteEnd(&ch, &key):
				fmt.Printf("ending note %s on channel %v\n", midi.Note(key), ch)

				if midi.Note(key).Name() == "B" {
					outSender <- midi.NoteOff(3, key)
				} else {
					outSender <- midi.NoteOff(0, key)
				}

			case msg.GetControlChange(&ch, &control, &value):
				fmt.Println(midi.ControlChange(ch, control, value).String())

				outSender <- midi.ControlChange(0, control, value)
				outSender <- midi.ControlChange(3, control, value)

			default:
				outSender <- msg
			}
		}, midi.UseSysEx())

		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			return
		}
		<-quit
		outQuit <- true
		stop()
	}()
	return quit
}

func main() {
	defer midi.CloseDriver()

	inPorts := midi.GetInPorts()

	if len(inPorts) == 0 {
		fmt.Println("Can't find a MIDI input port!")
		os.Exit(1)
	}
	in := inPorts[0]

	outPorts := midi.GetOutPorts()
	fmt.Println(outPorts.String())
	if len(outPorts) == 0 {
		fmt.Println("Can't find a MIDI output port!")
		os.Exit(1)
	}
	out := outPorts[0]

	_ = relay(in, out)
	select {}
}
