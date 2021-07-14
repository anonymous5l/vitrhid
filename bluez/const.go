package bluez

import (
	"bytes"
	"errors"
)

var (
	ErrInvalidMacAddress = errors.New("invalid mac address")
)

const (
	BluezInterface          = "org.bluez"
	BluezPath               = "/org/bluez"
	AdapterInterface        = "org.bluez.Adapter1"
	DeviceInterface         = "org.bluez.Device1"
	AgentManagerInterface   = "org.bluez.AgentManager1"
	AgentInterface          = "org.bluez.Agent1"
	ProfileInterface        = "org.bluez.Profile1"
	ProfileManagerInterface = "org.bluez.ProfileManager1"
)

const (
	Introspectable = "org.freedesktop.DBus.Introspectable"
)

const hextable = "0123456789abcdef"

func addressToPath(address []byte) string {
	if len(address) != 6 {
		return ""
	}

	buf := bytes.NewBufferString("")
	buf.WriteString("dev_") // prefix
	for i := 0; i < len(address)-1; i++ {
		buf.WriteByte(hextable[address[i]>>4])
		buf.WriteByte(hextable[address[i]&0x0f])
		buf.WriteByte('_')
	}
	buf.WriteByte(hextable[address[len(address)-1]>>4])
	buf.WriteByte(hextable[address[len(address)-1]&0x0f])
	return buf.String()
}
