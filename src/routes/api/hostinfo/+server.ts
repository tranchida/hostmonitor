import { json } from '@sveltejs/kit';
import si from 'systeminformation';
import os from 'os';
import fs from 'fs';

/** Format bytes to human-readable string (e.g. "3.8 GiB") */
function formatBytes(bytes: number): string {
	const units = ['B', 'KiB', 'MiB', 'GiB', 'TiB'];
	let value = bytes;
	let unit = 0;
	while (value >= 1024 && unit < units.length - 1) {
		value /= 1024;
		unit++;
	}
	return `${value.toFixed(1)} ${units[unit]}`;
}

/** Format uptime seconds to "X days, Y hours, Z minutes" */
function formatUptime(seconds: number): string {
	const days = Math.floor(seconds / 86400);
	const hours = Math.floor((seconds % 86400) / 3600);
	const minutes = Math.floor((seconds % 3600) / 60);
	const parts: string[] = [];
	if (days > 0) parts.push(`${days} day${days !== 1 ? 's' : ''}`);
	if (hours > 0) parts.push(`${hours} hour${hours !== 1 ? 's' : ''}`);
	if (minutes > 0 || parts.length === 0) parts.push(`${minutes} minute${minutes !== 1 ? 's' : ''}`);
	return parts.join(', ');
}

/** Detect if running inside a container */
function isRunningInContainer(): boolean {
	if (process.env.container || process.env.KUBERNETES_SERVICE_HOST) return true;
	try {
		fs.accessSync('/.dockerenv');
		return true;
	} catch {
		// not a docker container
	}
	try {
		const cgroup = fs.readFileSync('/proc/1/cgroup', 'utf8');
		if (cgroup.includes('docker') || cgroup.includes('lxc') || cgroup.includes('containerd')) {
			return true;
		}
	} catch {
		// /proc not available
	}
	return false;
}

export async function GET() {
	const [osInfo, cpuInfo, memInfo, fsSize, cpuLoad, cpuTemp, procs, netIfaces, time] =
		await Promise.all([
			si.osInfo(),
			si.cpu(),
			si.mem(),
			si.fsSize(),
			si.currentLoad(),
			si.cpuTemperature(),
			si.processes(),
			si.networkInterfaces(),
			si.time()
		]);

	// Disk: use the partition mounted at /  (fallback to first entry)
	const rootDisk =
		(Array.isArray(fsSize) ? fsSize : []).find((d) => d.mount === '/') ??
		(Array.isArray(fsSize) ? fsSize[0] : null);

	// CPU temperature
	let cpuTemperature = 'N/A';
	if (cpuTemp.main && cpuTemp.main > 0) {
		cpuTemperature = `${cpuTemp.main.toFixed(0)}°C`;
	} else if (cpuTemp.cores && cpuTemp.cores.length > 0) {
		const avg = cpuTemp.cores.reduce((a, b) => a + b, 0) / cpuTemp.cores.length;
		if (avg > 0) cpuTemperature = `${avg.toFixed(0)}°C`;
	}

	// Load averages (Node os.loadavg() gives [1, 5, 15])
	const [la1, la5, la15] = os.loadavg();

	// Network interface names
	const ifaceNames: string[] = Array.isArray(netIfaces)
		? netIfaces.map((n) => n.iface).filter(Boolean)
		: typeof netIfaces === 'object'
			? Object.values(netIfaces).map((n: unknown) => (n as { iface: string }).iface)
			: [];

	// Boot time
	const uptimeSec: number =
		typeof time.uptime === 'number' ? time.uptime : parseFloat(String(time.uptime)) || 0;
	const bootDate = new Date(Date.now() - uptimeSec * 1000);

	const hostInfo = {
		CurrentTime: new Date().toISOString(),
		Hostname: osInfo.hostname ?? os.hostname(),
		Uptime: formatUptime(uptimeSec),
		OS: osInfo.os ?? osInfo.distro ?? os.type(),
		Platform: osInfo.platform ?? os.platform(),
		PlatformVersion: osInfo.release ?? '',
		KernelVersion: osInfo.kernel ?? os.release(),
		CPUP: cpuInfo.physicalCores ?? cpuInfo.cores ?? os.cpus().length,
		CPUV: cpuInfo.cores ?? os.cpus().length,
		TotalMemory: formatBytes(memInfo.total),
		FreeMemory: formatBytes(memInfo.free),
		CacheMemory: formatBytes(memInfo.cached ?? 0),
		TotalDiskSpace: rootDisk ? formatBytes(rootDisk.size) : 'N/A',
		FreeDiskSpace: rootDisk ? formatBytes(rootDisk.available) : 'N/A',
		CPUTemperature: cpuTemperature,
		CPUUsage: `${cpuLoad.currentLoad.toFixed(2)}%`,
		LoadAverage1: la1.toFixed(2),
		LoadAverage5: la5.toFixed(2),
		LoadAverage15: la15.toFixed(2),
		TotalSwap: formatBytes(memInfo.swaptotal),
		FreeSwap: formatBytes(memInfo.swapfree),
		NetworkInterfaces: ifaceNames,
		RunningProcesses: procs.all ?? 0,
		BootTime: bootDate.toISOString(),
		IsRunningInContainer: isRunningInContainer()
	};

	return json(hostInfo);
}
