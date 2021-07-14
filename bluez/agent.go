package bluez

import (
	"log"
	"math/rand"
	"time"

	"github.com/godbus/dbus"
)

type Agent interface {
	RequestPinCode(device dbus.ObjectPath) (string, *dbus.Error)
	DisplayPinCode(device dbus.ObjectPath, pinCode string) *dbus.Error
	RequestPasskey(device dbus.ObjectPath) (uint32, *dbus.Error)
	DisplayPasskey(device dbus.ObjectPath, passkey uint32, entered uint16) *dbus.Error
	RequestConfirmation(device dbus.ObjectPath, passkey uint32) *dbus.Error
	RequestAuthorization(device dbus.ObjectPath) *dbus.Error
	AuthorizeService(device dbus.ObjectPath, uuid string) *dbus.Error
	Release() *dbus.Error
	Cancel() *dbus.Error
}

var (
	ErrRejected = &dbus.Error{
		Name: "org.bluez.Error.Rejected",
		Body: []interface{}{"Rejected"},
	}
	ErrCanceled = &dbus.Error{
		Name: "org.bluez.Error.Canceled",
		Body: []interface{}{"Canceled"},
	}
)

type SimpleAgent struct {
	path    dbus.ObjectPath
	passKey uint32
	pinCode string
	conn    *dbus.Conn
}

func NewSimpleAgent(path dbus.ObjectPath) (*SimpleAgent, error) {
	agent := &SimpleAgent{}
	agent.path = path

	conn, err := ExportInterface(agent, agent.path, AgentInterface)
	if err != nil {
		return nil, err
	}
	agent.conn = conn

	return agent, nil
}

func (a *SimpleAgent) rand(min, max int) int {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Intn(max-min) + min
}

func (a *SimpleAgent) RequestPinCode(device dbus.ObjectPath) (string, *dbus.Error) {
	log.Printf("RequestPinCode")
	pinCode := string(rune(a.rand(0, 9) + 48))
	for i := 0; i < 3; i++ {
		pinCode += string(rune(a.rand(0, 9) + 48))
	}
	a.pinCode = pinCode
	return a.pinCode, nil
}

func (a *SimpleAgent) DisplayPinCode(device dbus.ObjectPath, pinCode string) *dbus.Error {
	log.Printf("DisplayPinCode %s", pinCode)
	return nil
}

func (a *SimpleAgent) RequestPasskey(device dbus.ObjectPath) (uint32, *dbus.Error) {
	a.passKey = uint32(a.rand(0, 999999))
	log.Printf("RequestPasskey %06d", a.passKey)
	return a.passKey, nil
}

func (a *SimpleAgent) RequestConfirmation(device dbus.ObjectPath, passkey uint32) *dbus.Error {
	log.Printf("RequestConfirmation %d", passkey)
	d, err := NewDeviceWithFullPath(device)
	if err != nil {
		return ErrCanceled
	}
	if err := d.SetTrusted(true); err != nil {
		return ErrCanceled
	}
	return nil
}

func (a *SimpleAgent) RequestAuthorization(device dbus.ObjectPath) *dbus.Error {
	log.Printf("RequestAuthorization")
	return nil
}

func (a *SimpleAgent) AuthorizeService(device dbus.ObjectPath, uuid string) *dbus.Error {
	log.Printf("AuthorizeService %s", uuid)
	return nil
}

func (a *SimpleAgent) Release() *dbus.Error {
	return nil
}

func (a *SimpleAgent) Cancel() *dbus.Error {
	return nil
}
