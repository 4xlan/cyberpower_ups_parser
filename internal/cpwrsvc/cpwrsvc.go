package cpwrsvc

import (
	"cyberpower_service/internal/cpwrsvc_cfg"
	"cyberpower_service/internal/upsoperator"
	"fmt"
	"log"
	"net/http"
	"sync"
)

var (
	wgLocal  *sync.WaitGroup
	config   *cpwrsvc_cfg.UPSConfig
	upsState *upsoperator.UPSState
)

func Init(extCfg *cpwrsvc_cfg.UPSConfig, wg *sync.WaitGroup, ups *upsoperator.UPSState) error {
	wgLocal = wg
	config = extCfg
	upsState = ups
	return nil
}

func StartHttpServer() {

	server := &http.Server{
		Addr: fmt.Sprintf("%s:%s", config.ServerIP, config.ServerPort),
	}

	log.Println(fmt.Sprintf("Server was started on %s:%s", config.ServerIP, config.ServerPort))

	http.HandleFunc("/getState", func(w http.ResponseWriter, r *http.Request) {

		state := upsState.GetState()

		for key, value := range *state {
			_, err := fmt.Fprintf(w, "%v = %v\n", key, value)
			if err != nil {
				log.Println(err)
			}
		}

	})

	go func() {
		defer wgLocal.Done()
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("error during websrv execution: %v", err)
		}
	}()
}
