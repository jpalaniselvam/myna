<script lang="ts">
	import { FileText, FileCode } from '@lucide/svelte';

	let { data = $bindable(''), file = $bindable('') } = $props<{
		data: string;
		file: string;
	}>();

	let mode = $state<'text' | 'file'>('text');
</script>

<div class="space-y-4">
	<div class="flex items-center justify-between">
		<div class="flex items-center gap-2">
			<FileText size={18} class="text-primary-500" />
			<span class="font-semibold">Payload</span>
		</div>
		<div class="flex gap-2">
			<button
				class="btn btn-sm {mode === 'text' ? 'preset-filled-primary-500' : 'preset-tonal'}"
				onclick={() => (mode = 'text')}
			>
				Text
			</button>
			<button
				class="btn btn-sm {mode === 'file' ? 'preset-filled-primary-500' : 'preset-tonal'}"
				onclick={() => (mode = 'file')}
			>
				File
			</button>
		</div>
	</div>

	{#if mode === 'text'}
		<div class="relative">
			<textarea
				bind:value={data}
				class="textarea min-h-[200px] font-mono text-sm"
				placeholder="Enter payload (JSON, XML, text...)"
			></textarea>
			<div class="absolute right-2 bottom-2 flex gap-2">
				<span class="preset-filled-surface-500/50 badge">Raw</span>
			</div>
		</div>
	{:else}
		<div class="variant-soft-surface card p-4">
			<label class="label">
				<span class="label-text">Select Path to File</span>
				<div class="flex gap-2">
					<input type="text" bind:value={file} class="input" placeholder="path/to/payload.json" />
					<button class="btn preset-tonal">Browse</button>
				</div>
			</label>
			<p class="mt-2 text-xs text-surface-400">
				The payload will be read from this file during execution.
			</p>
		</div>
	{/if}
</div>
