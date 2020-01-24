package share

import "sync"

type Procedure struct {
	Callback func([]byte)
}

var mutex = sync.Mutex{}
var procedures = map[string][]Procedure{}

func AddProcedure(key string, procedure Procedure) {
  mutex.Lock()
  defer mutex.Unlock()
	if _, ok := procedures[key]; !ok {
		procedures[key] = []Procedure{}
	}
	procedures[key] = append(procedures[key], procedure)
}

func ProceedProcedure(key string, data []byte) {
  mutex.Lock()
  defer mutex.Unlock()
	if _, ok := procedures[key]; !ok {
		return
	}
	for _, procedure := range procedures[key] {
		procedure.Callback(data)
	}

	// To empty
	procedures[key] = []Procedure{}
}
