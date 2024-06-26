package edge

import (
	"context"
	"net"

	"github.com/singchia/geminio"
	"github.com/singchia/geminio/client"
)

// RPCer is method oriented
type RPCer interface {
	NewRequest(data []byte) geminio.Request

	Call(ctx context.Context, method string, req geminio.Request) (geminio.Response, error)
	CallAsync(ctx context.Context, method string, req geminio.Request, ch chan *geminio.Call) (*geminio.Call, error)
	Register(ctx context.Context, method string, rpc geminio.RPC) error
}

// Messager is topic oriented
type Messager interface {
	NewMessage(data []byte) geminio.Message

	// Publish a message to specific topic
	Publish(ctx context.Context, topic string, msg geminio.Message) error
	// Publish async a message to specific topic
	PublishAsync(ctx context.Context, topic string, msg geminio.Message, ch chan *geminio.Publish) (*geminio.Publish, error)
	Receive(ctx context.Context) (geminio.Message, error)
}

type RPCMessager interface {
	RPCer
	Messager
}

// Stream multiplexer
type Multiplexer interface {
	// Open a stream to a specific service
	OpenStream(serviceName string) (geminio.Stream, error)
	AcceptStream() (geminio.Stream, error)
	ListStreams() []geminio.Stream
}

type Edge interface {
	// Edge can directly Publish Message or Call RPC
	RPCMessager

	// Edge can manage(create, list...) streams from or to a Service
	Multiplexer

	// Edge is a net.Listener, actually it's wrapper of Multiplexer
	// The Accept function is a wrapper from AccetpStream
	// The Addr is a wrapper from LocalAddr
	net.Listener

	// Meta
	EdgeID() uint64

	Close() error
}

type Dialer func() (net.Conn, error)

func NewEdge(dialer Dialer, opts ...EdgeOption) (Edge, error) {
	return newRetryEdgeEnd(client.Dialer(dialer), opts...)
}

func NewNoRetryEdge(dialer Dialer, opts ...EdgeOption) (Edge, error) {
	return newEdgeEnd(client.Dialer(dialer), opts...)
}
