export interface HostInfo {
	CurrentTime: string;
	Hostname: string;
	Uptime: string;
	OS: string;
	Platform: string;
	PlatformVersion: string;
	CPUP: number;
	CPUV: number;
	TotalMemory: string;
	CacheMemory: string;
	FreeMemory: string;
	TotalDiskSpace: string;
	FreeDiskSpace: string;
	CPUTemperature: string;
	CPUUsage: string;
	LoadAverage1: string;
	LoadAverage5: string;
	LoadAverage15: string;
	TotalSwap: string;
	FreeSwap: string;
	NetworkInterfaces: string[];
	RunningProcesses: number;
	KernelVersion: string;
	BootTime: string;
	IsRunningInContainer: boolean;
}
