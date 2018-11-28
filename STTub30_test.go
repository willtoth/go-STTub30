// +build windows

package sttub30

import (
	"testing"

	"github.com/willtoth/setupapi"
)

func TestWithDevice(t *testing.T) {
	//GUID of STM32F3 DFU Driver
	guid := setupapi.Guid{0x3fe809ab, 0xfb91, 0x4cb5, [8]byte{0xa6, 0x43, 0x69, 0x67, 0x0d, 0x52, 0x36, 0x6e}}
	devInfo, err := setupapi.SetupDiGetClassDevsEx(guid, "", 0, setupapi.Present|setupapi.InterfaceDevice, 0, "", 0)
	if err != nil {
		t.Errorf("Error get class devs ex: %v", err)
		return
	}

	devPath, err := devInfo.DevicePath(guid)
	if err != nil {
		t.Errorf("Error device path: %s", err.Error())
		return
	}

	dev, err := Open(devPath)
	if err != nil {
		t.Errorf("Failed to open device: %v", err)
		return
	}
	defer func() {
		err = dev.Close()

		if err != nil {
			t.Errorf("%v", err)
		}
	}()

	desc, err := dev.GetDeviceDescriptor()

	if err != nil {
		t.Errorf("Failed to get device descriptor: %v", err)
		return
	}

	manStr, err := dev.GetStringDescriptor(uint(desc.iManufacturer))

	if err != nil {
		t.Errorf("Failed to get device string descriptor: %v", err)
		return
	}

	if manStr != "STMicroelectronics" {
		t.Errorf("iManufacturer was incorrect: %s\n", manStr)
	}
}
