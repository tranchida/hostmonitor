import type { HostInfo } from './types';

const POLL_INTERVAL = 10_000;

let hostInfo = $state<HostInfo | null>(null);
let error = $state<string | null>(null);
let loading = $state(true);

async function fetchHostInfo() {
	try {
		const res = await fetch('/api/hostinfo');
		if (!res.ok) throw new Error(`HTTP ${res.status}`);
		hostInfo = await res.json();
		error = null;
	} catch (e) {
		error = e instanceof Error ? e.message : 'Unknown error';
	} finally {
		loading = false;
	}
}

let intervalId: ReturnType<typeof setInterval> | null = null;

export function startPolling() {
	fetchHostInfo();
	intervalId = setInterval(fetchHostInfo, POLL_INTERVAL);
}

export function stopPolling() {
	if (intervalId) {
		clearInterval(intervalId);
		intervalId = null;
	}
}

export function getHostInfo() {
	return {
		get data() { return hostInfo; },
		get error() { return error; },
		get loading() { return loading; }
	};
}
