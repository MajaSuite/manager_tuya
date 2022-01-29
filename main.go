package main

import (
	"flag"
	"log"
	"manager_tuya/tuya"
	"os"
	"os/signal"
	"syscall"
)

var (
	key   = flag.String("key", "", "tuya account key")
	fetch = flag.String("fetch", "openapi.tuyaeu.com", "to fetch devices from the cloud")
	reg   = flag.Bool("reg", false, "to register new device")
	debug = flag.Bool("debug", false, "debugging mode")
)

func main() {
	flag.Parse()

	if *reg {
		log.Println("new tuya device registration")
		return
	}

	if *fetch != "" {
		log.Println("fetching devices from the cloud")
	}

	engine := tuya.NewEngine()

	listener := tuya.NewListener(*debug)
	if listener == nil {
		log.Panic("error start listener")
	}
	go listener.Receiver(engine)

	// todo
	a := tuya.NewAlarmSensor()
	a.SetIP("a")
	log.Println("main:", a)
	log.Println("start:", a.Start())

	finish := make(chan os.Signal, 1)
	signal.Notify(finish, syscall.SIGINT, syscall.SIGTERM)
	<-finish
	log.Println("finished")
}
