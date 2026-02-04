<script lang="ts">
	import { collectionStore } from '$lib/stores/collections.svelte';
	import { GetCollection } from 'wailsjs/go/main/App';
	import type { collection } from 'wailsjs/go/models';
	import EditableDescription from './EditableDescription.svelte';
	import CollectionVariables from '$lib/components/dashboard/CollectionVariables.svelte';
	import EnvironmentManager from '$lib/components/dashboard/EnvironmentManager.svelte';

	let collectionData = $state<collection.CollectionResponse | null>(null);
	let loading = $state(false);
	let error = $state('');

	// React to activeCollection changes
	$effect(() => {
		loadCollection(collectionStore.activeCollection);
	});

	async function loadCollection(path: string | null) {
		if (!path) {
			collectionData = null;
			return;
		}

		loading = true;
		error = '';
		try {
			collectionData = await GetCollection(path);
		} catch (e: any) {
			error = e.toString();
			collectionData = null;
		} finally {
			loading = false;
		}
	}

	function refresh() {
		if (collectionStore.activeCollection) {
			loadCollection(collectionStore.activeCollection);
		}
	}

	// Derived name from path
	let collectionName = $derived(
		collectionStore.activeCollection ? collectionStore.activeCollection.split(/[\\/]/).pop() : ''
	);

	let activeTab = $state('overview');
</script>

<div class="flex h-full w-full flex-col">
	<!-- Header -->
	<header class="bg-surface-100-800-token flex flex-col border-b border-surface-500/30">
		<div class="flex items-center justify-between p-4 pb-0">
			<div class="mr-4 min-w-0 flex-1">
				<h2 class="truncate h2 font-bold text-primary-500" title={collectionName}>
					{collectionName}
				</h2>
				{#if collectionData && collectionStore.activeCollection}
					<div class="mt-1">
						<EditableDescription
							collectionPath={collectionStore.activeCollection}
							description={collectionData.description || ''}
							onupdate={(newDesc) => {
								if (collectionData) collectionData.description = newDesc;
							}}
						/>
					</div>
				{/if}
			</div>

			<div class="flex shrink-0 items-center gap-4">
				<!-- Environment Selector moved to Environments tab -->
			</div>
		</div>

		{#if collectionData}
			<div class="flex px-4 pt-4">
				<button
					class="hover:bg-surface-200-700-token border-b-2 px-4 py-2 font-medium transition-colors {activeTab ===
					'overview'
						? 'border-primary-500 text-primary-500'
						: 'text-surface-600-300-token border-transparent'}"
					onclick={() => (activeTab = 'overview')}
				>
					Overview
				</button>
				<button
					class="hover:bg-surface-200-700-token border-b-2 px-4 py-2 font-medium transition-colors {activeTab ===
					'settings'
						? 'border-primary-500 text-primary-500'
						: 'text-surface-600-300-token border-transparent'}"
					onclick={() => (activeTab = 'settings')}
				>
					Variables & Settings
				</button>
				<button
					class="hover:bg-surface-200-700-token border-b-2 px-4 py-2 font-medium transition-colors {activeTab ===
					'environments'
						? 'border-primary-500 text-primary-500'
						: 'text-surface-600-300-token border-transparent'}"
					onclick={() => (activeTab = 'environments')}
				>
					Environments
				</button>
			</div>
		{/if}
	</header>

	<!-- Content -->
	<div class="flex-1 overflow-auto p-6">
		{#if loading}
			<div class="flex h-full items-center justify-center">
				<p>Loading...</p>
			</div>
		{:else if error}
			<div class="alert variant-filled-error">
				{error}
			</div>
		{:else if collectionData}
			{#if activeTab === 'overview'}
				<!-- Placeholder for Collection Content -->
				<div class="card p-4">
					<h3 class="mb-2 h3">Details</h3>
					<p>Collection path: {collectionStore.activeCollection}</p>
					<div class="mt-4">
						<h4 class="h4">Stats</h4>
						<p>Environments: {collectionData.environments?.length || 0}</p>
						<p>Actions loaded.</p>
					</div>
				</div>
			{:else if activeTab === 'settings'}
				<!-- Collection Variables (Settings) -->
				<div class="card p-4">
					<CollectionVariables
						collectionPath={collectionStore.activeCollection || ''}
						variables={collectionData.pre}
						onupdate={(newVars) => {
							if (collectionData) collectionData.pre = newVars;
						}}
					/>
				</div>
			{:else if activeTab === 'environments'}
				<!-- Environment Manager -->
				<div class="h-full card p-4">
					<EnvironmentManager
						collectionPath={collectionStore.activeCollection || ''}
						environments={collectionData.environments}
						onrefresh={refresh}
					/>
				</div>
			{/if}
		{:else}
			<!-- Should not happen if path is set and no error/loading -->
			<p>No data found.</p>
		{/if}
	</div>
</div>
