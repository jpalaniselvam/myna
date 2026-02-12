<script lang="ts">
	import { collectionStore } from '$lib/stores/collections.svelte';
	import { GetCollection } from 'wailsjs/go/main/App';
	import type { collection } from 'wailsjs/go/models';
	import EditableDescription from './EditableDescription.svelte';
	import CollectionVariables from '$lib/components/dashboard/CollectionVariables.svelte';
	import EnvironmentManager from '$lib/components/dashboard/EnvironmentManager.svelte';
	import ActionList from '$lib/components/dashboard/ActionList.svelte';
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

	// Derived name from path
	let collectionName = $derived(
		collectionStore.activeCollection ? collectionStore.activeCollection.split(/[\\/]/).pop() : ''
	);

	let activeTab = $state('overview');
</script>

<div class="flex h-full w-full flex-col">
	<Tabs value={activeTab} onValueChange={(e) => (activeTab = e.value)} class="flex h-full flex-col">
		<!-- Header -->
		<header class="bg-surface-100-800-token flex flex-col">
			<div class="flex items-center justify-between p-4 pb-0">
				<div class="mr-4 min-w-0 flex-1">
					<h2 class="truncate h2 font-bold text-primary-500" title={collectionName}>
						{collectionName}
					</h2>
				</div>

				<div class="flex shrink-0 items-center gap-4">
					<!-- Environment Selector moved to Environments tab -->
				</div>
			</div>

			{#if collectionData}
				<Tabs.List class="px-4">
					<Tabs.Trigger value="overview">Overview</Tabs.Trigger>
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
		<div class="flex-1 overflow-auto p-2">
			{#if loading}
				<div class="flex h-full items-center justify-center">
					<p>Loading...</p>
				</div>
			{:else if error}
				<div class="alert variant-filled-error">
					{error}
				</div>
			{:else if collectionData}
				<Tabs.Content value="overview">
					<!-- Placeholder for Collection Content -->
					<div class="card">
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
						<div class="mt-2">
							<p>Collection path: {collectionStore.activeCollection}</p>
						</div>
						<div class="mt-4">
							<h4 class="h4">Stats</h4>
							<p>Environments: {collectionData.environments?.length || 0}</p>
							<p>Actions: {Object.keys(collectionData.actions || {}).length} root items</p>
						</div>
					</div>
				</Tabs.Content>
				<Tabs.Content value="actions">
					<div class="h-full card p-4">
						<ActionList
							actions={collectionData.actions}
							collectionPath={collectionStore.activeCollection || ''}
						/>
					</div>
				</Tabs.Content>
				<Tabs.Content value="settings">
					<!-- Collection Variables (Settings) -->
					<div class="card">
						<CollectionVariables
							collectionPath={collectionStore.activeCollection || ''}
							variables={collectionData.pre}
							onupdate={(newVars) => {
								if (collectionData) collectionData.pre = newVars;
							}}
						/>
					</div>
				</Tabs.Content>
				<Tabs.Content value="environments">
					<!-- Environment Manager -->
					<div class="h-full card">
						<EnvironmentManager
							collectionPath={collectionStore.activeCollection || ''}
							environments={collectionData.environments}
							onrefresh={refresh}
						/>
					</div>
				</Tabs.Content>
			{:else}
				<!-- Should not happen if path is set and no error/loading -->
				<p>No data found.</p>
			{/if}
		</div>
	</Tabs>
</div>
