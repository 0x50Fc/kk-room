package room

import (
	"github.com/gorilla/websocket"
	"github.com/golang/protobuf/proto"
	"github.com/hailongz/kk-room/proto/golang/kk"
	"errors"
)

const (
	WSChannelReadClosed = 1
	WSChannelWriteClosed = 2
)

type WSChannel struct {
	id    int64
	conn  *websocket.Conn
}

func NewWSChannel(id int64, conn *websocket.Conn) *WSChannel {
	v := WSChannel{}
	v.id = id
	v.conn = conn
	return &v
}

/**
 * 发送消息
 */
func (C *WSChannel) Send(data []byte) error {
	if C.conn != nil {
		return C.conn.WriteMessage(websocket.BinaryMessage, data)
	}
	return errors.New("Not Found Connection")
}

/**
* 发送消息
*/
func (C *WSChannel) SendMessage(message *kk.Message) error {

	if C.conn != nil {

		data,err := proto.Marshal(message)

		if( err != nil) {
			return err
		} 
		return C.conn.WriteMessage(websocket.BinaryMessage, data)
	}

	return errors.New("Not Found Connection")

}


func (C *WSChannel) Close() {
	if(C.conn != nil) {
		conn := C.conn
		C.conn = nil;
		conn.Close();
	}
}

/**
* 通道ID
*/
func (C *WSChannel) GetId() int64 {
	return C.id
}
