<script lang="ts">
	import { onMount } from 'svelte';

	let { creating, onstart, onsubmit, oncancel } = $props<{
		creating: boolean;
		onstart: () => void;
		onsubmit: (name: string) => void;
		oncancel: () => void;
	}>();

	let name = $state('');
	let inputRef: HTMLInputElement | undefined = $state();

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Enter') {
			if (name.trim()) {
				onsubmit(name.trim());
			}
		} else if (e.key === 'Escape') {
			oncancel();
			name = '';
		}
	}

	$effect(() => {
		if (creating && inputRef) {
			inputRef.focus();
		}
	});

	// Reset name when creating mode ends
	$effect(() => {
		if (!creating) {
			name = '';
		}
	});
</script>

<div class="flex flex-col border-b border-surface-500/30 p-4">
	<div class="flex h-10 items-center">
		{#if !creating}
			<button class="variant-filled-primary btn w-full gap-2" onclick={onstart}>
				<!-- Plus Icon -->
				<span class="text-xl leading-none">+</span>
				<span>Create Collection</span>
			</button>
		{:else}
			<input
				bind:this={inputRef}
				type="text"
				class="input"
				placeholder="Enter collection name..."
				bind:value={name}
				onkeydown={handleKeydown}
				onblur={() => {
					// Optional: close on blur delay to allow clicks
					setTimeout(() => oncancel(), 200);
				}}
			/>
		{/if}
	</div>
</div>
