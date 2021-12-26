package upsoperator

import (
	"cyberpower_service/internal/cpwrsvc_cfg"
	"os/exec"
	"strings"
	"sync"
	"time"
	"unicode"
)

func (f *Freq) SetFreq(freq int) {
	f.freqMtx.Lock()
	defer f.freqMtx.Unlock()

	f.freq = freq
}

func (f *Freq) GetFreq() int {
	return f.freq
}

func (ups *UPSState) Init(cfg *cpwrsvc_cfg.UPSConfig, wg *sync.WaitGroup) error {
	ups.config = cfg
	ups.wg = wg
	ups.dataChan = make(chan map[string]string, 1)

	return nil
}

func (ups *UPSState) GetState() *map[string]string {
	return &ups.currentState
}

func (ups *UPSState) Listen() {
	defer ups.wg.Done()

	for {
		ups.currentState = <-ups.dataChan
		time.Sleep(time.Duration(ups.freqVal.freq/2) * time.Second)
	}

}

func (ups *UPSState) Read() error {
	defer ups.wg.Done()

	for {

		tmp := map[string]string{}

		err := ups.getInfo(&tmp)
		if err != nil {
			return err
		}

		ups.dataChan <- tmp
		time.Sleep(time.Duration(ups.freqVal.freq) * time.Second)
	}

	return nil
}

func (ups *UPSState) getInfo(tmp *map[string]string) error {

	cmd := exec.Command(command, arg1)
	resp, err := cmd.Output()

	if err != nil {
		return err
	}

	convertedResp := strings.Split(string(resp), "\n")

	err = ups.parseOutput(&convertedResp, tmp)
	if err != nil {
		return err
	}

	return nil
}

func (ups *UPSState) parseOutput(output *[]string, tmp *map[string]string) error {

	for _, line := range *output {
		parsedString := strings.FieldsFunc(line, func(r rune) bool {
			return r == '.'
		})

		// If we have 2 strings, then start to parse: first arg expected to be a key for map
		// (if it exists in CyberPowerResponse struct, and the sophomore - should be a value for this key)

		if len(parsedString) > 1 {
			key := convertKey(parsedString[0])
			index := ups.isKeyExists(key)

			if index != -1 {
				(*tmp)[ups.config.UPSResponse[index].PrettyName] = strings.TrimSpace(parsedString[1])
			}
		}
	}

	return nil
}

func (ups *UPSState) isKeyExists(key string) int {

	for i, item := range ups.config.UPSResponse {
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
