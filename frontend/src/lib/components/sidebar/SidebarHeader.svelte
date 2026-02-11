<script lang="ts">
	import { Plus, FolderOpen } from '@lucide/svelte';

	let { onadd, onopen } = $props<{ onadd: () => void; onopen: () => void }>();
</script>

<div class="flex flex-col border-b border-surface-500/30 p-4">
	<div class="flex items-center justify-between">
		<h2 class="flex-1">Collections</h2>
		<!-- open icon to open existing collections -->
		<button class="btn-icon btn-sm" onclick={onopen} title="Open collection">
			<FolderOpen size={20} />
		</button>
		<!-- plus icon to add collections -->
		<button class="btn-icon btn-sm" onclick={onadd} title="Create collection">
			<Plus size={20} />
		</button>
	</div>
</div>
