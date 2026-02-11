<script lang="ts">
	import { FolderOpen, X } from '@lucide/svelte';
	import { SelectDirectory } from '../../../../wailsjs/go/main/App';

	let { open = $bindable(false), onsubmit } = $props<{
		open: boolean;
		onsubmit: (data: { name: string; baseDir: string }) => void;
	}>();

	let name = $state('');
	let baseDir = $state('');
	let error = $state('');

	function handleClose() {
		open = false;
		name = '';
		baseDir = '';
		error = '';
	}

	async function handleBrowse() {
		try {
			const selected = await SelectDirectory();
			if (selected) {
				baseDir = selected;
			}
		} catch (e: any) {
			console.error('Failed to select directory:', e);
		}
	}

	function handleSubmit() {
		if (!name.trim()) {
			error = 'Name is required';
			return;
		}
		if (!baseDir.trim()) {
			error = 'Base directory is required';
			return;
		}
		onsubmit({ name: name.trim(), baseDir: baseDir.trim() });
		handleClose();
	}

	function handleKeyDown(e: KeyboardEvent) {
		if (e.key === 'Escape') {
			handleClose();
		} else if (e.key === 'Enter') {
			handleSubmit();
		}
	}
</script>

{#if open}
	<!-- Backdrop -->
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-surface-900/60 p-4 backdrop-blur-sm"
		role="dialog"
		aria-modal="true"
		tabindex="-1"
		onkeydown={handleKeyDown}
	>
		<!-- Dialog Card -->
		<div
			class="bg-surface-100-800-token w-full max-w-lg overflow-hidden card border border-surface-500/20 shadow-2xl"
		>
			<!-- Header -->
			<header class="flex items-center justify-between border-b border-surface-500/20 p-6">
				<h2 class="h2 font-bold">New Collection</h2>
				<button
					class="btn-icon btn-sm transition-colors hover:bg-surface-500/20"
					onclick={handleClose}
				>
					<X size={20} />
				</button>
			</header>

			<!-- Content -->
			<div class="space-y-6 p-6">
				<div class="text-surface-600-300-token space-y-2">
					<p>Create a new Myna collection to organize your AWS serverless requests.</p>
				</div>

				<div class="grid gap-4">
					<label class="label">
						<span class="font-semibold">Collection Name</span>
						<input
							class="input"
							type="text"
							placeholder="e.g. my-awesome-project"
							bind:value={name}
						/>
					</label>

					<label class="label">
						<span class="font-semibold">Base Directory</span>
						<div class="flex gap-2">
							<input
								class="input"
								type="text"
								placeholder="/path/to/projects"
								bind:value={baseDir}
							/>
							<button class="btn shrink-0 preset-tonal" onclick={handleBrowse}>
								<FolderOpen size={18} class="mr-2" />
								Browse
							</button>
						</div>
					</label>
				</div>

				{#if error}
					<p class="text-sm text-error-500">{error}</p>
				{/if}
			</div>

			<!-- Footer -->
			<footer class="flex justify-end gap-3 border-t border-surface-500/20 bg-surface-500/5 p-6">
				<button class="preset-outline-surface-500 btn" onclick={handleClose}> Cancel </button>
				<button class="btn preset-filled-primary-500 shadow-md" onclick={handleSubmit}>
					Create Collection
				</button>
			</footer>
		</div>
	</div>
{/if}
