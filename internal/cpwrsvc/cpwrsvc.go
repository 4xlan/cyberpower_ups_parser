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
		showCurrentState(&w)
	})

	http.HandleFunc("/runningConf", func(w http.ResponseWriter, r *http.Request) {
		showRunningConf(&w)
	})

	http.HandleFunc("/forceSet", func(w http.ResponseWriter, r *http.Request) {
		//TODO: Get 1/0 as argument of request
		// 1: Battery
		// 0: On-site
		// Notification: (has been manually moved to state)
		// Will be done after main logic
	})

	go func() {
		defer wgLocal.Done()
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("error during websrv execution: %v", err)
		}
	}()
}

func showCurrentState(w *http.ResponseWriter) {
	state := upsState.GetState()

	_, err := fmt.Fprintf(*w, "Last capture: %v\n---\n%v", upsState.GetDate(), state)
	if err != nil {
		log.Println(err)
	}

	if state == nil {
		_, err := fmt.Fprintf(*w, "No data.")
		if err != nil {
			log.Println(err)
		}
	}
}

func showRunningConf(w *http.ResponseWriter) {
	_, err := fmt.Fprintf(*w, "%+v", config.UPSResponse)
	if err != nil {
		log.Println(err)
	}
}
