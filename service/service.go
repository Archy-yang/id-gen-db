package service

import (
	"encoding/json"
	"fmt"
	"id-gen-db/gen"
	"net/http"
	"time"
)

func Start() {
	http.HandleFunc("/getID", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("request in: %v\n", r.URL.Query())
		defer func() {
			fmt.Printf("request out: %v\n", r.URL.Query())
		}()
		o := time.After(2 * time.Millisecond)
		d := make(chan struct{})
		go func() {
			select {
			case <-o:
				fmt.Println("超时 http")
			case <-d:
			}
		}()
		id, err := gen.GetIDGen().NextID()
		close(d)

		errStr := ""
		if err != nil {
			errStr = err.Error()
		}
		res := struct {
			ID int64 `json:"id"`
			Err string `json:"err"`
		}{
			ID: id,
			Err: errStr,
		}

		re, err := json.Marshal(res)
		if err != nil {
			fmt.Println(err)
		}

		w.Write([]byte(re))
	})
	fmt.Println("start listen on 127.0.0.1:9999")
	http.ListenAndServe(":9999", nil)
}
