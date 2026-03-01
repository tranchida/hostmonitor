<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { startPolling, stopPolling, getHostInfo } from '$lib/hostinfo.svelte';

	const state = getHostInfo();

	onMount(() => startPolling());
	onDestroy(() => stopPolling());
</script>

<div class="font-sans bg-gray-900 min-h-screen p-4">
	{#if state.loading && !state.data}
		<div class="bg-sky-600 p-8 rounded-md text-white text-center text-xl">
			Loading host information...
		</div>
	{:else if state.error && !state.data}
		<div class="bg-red-600 p-8 rounded-md text-white text-center text-xl">
			Error: {state.error}
		</div>
	{:else if state.data}
		{@const info = state.data}
		<div class="bg-sky-600 p-4 rounded-md">
			<div class="flex flex-row">
				<div class="flex items-center justify-center w-15 -rotate-90 whitespace-nowrap text-white uppercase text-2xl font-bold">
					System information
				</div>
				<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 w-full">
					<!-- Host -->
					<div class="card">
						<span class="card-label">Host</span>
						<span class="card-value">{info.Hostname}</span>
					</div>

					<!-- OS -->
					<div class="card">
						<span class="card-label">OS</span>
						<span class="card-value">{info.OS}</span>
					</div>

					<!-- Platform -->
					<div class="card">
						<span class="card-label">Platform</span>
						<span class="card-value">{info.Platform} {info.PlatformVersion}</span>
					</div>

					<!-- Kernel -->
					<div class="card">
						<span class="card-label">Kernel</span>
						<span class="card-value">{info.KernelVersion}</span>
					</div>

					<!-- Boot time -->
					<div class="card">
						<span class="card-label">Boot time</span>
						<span class="card-value">{info.BootTime}</span>
					</div>

					<!-- Uptime -->
					<div class="card">
						<span class="card-label">Uptime</span>
						<span class="card-value">{info.Uptime}</span>
					</div>

					<!-- CPU -->
					<div class="card">
						<span class="card-label">CPU (Physical / Virtual)</span>
						<span class="card-value">{info.CPUP} / {info.CPUV}</span>
					</div>

					<!-- CPU Usage -->
					<div class="card">
						<span class="card-label">CPU usage</span>
						<span class="card-value">{info.CPUUsage}</span>
					</div>

					<!-- CPU Temperature -->
					<div class="card">
						<span class="card-label">CPU temperature</span>
						<span class="card-value">{info.CPUTemperature}</span>
					</div>

					<!-- Running Processes -->
					<div class="card">
						<span class="card-label">Running processes</span>
						<span class="card-value">{info.RunningProcesses}</span>
					</div>

					<!-- Load Average -->
					<div class="card">
						<span class="card-label">Load average</span>
						<span class="card-value">
							1 min: {info.LoadAverage1} &nbsp; 5 min: {info.LoadAverage5} &nbsp; 15 min: {info.LoadAverage15}
						</span>
					</div>

					<!-- Memory -->
					<div class="card">
						<span class="card-label">Memory (Total / Free / Cache)</span>
						<span class="card-value">{info.TotalMemory} / {info.FreeMemory} / {info.CacheMemory}</span>
					</div>

					<!-- Storage -->
					<div class="card">
						<span class="card-label">Storage (Total / Free)</span>
						<span class="card-value">{info.TotalDiskSpace} / {info.FreeDiskSpace}</span>
					</div>

					<!-- Swap -->
					<div class="card">
						<span class="card-label">Swap (Total / Free)</span>
						<span class="card-value">{info.TotalSwap} / {info.FreeSwap}</span>
					</div>

					<!-- Container -->
					<div class="card">
						<span class="card-label">Running in container</span>
						<span class="card-value">
							{#if info.IsRunningInContainer}
								<span class="inline-flex items-center gap-1 text-green-600 font-semibold">
									<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/></svg>
									Yes
								</span>
							{:else}
								<span class="inline-flex items-center gap-1 text-gray-500 font-semibold">
									<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/></svg>
									No
								</span>
							{/if}
						</span>
					</div>

					<!-- Network Interfaces -->
					<div class="card lg:col-span-4 sm:col-span-2">
						<span class="card-label">Network interfaces</span>
						<span class="card-value">
							{#if info.NetworkInterfaces && info.NetworkInterfaces.length > 0}
								<div class="flex flex-wrap gap-2">
									{#each info.NetworkInterfaces as iface}
										<span class="inline-block bg-sky-100 text-sky-800 text-sm font-mono px-2 py-0.5 rounded">{iface}</span>
									{/each}
								</div>
							{:else}
								<span class="text-gray-400">None</span>
							{/if}
						</span>
					</div>
				</div>
			</div>

			{#if state.error}
				<div class="mt-2 text-yellow-200 text-sm text-center">
					⚠ Polling error: {state.error} — showing last known data
				</div>
			{/if}
		</div>

		<!-- Footer -->
		<div class="text-center text-gray-500 text-sm mt-4">
			Last updated: {info.CurrentTime} · Polling every 10s
		</div>
	{/if}
</div>

<style>
	@reference '../app.css';

	:global(.card) {
		@apply p-4 bg-white border border-solid rounded-md shadow-md shadow-sky-950 hover:ring-2 hover:ring-inset hover:ring-sky-800 transition;
	}
	:global(.card-label) {
		@apply font-bold block mb-1.5;
	}
	:global(.card-value) {
		@apply break-all;
	}
</style>
