package admin

import (
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/hailongz/kk-room/room"
	"github.com/shirou/gopsutil/process"
)

const VERSION = "1.0.0"

/**
 * data.CPU.threads 线程数
 * data.CPU.used	CPU 使用占比 %
 * data.MEM.used	内存使用占比 %
 * data.MEM.rss		使用的物理内存 b
 * data.MEM.vms		虚拟内存 b
 * data.NET.count	当前连接数
 * data.goroutine	当前协程数
 * data.GC.last		最后GC时间 毫秒
 * data.GC.count	GC次数
 */
func RuntimeState(server room.IServer) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		log.Println(r.URL)

		data := map[string]interface{}{}

		data["version"] = VERSION

		data["roomCount"] = server.GetRoomCount()

		{
			proc, _ := process.NewProcess(int32(os.Getpid()))

			if proc != nil {

				{
					used, _ := proc.CPUPercent()
					threads, _ := proc.NumThreads()

					data["CPU"] = map[string]interface{}{"used": used, "threads": threads}
				}
				{
					mem, _ := proc.MemoryInfo()
					used, _ := proc.MemoryPercent()
					data["MEM"] = map[string]interface{}{"used": used, "rss": mem.RSS, "vms": mem.VMS}
				}
				{
					conns, _ := proc.Connections()
					data["NET"] = map[string]interface{}{"count": len(conns)}
				}
			}

		}

		{
			data["goroutine"] = runtime.NumGoroutine()
		}

		{
			m := runtime.MemStats{}
			runtime.ReadMemStats(&m)
			data["GC"] = map[string]interface{}{"last": m.LastGC / 1000000, "count": m.NumGC}
		}

		SetOutputData(map[string]interface{}{"data": data, "errcode": 200}, w)
	}
}

/**
 * 当前连接
 * data[].remoteaddr.ip		远程IP
 * data[].remoteaddr.port	远程端口
 * data[].localaddr.ip		本地IP
 * data[].localaddr.port	本地端口
 */
func RuntimeConns(server room.IServer) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		log.Println(r.URL)

		data := map[string]interface{}{}

		data["errcode"] = 200

		{
			proc, _ := process.NewProcess(int32(os.Getpid()))

			if proc != nil {

				conns, _ := proc.Connections()

				data["data"] = conns
			}

		}

		SetOutputData(data, w)
	}
}
