package bluez

import "github.com/godbus/dbus"

type Device struct {
	client *Client
}

func NewDeviceWithFullPath(path dbus.ObjectPath) (*Device, error) {
	devClient, err := NewClientWithFullPath(BluezInterface, DeviceInterface, path)
	if err != nil {
		return nil, err
	}
	return &Device{
		client: devClient,
	}, nil
}

func NewDevice(adapter string, address []byte) (*Device, error) {
	ap := addressToPath(address)
	if ap == "" {
		return nil, ErrInvalidMacAddress
	}
	return NewDeviceWithFullPath(dbus.ObjectPath(BluezPath + "/" + adapter + "/" + ap))
}

func (d *Device) GetTrusted() (bool, error) {
	v, err := d.client.GetProperty("Trusted")
	if err != nil {
		return false, err
	}
	return v.Value().(bool), nil
}

func (d *Device) SetTrusted(trusted bool) error {
	return d.client.SetProperty("Trusted", trusted)
}
