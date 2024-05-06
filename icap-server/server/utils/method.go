package utils

type Method int

const (
	MethodUnknown Method = iota
	MethodREQMOD
	MethodRESPMOD
)

func ParseMethod(str string) Method {
	switch str {
	case "REQMOD":
		return MethodREQMOD
	case "RESPMOD":
		return MethodRESPMOD
	default:
		return MethodUnknown
	}
}
