package upsoperator

import (
	"cyberpower_service/internal/cpwrsvc_cfg"
	"sync"
	"time"
)

// Here we just keep parsed text and change it only through channel

const (
	command         = "pwrstat"
	arg1            = "-status"
	onSitePower     = "Utility Power"
	checkPowerField = "powersupplyby"
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

type UPSState struct {
	date         time.Time
	currentState map[string]UPSCurrent
	freq         int
	freqMtx      sync.Mutex
	config       *cpwrsvc_cfg.UPSConfig
	wg           *sync.WaitGroup
	dataChan     chan map[string]UPSCurrent
	lastState    string
}

type UPSCurrent struct {
	Pretty string
	Value  string
}
