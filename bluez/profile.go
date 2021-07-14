package bluez

import "github.com/godbus/dbus"

type Profile interface {
	NewConnection(device dbus.ObjectPath, fd int32, properties map[string]interface{}) *dbus.Error
	RequestDisconnection(device dbus.ObjectPath) *dbus.Error
	Release()
}
