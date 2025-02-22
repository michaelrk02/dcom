package dcom

import (
	"bytes"
	"errors"
	"io"
	"net"

	"github.com/google/uuid"
)

var (
	ErrClientConnectionInvoke error = errors.New("client connection invoke error")
)

type ProxyConnection struct {
	address string
}

func NewProxyConnection(address string) *ProxyConnection {
	return &ProxyConnection{
		address: address,
	}
}

func (conn *ProxyConnection) Invoke(message string, body io.Reader) (io.Reader, error) {
	tcp, err := net.Dial("tcp", conn.address)
	if err != nil {
		return nil, errors.Join(ErrClientConnectionInvoke, err)
	}
	defer tcp.Close()

	m := NewDefaultMarshaler(tcp)

	err = m.WriteString(message)
	if err != nil {
		return nil, errors.Join(ErrClientConnectionInvoke, err)
	}

	var bodyBuffer bytes.Buffer
	_, err = bodyBuffer.ReadFrom(body)
	if err != nil {
		return nil, errors.Join(ErrClientConnectionInvoke, err)
	}

	_, err = bodyBuffer.WriteTo(tcp)
	if err != nil {
		return nil, errors.Join(ErrClientConnectionInvoke, err)
	}

	var respBuffer bytes.Buffer
	_, err = respBuffer.ReadFrom(tcp)
	if err != nil {
		return nil, errors.Join(ErrClientConnectionInvoke, err)
	}

	return bytes.NewReader(respBuffer.Bytes()), nil
}

func (conn *ProxyConnection) InvokeObject(clsid, instanceID uuid.UUID, method string, params io.Reader) (io.Reader, error) {
	var bodyBuffer bytes.Buffer

	m := NewDefaultMarshaler(&bodyBuffer)

	err := m.WriteUUID(clsid)
	if err != nil {
		return nil, errors.Join(ErrClientConnectionInvoke, err)
	}

	err = m.WriteUUID(instanceID)
	if err != nil {
		return nil, errors.Join(ErrClientConnectionInvoke, err)
	}

	err = m.WriteString(method)
	if err != nil {
		return nil, errors.Join(ErrClientConnectionInvoke, err)
	}

	_, err = bodyBuffer.ReadFrom(params)
	if err != nil {
		return nil, errors.Join(ErrClientConnectionInvoke, err)
	}

	return conn.Invoke("Object.invoke", bytes.NewReader(bodyBuffer.Bytes()))
}
