package main

import (
	"io/ioutil"
)

func configParse() []byte {
	data, err := ioutil.ReadFile("settings.json")
	if err != nil {
		panic(err)
	}
	return data
}
