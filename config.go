package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

const configFile = "./config.json"

func readConfig() *config {
	configFileBytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal("unable to read config file", err)
	}
	c := &config{}
	if err := json.Unmarshal(configFileBytes, c); err != nil {
		log.Fatal("unable to read config file", err)
	}
	return c
}

type config struct {
	BindAddress          string          `json:"bind_address"`
	FolderRescanInterval int             `json:"folder_rescan_interval"`
	FolderSizes          []folderSizeCfg `json:"folder_sizes"`
	HornetNodes          []hornetcfg     `json:"hornet_nodes"`
	IRINodes             []iricfg        `json:"iri_nodes"`
}

type folderSizeCfg struct {
	ID   string `json:"id"`
	Path string `json:"path"`
}

type hornetcfg struct {
	ID   string `json:"id"`
	Host string `json:"host"`
}

type iricfg struct {
	ID   string `json:"id"`
	Host string `json:"host"`
}
