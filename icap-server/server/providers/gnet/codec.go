package gnet

import (
	"bytes"
	"github.com/evanphx/wildcat"
	"igops.me/icap/server/utils"
	"io"
	"net/http"
)

type Codec struct {
	parser *wildcat.HTTPParser
}

func NewCodec() *Codec {
	return &Codec{parser: wildcat.NewHTTPParser()}
}

func (c *Codec) parseBody(data []byte) (body []byte, err error) {
	headerOffset, err := c.parser.Parse(data)
	bodyLen := int(c.parser.ContentLength())
	if bodyLen == -1 {
		bodyLen = 0
	}
	return data[headerOffset : headerOffset+bodyLen], err
}

func (c *Codec) getMethod() utils.Method {
	return utils.ParseMethod(string(c.parser.Method))
}

func (c *Codec) buildResponse(data []byte) []byte {
	t := &http.Response{
		Status:        "200 OK",
		StatusCode:    200,
		Proto:         "HTTP/1.1",
		Body:          io.NopCloser(bytes.NewReader(data)),
		ContentLength: int64(len(data)),
		//Header:        make(http.Header),
	}

	buff := bytes.NewBuffer(nil)
	t.Write(buff)
	return buff.Bytes()
}
