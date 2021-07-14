package main

import (
	"errors"
	"log"
	"net/http"
	"syscall"
	"vitrhid/bluez"
	"vitrhid/growcastle"
	"vitrhid/mgmt"

	"golang.org/x/sys/unix"
)

var (
	controlListenFd   int
	interruptListenFd int
)

func l2capListen(psm uint16) (int, error) {
	fd, err := unix.Socket(syscall.AF_BLUETOOTH, syscall.SOCK_SEQPACKET, unix.BTPROTO_L2CAP)
	if err != nil {
		return -1, err
	}

	if err := syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1); err != nil {
		return -1, err
	}

	if err := unix.Bind(fd, &unix.SockaddrL2{
		PSM: psm,
		//Addr: [6]byte{0xB8, 0x27, 0xEB, 0xCC, 0x85, 0x7A},
	}); err != nil {
		return -1, err
	}

	if err := unix.Listen(fd, 5); err != nil {
		return -1, err
	}

	return fd, nil
}

func initLowLevelBluetooth() (uint16, error) {
	ll := mgmt.NewBluetoothLowLevel()
	if err := ll.Connect(); err != nil {
		return 0, err
	}

	list, err := ll.ReadControllerIndexList()
	if err != nil {
		return 0, err
	}
	if len(list.Controllers) == 0 {
		return 0, errors.New("no controller")
	}

	index := list.Controllers[0]

	if _, err := ll.SetPowered(index, mgmt.On); err != nil {
		return 0, err
	}
	log.Printf("Bluetooth Powered On")

	if _, err := ll.SetConnectable(index, mgmt.On); err != nil {
		return 0, err
	}
	log.Printf("Bluetooth Connectable On")

	if err := ll.SetLocalName(index, "AnonymousCheat", "AC"); err != nil {
		return 0, err
	}
	log.Printf("Bluetooth SetLocalName")

	if _, err := ll.SetSecureSimplePairing(index, mgmt.On); err != nil {
		return 0, err
	}
	log.Printf("Bluetooth Set Secure Simple Pairing")

	if _, err := ll.SetDiscoverable(index, mgmt.On, 0x500); err != nil {
		return 0, err
	}
	log.Printf("Bluetooth Set Discovereable")

	if _, err := ll.SetDeviceClass(index, 5, 64); err != nil {
		return 0, err
	}
	log.Printf("Bluetooth Set Device Class")

	if err := ll.SetAppearance(index, 0x03C0); err != nil {
		return 0, err
	}
	log.Printf("Bluetooth Set Appearance")

	return index, nil
}

func initBluez(index uint16) error {
	am, err := bluez.NewAgentManager()
	if err != nil {
		return err
	}

	_, err = bluez.NewSimpleAgent(growcastle.AgentPath)
	if err != nil {
		return err
	}

	err = am.RegisterAgent(growcastle.AgentPath, bluez.AgentCapabilityDisplayOnly)
	if err != nil {
		return err
	}

	err = am.RequestDefaultAgent(growcastle.AgentPath)
	if err != nil {
		return err
	}

	pm, err := bluez.NewProfileManager()
	if err != nil {
		return err
	}

	_, err = growcastle.NewProfile()
	if err != nil {
		return err
	}

	var descriptor [][]byte
	descriptor = append(descriptor, growcastle.MouseDescriptor())

	record, err := growcastle.SDPRecord(descriptor)
	if err != nil {
		return err
	}

	service := "00001124-0000-1000-8000-00805f9b34fb"

	opts := make(map[string]interface{})
	opts["Name"] = "GrowCastle"
	opts["Role"] = "server"
	opts["AutoConnect"] = true
	opts["ServiceRecord"] = record

	err = pm.RegisterProfile(
		growcastle.ProfilePath,
		service, // HumanInterfaceDeviceServiceClass
		opts,
	)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	var err error

	controlListenFd, err = l2capListen(0x11)
	if err != nil {
		log.Fatalf("l2cap: listen control")
	}
	interruptListenFd, err = l2capListen(0x13)
	if err != nil {
		log.Fatalf("l2cap: listen interrupt")
	}

	index, err := initLowLevelBluetooth()
	if err != nil {
		log.Fatalf("bluetooth: %s\n", err)
	}

	if err := initBluez(index); err != nil {
		log.Fatalf("bluez: %s\n", err)
	}

	s := NewServices()
	go s.AcceptControl()
	go s.AcceptInterrupt()

	if err := http.ListenAndServe(":8080", s); err != nil {
		log.Fatalf("http: %s\n", err)
	}
}
