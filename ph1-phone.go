package main

import (
	"amf_common"
	"github.com/tarm/goserial"

	"log"
	"fmt"
	"io"
	"bufio"
	"strings"
	"time"
)


var appsStates = make(map[string]string)

type GsmModem struct {
	port       string
	readWriter io.ReadWriteCloser
}

//======================================
//PARSERS

func parseClip(clip string) {
	//+CLIP: "+79055807231",145,"",0,"",0
}

func main() {
	appsStates["on"] = "-1"
	amf_common.AppHeader = amf_common.TAppHeader{
		AppName: "ph1_phone",
		AppVersion: "1.0.0",
	}
	amf_common.OnAppStart()
	defer amf_common.OnAppEnd()

	port := "/dev/ttyAMA0"

	modem := NewGsmModem(port)
	modem.Connect()
	lines, resp, isOK := modem.SendCommand("AT+CLIP=1\r\n")
	fmt.Println(resp, lines, isOK)
	if isOK != true {
		log.Fatal("failed to set modem on text mode")
	}

	linezz := ""
	for {
		linezz, resp, isOK = modem.Read()

		switch(resp) {
		case "CLIP" : {
			amf_common.WriteAppState(map[string]string{
				"INCOMING": linezz,
			})
			break
		}
		case "HANGUP" : {
			amf_common.WriteAppState(map[string]string{
				"HANGUP": linezz,
			})
			break
		}
		}

		fmt.Printf("LI=%v\nRes=%v\nIS=%v\n\n", linezz, resp, isOK)
		time.Sleep(time.Millisecond * 300)
	}


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
			line = amf_common.PH_parseClip(strread)
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





