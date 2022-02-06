package main

import (
	"encoding/hex"
	"flag"
	"log"
	"manager_tuya/tuya"
	"os"
	"os/signal"
	"syscall"
)

var (
	key   = flag.String("key", "", "tuya account key")
	reg   = flag.Bool("reg", false, "to register new device")
	debug = flag.Bool("debug", false, "debugging mode")
)

func main() {
	flag.Parse()

	if *reg {
		log.Println("new tuya device registration")
		return
	}

	if *key == "" {
		panic("key is not set")
	}
	deviceKey, err := hex.DecodeString(*key)
	if err != nil {
		panic(err)
	}

	log.Println("starting manager_tuya ...")

	devices := make(map[string]tuya.Device)
	log.Println("start broadcast discovery")
	discovery := make(chan tuya.Discovered)
	go tuya.NewDiscovery(*debug, discovery)

	for d := range discovery {
		if devices[d.GwId] == nil {
			if *debug {
				log.Printf("new device: %s", d)
			}

			dev := tuya.NewDevice(*debug, d.Ip, d.GwId, deviceKey)
			if dev != nil {
				devices[d.GwId] = dev
			}
		} else {
			if *debug {
				log.Println("reconnect device", devices[d.GwId])
			}
			devices[d.GwId].Connect(d.Ip)
		}
	}

	finish := make(chan os.Signal, 1)
	signal.Notify(finish, syscall.SIGINT, syscall.SIGTERM)
	<-finish
	log.Println("finished")
}
