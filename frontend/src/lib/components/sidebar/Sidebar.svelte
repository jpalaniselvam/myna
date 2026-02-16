<script lang="ts">
	import { collectionStore } from '$lib/stores/collections.svelte';
	import SidebarHeader from './SidebarHeader.svelte';
	import SidebarList from './SidebarList.svelte';
	import CreateCollectionDialog from '$lib/components/dialogs/CreateCollectionDialog.svelte';
	import {
		CreateCollection,
		SelectDirectory,
		GetCollection
	} from '../../../../wailsjs/go/main/App';
	import ErrorBanner from './ErrorBanner.svelte';

	import { workspaceStore } from '$lib/stores/workspace.svelte';
	import { layoutStore } from '$lib/stores/layout.svelte';

	let showCreateDialog = $state(false);
	let error = $state('');
	let isResizing = $state(false);

	function handleSelectCollection(path: string) {
		const name = path.split(/[\\/]/).pop() || 'Collection';
		workspaceStore.openTab('collection', path, name, { path });
		collectionStore.setActive(path);
	}

	function handleAddCollection() {
		showCreateDialog = true;
	}

	async function handleOpenCollection() {
		try {
			error = '';
			const path = await SelectDirectory();
			if (!path) return;

			// Verify it's a valid collection
			await GetCollection(path);

			const name = path.split(/[\\/]/).pop() || 'Collection';
			collectionStore.add(path);
			workspaceStore.openTab('collection', path, name, { path });
			collectionStore.setActive(path);
		} catch (e: any) {
			error = 'This directory does not appear to be a valid Myna collection.';
			console.error(e);
		}
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
			workspaceStore.openTab('collection', fullPath, data.name, { path: fullPath });
			collectionStore.setActive(fullPath);
			showCreateDialog = false;
		} catch (e: any) {
			error = e.toString() || 'Failed to create collection';
			console.error(error);
		}
	}

	function startResizing(e: MouseEvent) {
		isResizing = true;
		document.body.style.cursor = 'col-resize';
		document.body.style.userSelect = 'none';
		e.preventDefault();
	}

	function stopResizing() {
		isResizing = false;
		document.body.style.cursor = '';
		document.body.style.userSelect = '';
	}

	function onMouseMove(e: MouseEvent) {
		if (!isResizing) return;
		layoutStore.setSidebarWidth(e.clientX);
	}
</script>

<svelte:window onmousemove={onMouseMove} onmouseup={stopResizing} />

<!-- svelte-ignore a11y_no_static_element_interactions -->
<aside
	class="bg-surface-50-900-token group relative flex h-full shrink-0 flex-col border-r border-surface-500/30"
	style="width: {layoutStore.sidebarWidth}px;"
>
	{#if error}
		<div class="fixed top-4 right-4 z-50 w-full max-w-xs transition-all">
			<ErrorBanner message={error} ondismiss={() => (error = '')} />
		</div>
	{/if}

	<SidebarHeader onadd={handleAddCollection} onopen={handleOpenCollection} />

	<hr class="border-t-2!" />

	<SidebarList
		collections={collectionStore.collections}
		activeCollection={collectionStore.activeCollection}
		onselect={handleSelectCollection}
	/>

	<CreateCollectionDialog bind:open={showCreateDialog} onsubmit={onCreateSubmit} />

	<!-- Resize Handle -->
	<!-- svelte-ignore a11y_no_static_element_interactions -->
	<div
		onmousedown={startResizing}
		class="absolute top-0 -right-1 z-10 h-full w-2 cursor-col-resize transition-colors hover:bg-primary-500/30 {isResizing
			? 'bg-primary-500/50'
			: ''}"
	></div>
</aside>
