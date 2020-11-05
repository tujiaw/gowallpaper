package winapi

import (
	"log"
	"strings"
	"syscall"
	"unsafe"
)

var ApiList = map[string][]string{
	"user32.dll": {
		"MessageBoxW",
		"SystemParametersInfoW",
	},
	"kernel32.dll": {},
}

var ProcCache map[string]*syscall.Proc

func init() {
	ProcCache = make(map[string]*syscall.Proc)
	for dllName, apiList := range ApiList {
		d, err := syscall.LoadDLL(dllName)
		if err != nil {
			panic(err)
		}
		for _, name := range apiList {
			api, err := d.FindProc(name)
			if err != nil {
				log.Println(name, err)
			}
			ProcCache[name] = api
		}
		_ = syscall.FreeLibrary(d.Handle)
	}
}

func WinCall(name string, a ...uintptr) {
	if api, ok := ProcCache[name]; ok {
		_, _, err := api.Call(a...)
		if err != nil && !strings.Contains(err.Error(), "timeout period expired") {
			log.Println("api.Call", err)
		}
	} else {
		log.Println("api not found, name:", name)
	}
}

func IntPtr(n int) uintptr {
	return uintptr(n)
}

func StrPtr(s string) uintptr {
	p, _ := syscall.UTF16PtrFromString(s)
	return uintptr(unsafe.Pointer(p))
}

func ShowMessage(title, text string) {
	WinCall("MessageBoxW", IntPtr(0), StrPtr(text), StrPtr(title), IntPtr(0))
}

func SetWallpaper(bmpPath string) {
	WinCall("SystemParametersInfoW", IntPtr(20), IntPtr(0), StrPtr(bmpPath), IntPtr(3))
}
