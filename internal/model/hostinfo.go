package model

import (
	"fmt"
	"os"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/process"
	"github.com/shirou/gopsutil/v4/sensors"
)

// HostInfo struct
type HostInfo struct {
	CurrentTime       string
	Hostname          string
	Uptime            string
	OS                string
	Platform          string
	PlatformVersion   string
	CPUP              int
	CPUV              int
	TotalMemory       string
	CacheMemory       string
	FreeMemory        string
	TotalDiskSpace    string
	FreeDiskSpace     string
	CPUTemperature    string
	CPUUsage          string
	LoadAverage1      string
	LoadAverage5      string
	LoadAverage15     string
	TotalSwap         string
	FreeSwap          string
	NetworkInterfaces []string
	RunningProcesses  int
	KernelVersion     string
	BootTime          string
	IsRunningInContainer bool
}

// GetHostInfo retrieves and formats host information.
func GetHostInfo() (HostInfo, error) {
	hostname, _ := os.Hostname()

	// extract information from the current host using gopsutil
	uptime, _ := host.Uptime()
	info, _ := host.Info()
	cpuCountP, _ := cpu.Counts(false)
	cpuCountV, _ := cpu.Counts(true)
	memory, _ := mem.VirtualMemory()
	disk, _ := disk.Usage("/")

	// Get CPU temperature
	temps, err := sensors.SensorsTemperatures()
	cpuTemp := "N/A"
	if err == nil && len(temps) > 0 {
		// Find the first CPU temperature sensor
		for _, temp := range temps {
			if temp.SensorKey == "coretemp_core_0" {
				cpuTemp = fmt.Sprintf("%.0f°C", temp.Temperature)
				break
			}
		}
	}

	// Get CPU usage
	cpuPercent, _ := cpu.Percent(0, false)
	cpuUsage := fmt.Sprintf("%.2f%%", cpuPercent[0])

	// Get load average
	loadAvg, err := load.Avg()
	if err != nil {
		return HostInfo{}, err
	}

	// Get swap memory
	swap, err := mem.SwapMemory()
	if err != nil {
		return HostInfo{}, err
	}

	// Get network interfaces
	netInterfaces, _ := net.Interfaces()
	var interfaces []string
	for _, iface := range netInterfaces {
		interfaces = append(interfaces, iface.Name)
	}

	// Get number of running processes
	processes, _ := process.Pids()

	// Get kernel version
	kernelVersion := info.KernelVersion

	// Get boot time
	bootTime := time.Unix(int64(info.BootTime), 0).Format(time.RFC3339)

	hostInfo := HostInfo{
		CurrentTime:       time.Now().Format(time.RFC3339),
		Hostname:          hostname,
		Uptime:            formatDuration(time.Duration(uptime) * time.Second),
		OS:                info.OS,
		Platform:          info.Platform,
		PlatformVersion:   info.PlatformVersion,
		CPUP:              cpuCountP,
		CPUV:              cpuCountV,
		TotalMemory:       humanize.IBytes(memory.Total),
		FreeMemory:        humanize.IBytes(memory.Free),
		CacheMemory:       humanize.IBytes(memory.Cached),
		TotalDiskSpace:    humanize.IBytes(disk.Total),
		FreeDiskSpace:     humanize.IBytes(disk.Free),
		CPUTemperature:    cpuTemp,
		CPUUsage:          cpuUsage,
		LoadAverage1:      fmt.Sprintf("%.2f", loadAvg.Load1),
		LoadAverage5:      fmt.Sprintf("%.2f", loadAvg.Load5),
		LoadAverage15:     fmt.Sprintf("%.2f", loadAvg.Load15),
		TotalSwap:         humanize.IBytes(swap.Total),
		FreeSwap:          humanize.IBytes(swap.Free),
		NetworkInterfaces: interfaces,
		RunningProcesses:  len(processes),
		KernelVersion:     kernelVersion,
		BootTime:          bootTime,
	}

	hostInfo.IsRunningInContainer = isRunningInContainer()

	return hostInfo, nil
}

func formatDuration(d time.Duration) string {
	days := d / (24 * time.Hour)
	d -= days * 24 * time.Hour
	hours := d / time.Hour
	d -= hours * time.Hour
	minutes := d / time.Minute

	return fmt.Sprintf("%d days, %d hours, %d minutes", days, hours, minutes)
}

func isRunningInContainer() bool {
    // docker creates a .dockerenv file at the root
    // of the directory tree inside the container.
    // if this file exists then the viewer is running
    // from inside a container so return true
        
    if _, err := os.Stat("/.dockerenv"); err == nil {
        return true
    }

	// podman create a container with the environment variable container=podman
	if os.Getenv("container") == "podman" {
		return true
	}

	// check if running in kubernetes
	if os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
		return true
	}

	return false
}