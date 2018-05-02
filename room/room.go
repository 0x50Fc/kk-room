package room

import (
	"github.com/hailongz/kk-room/proto/golang/kk"
)

type IRoom interface {
	GetId() int64
	AddChannel(channel IChannel)
	RemoveChannel(channel IChannel)
	Send(message *kk.Message)
	Exit()
}

type Room struct {
	id       int64
	channels map[int64]IChannel
	ch       chan func()
	run      bool
}

func NewRoom(id int64, size int) IRoom {
	v := Room{}
	v.channels = map[int64]IChannel{}
	v.id = id
	v.ch = make(chan func(), size)
	v.run = true

	go func() {

		for v.run {
			select {
			case fn := <-v.ch:
				fn()
			}
		}

		for _, channel := range v.channels {
			channel.Close()
		}

		close(v.ch)

	}()

	return &v
}

func (R *Room) GetId() int64 {
	return R.id
}

func (R *Room) AddChannel(channel IChannel) {
	if !R.run {
		return
	}
	R.ch <- func() {
		R.channels[channel.GetId()] = channel
	}
}

func (R *Room) RemoveChannel(channel IChannel) {
	if !R.run {
		return
	}
	R.ch <- func() {
		delete(R.channels, channel.GetId())
	}
}

func (R *Room) Exit() {
	if R.run {
		R.run = false
		R.ch <- func() {}
	}
}

func (R *Room) Send(message *kk.Message) {

	if !R.run {
		return
	}

	R.ch <- func() {

		for _, channel := range R.channels {
			channel.Send(message)
		}

	}
}
