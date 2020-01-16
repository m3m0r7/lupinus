package share

type Procedure struct {
	Callback func([]byte)
}

var procedures = map[string][]Procedure{}

func AddProcedure(key string, procedure Procedure) {
	if _, ok := procedures[key]; !ok {
		procedures[key] = []Procedure{}
	}
	procedures[key] = append(procedures[key], procedure)
}

func ProceedProcedure(key string, data []byte) {
	if _, ok := procedures[key]; !ok {
		return
	}
	for _, procedure := range procedures[key] {
		procedure.Callback(data)
	}

	// To empty
	procedures[key] = []Procedure{}
}
