package cpwrsvc_cfg

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
	IsShown    bool   `yaml:"isShown"`
}

type Action struct {
	Name        string `yaml:"name"`
	ExecCommand string `yaml:"exec"`
	PowerLevel  int    `yaml:"on"`
}
