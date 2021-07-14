package mgmt

import (
	"encoding/binary"
	"errors"
	"io"
)

func simpleTo(r io.Reader, v interface{}) error {
	return binary.Read(r, binaryOrder, v)
}

func autoTransResponse(r io.Reader, base *CommandComplete) error {
	switch base.OpCode {
	case OpReadManagementVersionInformation:
		base.Response = &ReadVersion{}
		return simpleTo(r, base.Response)
	case OpReadManagementSupporteds:
		var (
			numCommands uint16
			numEvents   uint16
		)
		if err := binary.Read(r, binaryOrder, &numCommands); err != nil {
			return err
		}
		if err := binary.Read(r, binaryOrder, &numEvents); err != nil {
			return err
		}

		commands := &ReadCommands{
			Commands: make([]uint16, numCommands),
			Events:   make([]uint16, numEvents),
		}

		for i := uint16(0); i < numCommands; i++ {
			if err := binary.Read(r, binaryOrder, &commands.Commands[i]); err != nil {
				return err
			}
		}
		for i := uint16(0); i < numEvents; i++ {
			if err := binary.Read(r, binaryOrder, &commands.Events[i]); err != nil {
				return err
			}
		}

		base.Response = commands
		return nil
	case OpReadControllerIndexList:
		var numControllers uint16
		if err := binary.Read(r, binaryOrder, &numControllers); err != nil {
			return err
		}
		list := &ReadControllerIndexList{
			Controllers: make([]uint16, numControllers),
		}
		for i := uint16(0); i < numControllers; i++ {
			if err := binary.Read(r, binaryOrder, &list.Controllers[i]); err != nil {
				return err
			}
		}
		base.Response = list
		return nil
	case OpReadControllerInformation:
		base.Response = &ReadControllerInformation{}
		return simpleTo(r, base.Response)
	case OpSetPowered,
		OpSetDiscoverable,
		OpSetConnectable,
		OpSetFastConnectable,
		OpSetBondable,
		OpSetLinkSecurity,
		OpSetSecureSimplePairing,
		OpSetHighSpeed,
		OpSetLowEnergy:
		var cs uint32
		if err := simpleTo(r, &cs); err != nil {
			return err
		}
		base.Response = cs
		return nil
	case OpSetDeviceClass:
		base.Response = make([]byte, 3)
		return simpleTo(r, base.Response)
	case OpSetLocalName:
		base.Response = &LocalName{}
		return simpleTo(r, base.Response)
	}

	return errors.New("not support to trans")
}
