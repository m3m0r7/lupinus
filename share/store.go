package share

type Procedure struct {
	Callback func(string)
}

var procedures = []Procedure{}

func AddProcedure(procedure Procedure) {
	procedures = append(procedures, procedure)
}

func ProceedProcedure(data string) {
	for _, procedure := range procedures {
		procedure.Callback(data)
	}

	// To empty
	procedures = []Procedure{}
}