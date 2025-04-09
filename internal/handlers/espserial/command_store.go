package espserial

import (
	"encoding/json"
	"sync"
)

type commandStore struct {
	cmdMap map[string]*espCommand
	mu     sync.RWMutex
}

func newCommandStore() *commandStore {
	return &commandStore{
		cmdMap: make(map[string]*espCommand),
	}
}

func (r *commandStore) GetCommand(id string) (espCommand, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	cmd, ok := r.cmdMap[id]
	return *cmd, ok
}

func (r *commandStore) AddCommand(cmd espCommand) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.cmdMap[cmd.ID] = &cmd
}

func (r *commandStore) RemoveCommand(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.cmdMap, id)
}

type espCommand struct {
	ID   string         `json:"id"`
	Type espCommandType `json:"type"`
	Data espData        `json:"data"`
}

type espCommandType uint8

func (t espCommandType) MarshalJSON() ([]byte, error) {
	return json.Marshal(uint8(t))
}

const (
	espCommandTypeCargoDoorMotor espCommandType = 0
)

type espData interface {
	isEspData()
}

type espCargoDoorMotorData struct {
	Direction doorDirection
	Speed     uint8
	Enable    bool
}

func (espCargoDoorMotorData) isEspData() {}

func (d espCargoDoorMotorData) MarshalJSON() ([]byte, error) {
	var data struct {
		State  uint8 `json:"state"`
		Speed  uint8 `json:"speed"`
		Enable uint8 `json:"enable"`
	}

	data.State = uint8(d.Direction)
	data.Speed = d.Speed
	data.Enable = boolToUint8(d.Enable)
	return json.Marshal(data)
}

type doorDirection uint8

const (
	doorDirectionClose doorDirection = 0
	doorDirectionOpen  doorDirection = 1
)

func boolToUint8(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}
