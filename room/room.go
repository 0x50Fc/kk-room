package room

import (
	"log"
	"strconv"
	"strings"

	"github.com/golang/protobuf/proto"
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

		log.Println("[ROOM] [RUN]", id)

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

		log.Println("[ROOM] [EXIT]", id)
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

		data, err := proto.Marshal(message)

		if err != nil {
			log.Println("[ROOM] [ERROR]", err)
		} else {

			var ids map[int64]bool = nil

			if message.To != "" {
				ids = map[int64]bool{}
				for _, v := range strings.Split(message.To, ",") {
					id, _ := strconv.ParseInt(v, 10, 64)
					ids[id] = true
				}
			}

			for _, channel := range R.channels {
				if ids == nil || ids[channel.GetId()] {
					channel.Send(data)
				}
			}
		}

	}
}
