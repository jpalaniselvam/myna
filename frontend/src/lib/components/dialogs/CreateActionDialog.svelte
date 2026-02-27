<script lang="ts">
	import { X } from '@lucide/svelte';
	import { GetActionKinds } from 'wailsjs/go/main/App';

	let { open = $bindable(false), onsubmit } = $props<{
		open: boolean;
		onsubmit: (data: { kind: string; name: string }) => void;
	}>();

	let actionKinds = $state<string[]>([]);
	let selectedKind = $state('');
	let actionName = $state('');
	let error = $state('');

	$effect(() => {
		if (open) {
			GetActionKinds().then((kinds) => {
				actionKinds = kinds;
				if (kinds.length > 0 && !selectedKind) {
					selectedKind = kinds[0];
				}
			});
		}
	});

	function handleClose() {
		open = false;
		actionName = '';
		error = '';
	}

	function handleSubmit() {
		if (!actionName.trim()) {
			error = 'Name is required';
			return;
		}
		if (!selectedKind) {
			error = 'Action kind is required';
			return;
		}
		onsubmit({ kind: selectedKind, name: actionName.trim() });
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
			class="bg-surface-100-800-token w-full max-w-lg card border border-surface-500/20 shadow-2xl"
		>
			<!-- Header -->
			<header class="flex items-center justify-between border-b border-surface-500/20 p-6">
				<h2 class="h2 font-bold">New Action</h2>
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
					<p>Create a new action to perform tasks like invoking Lambdas or sending SQS messages.</p>
				</div>

				<div class="grid gap-4">
					<label class="label">
						<span class="font-semibold">Action Name</span>
						<input
							class="input"
							type="text"
							placeholder="e.g. process-order"
							bind:value={actionName}
						/>
					</label>

					<label class="label">
						<span class="font-semibold">Action Kind</span>
						<select bind:value={selectedKind} class="select">
							{#each actionKinds as kind}
								<option value={kind}>{kind}</option>
							{/each}
						</select>
					</label>
				</div>

				{#if error}
					<p class="text-sm text-error-500">{error}</p>
				{/if}
			</div>

			<!-- Footer -->
			<footer class="flex justify-end gap-3 border-t border-surface-500/20 bg-surface-500/5 p-6">
				<button class="preset-outline-surface-500 btn" onclick={handleClose}> Cancel </button>
				<button
					class="btn preset-filled-primary-500 shadow-md"
					onclick={handleSubmit}
					disabled={!actionName.trim() || !selectedKind}
				>
					Create Action
				</button>
			</footer>
		</div>
	</div>
{/if}
