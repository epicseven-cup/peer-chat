package client

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"log/slog"
	"net"
	"os"
	"peerchat/internal/protobuf/pb"
	"sync"
	"time"
)

type Client struct {
	addr     net.Addr
	messages []*pb.Message
	status   chan bool
	mutex    sync.Mutex
	slog     slog.Logger
}

func NewClient(addr net.Addr) *Client {
	return &Client{
		addr:     addr,
		messages: []*pb.Message{},
	}
}

func (c *Client) Start() {
	listen, err := net.Listen("tcp", c.addr.String())

	if err != nil {
		c.slog.Error(err.Error())
		return
	}

	for {
		// go routine
		go func() {

			select {
			case <-c.status:
				c.Close()
				return
			default:
				conn, err := listen.Accept()

				// handel error but continue
				if err != nil {
					c.slog.Error(err.Error())
					return
				}

				var unmarshalMessage []byte
				bytesRead, err := conn.Read(unmarshalMessage)

				if err != nil {
					c.slog.Error(err.Error())
					return
				}

				c.slog.Info(fmt.Sprintf("received message from %s", conn.RemoteAddr().String()))
				c.slog.Info(fmt.Sprintf("message size: %d", bytesRead))

				message := &pb.Message{}
				err = proto.Unmarshal(unmarshalMessage, message)
				if err != nil {
					c.slog.Error(err.Error())
					return
				}

				// locking mutex when accessing the slice
				c.mutex.Lock()
				defer c.mutex.Unlock()
				c.messages = append(c.messages, message)
			}
		}()
	}
}

func (c *Client) GetMessages() []*pb.Message {
	return c.messages
}

func (c *Client) AddMessage(m *pb.Message) {
	c.messages = append(c.messages, m)
}

func (c *Client) SendMessage(remoteAddress net.Addr, m *pb.Message) {
	conn, err := net.DialTimeout("tcp", remoteAddress.String(), time.Second*15)

	// Closing connection after finishing dialing
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			c.slog.Error(err.Error())
		}
	}(conn)

	if err != nil {
		c.slog.Error(err.Error())
		return
	}

	marshalData, err := proto.Marshal(m)
	if err != nil {
		c.slog.Error(err.Error())
		return
	}

	writtenBytes, err := conn.Write(marshalData)
	if err != nil {
		c.slog.Error(err.Error())
		return
	}

	c.slog.Info(fmt.Sprintf("Sent [%d] bytes", writtenBytes), remoteAddress.String())
	// There should be a way to send back request when the request was successfully from the other client
}

func (c *Client) Close() {
	os.Exit(0)
}
