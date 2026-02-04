<script lang="ts">
	import { slide } from 'svelte/transition';

	let { collections, activeCollection, onselect } = $props<{
		collections: string[];
		activeCollection: string | null;
		onselect: (path: string) => void;
	}>();
</script>

<div class="flex-1 overflow-y-auto p-2">
	<nav class="list-nav">
		<ul>
			{#each collections as col}
				<li>
					<button
						class="rounded-container-token hover:variant-soft-primary w-full px-4 py-2 text-left {activeCollection ===
						col
							? 'variant-filled-secondary'
							: ''}"
						onclick={() => onselect(col)}
					>
						<span class="truncate">{col.split(/[\\/]/).pop()}</span>
					</button>
				</li>
			{/each}
		</ul>
	</nav>
</div>
