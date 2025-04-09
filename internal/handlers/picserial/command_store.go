package picserial

import (
	"encoding/json"
	"sync"
)

type commandStore struct {
	cmdMap map[string]*picCommand
	mu     sync.RWMutex
}

func newCommandStore() *commandStore {
	return &commandStore{
		cmdMap: make(map[string]*picCommand),
	}
}

func (r *commandStore) GetCommand(id string) (picCommand, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	cmd, ok := r.cmdMap[id]
	return *cmd, ok
}

func (r *commandStore) AddCommand(cmd picCommand) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.cmdMap[cmd.ID] = &cmd
}

func (r *commandStore) RemoveCommand(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.cmdMap, id)
}

type picCommand struct {
	ID   string         `json:"id"`
	Type picCommandType `json:"type"`
	Data picCommandData `json:"data"`
}

type picCommandType uint8

func (t picCommandType) MarshalJSON() ([]byte, error) {
	return json.Marshal(uint8(t))
}

const (
	picCommandTypeBatteryCharge    picCommandType = 0
	picCommandTypeBatteryDischarge picCommandType = 1
	picCommandTypeLiftMotor        picCommandType = 2
	picCommandTypeDriveMotor       picCommandType = 3
)

type picCommandData interface {
	isPICCommandData()
}

type picCommandBatteryChargeData struct {
	CurrentLimit uint16
	Enable       bool
}

func (d picCommandBatteryChargeData) MarshalJSON() ([]byte, error) {
	var temp struct {
		CurrentLimit uint16 `json:"current_limit"`
		Enable       uint8  `json:"enable"`
	}

	temp.CurrentLimit = d.CurrentLimit
	temp.Enable = boolToUint8(d.Enable)
	return json.Marshal(temp)
}

func (picCommandBatteryChargeData) isPICCommandData() {}

type picCommandBatteryDischargeData struct {
	CurrentLimit uint16
	Enable       bool
}

func (d picCommandBatteryDischargeData) MarshalJSON() ([]byte, error) {
	var temp struct {
		CurrentLimit uint16 `json:"current_limit"`
		Enable       uint8  `json:"enable"`
	}

	temp.CurrentLimit = d.CurrentLimit
	temp.Enable = boolToUint8(d.Enable)
	return json.Marshal(temp)
}

func (picCommandBatteryDischargeData) isPICCommandData() {}

type picCommandLiftMotorData struct {
	TargetPosition uint16
	Enable         bool
}

func (d picCommandLiftMotorData) MarshalJSON() ([]byte, error) {
	var temp struct {
		TargetPosition uint16 `json:"target_position"`
		Enable         uint8  `json:"enable"`
	}

	temp.TargetPosition = d.TargetPosition
	temp.Enable = boolToUint8(d.Enable)
	return json.Marshal(temp)
}

func (picCommandLiftMotorData) isPICCommandData() {}

type picCommandDriveMotorData struct {
	Direction moveDirection
	Speed     uint8
	Enable    bool
}

func (d picCommandDriveMotorData) MarshalJSON() ([]byte, error) {
	var temp struct {
		Direction uint8 `json:"direction"`
		Speed     uint8 `json:"speed"`
		Enable    uint8 `json:"enable"`
	}

	temp.Direction = uint8(d.Direction)
	temp.Speed = d.Speed
	temp.Enable = boolToUint8(d.Enable)
	return json.Marshal(temp)
}

func (picCommandDriveMotorData) isPICCommandData() {}

type moveDirection uint8

const (
	moveDirectionForward  moveDirection = 0
	moveDirectionBackward moveDirection = 1
)

func boolToUint8(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}
