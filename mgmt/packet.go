package mgmt

import (
	"bytes"
	"encoding/binary"
	"io"
)

const (
	MaxControllerId = 0xFFFE + iota
	NonController
)

const (
	Success = iota
	ErrUnknownCommand
	ErrNotConnect
	ErrFailed
	ErrConnectFailed
	ErrAuthenticationFailed
	ErrNotPaired
	ErrNoResources
	ErrTimeout
	ErrAlreadyConnected
	ErrBusy
	ErrRejected
	ErrNotSupported
	ErrInvalidParameters
	ErrDisconnected
	ErrNotPowered
	ErrCancelled
	ErrInvalidIndex
	ErrRFKilled
	ErrAlreadyPaired
	ErrPermissionDenied
)

type Command struct {
	pkt        chan *CommandComplete
	OpCode     uint16
	Controller uint16
	Data       []byte
}

func (d *Command) Serialize() []byte {
	buf := &bytes.Buffer{}
	binary.Write(buf, binaryOrder, d.OpCode)
	binary.Write(buf, binaryOrder, d.Controller)
	binary.Write(buf, binaryOrder, uint16(len(d.Data)))
	if d.Data != nil {
		buf.Write(d.Data)
	}
	return buf.Bytes()
}

func (d *Command) Deserialize(r io.Reader) {
	pLen := uint16(0)
	binary.Read(r, binaryOrder, &d.OpCode)
	binary.Read(r, binaryOrder, &d.Controller)
	binary.Read(r, binaryOrder, &pLen)
	if pLen > 0 {
		d.Data = make([]byte, pLen)
		r.Read(d.Data)
	}
}

type CommandComplete struct {
	OpCode   uint16
	Status   byte
	Response interface{}
}

const (
	OpReadManagementVersionInformation uint16 = iota + 1
	OpReadManagementSupporteds
	OpReadControllerIndexList
	OpReadControllerInformation
	OpSetPowered
	OpSetDiscoverable
	OpSetConnectable
	OpSetFastConnectable
	OpSetBondable
	OpSetLinkSecurity
	OpSetSecureSimplePairing
	OpSetHighSpeed
	OpSetLowEnergy
	OpSetDeviceClass
	OpSetLocalName
	OpAddUUID
	OpRemoveUUID
	OpLoadLinkKeys
	OpLoadLongTermKeys
	OpDisconnect
	OpGetConnections
	OpPINCodeReply
	OpPINCodeNegativeReply
	OpSetIOCapability
	OpPairDevice
	OpCancelPairDevice
	OpUnpairDevice
	OpUserConfirmationReply
	OpUserConfirmationNegativeReply
	OpUserPasskeyReply
	OpUserPasskeyNegativeReply
	OpReadLocalOutOfBandData
	OpAddRemoteOutOfBandData
	OpRemoveRemoteOutOfBandData
	OpStartDiscovery
	OpStopDiscovery
	OpConfirmName
	OpBlockDevice
	OpUnblockDevice
	OpSetDeviceID
	OpSetAdvertising
	OpSetBREDR
	OpSetStaticAddress
	OpSetScanParameters
	OpSetSecureConnections
	OpSetDebugKeys
	OpSetPrivacy
	OpLoadIdentityResolvingKeys
	OpGetConnectionInformation
	OpGetClockInformation
	OpAddDevice
	OpRemoveDevice
	OpLoadConnectionParameters
	OpReadUnconfiguredControllerIndexList
	OpReadControllerConfigurationInformation
	OpSetExternalConfiguration
	OpSetPublicAddress
	OpStartServiceDiscovery
	OpReadLocalOutOfBandExtendedData
	OpReadExtendedControllerIndexList
	OpReadAdvertisingFeatures
	OpAddAdvertising
	OpRemoveAdvertising
	OpGetAdvertisingSizeInformation
	OpStartLimitedDiscovery
	OpReadExtendedControllerInformation
	OpSetAppearance
	OpGetPHYConfiguration
	OpSetPHYConfiguration
	OpLoadBlockedKeys
	OpSetWidebandSpeech
	OpReadControllerCapabilities
	OpReadExperimentalFeaturesInformation
	OpSetExperimentalFeature
	OpReadDefaultSystemConfiguration
	OpSetDefaultSystemConfiguration
	OpReadDefaultRuntimeConfiguration
	OpSetDefaultRuntimeConfiguration
	OpGetDeviceFlags
	OpSetDeviceFlags
	OpReadAdvertisementMonitorFeatures
	OpAddAdvertisementPatternsMonitor
	OpRemoveAdvertisementMonitor
	OpAddExtendedAdvertisingParameters
	OpAddExtendedAdvertisingData
	OpAddAdvertisementPatternsMonitorWithRSSIThreshold
	EvComplete
	EvStatus
	EvControllerError
	EvIndexAdded
	EvIndexRemoved
	EvNewSettings
	EvClassOfDeviceChanged
	EvLocalNameChanged
	EvNewLinkKey
	EvNewLongTermKey
	EvDeviceConnected
	EvDeviceDisconnected
	EvConnectFailed
	EvPINCodeRequest
	EvUserConfirmationRequest
	EvUserPasskeyRequest
	EvAuthenticationFailed
	EvDeviceFound
	EvDiscovering
	EvDeviceBlocked
	EvDeviceUnblocked
	EvDeviceUnpaired
	EvPasskeyNotify
	EvNewIdentityResolvingKey
	EvNewSignatureResolvingKey
	EvDeviceAdded
	EvDeviceRemoved
	EvNewConnectionParameter
	EvUnconfiguredIndexAdded
	EvUnconfiguredIndexRemoved
	EvNewConfigurationOptions
	EvExtendedIndexAdded
	EvExtendedIndexRemoved
	EvLocalOutOfBandExtendedDataUpdated
	EvAdvertisingAdded
	EvAdvertisingRemoved
	EvExtendedControllerInformationChanged
	EvPHYConfigurationChanged
	EvExperimentalFeatureChanged
	EvDefaultSystemConfigurationChanged
	EvDefaultRuntimeConfigurationChanged
	EvDeviceFlagsChanged
	EvAdvertisementMonitorAdded
	EvAdvertisementMonitorRemoved
	EvControllerSuspend
	EvControllerResume
)

type ReadVersion struct {
	Version  uint8
	Revision uint16
}

type ReadCommands struct {
	Commands []uint16
	Events   []uint16
}

type ReadControllerIndexList struct {
	Controllers []uint16
}

const (
	SettingPowered                 = 1
	SettingConnectable             = 1 << 1
	SettingFastConnectable         = 1 << 2
	SettingDiscoverable            = 1 << 3
	SettingBondable                = 1 << 4
	SettingLinkLevelSecurity       = 1 << 5
	SettingSecureSimplePairing     = 1 << 6
	SettingBREDR                   = 1 << 7 //Basic Rate/Enhanced Data Rate
	SettingHighSpeed               = 1 << 8
	SettingLowEnergy               = 1 << 9
	SettingAdvertising             = 1 << 10
	SettingSecureConnections       = 1 << 11
	SettingDebugKeys               = 1 << 12
	SettingPrivacy                 = 1 << 13
	SettingControllerConfiguration = 1 << 14
	SettingStaticAddress           = 1 << 15
	SettingPHYConfiguration        = 1 << 16
	SettingWidebandSpeech          = 1 << 17
)

type ReadControllerInformation struct {
	Address           [6]byte
	BluetoothVersion  byte
	Manufacturer      uint16
	SupportedSettings uint32
	CurrentSettings   uint32
	ClassOfDevice     [3]byte
	Name              [249]byte
	ShortName         [11]byte
}

type LocalName struct {
	Name      [249]byte
	ShortName [11]byte
}
