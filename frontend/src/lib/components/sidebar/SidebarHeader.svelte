<script lang="ts">
	import { Plus } from '@lucide/svelte';

	let { onadd } = $props<{ onadd: () => void }>();
</script>

<div class="flex flex-col border-b border-surface-500/30 p-4">
	<div class="flex h-10 items-center justify-between">
		<h2 class="h4 font-bold">Collections</h2>
		<!-- plus icon to add collections -->
		<button class="btn-icon btn-sm" onclick={onadd}>
			<Plus size={20} />
		</button>
	</div>
</div>
