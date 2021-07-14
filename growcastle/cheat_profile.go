package growcastle

import (
	"encoding/xml"
	"fmt"
	"vitrhid/bluez"

	"github.com/godbus/dbus"
)

type Profile struct {
	conn *dbus.Conn
}

func NewProfile() (*Profile, error) {
	profile := &Profile{}
	conn, err := bluez.ExportInterface(profile, ProfilePath, bluez.ProfileInterface)
	if err != nil {
		return nil, err
	}
	profile.conn = conn
	return profile, nil
}

func (p *Profile) NewConnection(device dbus.ObjectPath, fd int32, properties map[string]interface{}) *dbus.Error {
	fmt.Println("GrowCastle", "GotConnection", device, fd, properties)
	return nil
}

func (p *Profile) RequestDisconnection(device dbus.ObjectPath) *dbus.Error {
	fmt.Println("GrowCastle", "RequestDisconnection", device)
	return nil
}

func (p *Profile) Release() {
}

type UUID struct {
	XMLName xml.Name `xml:"uuid"`
	Value   string   `xml:"value,attr"`
}

type Sequence struct {
	XMLName xml.Name `xml:"sequence"`
	Value   []interface{}
}

type UInt16 struct {
	XMLName xml.Name `xml:"uint16"`
	Value   string   `xml:"value,attr"`
}

type UInt8 struct {
	XMLName xml.Name `xml:"uint8"`
	Value   string   `xml:"value,attr"`
}

type Text struct {
	XMLName  xml.Name `xml:"text"`
	Encoding string   `xml:"encoding,attr,omitempty"`
	Value    string   `xml:"value,attr"`
}

type Boolean struct {
	XMLName xml.Name `xml:"boolean"`
	Value   string   `xml:"value,attr"`
}

type Attribute struct {
	XMLName xml.Name `xml:"attribute"`
	Id      string   `xml:"id,attr"`
	Value   interface{}
}

type Record struct {
	XMLName xml.Name `xml:"record"`
	Records []interface{}
}
