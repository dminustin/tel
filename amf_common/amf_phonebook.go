package amf_common

import (
	"fmt"
	"strings"
	"strconv"
)

type TPhoneBookRec struct {
	Phone string
	Name string
	Button int
}

var PhoneBook map[string]TPhoneBookRec
var PhonesOnButtons map[int]string

func LoadPhoneBook() {
	PhoneBook = make(map[string]TPhoneBookRec)
	PhonesOnButtons = make(map[int]string)
	var arr, err = ReadLines("./phonebook.ini")
	if (err == nil) {
		for i,_ := range(arr) {
			var line = strings.Split(arr[i],";")
			if (len(line[0])>0) {
				if (line[2] == "") {
					line[2] = "0"
				}
				var btn,_ = strconv.ParseInt(line[2], 10, 8)
				PhoneBook[line[0]] = TPhoneBookRec{
					Name: line[1],
					Button: int(btn),
					Phone: line[0],
				}

				if (int(btn) >0) {
					PhonesOnButtons[int(btn)] = line[0]
				}

			}
		}
	}
	fmt.Println(PhoneBook, PhonesOnButtons)

}
