package bluez

import "github.com/godbus/dbus"

type AgentManager struct {
	client *Client
}

func NewAgentManager() (*AgentManager, error) {
	client, err := NewClient(BluezInterface, AgentManagerInterface, BluezPath)
	if err != nil {
		return nil, err
	}
	return &AgentManager{client: client}, nil
}

type AgentCapability string

const (
	AgentCapabilityDisplayOnly     AgentCapability = "DisplayOnly"
	AgentCapabilityDisplayYesNo    AgentCapability = "DisplayYesNo"
	AgentCapabilityKeyboardOnly    AgentCapability = "KeyboardOnly"
	AgentCapabilityNoInputNoOutput AgentCapability = "NoInputNoOutput"
	AgentCapabilityKeyboardDisplay AgentCapability = "KeyboardDisplay"
)

func (a *AgentManager) RegisterAgent(agent dbus.ObjectPath, capability AgentCapability) error {
	call, err := a.client.Call("RegisterAgent", 0, agent, capability)
	if err != nil {
		return err
	}
	return call.Store()
}

func (a *AgentManager) UnregisterAgent(agent dbus.ObjectPath) error {
	call, err := a.client.Call("UnregisterAgent", 0, agent)
	if err != nil {
		return err
	}
	return call.Store()
}

func (a *AgentManager) RequestDefaultAgent(agent dbus.ObjectPath) error {
	call, err := a.client.Call("RequestDefaultAgent", 0, agent)
	if err != nil {
		return err
	}
	return call.Store()
}
