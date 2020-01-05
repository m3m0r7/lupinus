package share

var (
	cameraEnv Env
)

type Env struct {
	Temp     float64
	Humidity float64
	CpuTemp  float64
	Pressure float64
}

func SetCameraEnv(temp float64, humidity float64, cpuTemp float64, pressure float64) {
	cameraEnv = Env{
		Temp:     temp,
		Humidity: humidity,
		CpuTemp:  cpuTemp,
		Pressure: pressure,
	}
}

func GetCameraEnv() *Env {
	return &cameraEnv
}
