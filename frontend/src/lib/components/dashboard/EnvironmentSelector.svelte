<script lang="ts">
	import { CreateEnvironment } from 'wailsjs/go/main/App';
	import { environmentStore } from '$lib/stores/environment.svelte';
	import { slide } from 'svelte/transition';

	let { collectionPath, environments, onrefresh } = $props<{
		collectionPath: string;
		environments: string[];
		onrefresh: () => void;
	}>();

	let creating = $state(false);
	let newEnvName = $state('');
	let error = $state('');
	let inputRef: HTMLInputElement | undefined = $state();

	async function handleCreate() {
		if (!newEnvName.trim()) return;

		try {
			await CreateEnvironment(collectionPath, newEnvName.trim());
			error = '';
			creating = false;
			newEnvName = '';
			onrefresh();
			// Auto select new environment?
			environmentStore.setActive(newEnvName.trim() + '.toml'); // Assuming file name is env name from GetCollection
			// Actually GetCollection returns file names e.g. "dev.toml".
			// CreateEnvironment takes "dev" (logic check needed).
			// backend: func CreateEnvironment(collectionPath, envName string) -> creates envName.toml
		} catch (e: any) {
			error = e.toString();
		}
	}

	function onKeyDown(e: KeyboardEvent) {
		if (e.key === 'Enter') handleCreate();
		if (e.key === 'Escape') {
			creating = false;
			error = '';
		}
	}

	$effect(() => {
		if (creating && inputRef) {
			inputRef.focus();
		}
	});

	function handleSelect(event: Event) {
		const select = event.target as HTMLSelectElement;
		environmentStore.setActive(select.value);
	}
</script>

<div class="relative flex items-center gap-2">
	{#if error}
		<div
			class="absolute top-full right-0 z-10 mt-2 w-48 card bg-error-500 p-2 text-xs text-white"
			transition:slide
		>
			{error}
			<button class="ml-2 font-bold" onclick={() => (error = '')}>X</button>
		</div>
	{/if}

	{#if !creating}
		<select class="select w-48" value={environmentStore.activeEnvironment} onchange={handleSelect}>
			<option value={null}>Select Environment</option>
			{#each environments as env}
				<!-- Remove .toml extension for display? Or keep it? User sees files.
                     Backend Get returns file names "dev.toml".
                     Let's display clean name but value is filename.
                -->
				<option value={env}>{env.replace('.toml', '')}</option>
			{/each}
		</select>
		<button
			class="variant-filled-secondary btn-icon"
			onclick={() => (creating = true)}
			title="Create Environment"
		>
			+
		</button>
	{:else}
		<div class="flex items-center gap-1">
			<input
				bind:this={inputRef}
				class="input w-48"
				type="text"
				placeholder="New Environment..."
				bind:value={newEnvName}
				onkeydown={onKeyDown}
				onblur={() => {
					// Check if we want to cancel on blur.
					// A slight delay to allow button click if needed, or just let them cancel via Escape.
					// Sidebar behavior was cancel on blur (with delay).
					setTimeout(() => (creating = false), 200);
				}}
			/>
		</div>
	{/if}
</div>
