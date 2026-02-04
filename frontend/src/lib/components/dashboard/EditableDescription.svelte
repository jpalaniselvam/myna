<script lang="ts">
	import { UpdateCollection } from 'wailsjs/go/main/App';
	import { slide } from 'svelte/transition';

	let { collectionPath, description, onupdate } = $props<{
		collectionPath: string;
		description: string;
		onupdate: (newDesc: string) => void;
	}>();

	let editing = $state(false);
	let editValue = $state(description);
	let error = $state('');
	let loading = $state(false);

	// Sync if prop changes and not editing
	$effect(() => {
		if (!editing) editValue = description;
	});

	async function handleSave() {
		if (editValue === description) {
			editing = false;
			return;
		}

		loading = true;
		error = '';
		try {
			// UpdateCollection(collectionPath, name, desc, pre)
			// We only want to update description. Name and pre usage:
			// Backend logic: name="" -> no rename. pre=nil -> no update to pre.
			await UpdateCollection(collectionPath, '', editValue, new Map<string, string>());
			onupdate(editValue);
			editing = false;
		} catch (e: any) {
			error = e.toString();
		} finally {
			loading = false;
		}
	}

	function focus(el: HTMLElement) {
		el.focus();
	}

	function onKeyDown(e: KeyboardEvent) {
		if (e.key === 'Enter' && !e.shiftKey) {
			e.preventDefault();
			handleSave();
		}
		if (e.key === 'Escape') {
			editing = false;
			editValue = description;
			error = '';
		}
	}
</script>

<div class="group relative">
	{#if error}
		<div class="alert variant-filled-error mb-2" transition:slide>
			{error}
			<button class="ml-2 font-bold" onclick={() => (error = '')}>X</button>
		</div>
	{/if}

	{#if !editing}
		<div
			role="button"
			tabindex="0"
			class="hover:bg-surface-200-700-token -ml-2 min-h-[1.5em] cursor-text rounded p-2"
			onclick={() => {
				editing = true;
				editValue = description;
			}}
			onkeydown={(e) => {
				if (e.key === 'Enter') editing = true;
			}}
			title="Click to edit description"
		>
			<p class="text-surface-600-300-token text-sm whitespace-pre-wrap">
				{description || 'No description provided.'}
			</p>
		</div>
	{:else}
		<textarea
			class="textarea w-full p-2 text-sm"
			rows="2"
			bind:value={editValue}
			use:focus
			onkeydown={onKeyDown}
			onblur={handleSave}
			disabled={loading}
		></textarea>
		{#if loading}
			<div class="absolute right-2 bottom-2">
				<!-- Simple spinner or text -->
				<span class="text-xs text-surface-400">Saving...</span>
			</div>
		{/if}
	{/if}
</div>
