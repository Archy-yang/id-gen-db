package main

import (
	"id-gen-db/client"
	"id-gen-db/gen"
	"id-gen-db/service"
	"time"
)

func main()  {
	client.Init()
	gen.GenInit()

	go func() {
		c := time.Tick(time.Second)
		for {
			select {
				case <- c:
					//fmt.Println(runtime.NumGoroutine())
			}
		}
	}()
	service.Start()
}
