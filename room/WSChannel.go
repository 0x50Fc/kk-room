package room

import (
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/hailongz/kk-room/proto/golang/kk"
)

type WSChannel struct {
	id     int64
	status ChannelStatus
	err    error
	w      chan *kk.Message
	r      chan *kk.Message
	c      chan bool
}

func NewWSChannel(id int64, conn *websocket.Conn, size int32) *WSChannel {
	v := WSChannel{}
	v.id = id
	v.status = ChannelStatusConnected
	v.w = make(chan *kk.Message, size)
	v.r = make(chan *kk.Message, size)
	v.c = make(chan bool)

	go func() {

		run := true

		for run {

			select {
			case message := <-v.w:

				if message == nil {
					continue
				}

				data, err := proto.Marshal(message)

				if err != nil {
					v.status = ChannelStatusFail
					v.err = err
					run = false
					close(v.c)
					break
				}

				err = conn.WriteMessage(websocket.BinaryMessage, data)

				if err != nil {
					if v.status == ChannelStatusConnected {
						v.status = ChannelStatusFail
						v.err = err
						close(v.c)
					}
					run = false
					close(v.c)
					break
				}

			case <-v.c:
				run = false
				break
			}

		}

		close(v.w)

		conn.Close()
	}()

	go func() {

		run := true

		for run {

			mType, data, err := conn.ReadMessage()

			if err != nil {
				v.status = ChannelStatusFail
				v.err = err
				run = false
				close(v.c)
				break
			}

			if mType == websocket.BinaryMessage {

				message := kk.Message{}

				err = proto.Unmarshal(data, &message)

				if err != nil {
					v.status = ChannelStatusFail
					v.err = err
					run = false
					close(v.c)
					break
				}

				select {
				case v.r <- &message:
					break
				case <-v.c:
					run = false
					break
				}
			}

		}

		close(v.r)

	}()

	return &v
}

/**
 * 发送消息
 */
func (C *WSChannel) Send(message *kk.Message) {
	if C.status == ChannelStatusConnected {
		C.w <- message
	}
}

/**
 * 读取消息通道
 */
func (C *WSChannel) ReadChannel() chan *kk.Message {
	return C.r
}

/**
 * 关闭通道
 */
func (C *WSChannel) CloseChannel() chan bool {
	return C.c
}

/**
 * 通道状态
 */
func (C *WSChannel) GetStatus() ChannelStatus {
	return C.status
}

/**
 * 获取错误
 */
func (C *WSChannel) GetError() error {
	return C.err
}

/**
 * 关闭
 */
func (C *WSChannel) Close() {
	if C.status == ChannelStatusConnected {
		C.status = ChannelStatusDisconnected
		close(C.c)
	}
}

/**
* 通道ID
 */
func (C *WSChannel) GetId() int64 {
	return C.id
}
