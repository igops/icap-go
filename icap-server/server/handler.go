package server

type Handler interface {
	OnREQMOD(request []byte) (response []byte, err error)
	OnRESPMOD(request []byte) (response []byte, err error)
}
