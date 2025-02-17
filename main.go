package main

import (
	"embed"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/process"
	"github.com/shirou/gopsutil/v4/sensors"
)

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
	LoadAverage       string
	TotalSwap         string
	FreeSwap          string
	NetworkInterfaces []string
	RunningProcesses  int
	KernelVersion     string
	BootTime          string
}

//go:embed static templates
var contentFS embed.FS

func newEcho() (*echo.Echo, error) {

	e := echo.New()
	e.HideBanner = false

	e.Renderer = &echo.TemplateRenderer{
		Template: template.Must(template.ParseFS(contentFS, "templates/*")),
	}

	e.Use(middleware.Gzip())

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           "${remote_ip} - - [${time_custom}] \"${method} ${uri} ${protocol}\" ${status} ${bytes_in} ${bytes_out} ${latency_human}\n",
		CustomTimeFormat: "02/Jan/2006:15:04:05 -0700",
	}))

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:       "static",
		Filesystem: http.FS(contentFS),
	}))

	e.GET("/host", func(c echo.Context) error {
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
		loadAvg, _ := load.Avg()
		loadAverage := fmt.Sprintf("%.2f, %.2f, %.2f", loadAvg.Load1, loadAvg.Load5, loadAvg.Load15)

		// Get swap memory
		swap, _ := mem.SwapMemory()

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
			LoadAverage:       loadAverage,
			TotalSwap:         humanize.IBytes(swap.Total),
			FreeSwap:          humanize.IBytes(swap.Free),
			NetworkInterfaces: interfaces,
			RunningProcesses:  len(processes),
			KernelVersion:     kernelVersion,
			BootTime:          bootTime,
		}

		return c.Render(http.StatusOK, "main", hostInfo)

	})

	return e, nil
}

func main() {

	e, err := newEcho()
	if err != nil {
		panic(err)
	}

	fmt.Println("open browser on : http://localhost:8080")

	if err := e.Start(":8080"); err != nil {
		panic(err)
	}

}

func formatDuration(d time.Duration) string {
	days := d / (24 * time.Hour)
	d -= days * 24 * time.Hour
	hours := d / time.Hour
	d -= hours * time.Hour
	minutes := d / time.Minute

	return fmt.Sprintf("%d jours, %d heures, %d minutes", days, hours, minutes)
}


