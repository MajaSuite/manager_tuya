package tuya

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"
	"log"
	"manager_tuya/utils"
	"net"
	"strconv"
	"time"
)

const (
	NO_TYPE Type = iota
	ALARM_SENSOR
	BULB
	CEILING_LIGHT
)

var (
	ErrWrongHeader = errors.New("wrong response header")
	ErrWrongSize   = errors.New("wrong payload size")
	ErrWrongCRC    = errors.New("wrong crc")

	timeout = time.Second * 3
)

type Device interface {
	Type() Type
	Name() string
	IP() string
	Connect(ip string) error
	Close()
	SendPacket(command Command, data []byte, addVersion bool) error
	ReceivePacket() ([]byte, uint32, error)
	String() string
}

type TuyaDevice struct {
	debug      bool
	deviceType Type
	deviceName string
	deviceIp   string
	deviceId   string
	deviceKey  []byte
	deviceSeq  uint32
	conn       net.Conn
}

type Type byte

func (t Type) String() string {
	switch t {
	case ALARM_SENSOR:
		return "Alarm sensor"
	case BULB:
		return "Bulb"
	case CEILING_LIGHT:
		return "Ceiling light"
	default:
		return "n/a"
	}
}

func NewDevice(debug bool, ip string, id string, key []byte) Device {
	d := &TuyaDevice{
		debug:     debug,
		deviceId:  id,
		deviceKey: key,
	}

	if err := d.Connect(ip); err != nil {
		log.Println("error establish connection", err)
		return nil
	}
	return d

	// todo
	//switch expr {
	//case ALARM_SENSOR:
	//return NewAlarmSensor()
	//default:
	//return nil
	//}
}

func (td *TuyaDevice) Type() Type {
	return td.deviceType
}

func (td *TuyaDevice) Name() string {
	return td.deviceName
}

func (td *TuyaDevice) IP() string {
	return td.deviceIp
}

func (td *TuyaDevice) Connect(ip string) error {
	td.deviceIp = ip

	var err error
	log.Printf("connecting to device %s (%s)", td.deviceId, td.deviceIp)
	td.conn, err = net.DialTimeout("tcp", net.JoinHostPort(td.deviceIp, "6668"), timeout)
	if err != nil {
		return err
	}

	request, err := json.Marshal(Request{
		//GwId:     td.deviceId,
		DevId:    td.deviceId,
		Uid:      td.deviceId,
		DateTime: strconv.Itoa(int(time.Now().Unix())),
	})
	if err != nil {
		return err
	}

	if td.debug {
		log.Println("prepared request", string(request))
	}

	err = td.SendPacket(DP_QUERY, request, false)
	if err != nil {
		td.conn.Close()
		return err
	}
	td.deviceSeq++

	payload, retCode, err := td.ReceivePacket()
	log.Printf("payload '%v' retcode:%d", payload, retCode)

	///////////////////////////
	/*
			req, err := json.Marshal(Request{
				//	GwId:  td.deviceId,
				DevId:    td.deviceId,
				Uid:      td.deviceId,
				DateTime: strconv.Itoa(int(time.Now().Unix())),
			})
			err = td.SendPacket(DP_QUERY_NEW, req, true)
			if err != nil {
				td.conn.Close()
				return err
			}
			td.deviceSeq++


		pl, err := td.ReceivePacket()
		log.Println("payload", pl)
	*/

	go func() {
		for {
			time.Sleep(time.Second * 10)
			err = td.SendPacket(PING, []byte("{}"), false)
			if err != nil {
				td.conn.Close()
				return
			}

			payload, retCode, err := td.ReceivePacket()
			if err != nil {
				log.Println("error", err)
			}
			log.Printf("payload '%v' retcode:%d", payload, retCode)

		}
	}()
	return nil
}

func (td *TuyaDevice) Close() {
	if td.conn != nil {
		td.conn.Close()
	}
	td.conn = nil
	td.deviceIp = ""
}

func (td *TuyaDevice) SendPacket(command Command, data []byte, addVersion bool) error {
	if td.debug {
		log.Println("send command", command.String())
	}

	payload, err := aesEncrypt(data, td.deviceKey)
	if err != nil {
		return err
	}

	v := 0
	if addVersion {
		v = 4*4 + 3
	}

	buffer := make([]byte, len(payload)+(6*4)+v)
	offset := utils.WriteInt32(buffer, 0, uint32(0x55aa))
	offset = utils.WriteInt32(buffer, offset, uint32(td.deviceSeq))
	offset = utils.WriteInt32(buffer, offset, uint32(command))
	offset = utils.WriteInt32(buffer, offset, uint32(len(payload)+(2*4)+v))
	if addVersion {
		offset = utils.WriteInt32(buffer, offset, uint32(0x0000))
		offset = utils.WriteBytes(buffer, offset, []byte{0x33, 0x2e, 0x33})
		offset = utils.WriteInt32(buffer, offset, uint32(0x0000))
		offset = utils.WriteInt32(buffer, offset, uint32(0x0000))
		offset = utils.WriteInt32(buffer, offset, uint32(0x0000))
	}
	offset = utils.WriteBytes(buffer, offset, payload)
	//crcOffset := offset
	offset = utils.WriteInt32(buffer, offset, crc32.ChecksumIEEE(buffer[:len(payload)+(4*4)+v]))
	offset = utils.WriteInt32(buffer, offset, uint32(0xaa55))

	td.conn.SetWriteDeadline(time.Now().Add(timeout))
	if _, err := td.conn.Write(buffer); err != nil {
		return err
	}

	return nil
}

func (td *TuyaDevice) ReceivePacket() ([]byte, uint32, error) {
	header := make([]byte, 16)
	td.conn.SetReadDeadline(time.Now().Add(timeout))
	if _, err := td.conn.Read(header); err != nil {
		td.conn.Close()
		return nil, 0, err
	}

	if td.debug {
		log.Printf("receive header\n%s", hex.Dump(header))
	}

	startToken := uint32(header[3]) | uint32(header[2])<<8 | uint32(header[1])<<16 | uint32(header[0])<<24
	if startToken != 0x55aa {
		td.conn.Close()
		return nil, 0, ErrWrongHeader
	}

	td.deviceSeq = uint32(header[7]) | uint32(header[6])<<8 | uint32(header[5])<<16 | uint32(header[4])<<24
	// 11 10 9 8 - command
	size := int(uint32(header[15]) | uint32(header[14])<<8 | uint32(header[13])<<16 | uint32(header[12])<<24)

	/* payload format:
	4 byte - 0x0000 0x0000 - retcade
	3 byte optional - 0x33 0x2e 0x33 - version (3.3)
	4 byte optional - 0x0000 0x0000 - ?
	4 byte optional - 0x0000 0x000e - looks like connection uptime
	4 byte optional - 0x0000 0x0001 - ?
	x byte actually optional - ... - encrypted data block (json)
	4 byte - CRC
	4 byte end header
	*/
	payload := make([]byte, size)
	td.conn.SetReadDeadline(time.Now().Add(timeout))
	if n, err := td.conn.Read(payload); n < size || err != nil {
		td.conn.Close()
		return nil, 0, err
	}

	if td.debug {
		log.Printf("receive payload\n%s", hex.Dump(payload))
	}

	retCode := uint32(payload[3]) | uint32(payload[2])<<8 | uint32(payload[1])<<16 | uint32(payload[0])<<24

	orig := uint32(payload[size-5]) | uint32(payload[6])<<8 | uint32(payload[size-7])<<16 | uint32(payload[size-8])<<24
	calc := crc32.ChecksumIEEE(append(header, payload[:size-8]...))
	if orig != calc {
		if td.debug {
			log.Printf("wrong crc expect:%x calc:%x", orig, calc)
		}
		//rd.conn.Close()
		//return nil, retCode, ErrWrongCRC
	}

	if size < 14 {
		return nil, retCode, nil
	}

	data, _, err := utils.ReadBytes(payload, 4, int(size)-4-8)
	if err != nil {
		return nil, retCode, err
	}

	log.Printf("encrypted data:\n%s", hex.Dump(data))
	decrypted, err := aesDecrypt(data, td.deviceKey)
	if err != nil {
		return nil, retCode, err
	}

	if td.debug {
		log.Printf("decrypted data:\n%s", string(decrypted))
	}

	return decrypted, retCode, nil
}

func (td *TuyaDevice) String() string {
	return fmt.Sprintf(`{"debug":%b,"type":"%x","name":"%s","ip":"%s","id":"%x","key":"%x","seq":%d}`,
		td.deviceSeq, td.deviceType, td.deviceName, td.deviceIp, td.deviceId, td.deviceKey, td.deviceSeq)
}
