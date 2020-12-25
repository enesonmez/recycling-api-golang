package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func (conf *Config) ConfigRead() {
	byteValue, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	json.Unmarshal(byteValue, &conf)
}
