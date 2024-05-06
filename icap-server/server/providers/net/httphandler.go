package net

import (
	"igops.me/icap/server"
	"igops.me/icap/server/utils"
	"io"
	"net/http"
)

type HTTPHandler struct {
	ICAPHandler server.Handler
}

func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	method := utils.ParseMethod(r.Method)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var out []byte
	switch method {
	case utils.MethodREQMOD:
		out, err = h.ICAPHandler.OnREQMOD(body)
	case utils.MethodRESPMOD:
		out, err = h.ICAPHandler.OnRESPMOD(body)
	default:
		http.Error(w, "500 Error", http.StatusInternalServerError)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(out)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
