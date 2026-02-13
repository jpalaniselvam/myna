<script lang="ts">
	import { collectionStore } from '$lib/stores/collections.svelte';
	import { GetCollection } from 'wailsjs/go/main/App';
	import type { collection } from 'wailsjs/go/models';
	import CollectionVariables from '$lib/components/dashboard/CollectionVariables.svelte';
	import CollectionCredentials from '$lib/components/dashboard/CollectionCredentials.svelte';
	import EnvironmentManager from '$lib/components/dashboard/EnvironmentManager.svelte';
	import ActionList from '$lib/components/dashboard/ActionList.svelte';
	import CollectionOverview from '$lib/components/dashboard/CollectionOverview.svelte';
	import { Tabs } from '@skeletonlabs/skeleton-svelte';

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

	let activeTab = $state('overview');
</script>

<div class="flex h-full w-full flex-col">
	<Tabs value={activeTab} onValueChange={(e) => (activeTab = e.value)} class="flex h-full flex-col">
		<!-- Header -->
		<header class="bg-surface-100-800-token">
			{#if collectionData}
				<Tabs.List class="px-4">
					<Tabs.Trigger value="overview">Overview</Tabs.Trigger>
					<Tabs.Trigger value="credentials" class="flex items-center gap-2">
						Credentials
					</Tabs.Trigger>
					<Tabs.Trigger value="actions" class="flex items-center gap-2">
						Actions
						<span class="badge-icon preset-filled-primary-500">
							{Object.keys(collectionData.actions || {}).length}
						</span>
					</Tabs.Trigger>
					<Tabs.Trigger value="settings" class="flex items-center gap-2">
						Variables
						<span class="badge-icon preset-filled-primary-500">
							{Object.keys(collectionData.pre || {}).length}
						</span>
					</Tabs.Trigger>
					<Tabs.Trigger value="environments" class="flex items-center gap-2">
						Environments
						<span class="badge-icon preset-filled-primary-500">
							{collectionData.environments?.length || 0}
						</span>
					</Tabs.Trigger>
					<Tabs.Indicator />
				</Tabs.List>
			{/if}
		</header>

		<!-- Content -->
		<div class="flex-1 overflow-auto">
			{#if loading}
				<div class="flex h-full items-center justify-center">
					<p>Loading...</p>
				</div>
			{:else if error}
				<div class="alert variant-filled-error m-4">
					{error}
				</div>
			{:else if collectionData}
				<Tabs.Content value="overview">
					<CollectionOverview
						{collectionData}
						collectionPath={collectionStore.activeCollection || ''}
						onupdate={(newDesc) => {
							if (collectionData) collectionData.description = newDesc;
						}}
					/>
				</Tabs.Content>
				<Tabs.Content value="credentials">
					<div class="p-4">
						<div class="card p-4">
							<CollectionCredentials
								collectionPath={collectionStore.activeCollection || ''}
								credentials={collectionData.credentials}
								variables={collectionData.pre}
								onupdate={(newCreds) => {
									if (collectionData) collectionData.credentials = newCreds;
								}}
							/>
						</div>
					</div>
				</Tabs.Content>
				<Tabs.Content value="actions">
					<div class="h-full p-4">
						<div class="h-full card p-4">
							<ActionList
								actions={collectionData.actions}
								collectionPath={collectionStore.activeCollection || ''}
							/>
						</div>
					</div>
				</Tabs.Content>
				<Tabs.Content value="settings">
					<div class="p-4">
						<div class="card p-4">
							<CollectionVariables
								collectionPath={collectionStore.activeCollection || ''}
								variables={collectionData.pre}
								credentials={collectionData.credentials}
								onupdate={(newVars) => {
									if (collectionData) collectionData.pre = newVars;
								}}
							/>
						</div>
					</div>
				</Tabs.Content>
				<Tabs.Content value="environments">
					<div class="h-full p-4">
						<div class="h-full card">
							<EnvironmentManager
								collectionPath={collectionStore.activeCollection || ''}
								environments={collectionData.environments}
								onrefresh={refresh}
							/>
						</div>
					</div>
				</Tabs.Content>
			{:else}
				<p class="p-4 opacity-60">No data found.</p>
			{/if}
		</div>
	</Tabs>
</div>
