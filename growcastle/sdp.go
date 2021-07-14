package growcastle

import (
	"encoding/hex"
	"encoding/xml"
)

func KeyboardDescriptor(reportId byte) []byte {
	return []byte{
		0x05, 0x01, // Usage Page (Generic Desktop Ctrls)
		0x09, 0x06, // Usage (Keyboard)
		0xA1, 0x01, // Collection (Application)
		0x85, KeyboardReportId, //   Report ID (1)
		0x75, 0x01, //   Report Size (1)
		0x95, 0x08, //   Report Count (8)
		0x05, 0x07, //   Usage Page (Kbrd/Keypad)
		0x19, 0xE0, //   Usage Minimum (0xE0)
		0x29, 0xE7, //   Usage Maximum (0xE7)
		0x15, 0x00, //   Logical Minimum (0)
		0x25, 0x01, //   Logical Maximum (1)
		0x81, 0x02, //   Input (Data,Var,Abs,No Wrap,Linear,Preferred State,No Null Position)
		0x95, 0x01, //   Report Count (1)
		0x75, 0x08, //   Report Size (8)
		0x81, 0x03, //   Input (Const,Var,Abs,No Wrap,Linear,Preferred State,No Null Position)
		0x95, 0x05, //   Report Count (5)
		0x75, 0x01, //   Report Size (1)
		0x05, 0x08, //   Usage Page (LEDs)
		0x19, 0x01, //   Usage Minimum (Num Lock)
		0x29, 0x05, //   Usage Maximum (Kana)
		0x91, 0x02, //   Output (Data,Var,Abs,No Wrap,Linear,Preferred State,No Null Position,Non-volatile)
		0x95, 0x01, //   Report Count (1)
		0x75, 0x03, //   Report Size (3)
		0x91, 0x03, //   Output (Const,Var,Abs,No Wrap,Linear,Preferred State,No Null Position,Non-volatile)
		0x95, 0x06, //   Report Count (6)
		0x75, 0x08, //   Report Size (8)
		0x15, 0x00, //   Logical Minimum (0)
		0x26, 0xFF, 0x00, //   Logical Maximum (255)
		0x05, 0x07, //   Usage Page (Kbrd/Keypad)
		0x19, 0x00, //   Usage Minimum (0x00)
		0x29, 0xFF, //   Usage Maximum (0xFF)
		0x81, 0x00, //   Input (Data,Array,Abs,No Wrap,Linear,Preferred State,No Null Position)
		0xC0, // End Collection
	}
}

// TouchScreenDescriptor
// tip switch     byte
// press          byte
// x              int16
// y              int16
// width          touch width
// height         touch height
// contact identifier       int16
// contact count maximum    byte
func TouchScreenDescriptor() []byte {
	return []byte{
		0x05, 0x0D, // Usage Page (Digitizer)
		0x09, 0x04, // Usage (Touch Screen)
		0xA1, 0x01, // Collection (Application)
		0x85, TouchScreenReportId, //   Report ID (2)
		0x09, 0x22, //   Usage (Finger)
		0xA1, 0x00, //   Collection (Physical)
		0x09, 0x42, //     Usage (Tip Switch)
		0x15, 0x00, //     Logical Minimum (0)
		0x25, 0x01, //     Logical Maximum (1)
		0x75, 0x01, //     Report Size (1)
		0x95, 0x01, //     Report Count (1)
		0x81, 0x02, //     Input (Data,Var,Abs,No Wrap,Linear,Preferred State,No Null Position)
		0x95, 0x03, //     Report Count (3)
		0x81, 0x03, //     Input (Const,Var,Abs,No Wrap,Linear,Preferred State,No Null Position)
		0x09, 0x32, //     Usage (In Range)
		0x09, 0x47, //     Usage (0x47)
		0x95, 0x02, //     Report Count (2)
		0x81, 0x02, //     Input (Data,Var,Abs,No Wrap,Linear,Preferred State,No Null Position)
		0x95, 0x0A, //     Report Count (10)
		0x81, 0x03, //     Input (Const,Var,Abs,No Wrap,Linear,Preferred State,No Null Position)
		0x05, 0x01, //     Usage Page (Generic Desktop Ctrls)
		0x26, 0xFF, 0x7F, //     Logical Maximum (32767)
		0x75, 0x10, //     Report Size (16)
		0x95, 0x01, //     Report Count (1)
		0x09, 0x30, //     Usage (X)
		0x81, 0x02, //     Input (Data,Var,Abs,No Wrap,Linear,Preferred State,No Null Position)
		0x09, 0x31, //     Usage (Y)
		0x81, 0x02, //     Input (Data,Var,Abs,No Wrap,Linear,Preferred State,No Null Position)
		0x05, 0x0D, //     Usage Page (Digitizer)
		0x09, 0x48, //     Usage (0x48)
		0x09, 0x49, //     Usage (0x49)
		0x95, 0x02, //     Report Count (2)
		0x81, 0x02, //     Input (Data,Var,Abs,No Wrap,Linear,Preferred State,No Null Position)
		0x09, 0x51, //     Usage (0x51)
		0x95, 0x01, //     Report Count (1)
		0x81, 0x02, //     Input (Data,Var,Abs,No Wrap,Linear,Preferred State,No Null Position)
		0x09, 0x55, //     Usage (0x55)
		0x25, 0x08, //     Logical Maximum (8)
		0x75, 0x08, //     Report Size (8)
		0xB1, 0x02, //     Feature (Data,Var,Abs,No Wrap,Linear,Preferred State,No Null Position,Non-volatile)
		0xC0, //   End Collection
		0xC0, // End Collection
	}
}

func MouseDescriptor() []byte {
	return []byte{
		0x05, 0x01, // USAGE_PAGE (Generic Desktop)     0
		0x09, 0x02, // USAGE (Mouse)                    2
		0xa1, 0x01, // COLLECTION (Application)         4
		0x85, MouseReportId, //   REPORT_ID (Mouse)              6
		0x09, 0x01, //   USAGE (Pointer)                8
		0xa1, 0x00, //   COLLECTION (Physical)          10
		0x05, 0x09, //     USAGE_PAGE (Button)          12
		0x19, 0x01, //     USAGE_MINIMUM (Button 1)     14
		0x29, 0x02, //     USAGE_MAXIMUM (Button 2)     16
		0x15, 0x00, //     LOGICAL_MINIMUM (0)          18
		0x25, 0x01, //     LOGICAL_MAXIMUM (1)          20
		0x75, 0x01, //     REPORT_SIZE (1)              22
		0x95, 0x02, //     REPORT_COUNT (2)             24
		0x81, 0x02, //     INPUT (Data,Var,Abs)         26
		0x95, 0x06, //     REPORT_COUNT (6)             28
		0x81, 0x03, //     INPUT (Cnst,Var,Abs)         30
		0x05, 0x01, //     USAGE_PAGE (Generic Desktop) 32
		0x09, 0x30, //     USAGE (X)                    34
		0x09, 0x31, //     USAGE (Y)                    36
		0x15, 0x81, //     LOGICAL_MINIMUM (-127)       38
		0x25, 0x7f, //     LOGICAL_MAXIMUM (127)        40
		0x75, 0x08, //     REPORT_SIZE (8)              42
		0x95, 0x02, //     REPORT_COUNT (2)             44
		0x81, 0x06, //     INPUT (Data,Var,Rel)         46
		0xc0, //   END_COLLECTION                 48
		0xc0, // END_COLLECTION                   49/50
	}
}

// SDPRecord see https://btprodspecificationrefs.blob.core.windows.net/assigned-numbers/Assigned%20Number%20Types/Service%20Discovery.pdf
// section Human Interface Device Profile
func SDPRecord(descriptor [][]byte) (string, error) {
	var records []interface{}

	// ServiceClassIDList
	records = append(records, Attribute{
		Id: "0x0001",
		Value: Sequence{
			Value: []interface{}{
				UUID{Value: "0x1124"}, // HID profile
			},
		},
	})

	// ProtocolDescriptorList
	records = append(records, Attribute{
		Id: "0x0004",
		Value: Sequence{
			Value: []interface{}{
				Sequence{
					Value: []interface{}{
						UUID{Value: "0x0100"},   // means this profile is based on top of L2CAP
						UInt16{Value: "0x0011"}, // refers to the PSM for HID Control
					},
				},
				Sequence{
					Value: []interface{}{
						UUID{Value: "0x0011"}, // refers to the Protocol Identifier's UUID
					},
				},
			},
		},
	})

	// browse group visibility
	records = append(records, Attribute{
		Id: "0x0005",
		Value: Sequence{
			Value: []interface{}{
				UUID{Value: "0x1002"},
			},
		},
	})

	// LanguageBaseAttributeIDList
	records = append(records, Attribute{
		Id: "0x0006",
		Value: Sequence{
			Value: []interface{}{
				UInt16{Value: "0x656e"}, // "en" - English
				UInt16{Value: "0x006a"}, // UTF-8 encoding
				UInt16{Value: "0x0100"}, // PrimaryLanguageBaseId = 0
			},
		},
	})

	// BluetoothProfileDescriptorList
	records = append(records, Attribute{
		Id: "0x0009",
		Value: Sequence{
			Value: []interface{}{
				Sequence{
					Value: []interface{}{
						UUID{Value: "0x1124"},   // refers to the Protocol Identifier's UUID for the HID profile
						UInt16{Value: "0x0101"}, // is to define the version to 1.1
					},
				},
			},
		},
	})

	// AdditionalProtocolDescriptorLists
	records = append(records, Attribute{
		Id: "0x000d",
		Value: Sequence{
			Value: []interface{}{
				Sequence{
					Value: []interface{}{
						Sequence{
							Value: []interface{}{
								UUID{Value: "0x0100"},   // L2CAP
								UInt16{Value: "0x0013"}, // refers to the PSM for HID Interrupt
							},
						},
						Sequence{
							Value: []interface{}{
								UUID{Value: "0x0011"}, // refers to the Protocol Identifier's UUID for the HID profile
							},
						},
					},
				},
			},
		},
	})

	// HIDParserVersion
	records = append(records, Attribute{
		Id:    "0x0201",
		Value: UInt16{Value: "0x0111"},
	})

	// HIDDeviceSubclass
	records = append(records, Attribute{
		Id:    "0x0202",
		Value: UInt8{Value: "0xC0"},
	})

	// HIDCountryCode
	records = append(records, Attribute{
		Id:    "0x0203",
		Value: UInt8{Value: "0x00"},
	})

	// HIDVirtualCable
	records = append(records, Attribute{
		Id:    "0x0204",
		Value: Boolean{Value: "false"},
	})

	// HIDReconnectInitiate
	records = append(records, Attribute{
		Id:    "0x0205",
		Value: Boolean{Value: "false"},
	})

	var descriptorValue []interface{}
	for i := 0; i < len(descriptor); i++ {
		d := descriptor[i]
		descriptorValue = append(descriptorValue, Sequence{
			Value: []interface{}{
				UInt8{Value: "0x22"},
				Text{
					Encoding: "hex",
					Value:    hex.EncodeToString(d),
				},
			},
		})
	}

	// HIDDescriptorList
	records = append(records, Attribute{
		Id: "0x0206",
		Value: Sequence{
			Value: descriptorValue,
		},
	})

	// HIDLANGIDBaseList
	records = append(records, Attribute{
		Id: "0x0207",
		Value: Sequence{
			Value: []interface{}{
				Sequence{
					Value: []interface{}{
						UInt16{Value: "0x0409"}, // for en_US
						UInt16{Value: "0x0100"},
					},
				},
			},
		},
	})

	// HIDBatteryPower
	records = append(records, Attribute{
		Id:    "0x0209",
		Value: Boolean{Value: "false"},
	})

	// HIDRemoteWake
	records = append(records, Attribute{
		Id:    "0x020a",
		Value: Boolean{Value: "true"},
	})

	// HIDSupervisionTimeout
	records = append(records, Attribute{
		Id:    "0x020c",
		Value: UInt16{Value: "0x0c80"},
	})

	// HIDNormallyConnectable
	records = append(records, Attribute{
		Id:    "0x020d",
		Value: Boolean{Value: "false"},
	})

	// HIDBootDevice
	records = append(records, Attribute{
		Id:    "0x020e",
		Value: Boolean{Value: "false"},
	})

	// HIDSSRHostMaxLatency
	records = append(records, Attribute{
		Id:    "0x020e",
		Value: UInt16{Value: "0x0640"},
	})

	// HIDSSRHostMinTimeout
	records = append(records, Attribute{
		Id:    "0x0210",
		Value: UInt16{Value: "0x0320"},
	})

	record := Record{}
	record.Records = records

	x, err := xml.Marshal(record)
	if err != nil {
		return "", err
	}

	return xml.Header + string(x), nil
}
