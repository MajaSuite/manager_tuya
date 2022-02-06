package tuya

const (
	HEART_BEAT               Command = 0x0
	PRODUCT_INFO             Command = 0x1
	WIFI_WORK_MODE           Command = 0x2
	WIFI_STATUS              Command = 0x3
	WIFI_RESET               Command = 0x4
	WIFI_MODE                Command = 0x5
	DATA_QUERY               Command = 0x6
	DATA_SET                 Command = 0x7
	STATUS                   Command = 0x8
	PING                     Command = 0x9
	DP_QUERY                 Command = 0xA // get dp
	WIFI_QUERY               Command = 0xB // UPDATE_TRANS_CMD
	GET_ONLINE_TIME          Command = 0xC // token_query GET_ONLINE_TIME_CMD - system time (GMT)
	CONTROL_NEW              Command = 0xD // FACTORY_MODE_CMD
	TEST_WIFI                Command = 0xE // WIFI_TEST_CMD
	F                        Command = 0xF
	DP_QUERY_NEW             Command = 0x10
	SCENE_EXECUTE            Command = 0x11
	UPDATEDPS                Command = 0x12 // Request refresh of DPS
	BROADCAST_NEW            Command = 0x13
	AP_CONFIG_NEW            Command = 0x14
	GET_LOCAL_TIME           Command = 0x1C
	WEATHER_OPEN             Command = 0x20
	WEATHER_DATA             Command = 0x21
	STATE_UPLOAD_SYN         Command = 0x22
	STATE_UPLOAD_SYN_RECV    Command = 0x23
	HEART_BEAT_STOP          Command = 0x25
	STREAM_TRANS             Command = 0x26
	GET_WIFI_STATUS          Command = 0x2B
	WIFI_CONNECT_TEST        Command = 0x2C
	GET_MAC                  Command = 0x2D
	GET_IR_STATUS            Command = 0x2E
	IR_TX_RX_TEST            Command = 0x2F
	LAN_GW_ACTIVE            Command = 0xF0
	LAN_SUB_DEV_REQUEST      Command = 0xF1
	LAN_DELETE_SUB_DEV       Command = 0xF2
	LAN_REPORT_SUB_DEV       Command = 0xF3
	LAN_SCENE                Command = 0xF4
	LAN_PUBLISH_CLOUD_CONFIG Command = 0xF5
	LAN_PUBLISH_APP_CONFIG   Command = 0xF6
	LAN_EXPORT_APP_CONFIG    Command = 0xF7
	LAN_PUBLISH_SCENE_PANEL  Command = 0xF8
	LAN_REMOVE_GW            Command = 0xF9
	LAN_CHECK_GW_UPDATE      Command = 0xFA
	LAN_GW_UPDATE            Command = 0xFB
	LAN_SET_GW_CHANNEL       Command = 0xFC
)

type Command uint32

func (c Command) String() string {
	switch c {
	case HEART_BEAT:
		return "Heartbeat"
	case PRODUCT_INFO:
		return "Product info"
	case WIFI_WORK_MODE:
		return "Wifi work mode"
	case WIFI_STATUS:
		return "Wifi status"
	case WIFI_RESET:
		return "Wifi reset"
	case WIFI_MODE:
		return "Wifi mode"
	case DATA_QUERY:
		return "Data query"
	case DATA_SET:
		return "Set state"
	case STATUS:
		return "Status"
	case PING:
		return "Ping"
	case DP_QUERY:
		return "DP query"
	case 0xB:
		return "Wifi query"
	case 0xC:
		return "Get online time"
	case 0xD:
		return "Factory mode"
	case 0xE:
		return "Wifi test"
	case 0x10:
		return "Dp query new"
	case 0x11:
		return "Run scene"
	case 0x12:
		return "Refresh dps"
	case 0x13:
		return "Broadcast new"
	case 0x14:
		return "AP config"
	case 0x1C:
		return "Get local time"
	case 0x20:
		return "Weather open"
	case 0x21:
		return "Weather data"
	case 0x22:
		return "State upload syn"
	case 0x23:
		return "State upload syn receive"
	case 0x25:
		return "Heartbeat stop"
	case 0x26:
		return "Stream trans"
	case 0x2B:
		return "Wifi get status"
	case 0x2C:
		return "Wifi connect test"
	case 0x2D:
		return "Get mac"
	case 0x2E:
		return "Get IR status"
	case 0x2F:
		return "Test IR tx/rx"
	case 0xF0:
		return "Lan active gw"
	case 0xF1:
		return "Lan subdev request"
	case 0xF2:
		return "Lan delete subdev"
	case 0xF3:
		return "Lan report subdev"
	case 0xF4:
		return "Lan scene"
	case 0xF5:
		return "Lan pubish cloud config"
	case 0xF6:
		return "Lan publish app config"
	case 0xF7:
		return "Lan export app config"
	case 0xF8:
		return "Lan publish scene panel"
	case 0xF9:
		return "Lan remove gw"
	case 0xFA:
		return "Lan check gw update"
	case 0xFB:
		return "Lan update gw"
	case 0xFC:
		return "Lan set gw channel"
	default:
		return "Unknown"
	}
}

type Request struct {
	GwId     string `json:"gwId,omitempty"`
	DevId    string `json:"devId"`
	Uid      string `json:"uid,omitempty"`
	DateTime string `json:"t,omitempty"`
}
