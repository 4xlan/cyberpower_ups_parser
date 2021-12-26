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

	go func() {
		defer wgLocal.Done()
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("error during websrv execution: %v", err)
		}
	}()
}

func showCurrentState(w *http.ResponseWriter) {
	state := upsState.GetState()

	_, err := fmt.Fprintf(*w, "Last capture: %v\n", upsState.GetDate())
	if err != nil {
		log.Println(err)
	}

	if state == nil {
		_, err := fmt.Fprintf(*w, "No data.")
		if err != nil {
			log.Println(err)
		}
	} else {
		out := sortOutput(state, upsState.GetMaxOrder())
		_, err := fmt.Fprintf(*w, out)
		if err != nil {
			log.Println(err)
		}
	}
}

func showRunningConf(w *http.ResponseWriter) {
	_, err := fmt.Fprintf(*w, "%v", upsState.GetFreq())
	if err != nil {
		log.Println(err)
	}
}

func sortOutput(state *map[string]upsoperator.UPSCurrent, max int) string {
	tmp := ""
	iter := 1
	itsDone := false

	for itsDone {

		for _, value := range *state {
			if value.Order == iter {
				tmp += fmt.Sprintf("%v: %v\n", value.Pretty, value.Value)
				log.Printf("%v: %v: %v", iter, value.Pretty, value.Value)
				iter += 1
			}
		}

		if iter > max {
			itsDone = true
		}
	}

	return tmp
}
