package upsoperator

import (
	"cyberpower_service/internal/cpwrsvc_cfg"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"sync"
	"time"
	"unicode"
)

func (ups *UPSState) SetFreq(freq int) {
	ups.freqMtx.Lock()
	defer ups.freqMtx.Unlock()

	ups.freq = freq
}

func (ups *UPSState) GetFreq() int {
	return ups.freq
}

func (ups *UPSState) GetDate() string {
	return ups.date.Format("2006.01.02 15:04:05")
}

func (ups *UPSState) GetMaxOrder() int {
	return ups.maxOrder
}

func (ups *UPSState) Init(cfg *cpwrsvc_cfg.UPSConfig, wg *sync.WaitGroup) error {
	ups.config = cfg
	ups.wg = wg
	ups.dataChan = make(chan map[string]UPSCurrent, 1)

	tmp := map[string]UPSCurrent{}
	err := ups.getInfo(&tmp)
	if err != nil {
		return err
	}

	ups.findMaxOrder()
	fmt.Println(ups.maxOrder)

	ups.currentState = tmp
	ups.checkCurrentState()

	return nil
}

func (ups *UPSState) GetState() *map[string]UPSCurrent {
	return &ups.currentState
}

func (ups *UPSState) Listen() {
	defer ups.wg.Done()

	for {
		ups.currentState = <-ups.dataChan
		time.Sleep(time.Duration(ups.freq/2) * time.Second)
	}

}

func (ups *UPSState) Read() error {
	defer ups.wg.Done()

	for {

		tmp := map[string]UPSCurrent{}

		err := ups.getInfo(&tmp)
		if err != nil {
			return err
		}

		ups.checkCurrentState()

		ups.dataChan <- tmp
		time.Sleep(time.Duration(ups.freq) * time.Second)
	}

	return nil
}

func (ups *UPSState) getInfo(tmp *map[string]UPSCurrent) error {

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

	ups.date = time.Now()

	return nil
}

func (ups *UPSState) parseOutput(output *[]string, tmp *map[string]UPSCurrent) error {

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
				so := UPSCurrent{
					Pretty: strings.TrimSpace(parsedString[0]),
					Value:  strings.TrimSpace(parsedString[1]),
				}
				(*tmp)[ups.config.UPSResponse[index].Name] = so
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

func (ups *UPSState) checkCurrentState() {

	if ups.currentState[checkPowerField].Value != ups.lastState {
		log.Printf("Power source change detected\n- Prev: %v\n- Curr: %v",
			ups.lastState,
			ups.currentState[checkPowerField].Value)

		ups.setActualFreq()

		ups.lastState = ups.currentState[checkPowerField].Value
	}

}

func (ups *UPSState) setActualFreq() {

	if ups.currentState[checkPowerField].Value == onSitePower {
		ups.SetFreq(ups.config.ChFreqLine)
	} else {
		ups.SetFreq(ups.config.ChFreqBatt)
	}

	log.Printf("Frequency of state read has been changed to: %v\n", ups.freq)
}

func (ups *UPSState) findMaxOrder() {
	ups.maxOrder = 0

	for _, value := range ups.currentState {
		if ups.maxOrder < value.Order {
			ups.maxOrder = value.Order
		}
	}

}

func convertKey(key string) string {

	return strings.ReplaceAll(
		strings.TrimFunc(key, func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsNumber(r)
		}), " ", "")

}
