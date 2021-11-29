package cpwrsvc

type UPSSvc interface {
	getInfo(*map[string]string)
}

type CyberPowerResponse struct {
	ModelName        string
	FirmwareNumber   string
	RatingVoltage    string
	RatingPower      string
	State            string
	PowerSupplyby    string // don't change it, until utility output is parsing into this structure
	UtilityVoltage   string
	OutputVoltage    string
	BatteryCapacity  string
	RemainingRuntime string
	Load             string
	LineInteraction  string
	TestResult       string
	LastPowerEvent   string
}
