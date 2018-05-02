package room

import (
	"github.com/hailongz/kk-room/proto/golang/kk"
)

type ChannelStatus int32

const (
	ChannelStatusNone         ChannelStatus = 0
	ChannelStatusConnected    ChannelStatus = 1
	ChannelStatusDisconnected ChannelStatus = 2
	ChannelStatusFail         ChannelStatus = 3
)

type IChannel interface {

	/**
	 * 通道ID
	 */
	GetId() int64

	/**
	 * 通道状态
	 */
	GetStatus() ChannelStatus

	/**
	 * 获取错误
	 */
	GetError() error

	/**
	 * 发送消息
	 */
	Send(message *kk.Message)

	/**
	 * 读取消息通道
	 */
	ReadChannel() chan *kk.Message

	/**
	 * 关闭通道
	 */
	CloseChannel() chan bool
	/**
	 * 关闭
	 */
	Close()
}
