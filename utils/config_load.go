package utils

import (
	"io/ioutil"
)

func ConfigParse() []byte {
	data, err := ioutil.ReadFile("settings.json")
	if err != nil {
		panic(err)
	}
	return data
}
