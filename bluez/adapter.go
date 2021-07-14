package bluez

import (
	"github.com/godbus/dbus"
)

type Adapter struct {
	client  *Client
	adapter string
}

func NewAdapter(adapter string) (*Adapter, error) {
	client, err := NewClient(BluezInterface, AdapterInterface, BluezPath+"/"+adapter)
	if err != nil {
		return nil, err
	}
	return &Adapter{
		client:  client,
		adapter: adapter,
	}, nil
}

func (a *Adapter) StartDiscovery() error {
	cell, err := a.client.Call(
		AdapterInterface,
		0, "StartDiscovery")
	if err != nil {
		return err
	}
	return cell.Store()
}

func (a *Adapter) StopDiscovery() error {
	cell, err := a.client.Call(
		AdapterInterface,
		0, "StopDiscovery")
	if err != nil {
		return err
	}
	return cell.Store()
}

func (a *Adapter) RemoveDevice(path string) error {
	cell, err := a.client.Call(
		AdapterInterface,
		0, "StopDiscovery",
		dbus.ObjectPath(path))
	if err != nil {
		return err
	}
	return cell.Store()
}

func (a *Adapter) Adapter() string {
	return a.adapter
}
