// Copyright 2015, Homin Lee <homin.lee@suapapa.net>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/suapapa/go_devices/rpi/gpio"
	"github.com/suapapa/go_devices/tm1638"
	//"fmt"
	"fmt"
)

var (
	exitC = make(chan struct{})
	btns = []int{
		0:7,
		1:6,
		2:5,
		3:4,
		4:3,
		5:2,
		}
)

func main() {
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

	log.Println("Ctrl-C for exit")
	go func() {
		for {
			keys := m.GetButtons()

			var str string
			for i := 0; i < 8; i++ {
				if keys&(1<<byte(i)) == 0 {
					str += "0"
					//m.SetLed(i, tm1638.Off)
				} else {
					str += "1"
					m.SetString(fmt.Sprintf("%v", i))
					//m.SetLed(btns[i], tm1638.Green)
					go showBtn(m, i)
				}
			}



			time.Sleep(10 * time.Millisecond)
		}
	}()

	<-exitC
}

func showBtn(m *tm1638.Module, n int) {
	m.SetLed(btns[n], tm1638.Green)
	time.Sleep(time.Second)
	m.SetLed(btns[n], tm1638.Off)
}