<script lang="ts">
	import { GetAction, UpdateAction } from 'wailsjs/go/main/App';
	import { Save, CircleAlert, CircleCheck, Loader } from '@lucide/svelte';
	import PayloadEditor from './PayloadEditor.svelte';

	let { metadata } = $props<{
		metadata: {
			collectionPath: string;
			subPath: string;
			fileName: string;
		};
	}>();

	let actionData = $state<any>(null);
	let loading = $state(true);
	let saving = $state(false);
	let status = $state<{ type: 'success' | 'error'; message: string } | null>(null);

	// Load data when metadata changes
	$effect(() => {
		loadAction();
	});

	async function loadAction() {
		loading = true;
		status = null;
		try {
			const data = await GetAction({
				collection_path: metadata.collectionPath,
				sub_path: metadata.subPath,
				file_name: metadata.fileName
			});

			// Ensure structure exists
			if (!data.payload) {
				data.payload = { data: '', file: '' };
			} else {
				if (data.payload.data === undefined || data.payload.data === null) data.payload.data = '';
				if (data.payload.file === undefined || data.payload.file === null) data.payload.file = '';
			}

			const kind = data.kind || '';
			if (kind.startsWith('lambda')) {
				if (!data.lambda) {
					data.lambda = {
						function_name: '',
						invocation_type: 'RequestResponse',
						qualifier: '',
						client_context: {}
					};
				} else {
					if (data.lambda.function_name === undefined) data.lambda.function_name = '';
					if (data.lambda.invocation_type === undefined)
						data.lambda.invocation_type = 'RequestResponse';
					if (data.lambda.qualifier === undefined) data.lambda.qualifier = '';
					if (!data.lambda.client_context) data.lambda.client_context = {};
				}
			}

			if (data.description === undefined) data.description = '';
			if (data.timeout_ms === undefined || data.timeout_ms === 0) data.timeout_ms = 30000;

			actionData = data;
		} catch (e: any) {
			status = { type: 'error', message: 'Failed to load action: ' + e.message };
			console.error('Failed to load action:', e);
			actionData = null;
		} finally {
			loading = false;
		}
	}

	async function handleSave() {
		if (!actionData) return;
		saving = true;
		status = null;
		try {
			await UpdateAction({
				collection_path: metadata.collectionPath,
				sub_path: metadata.subPath,
				file_name: metadata.fileName,
				data: actionData
			});
			status = { type: 'success', message: 'Action saved successfully' };
			setTimeout(() => {
				if (status?.type === 'success') status = null;
			}, 3000);
		} catch (e: any) {
			status = { type: 'error', message: 'Failed to save action: ' + e.message };
		} finally {
			saving = false;
		}
	}
</script>

{#if loading}
	<div class="flex h-full items-center justify-center">
		<Loader size={32} class="animate-spin text-primary-500" />
	</div>
{:else if actionData}
	<div class="flex h-full flex-col gap-6">
		<div class="flex justify-end">
			<button
				class="btn flex items-center gap-2 preset-filled-primary-500 shadow-lg"
				onclick={handleSave}
				disabled={saving}
			>
				{#if saving}
					<Loader size={16} class="animate-spin" />
				{:else}
					<Save size={16} />
				{/if}
				<span>Save Action</span>
			</button>
		</div>

		{#if status}
			<div
				class="alert {status.type === 'success'
					? 'variant-filled-success'
					: 'variant-filled-error'} flex items-center gap-3"
			>
				{#if status.type === 'success'}
					<CircleCheck size={18} />
				{:else}
					<CircleAlert size={18} />
				{/if}
				<span>{status.message}</span>
			</div>
		{/if}

		<div class="grid flex-1 gap-6 overflow-auto pr-2">
			<!-- Common Metadata -->
			<section class="variant-soft-surface space-y-4 card">
				<div class="grid gap-4 sm:grid-cols-2">
					<label class="label">
						<span class="label-text font-semibold">Description</span>
						<input
							type="text"
							bind:value={actionData.description}
							class="input"
							placeholder="What does this action do?"
						/>
					</label>
				</div>
			</section>

			<!-- Kind Specific Configuration -->
			<h3>placeholder for {actionData.kind}</h3>

			<!-- Payload Section -->
			<section class="variant-soft-surface space-y-4">
				<PayloadEditor bind:data={actionData.payload.data} bind:file={actionData.payload.file} />
			</section>
		</div>
	</div>
{:else}
	<div class="flex h-full items-center justify-center text-error-500">
		<CircleAlert size={32} class="mb-2" />
		<p>Action data not found.</p>
	</div>
{/if}
