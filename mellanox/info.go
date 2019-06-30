package mellanox

// Opening up some functionality from: https://github.com/Mellanox/sriovnet/blob/master/sriovnet_helper.go
// This file also contains additional functionality

import (
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	NetSysDir    = "/sys/class/net"
	PcidevPrefix = "device"

	NetDevMaxVfCountFile     = "sriov_totalvfs"
	NetDevCurrentVfCountFile = "sriov_numvfs"
)

//GetNetDevDeviceDir returns the device dir for a given PF device name
func GetNetDevDeviceDir(netDevName string) string {
	devDirName := filepath.Join(NetSysDir, netDevName, PcidevPrefix)
	return devDirName
}

//GetMaxVfCount returns the maximum VFs for a given PF device
func GetMaxVfCount(pfNetdevName string) (int, error) {
	devDirName := GetNetDevDeviceDir(pfNetdevName)

	maxDevFile := FileObject{
		Path: filepath.Join(devDirName, NetDevMaxVfCountFile),
	}

	maxVfs, err := maxDevFile.ReadInt()
	if err != nil {
		return 0, err
	} else {
		log.Println("max_vfs = ", maxVfs)
		return maxVfs, nil
	}
}

//GetCurrentVfCount gets the current VF count for a given PF device name
func GetCurrentVfCount(pfNetdevName string) (int, error) {
	devDirName := GetNetDevDeviceDir(pfNetdevName)

	maxDevFile := FileObject{
		Path: filepath.Join(devDirName, NetDevCurrentVfCountFile),
	}

	curVfs, err := maxDevFile.ReadInt()
	if err != nil {
		return 0, err
	} else {
		return curVfs, nil
	}
}

//GetAllSriovEnabledDevices returns a list of devices on a physcial computer
//that have SRIOV enabled and the current VF count > 0. Note: this scans
//all the directories in NetSysDir variable, this should be used sparingly!!
func GetAllSriovEnabledDevices() (devices []string) {
	if !DirExists(NetSysDir) {
		return
	}

	systemDevices, err := LsFilesWithPrefix(NetSysDir, "", false)
	if err != nil {
		return
	}

	for _, deviceDir := range systemDevices {
		if currentVfCount, _ := GetCurrentVfCount(deviceDir); currentVfCount > 0 {
			devices = append(devices, deviceDir)
		}
	}
	return
}

//GetPfMaxSendingRate gets the maximum sending rate of a given PF device name
//the rate returned is in bits/second
func GetPfMaxSendingRate(pfNetdevName string) (rate uint, err error) {
	deviceDir := GetNetDevDeviceDir(pfNetdevName)
	infinibandDir := filepath.Join(deviceDir, "infiniband")
	infinibandDevices, err := LsFilesWithPrefix(infinibandDir, "", false)
	if err != nil || len(infinibandDevices) == 0 {
		err = fmt.Errorf("Failed to get any devices in infiniband dir[%s]", infinibandDir)
		return
	}

	portsDir := filepath.Join(infinibandDir, infinibandDevices[0], "ports")
	portsAvailable, err := LsFilesWithPrefix(portsDir, "", false)
	if err != nil || len(portsDir) == 0 {
		err = fmt.Errorf("Failed to get any ports from ports dir[%s]", portsDir)
		return
	}

	rateFile := FileObject{
		Path: filepath.Join(portsDir, portsAvailable[0], "rate"),
	}
	rateStr, err := rateFile.Read()
	if err != nil {
		err = fmt.Errorf("Failed to read any information from rate file[%s]: %s", rateFile.Path, err)
		return
	}
	// assuming string in rate file is in the following format: 100 Gb/sec (4X EDR)
	rateStrPieces := strings.Split(rateStr, " ")
	if len(rateStrPieces) < 2 {
		err = fmt.Errorf("Error rate file[%s] is not long enough, expected '100 Gb/sec'", rateFile.Path)
		return
	}

	rateStrNum := strings.TrimSpace(rateStrPieces[0])
	rateNum, err := strconv.ParseUint(rateStrNum, 10, 64)
	if err != nil {
		err = fmt.Errorf("Could not convert string rate[%s] to uint in rate file[%s]: %s", rateFile.Path, rateStrNum, err)
		return
	}

	rateSpeedPerSec := strings.TrimSpace(rateStrPieces[1])
	switch rateSpeedPerSec {
	case "Kb/sec":
		rate = uint(rateNum * 1000)
	case "Mb/sec":
		rate = uint(rateNum * 1000 * 1000)
	case "Gb/sec":
		rate = uint(rateNum * 1000 * 1000 * 1000)
	case "Tb/sec":
		rate = uint(rateNum * 1000 * 1000 * 1000 * 1000)
	default:
		err = fmt.Errorf("Unknown rate type[%s] for rate file[%s]", rateSpeedPerSec, rateFile.Path)
		return
	}
	return
}
