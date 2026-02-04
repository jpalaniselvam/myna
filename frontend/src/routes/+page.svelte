<script>
	import { collectionStore } from '$lib/stores/collections.svelte';
	import CollectionWorkspace from '$lib/components/dashboard/CollectionWorkspace.svelte';
</script>

<div class="h-full w-full">
	{#if collectionStore.activeCollection}
		<CollectionWorkspace />
	{:else}
		<div class="flex h-full flex-col items-center justify-center text-surface-500">
			<h1 class="mb-4 h1 font-bold">Welcome to Myna</h1>
			<p class="text-lg">Select a collection from the sidebar to get started.</p>
		</div>
	{/if}
</div>
