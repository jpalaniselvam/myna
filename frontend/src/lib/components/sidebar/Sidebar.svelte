<script lang="ts">
	import { collectionStore } from '$lib/stores/collections.svelte';
	import SidebarHeader from './SidebarHeader.svelte';
	import SidebarList from './SidebarList.svelte';
	import CreateCollectionDialog from '$lib/components/dialogs/CreateCollectionDialog.svelte';
	import { CreateCollection } from '../../../../wailsjs/go/main/App';
	import ErrorBanner from './ErrorBanner.svelte';

	let showCreateDialog = $state(false);
	let error = $state('');

	function handleSelectCollection(path: string) {
		collectionStore.setActive(path);
	}

	function handleAddCollection() {
		showCreateDialog = true;
	}

	async function onCreateSubmit(data: { name: string; baseDir: string }) {
		try {
			error = '';
			await CreateCollection(data.baseDir, data.name, '');

			// Success
			const separator = data.baseDir.includes('\\') ? '\\' : '/';
			const cleanBase = data.baseDir.endsWith(separator) ? data.baseDir.slice(0, -1) : data.baseDir;
			const fullPath = `${cleanBase}${separator}${data.name}`;

			collectionStore.add(fullPath);
			collectionStore.setActive(fullPath);
			showCreateDialog = false;
		} catch (e: any) {
			error = e.toString() || 'Failed to create collection';
			console.error(error);
		}
	}
</script>

<aside class="bg-surface-50-900-token flex h-full w-64 flex-col border-r border-surface-500/30">
	{#if error}
		<div class="fixed top-4 right-4 z-50 w-full max-w-xs transition-all">
			<ErrorBanner message={error} ondismiss={() => (error = '')} />
		</div>
	{/if}

	<SidebarHeader onadd={handleAddCollection} />

	<hr class="border-t-2!" />

	<SidebarList
		collections={collectionStore.collections}
		activeCollection={collectionStore.activeCollection}
		onselect={handleSelectCollection}
	/>

	<CreateCollectionDialog bind:open={showCreateDialog} onsubmit={onCreateSubmit} />
</aside>
