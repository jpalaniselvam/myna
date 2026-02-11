<script lang="ts">
	import { slide } from 'svelte/transition';

	let { collections, activeCollection, onselect } = $props<{
		collections: string[];
		activeCollection: string | null;
		onselect: (path: string) => void;
	}>();
</script>

<div class="flex-1 overflow-y-auto">
	<nav class="list-nav">
		<ul>
			{#each collections as col}
				<li>
					<button
						class="w-full px-4 py-2 text-left transition-colors {activeCollection &&
						activeCollection.replace(/\\/g, '/') === col.replace(/\\/g, '/')
							? 'preset-filled-primary-500'
							: 'hover:bg-surface-500'}"
						onclick={() => onselect(col)}
					>
						<span class="truncate">{col.split(/[\\/]/).pop()}</span>
					</button>
				</li>
			{/each}
		</ul>
	</nav>
</div>
