package main

import (
	"amf_common"
	//"fmt"
	"time"
	"fmt"
)


var appsStates = map[string]string{}

func main() {
	amf_common.AppHeader = amf_common.TAppHeader{
		AppName: "ph1_core",
		AppVersion: "1.0.0",
	}
	amf_common.OnAppStart()

	//load phonebook

	defer amf_common.OnAppEnd()

	for ;; {
		time.Sleep(time.Millisecond * 300)
		//check buttons state

		x := amf_common.ReadFromMem("ph1_buttons")
		var clicked = x["clicked"]

		if(len(clicked)>0) {
			if (appsStates["clicked"] != clicked) {
				appsStates["clicked"] = clicked
				fmt.Printf("%v\n\n", clicked)
				go amf_common.SendAppComand("ph1_buttons", map[string]string{
					"on": clicked,
				})
			}
		}



		//check phone state

		var phph = map[string]string{
			"+79055807231": "1",
			"+79672856116": "2",
		}

		phoneState := amf_common.ReadAppState("ph1_phone")
		if (len(phoneState["INCOMING"])>0) {
			go amf_common.SendAppComand("ph1_buttons", map[string]string{
				"on": phph[phoneState["INCOMING"]],
			})
		}else if(len(phoneState["HANGUP"])>0) {
			go amf_common.SendAppComand("ph1_buttons", map[string]string{
				"off": "all",
			})
		}

		//check mb temperature

		//check outside commands
	}


}


