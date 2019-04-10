package main

import (
	"amf_common"
	"fmt"
	"ini"

	"time"
)

type tIntervals struct {
	GetPhoneState int64 //Статус телефона
	GetSygnalLevel int64 //Уровень сигнала
	GetButtonsState int64 //Состояние кнопок
}

var AppIni *ini.File
var err error
var AppName = "Phone App"
var AppVersion = "1.0"
var Intervals = tIntervals{
	GetPhoneState: 30,
	GetSygnalLevel: 300,
	GetButtonsState: 3,
}


func main() {
	amf_common.IS_TEST = true
	fmt.Printf("%v v. %v starts\n", AppName, AppVersion)
	AppIni, err = ini.Load("./main.ini")
	amf_common.FailOnError(err,"Fail to read file")
	amf_common.LoadPhoneBook()
	amf_common.InitPhone("/dev/ttyAMA0")
	go Runner()
	go amf_common.LedRoutine()
	for {
		time.Sleep(time.Second * 60)
	}
}


func Runner() {

	var LastCalls = tIntervals{
		GetPhoneState: 0,
		GetSygnalLevel: 0,
		GetButtonsState: 0,
		}

	for {

		timeNow := time.Now().Unix()
		if timeNow - LastCalls.GetPhoneState > Intervals.GetPhoneState {
			LastCalls.GetPhoneState = timeNow
			go amf_common.GetPhoneState()
		}

		if timeNow - LastCalls.GetSygnalLevel > Intervals.GetSygnalLevel {
			LastCalls.GetSygnalLevel = timeNow
			go amf_common.GetSygnalLevel()
		}

		if timeNow - LastCalls.GetButtonsState > Intervals.GetButtonsState {
			LastCalls.GetButtonsState = timeNow
			buttons:= amf_common.GetButtonsState()
			pressed :=0
			for i:=1; i<9; i++ {
				if buttons[i] == amf_common.BUTTONSTATE_pressed {
					pressed = i
				}
			}

			if pressed>0 {
				for i:=1; i<9; i++ {
					if (i!=pressed){
					amf_common.SendLedCommand(i, amf_common.LEDSCOMMAND_idle)
					}
				}
				if (pressed == 8) {
					go amf_common.HangUp()
					continue
				}
				fmt.Println("Pressed ", pressed)
				phone := amf_common.PhonesOnButtons[pressed]
				if (len(phone)>0) {
					//go amf_common.CallNumber(phone)
					amf_common.SendLedCommand(pressed, amf_common.LEDSCOMMAND_idle)

				}
			}
		}

		time.Sleep(time.Millisecond * 100)
	}


}
