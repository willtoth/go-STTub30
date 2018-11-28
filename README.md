# go-STTub30
Go interface to STTubeDevice driver which implements driver for ST devices. This can be used instead of libusb to talk to devices using STTub30.sys driver since Windows 10 installs it by default. This driver is used for the ST based ROM bootloader and is used by several different projects. Instead of using libusb to implement the driver, which would require changing what driver Windows loads for generic ST DFU devices, we can use their driver directionly.

The ST source implementation can be found here: [DfuSe v3.0.6](https://www.st.com/en/development-tools/stsw-stm32080.html).

This project is Windows only

# Setup

Ideally these dependancies would be included with the repo so a single `go get` will pull in everything, but fir licensing concerns for the Windows header and DfuSe files.

1) For 64-bit builds, the DLL must be built in 64-bit mode and can be done with Visual Studio 2017 Community. The .dll and .lib file are placed in the root directory of the project
2) Path for minwindef.h added, or file copied to root directory
3) Sources/ folder from DfuSe V3.0.6 install copied to root directory, note: gcc/mingw don't like spaces in the include paths, cgo doesn't like spaces or '(' ')' characters for the default install location of that package.
