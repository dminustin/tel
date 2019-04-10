package main

import (
	"amf_common"
	//"fmt"
	"time"
	"github.com/suapapa/go_devices/tm1638"
	"github.com/suapapa/go_devices/rpi/gpio"

	"os"
	"os/signal"

	"fmt"
	"strconv"
)


var btns = []int{
	5:1,
	4:2,
	3:3,
	2:4,
	1:5,
	0:11,
}
var leds = []int{
	1:2,
	2:3,
	3:4,
	4:5,
	5:6,
	11:7,
}


var appsStates = make(map[string]string)
var exitC = make(chan struct{})

func main() {
	amf_common.AppHeader = amf_common.TAppHeader{
		AppName: "ph1_buttons",
		AppVersion: "1.0.0",
	}
	amf_common.OnAppStart()
	appsStates["on"] = "-1"
	//load phonebook

	defer amf_common.OnAppEnd()
	m, err := tm1638.Open(
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
	defer m.Close()

	for i := 0; i < 8; i++ {
		m.SetLed(i, tm1638.Off)
		//	time.Sleep(time.Millisecond * 100)
	}

	m.SetString("        ")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		exitC <- struct{}{}
	}()

	time.Sleep(time.Second)


	go func() {
		current_button := -1
		counter :=0
		for {
			keys := m.GetButtons()

//			var str string
			for i := 0; i < 8; i++ {
				if keys&(1<<byte(i)) == 0 {
//					str += "0"
					//m.SetLed(i, tm1638.Off)
				} else {
					if (current_button != btns[i]) {
						current_button = btns[i]
						go sendKeyPresed(current_button)
					}
//					str += "1"
					//m.SetLed(btns[i], tm1638.Green)
					//go showBtn(m, i)
				}
			}
			time.Sleep(10 * time.Millisecond)
			counter++
			if (counter > 20) {
				counter = 0
				cmd:= amf_common.ReadAppCommand("")

				if (len(cmd["off"])>0) {
					if (cmd["off"] == "all") {
						for rt := 0; rt < 8; rt++ {
							m.SetLed(rt, tm1638.Off)
						}
						continue
					}
					i,_ := strconv.ParseInt(cmd["off"], 10, 16)
					m.SetLed(leds[int(i)], tm1638.Off)
					appsStates["on"] = "-1"
				}

				if (len(cmd["on"])>0) {

					i,e := strconv.ParseInt(cmd["on"], 10, 16)
					if (e!=nil) {
						continue
					}
					if (i > -1) {

						if (cmd["on"]!=appsStates["on"]) {
							tmp,_ := strconv.ParseInt(appsStates["on"], 10, 16)
							if (tmp > -1) {
								m.SetLed(leds[int(tmp)], tm1638.Off)
								fmt.Printf("%v OFF\n", int(tmp))
							}
							appsStates["on"] = cmd["on"]
							fmt.Println(i)
							m.SetLed(leds[int(i)], tm1638.Green)
							appsStates["on"] = cmd["on"]
							fmt.Printf("%v ON\n", int(i))
						}
					}
				}
			}
		}
	}()

	<-exitC



}

func sendKeyPresed(k int) {
	pr:= make(map[string]string)
	pr["clicked"] = fmt.Sprintf("%v", k)
	amf_common.WriteToMem(pr, "ph1_buttons")
	fmt.Println(fmt.Sprintf("%v", k))
}
