// +build windows
package main

/*
#cgo CFLAGS: -I./Sources/STTubeDevice -I./Sources/Include
#cgo LDFLAGS: -L./ -lSTTubeDevice30

#include <stdlib.h>

#include "minwindef.h"

#include "usb100.h"
#include "STTubeDeviceErr30.h"
#include "STTubeDeviceTyp30.h"
#include "STTubeDeviceFun30.h"
*/
import "C"

import (
	"fmt"
	"unsafe"

	"github.com/willtoth/setupapi"
)

var STErrorStrings = map[int](string){
	C.STDEVICE_NOERROR:                 "No Error",
	C.STDEVICE_MEMORY:                  "Memory",
	C.STDEVICE_BADPARAMETER:            "Bad Parameter",
	C.STDEVICE_NOTIMPLEMENTED:          "Not Implemented",
	C.STDEVICE_ENUMFINISHED:            "Enum Finished",
	C.STDEVICE_OPENDRIVERERROR:         "Open Driver Error",
	C.STDEVICE_ERRORDESCRIPTORBUILDING: "Error Descriptor Building",
	C.STDEVICE_PIPECREATIONERROR:       "Pipe Creation Error",
	C.STDEVICE_PIPERESETERROR:          "Pipe Reset Error",
	C.STDEVICE_PIPEABORTERROR:          "Pipe Abort Error",
	C.STDEVICE_STRINGDESCRIPTORERROR:   "String Descriptor Error",
	C.STDEVICE_DRIVERISCLOSED:          "Driver is closed",
	C.STDEVICE_VENDOR_RQ_PB:            "Vendor RQ PB",
	C.STDEVICE_ERRORWHILEREADING:       "Error While Reading",
	C.STDEVICE_ERRORBEFOREREADING:      "Error Before Reading",
	C.STDEVICE_ERRORWHILEWRITING:       "Error While Writing",
	C.STDEVICE_ERRORBEFOREWRITING:      "Error Before Writing",
	C.STDEVICE_DEVICERESETERROR:        "Device Set Error",
	C.STDEVICE_CANTUSEUNPLUGEVENT:      "Cant Use Unplug Event",
	C.STDEVICE_INCORRECTBUFFERSIZE:     "Incorrect Buffer Size",
	C.STDEVICE_DESCRIPTORNOTFOUND:      "Descriptor Not Found",
	C.STDEVICE_PIPESARECLOSED:          "Pipes are closed",
	C.STDEVICE_PIPESAREOPEN:            "Pipes are open",
}

/* TODO: would be better to convert to pure go:
type DeviceDesc struct {
	Length             byte
	DescType           byte
	BcdUSB             uint16
	Class              byte
	SubClass           byte
	Protocol           byte
	MaxPacketSize      byte
	VendorID           uint16
	ProductID          uint16
	BcdDevice          uint16
	iManufacturer      byte
	iProduct           byte
	iSerialNumber      byte
	bNumConfigurations byte
}
*/

type DeviceDesc C.USB_DEVICE_DESCRIPTOR
type DeviceConfig C.USB_CONFIGURATION_DESCRIPTOR

func checkSTError(code int) error {
	if code != C.STDEVICE_NOERROR {
		return fmt.Errorf("%s", STErrorStrings[code])
	}
	return nil
}

type STDevice struct {
	handle C.HANDLE
}

func Open(devicePath string) (STDevice, error) {
	var device STDevice
	errno := C.STDevice_Open(C.CString(devicePath), &device.handle, C.LPHANDLE(C.NULL))
	return device, checkSTError(int(errno))
}

func (dev STDevice) OpenPipes(device STDevice) error {
	errno := C.STDevice_OpenPipes(dev.handle)
	return checkSTError(int(errno))
}

func (dev STDevice) ClosePipes() error {
	errno := C.STDevice_ClosePipes(dev.handle)
	return checkSTError(int(errno))
}

func (dev STDevice) Close() error {
	errno := C.STDevice_Close(dev.handle)
	return checkSTError(int(errno))
}

func (dev STDevice) GetStringDescriptor(nIndex uint) (string, error) {
	nStringLength := uint(512)
	var szStringBuf [512]byte
	errno := C.STDevice_GetStringDescriptor(dev.handle, C.UINT(nIndex), C.LPSTR(unsafe.Pointer(&szStringBuf[0])), C.UINT(nStringLength))
	return "", checkSTError(int(errno))
}

func (dev STDevice) GetDeviceDescriptor() (DeviceDesc, error) {
	var pDesc C.USB_DEVICE_DESCRIPTOR
	errno := C.STDevice_GetDeviceDescriptor(dev.handle, C.PUSB_DEVICE_DESCRIPTOR(&pDesc))
	desc := DeviceDesc(pDesc)
	return desc, checkSTError(int(errno))
}

func (dev STDevice) GetNbOfConfigurations() (uint, error) {
	var numConfigs C.UINT
	errno := C.STDevice_GetNbOfConfigurations(dev.handle, &numConfigs)
	return uint(numConfigs), checkSTError(int(errno))
}

func (dev STDevice) GetConfigurationDescriptor(nConfigIdx uint) (DeviceConfig, error) {
	var cfg C.USB_CONFIGURATION_DESCRIPTOR
	errno := C.STDevice_GetConfigurationDescriptor(dev.handle, C.UINT(nConfigIdx), &cfg)
	return DeviceConfig(cfg), checkSTError(int(errno))
}

func main() {
	//GUID of STM32F3 DFU Driver
	guid := setupapi.Guid{0x3fe809ab, 0xfb91, 0x4cb5, [8]byte{0xa6, 0x43, 0x69, 0x67, 0x0d, 0x52, 0x36, 0x6e}}
	devInfo, err := setupapi.SetupDiGetClassDevsEx(guid, "", 0, setupapi.Present|setupapi.InterfaceDevice, 0, "", 0)
	if err != nil {
		fmt.Printf("Error get class devs ex: %v", err)
		return
	}

	devPath, err := devInfo.DevicePath(guid)
	if err != nil {
		fmt.Printf("Error device path: %s", err.Error())
		return
	}

	dev, err := Open(devPath)

	if err != nil {
		fmt.Printf("Failed to open device: %v", err)
		return
	}

	desc, err := dev.GetDeviceDescriptor()

	if err != nil {
		fmt.Printf("Failed to get device descriptor: %v", err)
		return
	}

	fmt.Printf("%s\r\n", desc)

	err = dev.Close()

	if err != nil {
		fmt.Println(err)
	}

}
