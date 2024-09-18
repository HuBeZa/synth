//go:build windows
// +build windows

package keyboard

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	spi_getKeyboardDelay = 0x0016
	spi_setKeyboardDelay = 0x0017
)

var (
	user32DLL = windows.NewLazyDLL("user32.dll")

	// see https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-systemparametersinfoa
	//
	// BOOL SystemParametersInfoA(
	// 	[in]      UINT  uiAction,
	// 	[in]      UINT  uiParam,
	// 	[in, out] PVOID pvParam,
	// 	[in]      UINT  fWinIni
	//   );
	procSystemParamInfo = user32DLL.NewProc("SystemParametersInfoA")
)

func init() {
	origDelay, err := getKeyboardDelay()
	if err != nil || origDelay == 0 {
		return
	}

	restoreConfigOnAppClose(origDelay)
	setKeyboardDelay(0)
}

func restoreConfigOnAppClose(origKeyboardDelay int) {
	c := make(chan os.Signal, 1)
	signal.Notify(c,
		syscall.SIGTERM, // close gracefully
		syscall.SIGINT,  // Ctrl+C
		syscall.SIGQUIT, // Ctrl-\
		syscall.SIGKILL, // fatal
		syscall.SIGHUP,  // terminal disconnected
	)

	go func() {
		<-c
		fmt.Println("restore config")
		setKeyboardDelay(origKeyboardDelay)
		time.Sleep(10 * time.Second)
	}()

}

func getKeyboardDelay() (int, error) {
	var keyboardDelay int
	r1, _, lastErr := procSystemParamInfo.Call(spi_getKeyboardDelay, 0, uintptr(unsafe.Pointer(&keyboardDelay)), 0)
	if r1 == 1 && lastErr.(syscall.Errno) == windows.NOERROR {
		return keyboardDelay, nil
	}
	return -1, fmt.Errorf("windows error: %v", lastErr)
}

func setKeyboardDelay(keyboardDelay int) error {
	r1, _, lastErr := procSystemParamInfo.Call(spi_setKeyboardDelay, uintptr(keyboardDelay), 0, 0)
	if r1 == 1 && lastErr.(syscall.Errno) == windows.NOERROR {
		return nil
	}
	return fmt.Errorf("windows error: %v", lastErr)
}
