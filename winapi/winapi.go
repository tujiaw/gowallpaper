package winapi

import (
	"log"
	"syscall"
	"unsafe"
)

var ApiList = map[string][]string {
	"user32.dll": []string {
		"MessageBoxW",
		"SystemParametersInfoW",
	},
	"kernel32.dll": []string {

	},
}

var ProcCache map[string]*syscall.Proc

func init() {
	ProcCache = make(map[string]*syscall.Proc)
	for dllName, apiList := range ApiList {
		d, err := syscall.LoadDLL(dllName)
		if err != nil {
			panic(err)
		}
		defer syscall.FreeLibrary(d.Handle)
		for _, name := range apiList {
			api, err := d.FindProc(name)
			if err != nil {
				panic(err)
			}
			ProcCache[name] = api
		}
	}
}

func WinCall(name string, a ...uintptr) {
	if api, ok := ProcCache[name]; ok {
		_, _, err := api.Call(a...)
		if err != nil {
			log.Println(err)
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
