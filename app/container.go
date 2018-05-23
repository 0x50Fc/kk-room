package app

import "time"

type IContainer interface {
	Get(id int64, ch chan bool) int64
	Open(id int64, path string, query map[string]string, expires time.Duration, ch chan bool) int64
	Exit(id int64)
	RunCommand(id int64, command []byte)
	List(ch chan bool) []int64
}

type Container struct {
	apps   map[int64]IApplication
	ch     chan func()
	cb     chan int64
	autoId int64
}

func NewContainer() *Container {

	v := Container{}
	v.apps = map[int64]IApplication{}
	v.ch = make(chan func(), 204800)
	v.cb = make(chan int64, 2048)
	v.autoId = 0

	go func() {

		for {
			select {
			case fn := <-v.ch:
				fn()
				break
			case id := <-v.cb:
				delete(v.apps, id)
				break
			}
		}

	}()

	return &v

}

func (C *Container) Get(id int64, ch chan bool) int64 {

	var app IApplication = nil

	C.ch <- func() {

		app = C.apps[id]

		ch <- true
	}

	<-ch

	if app != nil {
		return id
	}

	return 0
}

func (C *Container) Open(id int64, path string, query map[string]string, expires time.Duration, ch chan bool) int64 {

	var app IApplication = nil

	C.ch <- func() {

		app = C.apps[id]

		if app != nil {
			ch <- true
			return
		}

		if id == 0 {
			for {
				C.autoId = C.autoId + 1
				id = C.autoId
				if C.apps[id] == nil {
					break
				}
			}
		}

		app = Open(id, path, query, expires, C.cb)

		C.apps[app.GetId()] = app

		ch <- true
	}

	<-ch

	return app.GetId()

}

func (C *Container) Exit(id int64) {

	C.ch <- func() {

		app := C.apps[id]

		if app != nil {
			app.Exit()
			delete(C.apps, id)
		}

	}

}

func (C *Container) RunCommand(id int64, command []byte) {

	C.ch <- func() {

		app := C.apps[id]

		if app != nil {
			app.RunCommand(command)
		}

	}

}

func (C *Container) List(ch chan bool) []int64 {

	ids := []int64{}

	C.ch <- func() {

		for id, _ := range C.apps {
			ids = append(ids, id)
		}

		ch <- true
	}

	<-ch

	return ids
}
