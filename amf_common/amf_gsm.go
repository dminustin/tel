package amf_common

import (
	"fmt"
	"io"
	"tarm/goserial"

	"bufio"
	"strings"
	"log"

	"gpio"
	"time"

)

type GsmModem struct {
	port       string
	readWriter io.ReadWriteCloser
}


var modem = &GsmModem{}

const PHONESTATE_idle = "idle"
const PHONESTATE_error = "error"
const PHONESTATE_calling = "calling"
const PHONESTATE_call = "call"

type TPhoneState struct {
	State string
	CurrentPhone string
	SygnalLevel int
	LastError string

}

var PhoneState = TPhoneState{
	State: PHONESTATE_idle,
	CurrentPhone: "",
	SygnalLevel: 0,
	LastError: "",
	}


func resetPhone() {

	_, resp, isOK := modem.SendCommand("AT\r\n")
	wtf := false
	if isOK != true {
	    wtf = true
	}
	if (resp!="OK"){
	    wtf = true
	}
	if (!wtf){
	    return
	}
	pin, err := gpio.OpenPin(gpio.GPIO23, gpio.ModeOutput)
	if err != nil {
		panic(err)
	}
	//defer pin.Close()
	fmt.Println("SET PIN23")
	pin.Set()
	time.Sleep(2 * time.Second)
	fmt.Println("UNSET PIN23")
	pin.Clear()

	time.Sleep(7 * time.Second)
	fmt.Println("CLOSE PIN23")

	pin.Close()
}

func InitPhone(port string) {
	modem = NewGsmModem(port)
	modem.Connect()
	resetPhone()
	lines, resp, isOK := modem.SendCommand("AT+CLIP=1\r\n")
	fmt.Println(resp, lines, isOK)
	if isOK != true {
		PhoneState.LastError = "Could not initialize modem"
		PhoneState.State = PHONESTATE_error
	} else {
		PhoneState.State = PHONESTATE_idle
	}

}

func GetPhoneState() {
	fmt.Println("Get Phone State")



	//on ring
	linezz, resp, isOK := modem.Read()
	if (!isOK) {
		PhoneState.State = PHONESTATE_error
		return
	}
	fmt.Println(linezz, resp, isOK)
	switch(resp) {
		case "CLIP" : {
			PhoneState.State = PHONESTATE_call
			PhoneState.CurrentPhone = linezz
			break
		}
		case "HANGUP" : {
			PhoneState.State = PHONESTATE_idle
			PhoneState.CurrentPhone = ""
			break
		}
	}


}

func GetSygnalLevel() {
	fmt.Println("Get Sygnal Level")
}

func CallNumber(num string) {
	if PhoneState.State != PHONESTATE_idle {
		if PhoneState.State == PHONESTATE_call {
			if PhoneState.CurrentPhone!=num {
				HangUp()
			} else {
				return
			}
		} else if PhoneState.State == PHONESTATE_calling {
			if PhoneState.CurrentPhone!=num {
				HangUp()
			} else {
				return
			}
		}
	}
	PhoneState.State = PHONESTATE_calling
	PhoneState.CurrentPhone=num
	modem.SendCommand("ATD" + num + ";")
}

func HangUp() {
	PhoneState.State = PHONESTATE_idle
	PhoneState.CurrentPhone=""
	modem.SendCommand("ATH")
}

func Answer() {
	PhoneState.State = PHONESTATE_call
	modem.SendCommand("ATA")
}


//=======================================

func NewGsmModem(port string) *GsmModem {
	return &GsmModem{port: port}
}
func (g *GsmModem) Connect() (io.ReadWriteCloser, error) {
	c := &serial.Config{Name: g.port, Baud: 115200}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	//defer s.Close()
	g.readWriter = s
	return g.readWriter, nil
}


func (g *GsmModem) Read() (line string, response string, isOK bool) {
	r := bufio.NewReader(g.readWriter)
	isOK = true
	line = ""
	response = "nil"
	for {
		read, _, err := r.ReadLine()
		if err != nil {
			isOK = false
			return
		}
		strread := string(read)

		if strings.Contains(strread, "+CLIP") {
			response = "CLIP"
			line = PH_parseClip(strread)
			return
		}
		if strings.Contains(strread, "NO CARRIER") {
			response = "HANGUP"
			line = strread
			return
		}

		if (len(strread)>0) {
			response = "UNKNOWN"
			line = strread
			return
		}


	}
	return
}

func (g *GsmModem) SendCommand(command string) (lines []string, response string, isOK bool) {
	_, err := g.readWriter.Write([]byte(command + "\n"))
	if err != nil {
		log.Fatal(err)
	}
	r := bufio.NewReader(g.readWriter)
	for {
		read, _, err := r.ReadLine()
		if err != nil {
			isOK = false
			return
		}
		strread := string(read)
		fmt.Println(strread)
		if strings.Contains(strread, "OK") {
			response = strread
			isOK = true
			return
		}
		if strings.Contains(strread, "ERROR") {
			response = strread
			isOK = false
			return
		}

		lines = append(lines, strread)
	}
	return
}





