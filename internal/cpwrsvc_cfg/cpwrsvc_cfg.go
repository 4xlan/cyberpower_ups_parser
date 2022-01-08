package cpwrsvc_cfg

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

func (uc *UPSConfig) Init(cfgPath string) error {

	err := uc.GetConfig(cfgPath)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UPSConfig) GetConfig(defaultConfigPath string) error {
	// Calculating the ABS path for config file
	absp, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}

	configFilePath := fmt.Sprintf("%s/%s", absp, defaultConfigPath)

	// Open file and read it
	data, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return err
	}

	// Reading cfg file
	err = yaml.Unmarshal([]byte(data), uc)
	if err != nil {
		return err
	}

	uc.findMaxOrder()

	return nil
}

func (uc *UPSConfig) findMaxOrder() {
	uc.MaxOrder = 0

	for _, value := range uc.UPSResponse {
		if uc.MaxOrder < value.Order {
			uc.MaxOrder = value.Order
		}
	}

}
