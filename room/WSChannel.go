package room

import (
	"errors"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/hailongz/kk-room/proto/golang/kk"
)

const (
	WSChannelReadClosed  = 1
	WSChannelWriteClosed = 2
)

type WSChannel struct {
	id   int64
	conn *websocket.Conn
	ch   chan func()
}

func NewWSChannel(id int64, conn *websocket.Conn, size int) *WSChannel {
	v := WSChannel{}
	v.id = id
	v.conn = conn
	v.ch = make(chan func(), size)

	go func() {

		for v.conn != nil {
			fn := <-v.ch
			fn()
		}

		close(v.ch)
	}()

	return &v
}

/**
 * 发送消息
 */
func (C *WSChannel) Send(data []byte) error {
	if C.conn != nil {
		C.ch <- func() {
			if C.conn != nil {
				C.conn.WriteMessage(websocket.BinaryMessage, data)
			}
		}
		return nil
	}
	return errors.New("Not Found Connection")
}

/**
* 发送消息
 */
func (C *WSChannel) SendMessage(message *kk.Message) error {

	if C.conn != nil {

		data, err := proto.Marshal(message)

		if err != nil {
			return err
		}

		return C.Send(data)
	}

	return errors.New("Not Found Connection")

}

func (C *WSChannel) Close() {
	if C.conn != nil {
		conn := C.conn
		C.conn = nil
		conn.Close()
		C.ch <- func() {}
	}
}

/**
* 通道ID
 */
func (C *WSChannel) GetId() int64 {
	return C.id
}
