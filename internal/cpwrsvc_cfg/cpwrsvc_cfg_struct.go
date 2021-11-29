package cpwrsvc_cfg

type UPSConfigSvc interface {
	GetConfig(*UPSConfig, string) error
}

type UPSConfig struct {
	ChFreqFail    int              `yaml:"check_freq_when_failed"`
	ChFreqNorm    int              `yaml:"check_freq_when_normal"`
	RemRunLimWarn int              `yaml:"remaining_runtime_limit_warn"`
	RemRunLinCrit int              `yaml:"remaining_runtime_limit_crit"`
	NotifWhenWarn bool             `yaml:"notify_when_warning"`
	NotifWhenCrit bool             `yaml:"notify_when_crit"`
	NotifWarnComm string           `yaml:"notification_warn_command"`
	NotifCritComm string           `yaml:"notification_crit_command"`
	ExecWarn      bool             `yaml:"execute_when_warning"`
	ExecCrit      bool             `yaml:"execute_when_crit"`
	ExecWarnComm  string           `yaml:"execute_warn_command"`
	ExecCritComm  string           `yaml:"execute_crit_command"`
	ServerIP      string           `yaml:"ip"`
	ServerPort    string           `yaml:"port"`
	UPSResponse   []UPSResponseMap `yaml:"cp_response"`
}

type UPSResponseMap struct {
	Name       string `yaml:"name"`
	PrettyName string `yaml:"prettyName"`
	IsShown    bool   `yaml:"isShown"`
	IsStatic   bool   `yaml:"isStatic"`
}
