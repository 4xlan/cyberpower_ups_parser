package main

import (
	"cyberpower_service/internal/cpwrsvc"
	"log"
	"sync"
)

func main() {

	wg := &sync.WaitGroup{}

	err := cpwrsvc.Init()
	if err != nil {
		log.Fatal(err)
	}

	wg.Add(1)
	go cpwrsvc.StartHttpServer(wg)
	wg.Wait()

	log.Println("Finished")
}
