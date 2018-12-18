// +build windows

package sttub30

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
	"strings"
	"unsafe"
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

const (
	descriptorConfigLevel         = C.DESCRIPTOR_CONFIGURATION_LEVEL
	descriptorIntfAltSettingLevel = C.DESCRIPTOR_INTERFACEALTSET_LEVEL
	descriptorEndpointLevel       = C.DESCRIPTOR_ENDPOINT_LEVEL

	urbVendorDevice    = C.URB_FUNCTION_VENDOR_DEVICE
	urbVendorInterface = C.URB_FUNCTION_VENDOR_INTERFACE
	urbVendorEndpoint  = C.URB_FUNCTION_VENDOR_ENDPOINT
	urbVendorOther     = C.URB_FUNCTION_VENDOR_OTHER

	urbClassDevice    = C.URB_FUNCTION_CLASS_DEVICE
	urbClassInterface = C.URB_FUNCTION_CLASS_INTERFACE
	urbClassEndpoint  = C.URB_FUNCTION_CLASS_ENDPOINT
	urbClassOther     = C.URB_FUNCTION_CLASS_OTHER

	vendorDirectionIn  = C.VENDOR_DIRECTION_IN
	vendorDirectionOut = C.VENDOR_DIRECTION_OUT
)

type ControlPipeRequest struct {
	Function  uint16
	Direction uint64
	Request   byte
	Value     uint16
	Index     uint16
	Length    uint64
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

//type InterfaceDesc C.USB_INTERFACE_DESCRIPTOR

type InterfaceDesc struct {
	Length               byte
	DescriptorType       byte
	InterfaceNumber      byte
	AlternateString      byte
	NumEndpoints         byte
	InterfaceClass       byte
	InterfaceSubClass    byte
	InterfaceProtocol    byte
	InterfaceStringIndex byte
}

type EndpointDesc C.USB_ENDPOINT_DESCRIPTOR

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
	return strings.TrimRight(string(szStringBuf[:]), "\x00"), checkSTError(int(errno))
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

func (dev STDevice) GetNbOfInterfaces(configIndex uint) (uint, error) {
	var numInterfaces C.UINT
	errno := C.STDevice_GetNbOfInterfaces(dev.handle, C.UINT(configIndex), &numInterfaces)
	return uint(numInterfaces), checkSTError(int(errno))
}

func (dev STDevice) GetNbOfAlternateInterfaces(configIndex, interfaceIndex uint) (uint, error) {
	var numAltInterfaces C.UINT
	errno := C.STDevice_GetNbOfAlternates(dev.handle, C.UINT(configIndex), C.UINT(interfaceIndex), &numAltInterfaces)
	return uint(numAltInterfaces), checkSTError(int(errno))
}

func (dev STDevice) GetInterfaceDescriptor(configIndex, interfaceIndex, nAltSetIdx uint) (InterfaceDesc, error) {
	var desc C.USB_INTERFACE_DESCRIPTOR
	var retval InterfaceDesc
	errno := C.STDevice_GetInterfaceDescriptor(dev.handle, C.UINT(configIndex), C.UINT(interfaceIndex), C.UINT(nAltSetIdx), &desc)
	retval.Length = byte(desc.bLength)
	retval.DescriptorType = byte(desc.bDescriptorType)
	retval.InterfaceNumber = byte(desc.bInterfaceNumber)
	retval.NumEndpoints = byte(desc.bNumEndpoints)
	retval.InterfaceClass = byte(desc.bInterfaceClass)
	retval.InterfaceSubClass = byte(desc.bInterfaceSubClass)
	retval.InterfaceProtocol = byte(desc.bInterfaceProtocol)
	retval.InterfaceStringIndex = byte(desc.iInterface)
	return retval, checkSTError(int(errno))
}

func (dev STDevice) GetNbOfEndPoints(configIndex, interfaceIndex, nAltSetIdx uint) (uint, error) {
	var numEndpoints C.UINT
	errno := C.STDevice_GetNbOfEndPoints(dev.handle,
		C.UINT(configIndex),
		C.UINT(interfaceIndex),
		C.UINT(nAltSetIdx),
		&numEndpoints)
	return uint(numEndpoints), checkSTError(int(errno))
}

func (dev STDevice) GetEndPointDescriptor(configIndex, interfaceIndex, nAltSetIdx, nEndPointIdx uint) (EndpointDesc, error) {
	var desc C.USB_ENDPOINT_DESCRIPTOR
	errno := C.STDevice_GetEndPointDescriptor(dev.handle,
		C.UINT(configIndex),
		C.UINT(interfaceIndex),
		C.UINT(nAltSetIdx),
		C.UINT(nEndPointIdx),
		&desc)
	return EndpointDesc(desc), checkSTError(int(errno))
}

func (dev STDevice) GetNbOfDescriptors(nLevel, nType byte, configIndex, interfaceIndex, nAltSetIdx, nEndPointIdx uint) (uint, error) {
	var numDescriptors C.UINT
	errno := C.STDevice_GetNbOfDescriptors(dev.handle,
		C.BYTE(nLevel),
		C.BYTE(nType),
		C.UINT(configIndex),
		C.UINT(interfaceIndex),
		C.UINT(nAltSetIdx),
		C.UINT(nEndPointIdx),
		&numDescriptors)
	return uint(numDescriptors), checkSTError(int(errno))
}

func (dev STDevice) GetDescriptor(nLevel, nType byte, configIndex, interfaceIndex, nAltSetIdx, nEndPointIdx, nIdx uint) (string, error) {
	szDesc := uint(512)
	var pDesc [512]byte
	errno := C.STDevice_GetDescriptor(dev.handle,
		C.BYTE(nLevel),
		C.BYTE(nType),
		C.UINT(configIndex),
		C.UINT(interfaceIndex),
		C.UINT(nAltSetIdx),
		C.UINT(nEndPointIdx),
		C.UINT(nIdx),
		C.PBYTE(unsafe.Pointer(&pDesc[0])),
		C.UINT(szDesc))
	return strings.TrimRight(string(pDesc[:]), "\x00"), checkSTError(int(errno))
}

func (dev STDevice) SelectCurrentConfiguration(configIndex, interfaceIndex, nAltSetIdx uint) error {
	errno := C.STDevice_SelectCurrentConfiguration(dev.handle,
		C.UINT(configIndex),
		C.UINT(interfaceIndex),
		C.UINT(nAltSetIdx))
	return checkSTError(int(errno))
}

func (dev STDevice) SetDefaultTimeout(timeout int) error {
	errno := C.STDevice_SetDefaultTimeOut(dev.handle, C.DWORD(timeout))
	return checkSTError(int(errno))
}

func (dev STDevice) SetSuspendModeBehaviour(allow bool) error {
	var doAllow C.BOOL
	if allow {
		doAllow = 1
	} else {
		doAllow = 0
	}
	errno := C.STDevice_SetSuspendModeBehaviour(dev.handle, doAllow)
	return checkSTError(int(errno))
}

const (
	endpointPipeReset     = C.PIPE_RESET
	endpointAbortTransfer = C.ABORT_TRANSFER
)

func (dev STDevice) EndPointControl(nEndPointIdx, nOperation uint) error {
	errno := C.STDevice_EndPointControl(dev.handle,
		C.UINT(nEndPointIdx),
		C.UINT(nOperation))
	return checkSTError(int(errno))
}

func (dev STDevice) Reset() error {
	errno := C.STDevice_Reset(dev.handle)
	return checkSTError(int(errno))
}

func (dev STDevice) ControlPipeRequest(Request ControlPipeRequest, Data []byte) error {
	var req C.CNTRPIPE_RQ

	req.Function = C.USHORT(Request.Function)
	req.Direction = C.ULONG(Request.Direction)
	req.Request = C.UCHAR(Request.Request)
	req.Value = C.USHORT(Request.Value)
	req.Index = C.USHORT(Request.Index)
	req.Length = C.ULONG(Request.Length)

	//var buffer [512]byte

	//copy(buffer[:], Data)

	if uint64(len(Data)) < Request.Length {
		return fmt.Errorf("Data buffer too small, must be at least as large as Request.Length.")
	}

	errno := C.STDevice_ControlPipeRequest(dev.handle, &req, C.PBYTE(unsafe.Pointer(&Data[0])))
	return checkSTError(int(errno))
}
