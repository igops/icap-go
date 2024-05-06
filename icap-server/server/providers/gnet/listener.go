package gnet

import (
	"github.com/panjf2000/gnet"
	"github.com/rs/zerolog"
	"igops.me/icap/server"
	"igops.me/icap/server/utils"
)

type Listener struct {
	*gnet.EventServer
	handler server.Handler
	log     zerolog.Logger
}

func (s *Listener) OnInitComplete(srv gnet.Server) (action gnet.Action) {
	s.log.Debug().Msgf("HTTP server is listening on %s (multi-cores: %t, event-loops: %d)\n",
		srv.Addr.String(), srv.Multicore, srv.NumEventLoop)
	return
}

func (s *Listener) OnOpened(c gnet.Conn) ([]byte, gnet.Action) {
	c.SetContext(NewCodec())
	return nil, gnet.None
}

func (s *Listener) React(data []byte, conn gnet.Conn) (out []byte, action gnet.Action) {
	codec := conn.Context().(*Codec)
	body, err := codec.parseBody(data)
	if err != nil {
		return []byte("500 Error"), gnet.Close
	}

	switch codec.getMethod() {
	case utils.MethodREQMOD:
		out, err = s.handler.OnREQMOD(body)
	case utils.MethodRESPMOD:
		out, err = s.handler.OnRESPMOD(body)
	default:
		return []byte("500 Error"), gnet.Close
	}

	return codec.buildResponse(out), gnet.Close
}
