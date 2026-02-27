<script lang="ts">
	import { workspaceStore } from '$lib/stores/workspace.svelte';
	import { collectionStore } from '$lib/stores/collections.svelte';
	import { Tabs } from '@skeletonlabs/skeleton-svelte';
	import CollectionWorkspace from './dashboard/CollectionWorkspace.svelte';
	import ActionWorkspace from './dashboard/ActionWorkspace.svelte';
	import IndexContent from './IndexContent.svelte';
	import CreateActionDialog from './dialogs/CreateActionDialog.svelte';
	import { Plus, X } from '@lucide/svelte';
	import { CreateAction } from 'wailsjs/go/main/App';

	interface ActionData {
		kind: string;
		name: string;
	}

	let showCreateActionDialog = $state(false);

	async function handleCreateAction(data: ActionData) {
		if (collectionStore.activeCollection) {
			const fileName = `${data.name}.toml`;
			const actionData = {
				version: 'v1',
				kind: data.kind,
				description: `Action for ${data.name}`
			};

			try {
				await CreateAction({
					collection_path: collectionStore.activeCollection,
					sub_path: '',
					file_name: fileName,
					data: actionData
				});

				const id = `${collectionStore.activeCollection}/${fileName}`;
				workspaceStore.openTab('action', id, data.name, {
					collectionPath: collectionStore.activeCollection,
					subPath: '',
					fileName: fileName,
					actionKind: data.kind
				});
			} catch (e: any) {
				console.error('Failed to create action:', e);
			}
		} else {
			alert('Please select a collection first');
		}
	}

	function handleValueChange(e: { value: string }) {
		if (e.value) {
			workspaceStore.setActiveTab(e.value);
			const tab = workspaceStore.tabs.find((t) => t.id === e.value);
			if (tab?.kind === 'collection') {
				collectionStore.setActive(tab.metadata.path);
			} else if (tab?.kind === 'action') {
				collectionStore.setActive(tab.metadata.collectionPath);
			}
		}
	}

	function closeTab(e: MouseEvent, id: string) {
		e.stopPropagation();
		workspaceStore.closeTab(id);
		// If we closed the last tab, clear active collection
		if (workspaceStore.tabs.length === 0) {
			collectionStore.setActive('');
		} else if (workspaceStore.activeTab?.kind === 'collection') {
			collectionStore.setActive(workspaceStore.activeTab.metadata.path);
		} else if (workspaceStore.activeTab?.kind === 'action') {
			collectionStore.setActive(workspaceStore.activeTab.metadata.collectionPath);
		}
	}
</script>

<div class="flex h-full flex-col pt-2">
	{#if workspaceStore.tabs.length > 0}
		<Tabs value={workspaceStore.activeTabId || ''} onValueChange={handleValueChange}>
			<Tabs.List class="px-4">
				{#each workspaceStore.tabs as tab (tab.id)}
					<Tabs.Trigger value={tab.id} class="flex w-64 items-center">
						<span class="flex-1 truncate text-left">{tab.title}</span>
						<button
							onclick={(e) => closeTab(e, tab.id)}
							aria-label="Close tab"
							class="btn-icon btn-icon-sm transition-colors hover:bg-surface-500/20"
						>
							<X size={14} />
						</button>
					</Tabs.Trigger>
				{/each}
				<button
					onclick={() => (showCreateActionDialog = true)}
					class="ml-2 btn-icon btn-icon-sm self-center transition-colors hover:bg-surface-500/20"
					title="Create New Action"
				>
					<Plus size={18} />
				</button>
				<Tabs.Indicator />
			</Tabs.List>

			<div class="flex-1 overflow-auto">
				{#if workspaceStore.activeTab}
					<Tabs.Content value={workspaceStore.activeTab.id} class="h-full">
						{#if workspaceStore.activeTab.kind === 'collection'}
							<CollectionWorkspace tab={workspaceStore.activeTab} />
						{:else if workspaceStore.activeTab.kind === 'action'}
							<ActionWorkspace tab={workspaceStore.activeTab} />
						{/if}
					</Tabs.Content>
				{/if}
			</div>
		</Tabs>
	{:else}
		<IndexContent />
	{/if}

	<CreateActionDialog bind:open={showCreateActionDialog} onsubmit={handleCreateAction} />
</div>

<style>
	:global(.no-scrollbar::-webkit-scrollbar) {
		display: none;
	}
	:global(.no-scrollbar) {
		-ms-overflow-style: none;
		scrollbar-width: none;
	}
</style>
