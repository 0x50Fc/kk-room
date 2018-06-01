package app

import (
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

type IApplication interface {
	GetId() int64
	RunCommand(command []byte)
	Exit()
}

type CommandStream struct {
	ch chan []byte
}

func NewCommandStream() *CommandStream {
	v := CommandStream{}
	v.ch = make(chan []byte, 2048)
	return &v
}

func (C *CommandStream) Read(data []byte) (int, error) {

	if C.ch == nil {
		return 0, io.EOF
	}

	b := <-C.ch

	if b != nil {
		var size int = len(data)
		var n int = len(b)
		if n > size {
			n = size
		}
		copy(data, b[:n])
		return n, nil
	} else {
		return 0, io.EOF
	}

}

func (C *CommandStream) Run(command []byte) {
	if C.ch != nil {
		C.ch <- command
	}
}

func (C *CommandStream) Close() {
	if C.ch != nil {
		ch := C.ch
		C.ch = nil
		close(ch)
	}
}

type Application struct {
	id      int64
	cmd     *exec.Cmd
	stdin   *CommandStream
	running bool
}

func (A *Application) Write(p []byte) (n int, err error) {
	log.Println(string(p))
	return len(p), nil
}

func Open(id int64, path string, query map[string]string, expires time.Duration, cb chan int64) *Application {
	v := Application{}
	v.id = id
	v.stdin = NewCommandStream()
	v.running = false

	a, _ := filepath.Abs("./bin/kk-app")
	p, _ := filepath.Abs(path)

	args := []string{"-id", strconv.FormatInt(id, 10)}

	for key, value := range query {
		args = append(args, "-"+key)
		args = append(args, value)
	}

	v.cmd = exec.Command(a, args...)
	v.cmd.Dir = p
	v.cmd.Stderr = os.Stderr
	v.cmd.Stdout = os.Stdout
	v.cmd.Stdin = v.stdin

	go func() {

		err := v.cmd.Start()

		if err != nil {
			log.Println("[APP] [FAIL]", id, err)
		} else {
			pid := v.cmd.Process.Pid
			log.Println("[APP] [RUN]", id, pid)
			err = v.cmd.Wait()
			log.Println("[APP] [EXIT]", id, pid)
		}

		v.cmd = nil

		cb <- id
	}()

	if expires != 0 {
		go func() {
			time.Sleep(expires)
			v.Exit()
		}()
	}

	return &v
}

func (A *Application) GetId() int64 {
	return A.id
}

func (A *Application) Exit() {
	if A.stdin != nil {
		stdin := A.stdin
		A.stdin = nil
		stdin.Run([]byte("exit\n"))
		stdin.Close()
	}

}

func (A *Application) RunCommand(command []byte) {
	if A.stdin != nil {
		log.Println("[APP] [COMMAND]", string(command))
		A.stdin.Run(append(command, '\n'))
	}
}
