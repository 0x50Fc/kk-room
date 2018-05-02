package room

type IServer interface {
	RoomCreate() IRoom
	RoomRemove(id int64)
	RoomGet(id int64) IRoom
	AutoId() int64
	GetRoomCount() int
	Exit()
}

type Server struct {
	id        int64
	rooms     map[int64]IRoom
	ch        chan func()
	run       bool
	roomCount int
}

func NewServer(size int) IServer {
	v := Server{}
	v.id = 0
	v.rooms = map[int64]IRoom{}
	v.ch = make(chan func(), size)
	v.run = true

	go func() {

		for v.run {
			select {
			case fn := <-v.ch:
				fn()
			}
		}

		for _, room := range v.rooms {
			room.Exit()
		}

		close(v.ch)
	}()

	return &v
}

func (S *Server) RoomCreate() IRoom {
	var room IRoom = nil

	ch := make(chan bool)

	S.ch <- func() {

		S.id = S.id + 1

		room = NewRoom(S.id, 20480)

		S.rooms[room.GetId()] = room
		S.roomCount = S.roomCount + 1

		ch <- true
	}

	<-ch

	close(ch)

	return room
}

func (S *Server) RoomRemove(id int64) {

	S.ch <- func() {

		room := S.rooms[id]

		if room != nil {
			room.Exit()
			delete(S.rooms, id)
			S.roomCount = S.roomCount - 1
		}

	}
}

func (S *Server) RoomGet(id int64) IRoom {
	var room IRoom = nil

	ch := make(chan bool)

	S.ch <- func() {
		room = S.rooms[id]
		ch <- true
	}

	<-ch

	close(ch)

	return room
}

func (S *Server) AutoId() int64 {

	var id int64 = 0

	ch := make(chan bool)

	S.ch <- func() {
		S.id = S.id + 1
		id = S.id
		ch <- true
	}

	<-ch

	close(ch)

	return id
}

func (S *Server) Exit() {
	S.ch <- func() {
		S.run = false
	}
}

func (S *Server) GetRoomCount() int {
	return S.roomCount
}
