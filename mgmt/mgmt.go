package mgmt

import (
	"bytes"
	"container/list"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"syscall"

	"golang.org/x/sys/unix"
)

const (
	EventCommandComplete = iota + 1
	EventCommandStatus
	EventControllerError
	EventIndexAdded
	EventIndexRemoved
	EventNewSettings
	EventClassOfDevChanged
	EventLocalNameChanged
	EventNewLinkKey
	EventNewLongTermKey
	EventDeviceConnected
)

var binaryOrder = binary.LittleEndian

type CommandError struct {
	Code byte
}

func (e *CommandError) Error() string {
	return e.String()
}

func (e *CommandError) String() string {
	switch e.Code {
	case ErrUnknownCommand:
		return "ErrUnknownCommand"
	case ErrNotConnect:
		return "ErrNotCommect"
	case ErrFailed:
		return "ErrFailed"
	case ErrConnectFailed:
		return "ErrConnectFailed"
	case ErrAuthenticationFailed:
		return "ErrAuthenticationFailed"
	case ErrNotPaired:
		return "ErrNotPaired"
	case ErrNoResources:
		return "ErrNoResources"
	case ErrTimeout:
		return "ErrTimeout"
	case ErrAlreadyConnected:
		return "ErrAlreadyConnected"
	case ErrBusy:
		return "ErrBusy"
	case ErrRejected:
		return "ErrRejected"
	case ErrNotSupported:
		return "ErrNotSupported"
	case ErrInvalidParameters:
		return "ErrInvalidParameters"
	case ErrDisconnected:
		return "ErrDisconnected"
	case ErrNotPowered:
		return "ErrNotPowered"
	case ErrCancelled:
		return "ErrCancelled"
	case ErrInvalidIndex:
		return "ErrInvalidIndex"
	case ErrRFKilled:
		return "ErrRFKilled"
	case ErrAlreadyPaired:
		return "ErrAlreadyPaired"
	case ErrPermissionDenied:
		return "ErrPermissionDenied"
	}
	return "ErrUnknown"
}

// BluetoothLowLevel detail docs https://github.com/bluez/bluez/blob/master/doc/mgmt-api.txt
type BluetoothLowLevel struct {
	fd      int
	pending *list.List
}

func NewBluetoothLowLevel() *BluetoothLowLevel {
	b := BluetoothLowLevel{}
	b.fd = -1
	b.pending = list.New()
	return &b
}

func (b *BluetoothLowLevel) commandComplete(cmd *Command) {
	if len(cmd.Data) < 3 {
		return
	}

	r := bytes.NewReader(cmd.Data)

	var (
		cmdCode   uint16
		cmdStatus uint8
	)
	binary.Read(r, binaryOrder, &cmdCode)
	binary.Read(r, binaryOrder, &cmdStatus)

	element := b.pending.Front()
	for element != nil {
		pendingCmd := element.Value.(*Command)
		if pendingCmd.OpCode != cmdCode {
			element = element.Next()
			continue
		}
		if pendingCmd.Controller != cmd.Controller {
			element = element.Next()
			continue
		}

		b.pending.Remove(element)

		base := &CommandComplete{
			OpCode: cmdCode,
			Status: cmdStatus,
		}
		if len(cmd.Data)-3 > 0 {
			if err := autoTransResponse(r, base); err != nil {
				base.Response = cmd.Data[3:]
			}
		}

		pendingCmd.pkt <- base
		return
	}
}

func (b *BluetoothLowLevel) eventLoop(epollFd int) {
	if b.fd == -1 {
		return
	}

	defer b.Close()

	var events [128]syscall.EpollEvent
	readBuf := make([]byte, 1024)

	for {
		numEvents, err := syscall.EpollWait(epollFd, events[:], -1)
		if err != nil {
			log.Println(err)
			return
		}
		for i := 0; i < numEvents; i++ {
			if (events[i].Events & syscall.EPOLLIN) == syscall.EPOLLIN {
				n, err := unix.Read(int(events[i].Fd), readBuf)
				if err != nil || n <= 0 {
					fmt.Printf("connection closed by peer\n")
					return
				} else if n < 6 {
					fmt.Printf("frame too short\n")
					break
				}

				cut := readBuf[:n]
				r := bytes.NewReader(cut)
				base := &Command{}
				base.Deserialize(r)

				switch base.OpCode {
				case EventCommandComplete, EventCommandStatus:
					b.commandComplete(base)
					break
				}
			}
		}
	}
}

func (b *BluetoothLowLevel) Connect() error {
	fd, err := unix.Socket(syscall.AF_BLUETOOTH,
		syscall.SOCK_RAW|syscall.SOCK_CLOEXEC|syscall.SOCK_NONBLOCK,
		unix.BTPROTO_HCI)
	if err != nil {
		return err
	}

	addr := unix.SockaddrHCI{
		Dev:     0xFFFF,
		Channel: unix.HCI_CHANNEL_CONTROL,
	}

	if err := unix.Bind(fd, &addr); err != nil {
		return err
	}

	epollFd, err := syscall.EpollCreate1(syscall.EPOLL_CLOEXEC)
	if err != nil {
		return err
	}

	ev := syscall.EpollEvent{
		Events: syscall.EPOLLIN,
		Fd:     int32(fd),
	}

	if err := syscall.EpollCtl(epollFd, syscall.EPOLL_CTL_ADD, fd, &ev); err != nil {
		return err
	}

	b.fd = fd

	go b.eventLoop(epollFd)

	return nil
}

func (b *BluetoothLowLevel) Send(cmd *Command) (*CommandComplete, error) {
	if b.fd == -1 {
		return nil, errors.New("not connect")
	}

	buf := cmd.Serialize()

	cmd.pkt = make(chan *CommandComplete, 1)
	defer close(cmd.pkt)

	b.pending.PushBack(cmd)

	if _, err := unix.Write(b.fd, buf); err != nil {
		return nil, err
	}

	pkt := <-cmd.pkt

	if pkt.Status != Success {
		return nil, &CommandError{Code: pkt.Status}
	}

	return pkt, nil
}

const (
	Off byte = 0
	On  byte = 1
)

func (b *BluetoothLowLevel) oneByteCommand(index, opcode uint16, e byte) (*CommandComplete, error) {
	return b.Send(&Command{
		OpCode:     opcode,
		Controller: index,
		Data:       []byte{e},
	})
}

func (b *BluetoothLowLevel) SetPowered(index uint16, powered byte) (uint32, error) {
	pkt, err := b.oneByteCommand(index, OpSetPowered, powered)
	if err != nil {
		return 0, err
	}
	return pkt.Response.(uint32), nil
}

const (
	DiscoverableLimit byte = 2
)

func (b *BluetoothLowLevel) SetDiscoverable(index uint16, discoverable byte, timeout uint16) (uint32, error) {
	buf := &bytes.Buffer{}
	binary.Write(buf, binaryOrder, discoverable)
	binary.Write(buf, binaryOrder, timeout)
	pkt, err := b.Send(&Command{
		OpCode:     OpSetDiscoverable,
		Controller: index,
		Data:       buf.Bytes(),
	})
	if err != nil {
		return 0, err
	}
	return pkt.Response.(uint32), nil
}

func (b *BluetoothLowLevel) SetConnectable(index uint16, connectable byte) (uint32, error) {
	pkt, err := b.oneByteCommand(index, OpSetConnectable, connectable)
	if err != nil {
		return 0, err
	}
	return pkt.Response.(uint32), nil
}

func (b *BluetoothLowLevel) SetFastConnectable(index uint16, enable byte) (uint32, error) {
	pkt, err := b.oneByteCommand(index, OpSetFastConnectable, enable)
	if err != nil {
		return 0, err
	}
	return pkt.Response.(uint32), nil
}

func (b *BluetoothLowLevel) SetBondable(index uint16, bondable byte) (uint32, error) {
	pkt, err := b.oneByteCommand(index, OpSetBondable, bondable)
	if err != nil {
		return 0, err
	}
	return pkt.Response.(uint32), nil
}

func (b *BluetoothLowLevel) SetLinkSecurity(index uint16, linkSecurity byte) (uint32, error) {
	pkt, err := b.oneByteCommand(index, OpSetLinkSecurity, linkSecurity)
	if err != nil {
		return 0, err
	}
	return pkt.Response.(uint32), nil
}

func (b *BluetoothLowLevel) SetSecureSimplePairing(index uint16, ssp byte) (uint32, error) {
	pkt, err := b.oneByteCommand(index, OpSetSecureSimplePairing, ssp)
	if err != nil {
		return 0, err
	}
	return pkt.Response.(uint32), nil
}

func (b *BluetoothLowLevel) SetHighSpeed(index uint16, highSpeed byte) (uint32, error) {
	pkt, err := b.oneByteCommand(index, OpSetHighSpeed, highSpeed)
	if err != nil {
		return 0, err
	}
	return pkt.Response.(uint32), nil
}

func (b *BluetoothLowLevel) SetLowEnergy(index uint16, lowEnergy byte) (uint32, error) {
	pkt, err := b.oneByteCommand(index, OpSetLowEnergy, lowEnergy)
	if err != nil {
		return 0, err
	}
	return pkt.Response.(uint32), nil
}

func (b *BluetoothLowLevel) SetDeviceClass(index uint16, majorDeviceClass, minorDeviceClass byte) ([]byte, error) {
	pkt, err := b.Send(&Command{
		OpCode:     OpSetDeviceClass,
		Controller: index,
		Data:       []byte{majorDeviceClass, minorDeviceClass},
	})
	if err != nil {
		return nil, err
	}
	return pkt.Response.([]byte), nil
}

func (b *BluetoothLowLevel) ReadControllerIndexList() (*ReadControllerIndexList, error) {
	pkt, err := b.Send(&Command{
		OpCode:     OpReadControllerIndexList,
		Controller: NonController,
	})
	if err != nil {
		return nil, err
	}
	return pkt.Response.(*ReadControllerIndexList), nil
}

func (b *BluetoothLowLevel) SetLocalName(index uint16, name, shortName string) error {
	bName := []byte(name)
	bShortName := []byte(shortName)

	if len(bName) >= 249 {
		return errors.New("name length not allow")
	}
	if len(bShortName) >= 11 {
		return errors.New("short name length not allow")
	}

	paramName := make([]byte, 249)
	paramShortName := make([]byte, 11)
	copy(paramName, bName)
	copy(paramShortName, bShortName)

	_, err := b.Send(&Command{
		OpCode:     OpSetLocalName,
		Controller: index,
		Data:       append(paramName, paramShortName...),
	})
	if err != nil {
		return err
	}
	return nil
}

func (b *BluetoothLowLevel) AddUUID(index uint16, uuid []byte, svcHint byte) error {
	_, err := b.Send(&Command{
		OpCode:     OpAddUUID,
		Controller: index,
		Data:       append(uuid, svcHint),
	})
	if err != nil {
		return err
	}
	return nil
}

func (b *BluetoothLowLevel) RemoveUUID(index uint16, uuid []byte) error {
	_, err := b.Send(&Command{
		OpCode:     OpRemoveUUID,
		Controller: index,
		Data:       append(uuid),
	})
	if err != nil {
		return err
	}
	return nil
}

func (b *BluetoothLowLevel) SetAppearance(index uint16, appearance uint16) error {
	data := make([]byte, 2)
	binaryOrder.PutUint16(data, appearance)
	_, err := b.Send(&Command{
		OpCode:     OpSetAppearance,
		Controller: index,
		Data:       data,
	})
	if err != nil {
		return err
	}
	return nil
}

func (b *BluetoothLowLevel) Close() error {
	if b.fd != -1 {
		if err := unix.Close(b.fd); err != nil {
			return err
		}
	}
	b.fd = -1
	return nil
}
