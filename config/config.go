package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var (
	// Public variables
	Token         string
	Masterchannel string
	KategoriId    string

	// Private variables
	config *configStruct
)

type configStruct struct {
	Token         string `json:"Token"`
	Masterchannel string `json:"Masterchannel"`
	KategoriId    string `json:"KategoriId"`
}

func ReadConfig() error {
	fmt.Println("Reading config file...")

	file, err := ioutil.ReadFile("./config.json")

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println(string(file))

	err = json.Unmarshal(file, &config)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	Token = config.Token
	Masterchannel = config.Masterchannel
	KategoriId = config.KategoriId

	return nil
}
