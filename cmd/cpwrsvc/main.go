package main

import (
	"cyberpower_service/internal/cpwrsvc"
	"cyberpower_service/internal/cpwrsvc_cfg"
	"cyberpower_service/internal/upsoperator"
	"log"
	"sync"
)

const defaultConfigPath string = "config/config.yaml"

func main() {

	wg := &sync.WaitGroup{}
	globalConfig := cpwrsvc_cfg.UPSConfig{}
	ups := upsoperator.UPSState{}

	err := globalConfig.Init(defaultConfigPath)

	if err != nil {
		log.Fatal(err)
	}

	err = ups.Init(&globalConfig, wg)
	if err != nil {
		log.Fatal(err)
	}

	wg.Add(1)
	go ups.Read()

	wg.Add(1)
	go ups.Listen()

	err = cpwrsvc.Init(&globalConfig, wg, &ups)
	if err != nil {
		log.Fatal(err)
	}

	wg.Add(1)
	go ups.Listen()

	wg.Add(1)
	go cpwrsvc.StartHttpServer()
	wg.Wait()

	log.Println("Finished")
}
