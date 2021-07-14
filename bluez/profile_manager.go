package bluez

import "github.com/godbus/dbus"

type ProfileManager struct {
	client *Client
}

func NewProfileManager() (*ProfileManager, error) {
	client, err := NewClient(BluezInterface, ProfileManagerInterface, BluezPath)
	if err != nil {
		return nil, err
	}
	return &ProfileManager{client: client}, nil
}

func (p *ProfileManager) RegisterProfile(profile dbus.ObjectPath, uuid string, option map[string]interface{}) error {
	call, err := p.client.Call("RegisterProfile", 0, profile, uuid, option)
	if err != nil {
		return err
	}
	return call.Store()
}

func (p *ProfileManager) UnregisterProfile(profile dbus.ObjectPath) error {
	call, err := p.client.Call("UnregisterProfile", 0, profile)
	if err != nil {
		return err
	}
	return call.Store()
}
