package upsoperator

import (
	"cyberpower_service/internal/cpwrsvc_cfg"
	"sync"
)

// Here we just keep parsed text and change it only through channel

const (
	command = "pwrstat"
	arg1    = "-status"
)

//type CyberPowerResponse struct {
//	ModelName        string
//	FirmwareNumber   string
//	RatingVoltage    string
//	RatingPower      string
//	State            string
//	PowerSupplyby    string // don't change it, until utility output is parsing into this structure
//	UtilityVoltage   string
//	OutputVoltage    string
//	BatteryCapacity  string
//	RemainingRuntime string
//	Load             string
//	LineInteraction  string
//	TestResult       string
//	LastPowerEvent   string
//}

type Freq struct {
	freq    int
	freqMtx sync.Mutex
}

type UPSState struct {
	currentState map[string]string
	freqVal      Freq
	config       *cpwrsvc_cfg.UPSConfig
	wg           *sync.WaitGroup
	dataChan     chan map[string]string
}
