<script lang="ts">
	import { collectionStore } from '$lib/stores/collections.svelte';
	import { CreateCollection, SelectDirectory } from '../../../../wailsjs/go/main/App';
	import SidebarHeader from './SidebarHeader.svelte';
	import SidebarList from './SidebarList.svelte';
	import ErrorBanner from './ErrorBanner.svelte';

	let creating = $state(false);
	let error = $state('');
	let picking = false;

	function handleStartCreation() {
		creating = true;
		error = '';
	}

	function handleCancelCreation() {
		if (picking) return;
		creating = false;
		error = '';
	}

	async function handleNameSubmit(name: string) {
		error = '';
		// Validation
		if (!name) {
			error = 'Name cannot be empty';
			return;
		}
		if (/[<>:"/\\|?*]/.test(name)) {
			error = 'Name contains invalid characters';
			return;
		}

		// Step 2: Location Selection
		try {
			picking = true;
			const baseDir = await SelectDirectory();
			picking = false;

			if (!baseDir) {
				// Return to input
				return;
			}

			// Step 3: Creation
			await CreateCollection(baseDir, name, '');

			// Success
			const separator = baseDir.includes('\\') ? '\\' : '/';
			const cleanBase = baseDir.endsWith(separator) ? baseDir.slice(0, -1) : baseDir;
			const fullPath = `${cleanBase}${separator}${name}`;

			collectionStore.add(fullPath);
			collectionStore.setActive(fullPath);

			// Success -> Idle
			creating = false;
			error = '';
		} catch (e: any) {
			picking = false;
			error = e.toString() || 'Failed to create collection';
		}
	}

	function handleSelectCollection(path: string) {
		collectionStore.setActive(path);
	}
</script>

<aside class="bg-surface-50-900-token flex h-full w-64 flex-col border-r border-surface-500/30">
	{#if error}
		<div class="p-2">
			<ErrorBanner message={error} ondismiss={() => (error = '')} />
		</div>
	{/if}

	<SidebarHeader
		{creating}
		onstart={handleStartCreation}
		onsubmit={handleNameSubmit}
		oncancel={handleCancelCreation}
	/>

	<hr class="!border-t-2" />

	<SidebarList
		collections={collectionStore.collections}
		activeCollection={collectionStore.activeCollection}
		onselect={handleSelectCollection}
	/>
</aside>
