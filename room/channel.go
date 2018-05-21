package room

import (
	"github.com/hailongz/kk-room/proto/golang/kk"
)

type IChannel interface {

	/**
	 * 通道ID
	 */
	GetId() int64

	/**
	 * 发送消息
	 */
	Send(data []byte) error

	/**
	 * 发送消息
	 */
	SendMessage(message *kk.Message) error

	/**
	 * 关闭
	 */
	Close()
}
