package go_server

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func ReadFile(filename string) ([]byte, error) {
	home, _ := os.UserHomeDir()
	flag.Parse()
	if _, err := os.Stat(home + "/" + filename); os.IsNotExist(err) {
		_, err = os.Create(home + "/" + filename)
		if err != nil {
			log.Println(err)
		}
	}
	fmt.Println(home)
	return ioutil.ReadFile(home + "/" + filename)
}

func ReadState(file []byte) {
	data := strings.Split(string(file), "\n")

	hashDec, err := base64.StdEncoding.DecodeString(data[0])
	if err != nil {
		log.Println(err)
		return
	}
	StateHash = string(hashDec)

	stateDec, err := base64.StdEncoding.DecodeString(data[1])
	err = json.Unmarshal(stateDec, &state)
	if err != nil {
		log.Println(err)
	}
}

func EncodeState() string {
	mapMutex.Lock()
	jsonstate, err := json.Marshal(state)
	mapMutex.Unlock()
	if err != nil {
		log.Fatal(err)
	}
	return base64.StdEncoding.EncodeToString(jsonstate)
}

func StateFile() *os.File {
	home, _ := os.UserHomeDir()
	mapMutex.Lock()
	file, err := os.Create(home + "/" + state.StatePath)
	mapMutex.Unlock()
	if err != nil {
		log.Println(err)
		return nil
	}
	return file
}

func WriteToDisk() {
	StateHashEnc := base64.StdEncoding.EncodeToString([]byte(StateHash))
	stateEnc := EncodeState()

	file := StateFile()
	if file == nil {
		return
	}

	_, err := file.WriteString(StateHashEnc + "\n" + stateEnc)
	if err != nil {
		log.Println(err)
		return
	}

	if err := file.Close(); err != nil {
		log.Println(err)
		return
	}
}
