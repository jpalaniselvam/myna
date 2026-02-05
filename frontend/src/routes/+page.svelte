<script>
	import { collectionStore } from '$lib/stores/collections.svelte';
	import CollectionWorkspace from '$lib/components/dashboard/CollectionWorkspace.svelte';
	import IndexContent from '$lib/components/IndexContent.svelte';
</script>

<div class="h-full w-full">
	{#if collectionStore.activeCollection}
		<CollectionWorkspace />
	{:else}
		<IndexContent />
	{/if}
</div>
