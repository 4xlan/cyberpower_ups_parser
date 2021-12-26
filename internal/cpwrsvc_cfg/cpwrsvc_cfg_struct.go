package cpwrsvc_cfg

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

type UPSConfigSvc interface {
	GetConfig(*UPSConfig, string) error
}

type UPSConfig struct {
	Mode        bool             `yaml:"read_only"`
	ChFreqBatt  int              `yaml:"check_freq_when_on_batt"`
	ChFreqLine  int              `yaml:"check_freq_when_on_line"`
	ServerIP    string           `yaml:"ip"`
	ServerPort  string           `yaml:"port"`
	UPSResponse []UPSResponseMap `yaml:"cp_response"`
	ActionsList []Action         `yaml:"actions_list"`
}

type UPSResponseMap struct {
	Name       string `yaml:"name"`
	PrettyName string `yaml:"prettyName"`
}

type Action struct {
	Name        string `yaml:"name"`
	ExecCommand string `yaml:"exec"`
	PowerLevel  int    `yaml:"on"`
}

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

	return nil
}
