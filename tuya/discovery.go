package tuya

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"log"
	"manager_tuya/utils"
	"net"
	"time"
)

var (
	udpKey = []byte("yGAdlopoPVldABfn")
)

type Discovered struct {
	Ip         string `json:"ip"`
	GwId       string `json:"gwId"`
	Active     int    `json:"active"`
	Ability    int    `json:"ability"`
	Mode       int    `json:"mode"`
	Encrypt    bool   `json:"encrypt"`
	ProductKey string `json:"productKey"`
	Version    string `json:"version"`
}

func (d Discovered) String() string {
	return fmt.Sprintf(`{"ip":"%s","gwid":"%s","active":%d,"ability":%d,"mode":%d,"encrypt":%v,"key":"%s","ver":"%s"}`,
		d.Ip, d.GwId, d.Active, d.Ability, d.Mode, d.Encrypt, d.ProductKey, d.Version)
}

func NewDiscovery(debug bool, report chan Discovered) error {
	conn, err := net.ListenPacket("udp4", ":6667")
	if err != nil {
		return err
	}

	key := md5.Sum(udpKey)
	for {
		buffer := make([]byte, 200)
		conn.SetReadDeadline(time.Now().Add(time.Second * 5))
		n, _, err := conn.ReadFrom(buffer)
		if err != nil || n < 24 {
			continue
		}

		if debug {
			log.Printf("receive request\n%s", hex.Dump(buffer))
		}

		startToken := uint32(buffer[3]) | uint32(buffer[2])<<8 | uint32(buffer[1])<<16 | uint32(buffer[0])<<24
		if startToken != 0x55aa {
			if debug {
				log.Println("wrong start token")
			}
			continue
		}

		// seq 7 6 5 4

		// command 11 10 9 8

		size := uint32(buffer[15]) | uint32(buffer[14])<<8 | uint32(buffer[13])<<16 | uint32(buffer[12])<<24
		if size > 180 {
			if debug {
				log.Println("buffer overflow")
			}
			continue
		}

		if (uint32(buffer[n-5]) | uint32(buffer[n-6])<<8 | uint32(buffer[n-7])<<16 | uint32(buffer[n-8])<<24) !=
			crc32.ChecksumIEEE(buffer[:n-8]) {
			if debug {
				log.Println("wrong crc in packet")
			}
			continue
		}

		data, _, err := utils.ReadBytes(buffer, 20, n-20-8)
		payload, err := aesDecrypt(data, key[:])
		if err != nil {
			log.Println("error decrypt", err)
		}

		if debug {
			log.Println("payload: ", string(payload))
		}

		if err == nil {
			var response Discovered
			if e := json.Unmarshal(payload, &response); e == nil {
				report <- response
			}
		}
	}
}
