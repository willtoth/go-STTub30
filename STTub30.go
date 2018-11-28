package STTub30

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

const (
	STDEVICE_ERROR_OFFSET            = 0x12340000
	STDEVICE_NOERROR                 = STDEVICE_ERROR_OFFSET
	STDEVICE_MEMORY                  = (STDEVICE_ERROR_OFFSET + 1)
	STDEVICE_BADPARAMETER            = (STDEVICE_ERROR_OFFSET + 2)
	STDEVICE_NOTIMPLEMENTED          = (STDEVICE_ERROR_OFFSET + 3)
	STDEVICE_ENUMFINISHED            = (STDEVICE_ERROR_OFFSET + 4)
	STDEVICE_OPENDRIVERERROR         = (STDEVICE_ERROR_OFFSET + 5)
	STDEVICE_ERRORDESCRIPTORBUILDING = (STDEVICE_ERROR_OFFSET + 6)
	STDEVICE_PIPECREATIONERROR       = (STDEVICE_ERROR_OFFSET + 7)
	STDEVICE_PIPERESETERROR          = (STDEVICE_ERROR_OFFSET + 8)
	STDEVICE_PIPEABORTERROR          = (STDEVICE_ERROR_OFFSET + 9)
	STDEVICE_STRINGDESCRIPTORERROR   = (STDEVICE_ERROR_OFFSET + 0xA)
	STDEVICE_DRIVERISCLOSED          = (STDEVICE_ERROR_OFFSET + 0xB)
	STDEVICE_VENDOR_RQ_PB            = (STDEVICE_ERROR_OFFSET + 0xC)
	STDEVICE_ERRORWHILEREADING       = (STDEVICE_ERROR_OFFSET + 0xD)
	STDEVICE_ERRORBEFOREREADING      = (STDEVICE_ERROR_OFFSET + 0xE)
	STDEVICE_ERRORWHILEWRITING       = (STDEVICE_ERROR_OFFSET + 0xF)
	STDEVICE_ERRORBEFOREWRITING      = (STDEVICE_ERROR_OFFSET + 0x10)
	STDEVICE_DEVICERESETERROR        = (STDEVICE_ERROR_OFFSET + 0x11)
	STDEVICE_CANTUSEUNPLUGEVENT      = (STDEVICE_ERROR_OFFSET + 0x12)
	STDEVICE_INCORRECTBUFFERSIZE     = (STDEVICE_ERROR_OFFSET + 0x13)
	STDEVICE_DESCRIPTORNOTFOUND      = (STDEVICE_ERROR_OFFSET + 0x14)
	STDEVICE_PIPESARECLOSED          = (STDEVICE_ERROR_OFFSET + 0x15)
	STDEVICE_PIPESAREOPEN            = (STDEVICE_ERROR_OFFSET + 0x16)
)

var STErrorStrings = map[int](string){
	STDEVICE_NOERROR:                 "No Error",
	STDEVICE_MEMORY:                  "Memory",
	STDEVICE_BADPARAMETER:            "Bad Parameter",
	STDEVICE_NOTIMPLEMENTED:          "Not Implemented",
	STDEVICE_ENUMFINISHED:            "Enum Finished",
	STDEVICE_OPENDRIVERERROR:         "Open Driver Error",
	STDEVICE_ERRORDESCRIPTORBUILDING: "Error Descriptor Building",
	STDEVICE_PIPECREATIONERROR:       "Pipe Creation Error",
	STDEVICE_PIPERESETERROR:          "Pipe Reset Error",
	STDEVICE_PIPEABORTERROR:          "Pipe Abort Error",
	STDEVICE_STRINGDESCRIPTORERROR:   "String Descriptor Error",
	STDEVICE_DRIVERISCLOSED:          "Driver is closed",
	STDEVICE_VENDOR_RQ_PB:            "Vendor RQ PB",
	STDEVICE_ERRORWHILEREADING:       "Error While Reading",
	STDEVICE_ERRORBEFOREREADING:      "Error Before Reading",
	STDEVICE_ERRORWHILEWRITING:       "Error While Writing",
	STDEVICE_ERRORBEFOREWRITING:      "Error Before Writing",
	STDEVICE_DEVICERESETERROR:        "Device Set Error",
	STDEVICE_CANTUSEUNPLUGEVENT:      "Cant Use Unplug Event",
	STDEVICE_INCORRECTBUFFERSIZE:     "Incorrect Buffer Size",
	STDEVICE_DESCRIPTORNOTFOUND:      "Descriptor Not Found",
	STDEVICE_PIPESARECLOSED:          "Pipes are closed",
	STDEVICE_PIPESAREOPEN:            "Pipes are open",
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

type STDevice C.HANDLE
type DeviceDesc C.USB_DEVICE_DESCRIPTOR
type DeviceConfig C.USB_CONFIGURATION_DESCRIPTOR

var stErrorString = map[int]string{}

func checkSTError(code int) error {
	if code != STDEVICE_NOERROR {
		return fmt.Errorf("%s", STErrorStrings[code])
	}
	return nil
}

func STDeviceOpen(devicePath string, device STDevice) error {
	handle := C.HANDLE(device)
	errno := C.STDevice_Open(C.CString(devicePath), C.LPHANDLE(&handle), C.LPHANDLE(C.NULL))
	return checkSTError(int(errno))
}

func STDeviceOpenPipes(device STDevice) error {
	errno := C.STDevice_OpenPipes(C.HANDLE(device))
	return checkSTError(int(errno))
}

func STDeviceClosePipes(device STDevice) error {
	errno := C.STDevice_ClosePipes(C.HANDLE(device))
	return checkSTError(int(errno))
}

func STDeviceClose(device STDevice) error {
	errno := C.STDevice_Close(C.HANDLE(device))
	return checkSTError(int(errno))
}

func STDeviceGetStringDescriptor(device STDevice, nIndex uint) (string, error) {
	nStringLength := uint(512)
	var szStringBuf [512]byte
	errno := C.STDevice_GetStringDescriptor(C.HANDLE(device), C.UINT(nIndex), C.LPSTR(unsafe.Pointer(&szStringBuf[0])), C.UINT(nStringLength))
	return "", checkSTError(int(errno))
}

func STDeviceGetDeviceDescriptor(device STDevice) (DeviceDesc, error) {
	var pDesc C.USB_DEVICE_DESCRIPTOR
	errno := C.STDevice_GetDeviceDescriptor(C.HANDLE(device), C.PUSB_DEVICE_DESCRIPTOR(&pDesc))
	desc := DeviceDesc(pDesc)
	return desc, checkSTError(int(errno))
}

func STDevice_GetNbOfConfigurations(device STDevice) (uint, error) {
	var numConfigs C.UINT
	errno := C.STDevice_GetNbOfConfigurations(C.HANDLE(device), &numConfigs)
	return uint(numConfigs), checkSTError(int(errno))
}

func STDeviceGetConfigurationDescriptor(device STDevice, nConfigIdx uint) (DeviceConfig, error) {
	var cfg C.USB_CONFIGURATION_DESCRIPTOR
	errno := C.STDevice_GetConfigurationDescriptor(C.HANDLE(device), C.UINT(nConfigIdx), &cfg)
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

	fmt.Printf("Test: %s\r\n", devPath)

	var handle C.HANDLE

	val1 := C.STDevice_Open(C.CString(devPath), C.LPHANDLE(&handle), C.LPHANDLE(C.NULL))

	C.STDevice_Close(C.HANDLE(handle))

	fmt.Println(val1)

}
