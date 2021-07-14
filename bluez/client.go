package bluez

import (
	"errors"

	"github.com/godbus/dbus"
)

type Client struct {
	iface      string
	conn       *dbus.Conn
	dbusObject dbus.BusObject
}

var ErrNotConnected = errors.New("not connected")

func NewClientWithFullPath(name, iface string, path dbus.ObjectPath) (*Client, error) {
	bus, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}
	obj := bus.Object(name, path)
	c := &Client{}
	c.dbusObject = obj
	c.conn = bus
	c.iface = iface
	return c, nil
}

func NewClient(name, iface, path string) (*Client, error) {
	return NewClientWithFullPath(name, iface, dbus.ObjectPath(path))
}

func (c *Client) fullDotName(name string) string {
	return c.iface + "." + name
}

func (c *Client) Call(method string, flags dbus.Flags, args ...interface{}) (*dbus.Call, error) {
	if c.conn == nil {
		return nil, ErrNotConnected
	}
	return c.dbusObject.Call(c.fullDotName(method), flags, args...), nil
}

func (c *Client) GetProperty(p string) (dbus.Variant, error) {
	return c.dbusObject.GetProperty(c.fullDotName(p))
}

func (c *Client) SetProperty(p string, v interface{}) error {
	return c.dbusObject.Call("org.freedesktop.DBus.Properties.Set", 0, c.iface, p, dbus.MakeVariant(v)).Store()
}

func (c *Client) Close() error {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return err
		}
		c.dbusObject = nil
		c.conn = nil
	}
	return nil
}
