package dcom

import (
	"bytes"
	"errors"
	"log"
	"net"
	"sync/atomic"
)

var (
	ErrStubListen error = errors.New("stub listen error")
	ErrStubAccept error = errors.New("stub accept error")
)

type StubConnection struct {
	l              *log.Logger
	address        string
	threadingModel ThreadingModel
	registry       *ServerRegistry

	running atomic.Bool
}

func NewStubConnection(
	l *log.Logger,
	address string,
	threadingModel ThreadingModel,
	reg *ServerRegistry,
) *StubConnection {
	conn := &StubConnection{
		l:              l,
		address:        address,
		threadingModel: threadingModel,
		registry:       reg,
	}
	conn.running.Store(true)
	return conn
}

func (conn *StubConnection) Listen() error {
	lis, err := net.Listen("tcp", conn.address)
	if err != nil {
		return errors.Join(ErrStubListen, err)
	}

	conn.l.Printf("listening at %s", conn.address)

	for conn.running.Load() {
		err = func() error {
			c, err := lis.Accept()
			if err != nil {
				return errors.Join(ErrStubAccept, err)
			}

			handle := func(tcp net.Conn) {
				var resp bytes.Buffer
				defer func() {
					_, err = tcp.Write(resp.Bytes())
					Assert(err)

					conn.l.Printf("sent %d bytes for the response", resp.Len())

					tcp.Close()
				}()

				in := NewDefaultUnmarshaler(tcp)

				message, err := in.ReadString()
				if err != nil {
					conn.l.Print(Describe(err))
					return
				}

				if message == "Object.invoke" {
					out := NewDefaultMarshaler(&resp)

					clsid, err := in.ReadUUID()
					if err != nil {
						conn.l.Print(Describe(err))
						Assert(out.WriteError(err))
						return
					}

					instanceID, err := in.ReadUUID()
					if err != nil {
						conn.l.Print(Describe(err))
						Assert(out.WriteError(err))
						return
					}

					method, err := in.ReadString()
					if err != nil {
						conn.l.Print(Describe(err))
						Assert(out.WriteError(err))
						return
					}

					obj, err := conn.registry.ResolveObject(clsid, instanceID)
					if err != nil {
						conn.l.Print(Describe(err))
						Assert(out.WriteError(err))
						return
					}

					stub, err := conn.registry.CreateStub(obj)
					if err != nil {
						conn.l.Print(Describe(err))
						Assert(out.WriteError(err))
						return
					}

					conn.l.Printf("got invoke message: class=%s instance=%s method=%s", clsid.String(), instanceID.String(), method)

					stub.Execute(method, in, out)
				}
			}

			if conn.threadingModel == ThreadingModelSingle {
				handle(c)
			} else if conn.threadingModel == ThreadingModelMultiple {
				go handle(c)
			}

			return nil
		}()
		if err != nil {
			conn.l.Print(Describe(err))
		}
	}

	return nil
}

func (conn *StubConnection) Close() {
	conn.running.Store(false)
}
