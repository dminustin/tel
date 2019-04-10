package amf_common

import (
	"fmt"
	"os"
	"ini"
	"memcache/memcache"
	"encoding/json"
	"bufio"
)

type TAppHeader struct {
	AppName string;
	AppVersion string;
	LockFile string;
	LogFile string;
	ControlPort string;
}

var IS_TEST = false
var AppIni *ini.File
var err error

var AppHeader TAppHeader
var MemClient *memcache.Client

func OnAppStart() {
	fmt.Printf("%v v. %v starts\n", AppHeader.AppName, AppHeader.AppVersion)
	AppIni, err = ini.Load("./main.ini")
	FailOnError(err,"Fail to read file")

	AppHeader.LockFile = AppIni.Section("paths").Key("lock").String() + AppHeader.AppName + ".lock"
	AppHeader.LogFile = AppIni.Section("paths").Key("log").String() + AppHeader.AppName + ".log"

	if IsFileExists(AppHeader.LockFile) {
		os.Remove(AppHeader.LockFile)
	}

	f, err := os.Create(AppHeader.LockFile)
	FailOnError(err,"Fail to create file")
	f.Close()


	//ListenPort
	AppHeader.ControlPort = AppIni.Section("control").Key("port").String()

	//Start Memcache connection
	MemClient = memcache.New(AppIni.Section("memcache").Key("server").String())

	MemClient.Set (&memcache.Item{
			Key: AppHeader.AppName+"_test",
			Value: []byte("test"),
			Expiration: 10,
			})
	it, err := MemClient.Get(AppHeader.AppName+"_test")
	FailOnError(err, "Cannot connect to Memcache")
	if (string(it.Value) != "test") {
		WriteLog("Error geting test value from Memcache")
		os.Exit(1)
	}

}

func OnAppEnd() {
	if IsFileExists(AppHeader.LockFile) {
		os.Remove(AppHeader.LockFile)
	}
}

func IsFileExists(filename string) (bool) {
	if _, err := os.Stat(filename);
	os.IsNotExist(err) {
		return false
	}
	return true
}

func FailOnError(err error, message string) {
	if err != nil {
		WriteLog(message)
		fmt.Printf("Error: %v\n%v\n\n", err, message)
		os.Exit(1)
	}
}

func WriteLog(message string) {
	f, err := os.OpenFile(AppHeader.LogFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(message); err != nil {
		panic(err)
	}
}


func MapToJson(m map[string]string) string {
	jsonString, _ := json.Marshal(m)
	return string(jsonString)
}

func JsonToMap(s string) (map[string]string) {
	var x = map[string]string{}
	json.Unmarshal([]byte(s), &x)
	return x
}




// readLines reads a whole file into memory
// and returns a slice of its lines.
func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// writeLines writes the lines to the given file.
func WriteLines(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}





//=============================================
// IO Functions
//=============================================

func WriteToMem(data map[string]string, key string) {
	m := MapToJson(data)
	MemClient.Set (&memcache.Item{
		Key: key,
		Value: []byte(m),
		Expiration: 180,
	})

}

func ReadFromMem(key string) (map[string]string) {
	result := map[string]string{}
	newres, e := MemClient.Get(key)
	if (e != nil) {
		fmt.Println(e)
		return result
	}
	return JsonToMap(string(newres.Value))

}

func WriteAppState(state map[string]string) {
	WriteToMem(state, AppHeader.AppName+"_state")
}

func ReadAppState(appName string) (map[string]string) {
	r := ReadFromMem(appName+"_state")
	return r
}

func ReadAppCommand(key string)  (map[string]string) {
	return ReadFromMem(AppHeader.AppName+"_cmd")
}

func SendAppComand(appName string, data map[string]string) {
	WriteToMem(data, appName + "_cmd")
}

