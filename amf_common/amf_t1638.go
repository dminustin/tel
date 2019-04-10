package amf_common

import (
	"suapapa/go_devices/tm1638"
	"suapapa/go_devices/rpi/gpio"

	"time"
	"fmt"
)

var btns = []int{
	5:1,
	4:2,
	3:3,
	2:4,
	1:5,
	0:8,
//unused
	6:8,
	7:8,
}
var leds = []int{
	1:2,
	2:3,
	3:4,
	4:5,
	5:6,
	8:7,
//unused
	6:7,
	7:7,
}

const BUTTONSTATE_idle = "idle"
const BUTTONSTATE_pressed = "pressed"

const LEDSSTATE_idle = "idle"
const LEDSSTATE_on = "on"

const LEDSCOMMAND_blink = "blink"
const LEDSCOMMAND_idle = "idle"
const LEDSCOMMAND_on = "on"


type ledChanStruct struct {
	ledID int
	ledState tm1638.Color
}

var ledChan = make(chan ledChanStruct,10)

type TT1638State struct {
	ButtonState map[int]string
	LedState map[int]string
	LedCommand map[int]string
}

var T1638State = TT1638State {
	ButtonState: map[int]string{1: BUTTONSTATE_idle,2: BUTTONSTATE_idle,3: BUTTONSTATE_idle,4: BUTTONSTATE_idle,5: BUTTONSTATE_idle,6: BUTTONSTATE_idle,7: BUTTONSTATE_idle,8: BUTTONSTATE_idle},
	LedState: map[int]string{1: LEDSSTATE_idle,2: LEDSSTATE_idle,3: LEDSSTATE_idle,4: LEDSSTATE_idle,5: LEDSSTATE_idle,6: LEDSSTATE_idle,7: LEDSSTATE_idle,8: LEDSSTATE_idle},
	LedCommand: map[int]string{1: LEDSCOMMAND_idle,2: LEDSCOMMAND_idle,3: LEDSCOMMAND_idle,4: LEDSCOMMAND_idle,5: LEDSCOMMAND_idle,6: LEDSCOMMAND_idle,7: LEDSCOMMAND_idle,8: LEDSCOMMAND_idle},
}
var defaultState = TT1638State {
	ButtonState: map[int]string{1: BUTTONSTATE_idle,2: BUTTONSTATE_idle,3: BUTTONSTATE_idle,4: BUTTONSTATE_idle,5: BUTTONSTATE_idle,6: BUTTONSTATE_idle,7: BUTTONSTATE_idle,8: BUTTONSTATE_idle},
	LedState: map[int]string{1: LEDSSTATE_idle,2: LEDSSTATE_idle,3: LEDSSTATE_idle,4: LEDSSTATE_idle,5: LEDSSTATE_idle,6: LEDSSTATE_idle,7: LEDSSTATE_idle,8: LEDSSTATE_idle},
	LedCommand: map[int]string{1: LEDSCOMMAND_idle,2: LEDSCOMMAND_idle,3: LEDSCOMMAND_idle,4: LEDSCOMMAND_idle,5: LEDSCOMMAND_idle,6: LEDSCOMMAND_idle,7: LEDSCOMMAND_idle,8: LEDSCOMMAND_idle},
}

var Module_ = &tm1638.Module{}
func LedRoutine() {

	Module_, err = tm1638.Open(
		&gpio.Sysfs{
			PinMap: map[string]int{
				tm1638.PinCLK:  27,
				tm1638.PinDATA: 22,
				tm1638.PinSTB:  17,
			},
		},
	)
	if err != nil {
		panic(err)
	}
	defer Module_.Close()
	
	
	for i:=0; i<8; i++ {
	    Module_.SetLed(leds[int(i)], tm1638.Off)
	}


	go func() {
		current_button := -1
		for {
			keys := Module_.GetButtons()

			//			var str string
			for i := 0; i < 8; i++ {
				if keys&(1<<byte(i)) == 0 {
					current_button = btns[i]
					T1638State.ButtonState[current_button] = BUTTONSTATE_idle
				} else {
					fmt.Println("PR ", i)
					current_button = btns[i]
					fmt.Println("BTN Pressed", current_button)
					T1638State.ButtonState[current_button] = BUTTONSTATE_pressed
				}
			}
			time.Sleep(10 * time.Millisecond)

			select {
			case tmp := <-ledChan: {
				fmt.Println("Set led routine", tmp.ledID, tmp.ledState)
				Module_.SetLed(leds[int(tmp.ledID)], tmp.ledState)
			}
			default: {
				//
			}
			}

		}
	}()



	blink :=0
	blink_interval := 1000
	for {
		for i := 1; i < 9; i++ {
			//новая команда?
			if (T1638State.LedCommand[i] != defaultState.LedCommand[i]) {
				defaultState.LedCommand[i] = T1638State.LedCommand[i]

				if (T1638State.LedState[i] == LEDSCOMMAND_idle) {
					if (T1638State.LedState[i] == LEDSCOMMAND_on) {
						SetLed(i, LEDSCOMMAND_on)
					}
				} else if (T1638State.LedState[i] == LEDSCOMMAND_on) {
					if (T1638State.LedState[i] == LEDSCOMMAND_idle) {
						SetLed(i, LEDSCOMMAND_idle)
					}
				}
			}

			if (T1638State.LedCommand[i] == LEDSCOMMAND_blink) {
				if (blink >= blink_interval) {
					if (T1638State.LedState[i] == LEDSSTATE_idle) {
						SetLed(i, LEDSCOMMAND_on)
					}
				} else {
					if (T1638State.LedState[i] == LEDSSTATE_on) {
						SetLed(i, LEDSCOMMAND_idle)
					}
				}
			}

		}
		if (blink>=blink_interval) {
			blink = 0
		} else {
			blink +=100
		}

		time.Sleep(time.Millisecond * 100)
	}
}

func SendLedCommand(num int, state string) {
	T1638State.LedCommand[num]=state

}

func SetLed(num int, state string) {
	T1638State.LedState[num] = state
	if (state == LEDSSTATE_idle) {
		ledChan <- ledChanStruct{ledID:num, ledState: tm1638.Off}
	} else {
		ledChan <- ledChanStruct{ledID:num, ledState: tm1638.Green}
	}
}

func GetButtonsState() (map[int]string){
	fmt.Println("GET STATE BTN")
	result := map[int]string{
		1: BUTTONSTATE_idle,
		2: BUTTONSTATE_idle,
		3: BUTTONSTATE_idle,
		4: BUTTONSTATE_idle,
		5: BUTTONSTATE_idle,
		6: BUTTONSTATE_idle,
		7: BUTTONSTATE_idle,
		8: BUTTONSTATE_idle,
	}
	for i:=1; i<9; i++ {
		if (T1638State.ButtonState[i] == BUTTONSTATE_pressed) {
			if (defaultState.ButtonState[i] == BUTTONSTATE_idle) {
				defaultState.ButtonState[i] = BUTTONSTATE_pressed
			}
		} else if (T1638State.ButtonState[i] == BUTTONSTATE_idle) {
			if (defaultState.ButtonState[i] == BUTTONSTATE_pressed) {
				defaultState.ButtonState[i] = BUTTONSTATE_idle
				result[i] = BUTTONSTATE_pressed
			}
		}

	}
	return result
}
