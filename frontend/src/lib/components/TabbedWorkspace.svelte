<script lang="ts">
	import { workspaceStore, type Tab } from '$lib/stores/workspace.svelte';
	import { collectionStore } from '$lib/stores/collections.svelte';
	import { X } from '@lucide/svelte';
	import { Tabs } from '@skeletonlabs/skeleton-svelte';
	import CollectionWorkspace from './dashboard/CollectionWorkspace.svelte';
	import IndexContent from './IndexContent.svelte';

	function handleValueChange(e: { value: string }) {
		if (e.value) {
			workspaceStore.setActiveTab(e.value);
			const tab = workspaceStore.tabs.find((t) => t.id === e.value);
			if (tab?.kind === 'collection') {
				collectionStore.setActive(tab.metadata.path);
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
		}
	}
</script>

<div class="flex h-full flex-col pt-2">
	{#if workspaceStore.tabs.length > 0}
		<Tabs value={workspaceStore.activeTabId || ''} onValueChange={handleValueChange}>
			<Tabs.List>
				{#each workspaceStore.tabs as tab (tab.id)}
					<Tabs.Trigger value={tab.id} class="flex w-64 items-center">
						<span class="flex-1 truncate">{tab.title}</span>
						<button
							onclick={(e) => closeTab(e, tab.id)}
							aria-label="Close tab"
							class="rounded-full p-2 hover:bg-surface-500/30"
						>
							<X size={14} />
						</button>
					</Tabs.Trigger>
				{/each}
				<Tabs.Indicator />
			</Tabs.List>

			<div class="flex-1 overflow-auto">
				{#if workspaceStore.activeTab}
					<Tabs.Content value={workspaceStore.activeTab.id} class="h-full">
						{#if workspaceStore.activeTab.kind === 'collection'}
							<CollectionWorkspace />
						{/if}
					</Tabs.Content>
				{/if}
			</div>
		</Tabs>
	{:else}
		<IndexContent />
	{/if}
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
