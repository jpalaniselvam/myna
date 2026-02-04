<script lang="ts">
	import {
		CreateEnvironment,
		GetEnvironment,
		AddEnvVar,
		UpdateEnvVar,
		DeleteEnvVar
	} from 'wailsjs/go/main/App';
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

	// Variables state
	let loadingVars = $state(false);
	let currentVars = $state<{ key: string; value: string }[]>([]);
	let selectedEnvName = $state<string | null>(null);
	let varError = $state('');

	// Effect to load variables when active environment changes
	$effect(() => {
		// If we have an active environment (filename like 'dev.toml'), load its vars
		if (environmentStore.activeEnvironment) {
			// Extract name without extension for backend calls if needed?
			// Backend uses envName for file lookup: envName + ".toml"
			// IF provided envName is already "dev.toml", backend does "dev.toml.toml".
			// We need to be careful.
			// Looking at backend:
			// func CreateEnvironment(collectionPath, envName string) -> filepath.Join(envDir, envName+".toml")
			// func GetCollection -> returns file names e.g. "dev.toml" in environments list.

			// So here we should strip .toml for display and backend calls that expect just name?
			// Or backend expects name?
			// AddEnvVar: envFilePath := filepath.Join(collectionPath, "environments", envName+".toml")
			// So if we pass "dev.toml", it looks for "dev.toml.toml".
			// So we MUST strip extension.

			const raw = environmentStore.activeEnvironment;
			const name = raw.endsWith('.toml') ? raw.slice(0, -5) : raw;
			selectedEnvName = name;
			loadVariables(name);
		} else {
			selectedEnvName = null;
			currentVars = [];
		}
	});

	async function loadVariables(envName: string) {
		loadingVars = true;
		varError = '';
		try {
			const vars = await GetEnvironment(collectionPath, envName);
			currentVars = Object.entries(vars || {}).map(([k, v]) => ({
				key: k,
				value: String(v)
			}));
		} catch (e: any) {
			varError = e.toString();
			currentVars = [];
		} finally {
			loadingVars = false;
		}
	}

	async function handleCreate() {
		if (!newEnvName.trim()) return;

		try {
			await CreateEnvironment(collectionPath, newEnvName.trim());
			error = '';
			creating = false;
			newEnvName = '';
			onrefresh();
			// Select new
			environmentStore.setActive(newEnvName.trim() + '.toml');
		} catch (e: any) {
			error = e.toString();
		}
	}

	// Variable operations
	async function addVariable() {
		if (!selectedEnvName) return;
		// Optimization: just UI add, save later? Or save immediately?
		// "Update, delete variables" usually implies direct action or bulk save.
		// Let's do direct action for simplicity per requirement "add, update, delete variables"
		// But table editing is usually deferred.
		// Let's try deferred save for UX, but requirement says "Update, delete variables".
		// Let's add an empty row. Saving row triggers AddEnvVar.
		currentVars.push({ key: '', value: '' });
	}

	async function saveVariable(index: number) {
		if (!selectedEnvName) return;
		const v = currentVars[index];
		if (!v.key) return; // Err?

		// Check if updating or adding?
		// Backend distinguishes Add vs Update.
		// Complexity: we don't know if this key existed before easily unless we track original state.
		// Fallback: Try Add, if fails (exists) -> Update.

		try {
			await AddEnvVar(collectionPath, selectedEnvName, v.key, v.value);
		} catch (e: any) {
			// Assume exists, try Update
			if (e.toString().includes('already exists')) {
				try {
					await UpdateEnvVar(collectionPath, selectedEnvName, v.key, v.value);
				} catch (e2: any) {
					varError = e2.toString();
				}
			} else {
				varError = e.toString();
			}
		}
	}

	async function deleteVariable(index: number) {
		if (!selectedEnvName) return;
		const v = currentVars[index];
		if (v.key) {
			try {
				await DeleteEnvVar(collectionPath, selectedEnvName, v.key);
			} catch (e: any) {
				varError = e.toString();
				return;
			}
		}
		currentVars.splice(index, 1);
	}

	function onKeyDown(e: KeyboardEvent) {
		if (e.key === 'Enter') handleCreate();
		if (e.key === 'Escape') {
			creating = false;
			error = '';
		}
	}

	function handleSelect(env: string) {
		environmentStore.setActive(env);
	}
</script>

<div class="flex h-full flex-col">
	<div class="mb-4 flex items-center justify-between">
		<h3 class="h3">Environments</h3>
		<!-- Create New -->
		{#if !creating}
			<button class="variant-filled-secondary btn btn-sm" onclick={() => (creating = true)}>
				+ New Environment
			</button>
		{:else}
			<div class="flex items-center gap-2">
				<input
					bind:this={inputRef}
					class="input-sm input w-48"
					type="text"
					placeholder="Name..."
					bind:value={newEnvName}
					onkeydown={onKeyDown}
				/>
				<button class="variant-filled-primary btn btn-sm" onclick={handleCreate}>OK</button>
				<button class="variant-ghost btn btn-sm" onclick={() => (creating = false)}>Cancel</button>
			</div>
		{/if}
	</div>

	{#if error}
		<div class="alert variant-filled-error mb-4" transition:slide>
			{error} <button class="ml-2 font-bold" onclick={() => (error = '')}>X</button>
		</div>
	{/if}

	<div class="flex h-full gap-4">
		<!-- List of Environments -->
		<div class="w-1/4 overflow-y-auto border-r border-surface-500/30 pr-4">
			<div class="list-nav">
				<ul>
					{#each environments as env}
						<li
							class={environmentStore.activeEnvironment === env ? 'bg-surface-200-700-token' : ''}
						>
							<button class="w-full rounded px-4 py-2 text-left" onclick={() => handleSelect(env)}>
								{env.replace('.toml', '')}
							</button>
						</li>
					{/each}
					{#if environments.length === 0}
						<li class="p-2 text-sm text-surface-500 italic">No environments created.</li>
					{/if}
				</ul>
			</div>
		</div>

		<!-- Selected Environment Variables -->
		<div class="flex flex-1 flex-col pl-4">
			{#if selectedEnvName}
				<div class="mb-2 flex items-center justify-between">
					<h4 class="h4 text-primary-500">{selectedEnvName} Variables</h4>
					<button class="variant-ghost-secondary btn btn-sm" onclick={addVariable}
						>+ Add Variable</button
					>
				</div>

				{#if varError}
					<div class="alert variant-filled-error mb-2" transition:slide>
						{varError} <button class="ml-2 font-bold" onclick={() => (varError = '')}>X</button>
					</div>
				{/if}

				{#if loadingVars}
					<div class="flex justify-center p-4">Loading...</div>
				{:else}
					<div class="table-container overflow-y-auto">
						<table class="table-hover table">
							<thead>
								<tr>
									<th>Key</th>
									<th>Value</th>
									<th class="w-16"></th>
								</tr>
							</thead>
							<tbody>
								{#each currentVars as v, i}
									<tr>
										<td>
											<input
												class="input-sm input"
												type="text"
												bind:value={v.key}
												placeholder="Key"
												onblur={() => saveVariable(i)}
											/>
										</td>
										<td>
											<input
												class="input-sm input"
												type="text"
												bind:value={v.value}
												placeholder="Value"
												onblur={() => saveVariable(i)}
											/>
										</td>
										<td>
											<button
												class="variant-filled-error btn-icon btn-icon-sm"
												onclick={() => deleteVariable(i)}>-</button
											>
										</td>
									</tr>
								{/each}
							</tbody>
						</table>
					</div>
				{/if}
			{:else}
				<div class="flex h-full items-center justify-center text-surface-500">
					Select an environment to view variables.
				</div>
			{/if}
		</div>
	</div>
</div>
