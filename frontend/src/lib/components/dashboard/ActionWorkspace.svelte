<script lang="ts">
	import { FileCode, Play } from '@lucide/svelte';
	import ActionManager from './ActionManager.svelte';

	let { tab } = $props<{
		tab: {
			id: string;
			title: string;
			metadata: {
				collectionPath: string;
				subPath: string;
				fileName: string;
				actionKind: string;
			};
		};
	}>();

	let metadata = $derived(tab?.metadata || {});
</script>

<div class="flex h-full flex-col p-6">
	<div class="mb-6 flex items-center justify-between">
		<div class="flex items-center gap-3">
			<div class="rounded-lg bg-primary-500/10 p-2">
				<FileCode size={24} class="text-primary-500" />
			</div>
			<div>
				<h1 class="text-2xl font-bold">{tab?.title}</h1>
				<p class="text-sm text-surface-400">{metadata.collectionPath}</p>
			</div>
		</div>

		<div class="flex gap-2">
			<button class="btn flex items-center gap-2 preset-filled-primary-500">
				<Play size={16} />
				<span>Run Action</span>
			</button>
		</div>
	</div>

	<div class="flex-1 overflow-hidden">
		{#if tab}
			<ActionManager {metadata} />
		{/if}
	</div>
</div>
