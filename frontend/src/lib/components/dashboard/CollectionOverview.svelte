<script lang="ts">
	import type { collection } from 'wailsjs/go/models';
	import { Folder, Globe, Zap, Share2, FileText, Pencil, Save, X } from '@lucide/svelte';
	import { UpdateCollection } from 'wailsjs/go/main/App';
	import { slide } from 'svelte/transition';

	let { collectionData, collectionPath, onupdate } = $props<{
		collectionData: collection.CollectionResponse;
		collectionPath: string;
		onupdate: (newDesc: string) => void;
	}>();

	let collectionName = $derived(collectionPath.split(/[\\/]/).pop() || 'Collection');
	let envCount = $derived(collectionData.environments?.length || 0);
	let actionCount = $derived(Object.keys(collectionData.actions || {}).length);

	let editing = $state(false);
	let editValue = $state('');
	let loading = $state(false);
	let error = $state('');

	$effect(() => {
		if (!editing) editValue = collectionData.description || '';
	});

	async function handleSave() {
		loading = true;
		error = '';
		try {
			await UpdateCollection(collectionPath, '', editValue, collectionData.pre);
			onupdate(editValue);
			editing = false;
		} catch (e: any) {
			error = e.toString();
		} finally {
			loading = false;
		}
	}

	function cancelEdit() {
		editing = false;
		editValue = collectionData.description || '';
		error = '';
	}
</script>

<div class="flex h-full w-full gap-8 p-6">
	<!-- Left Sidebar -->
	<div class="flex w-80 flex-col gap-8">
		<!-- Collection Name -->
		<div class="flex items-center gap-3">
			<div class="rounded preset-filled-primary-500 p-2">
				<Zap size={24} class="text-white" />
			</div>
			<h2 class="h2 font-bold">{collectionName}</h2>
		</div>

		<!-- Info List -->
		<div class="flex flex-col gap-6">
			<!-- Location -->
			<div class="flex gap-4">
				<div class="h-fit rounded-lg preset-tonal-surface p-3">
					<Folder size={20} class="text-blue-500" />
				</div>
				<div class="flex flex-col">
					<span class="text-sm font-bold opacity-80">Location</span>
					<span class="text-xs break-all opacity-60">{collectionPath}</span>
				</div>
			</div>

			<!-- Environments -->
			<div class="flex gap-4">
				<div class="h-fit rounded-lg preset-tonal-surface p-3">
					<Globe size={20} class="text-green-500" />
				</div>
				<div class="flex flex-col">
					<span class="text-sm font-bold opacity-80">Environments</span>
					<span class="text-xs opacity-60"
						>{envCount} environment{envCount === 1 ? '' : 's'} configured</span
					>
				</div>
			</div>

			<!-- Requests -->
			<div class="flex gap-4">
				<div class="h-fit rounded-lg preset-tonal-surface p-3">
					<Zap size={20} class="text-secondary-500" />
				</div>
				<div class="flex flex-col">
					<span class="text-sm font-bold opacity-80">Requests</span>
					<span class="text-xs opacity-60"
						>{actionCount} request{actionCount === 1 ? '' : 's'} in collection</span
					>
				</div>
			</div>

			<!-- Share -->
			<div class="flex gap-4">
				<div class="h-fit rounded-lg preset-tonal-surface p-3">
					<Share2 size={20} class="text-tertiary-500" />
				</div>
				<div class="flex flex-col">
					<span class="text-sm font-bold opacity-80">Share</span>
					<button class="text-left text-xs text-primary-500 hover:underline"
						>Share Collection</button
					>
				</div>
			</div>
		</div>
	</div>

	<!-- Right Content (Documentation) -->
	<div
		class="bg-surface-50-950-token border-surface-200-800-token flex flex-1 flex-col rounded-xl border p-6"
	>
		<div class="border-surface-200-800-token mb-6 flex items-center justify-between border-b pb-4">
			<div class="flex items-center gap-2">
				<FileText size={20} />
				<h3 class="h3 font-bold">Documentation</h3>
			</div>
			<div class="flex items-center gap-2">
				{#if editing}
					<button
						class="variant-filled-success btn-icon btn-icon-sm"
						onclick={handleSave}
						disabled={loading}
					>
						<Save size={16} />
					</button>
					<button
						class="variant-filled-surface btn-icon btn-icon-sm"
						onclick={cancelEdit}
						disabled={loading}
					>
						<X size={16} />
					</button>
				{:else}
					<button
						class="variant-soft-surface btn-icon btn-icon-sm"
						onclick={() => (editing = true)}
					>
						<Pencil size={16} />
					</button>
				{/if}
			</div>
		</div>

		{#if error}
			<div class="alert variant-filled-error mb-4" transition:slide>
				<p>{error}</p>
			</div>
		{/if}

		<div class="flex-1 overflow-auto">
			{#if editing}
				<textarea
					class="bg-surface-100-900-token textarea h-full w-full resize-none p-4 font-mono text-sm"
					bind:value={editValue}
					placeholder="Enter collection documentation (Supports Markdown)..."
					disabled={loading}
				></textarea>
			{:else}
				<div class="prose max-w-none whitespace-pre-wrap dark:prose-invert">
					{collectionData.description ||
						'Welcome to your collection documentation! Click the edit icon to get started.'}
				</div>
			{/if}
		</div>
	</div>
</div>
