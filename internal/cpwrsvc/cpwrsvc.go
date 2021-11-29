package cpwrsvc

import (
	"cyberpower_service/internal/cpwrsvc_cfg"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"unicode"
)

const (
	command                  = "pwrstat"
	arg1                     = "-status"
	defaultConfigPath string = "config/config.yaml"
)

var config cpwrsvc_cfg.UPSConfig

func Init() error {

	err := cpwrsvc_cfg.GetConfig(&config, defaultConfigPath)
	if err != nil {
		return err
	}

	return nil
}

func getInfo(state *map[string]string) error {

	cmd := exec.Command(command, arg1)
	resp, err := cmd.Output()

	if err != nil {
		return err
	}

	convertedResp := strings.Split(string(resp), "\n")

	err = parseOutput(&convertedResp, state)
	if err != nil {
		return err
	}

	return nil
}

func parseOutput(output *[]string, state *map[string]string) error {

	for _, line := range *output {
		parsedString := strings.FieldsFunc(line, func(r rune) bool {
			return r == '.'
		})

		// If we have 2 strings, then start to parse: first arg expected to be a key for map
		// (if it exists in CyberPowerResponse struct, and the sophomore - should be a value for this key)

		if len(parsedString) > 1 {
			key := convertKey(parsedString[0])
			index := isKeyExists(key)

			if index != -1 { // TODO: Add buffer for non-parsed values
				if config.UPSResponse[index].IsShown { // TODO: Add forced output through special parameter
					(*state)[config.UPSResponse[index].PrettyName] = strings.TrimSpace(parsedString[1])
				}
			}
		}
	}

	return nil
}

func isKeyExists(key string) int {

	for i, item := range config.UPSResponse {
		if strings.ToLower(item.Name) == strings.ToLower(key) {
			return i
		}
	}

	return -1
}

func convertKey(key string) string {

	return strings.ReplaceAll(
		strings.TrimFunc(key, func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsNumber(r)
		}), " ", "")

}

func StartHttpServer(wg *sync.WaitGroup) {

	server := &http.Server{
		Addr: fmt.Sprintf("%s:%s", config.ServerIP, config.ServerPort),
	}

	log.Println(fmt.Sprintf("Server was started on %s:%s", config.ServerIP, config.ServerPort))

	state := &map[string]string{}

	http.HandleFunc("/getState", func(w http.ResponseWriter, r *http.Request) {

		err := getInfo(state)

		if err != nil {
			log.Fatal(err)
		}

		for key, value := range *state {
			fmt.Fprintf(w, "%v = %v\n", key, value)
		}

	})

	go func() {
		defer wg.Done()
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("error during websrv execution: %v", err)
		}
	}()
}
