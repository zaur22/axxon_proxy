package main

import (
	"axxon_proxy/proxy"
	"axxon_proxy/router"
	"axxon_proxy/task"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	tm := task.NewTaskManager()
	proxy := proxy.NewProxy(tm)
	router := router.Router{
		Proxy: proxy,
	}

	http.HandleFunc("/fetch", router.FetchTask)
	http.HandleFunc("/get", router.GetTasks)
	http.HandleFunc("/delete", router.DeleteTask)

	go func() {
		log.Printf("starting server at 3000")
		err := http.ListenAndServe(":3001", nil)
		if err != nil {
			log.Fatalf("Server error: %v", err.Error())
		}
	}()

	go func(stop func()) {
		<-sigs
		fmt.Println()
		fmt.Println("Stopping program...")
		stop()
		done <- true
	}(proxy.StopWorkers)

	<-done
	print("Done\n")
}
