package bluez

import (
	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
)

func ExportInterface(i interface{}, path dbus.ObjectPath, interfaceName string) (*dbus.Conn, error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}

	err = conn.Export(i, path, interfaceName)
	if err != nil {
		return nil, err
	}

	node := &introspect.Node{
		Interfaces: []introspect.Interface{
			introspect.IntrospectData,
			{
				Name:    interfaceName,
				Methods: introspect.Methods(i),
			},
		},
	}

	err = conn.Export(introspect.NewIntrospectable(node), path, Introspectable)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
