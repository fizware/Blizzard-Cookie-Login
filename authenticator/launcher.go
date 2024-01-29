package authenticator

import (
	"errors"
	"github.com/mitchellh/go-ps"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
	"os/exec"
	"syscall"
)

var battleLocation string

func closeAndFindBattleNet() error {
	processes, err := ps.Processes()
	if err != nil {
		return err
	}
	for _, process := range processes {
		name := process.Executable()
		if name == "Battle.net.exe" {
			battleLocation, err = getExecutablePath(uint32(process.Pid()))
			if err != nil {
				return errors.New("unable to locate battle.net launcher")
			}
			hProc, _ := windows.OpenProcess(windows.PROCESS_TERMINATE|windows.PROCESS_QUERY_INFORMATION, false, uint32(process.Pid()))
			err = windows.TerminateProcess(hProc, 0)
			if err != nil {
				return errors.New("unable to terminate Battle.Net.exe")
			}
		}
		if name == "Agent.exe" {
			hProc, err := windows.OpenProcess(windows.PROCESS_TERMINATE, false, uint32(process.Pid()))
			if err != nil {
				return errors.New("unable to gain a handle to Agent.exe")
			}
			err = windows.TerminateProcess(hProc, 0)
			if err != nil {
				return errors.New("unable to terminate Agent.exe")
			}
			_ = windows.Close(hProc)
		}
	}
	if battleLocation == "" {
		return errors.New("please open battle.net")
	}
	return nil
}

func startBattleNet() error {
	return exec.Command(battleLocation).Start()
}

func logoutBattleNet() {
	_ = registry.DeleteKey(registry.CURRENT_USER, "SOFTWARE\\Blizzard Entertainment\\Battle.net\\UnifiedAuth")
}

func getExecutablePath(process uint32) (string, error) {
	hProc, err := windows.OpenProcess(windows.PROCESS_TERMINATE|windows.PROCESS_QUERY_INFORMATION, false, process)
	if err != nil {
		return "", err
	}
	defer func(fd windows.Handle) {
		_ = windows.Close(fd)
	}(hProc)
	pathPtr := make([]uint16, windows.MAX_PATH)
	length := uint32(windows.MAX_PATH)
	err = windows.QueryFullProcessImageName(hProc, 0, &pathPtr[0], &length)
	if err != nil {
		return "", err
	}
	path := syscall.UTF16ToString(pathPtr)
	return path, nil
}
