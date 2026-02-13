<script lang="ts">
	import { UpdateCollection } from 'wailsjs/go/main/App';
	import { types } from 'wailsjs/go/models';
	import { slide } from 'svelte/transition';
	import { Save, Lock } from '@lucide/svelte';

	let { collectionPath, credentials, variables, onupdate } = $props<{
		collectionPath: string;
		credentials: types.Credentials;
		variables: Record<string, any>;
		onupdate: (newCreds: types.Credentials) => void;
	}>();

	let creds = $state({
		region: '',
		profile: '',
		role_arn: ''
	});

	let loading = $state(false);
	let error = $state('');
	let isDirty = $state(false);

	$effect(() => {
		if (credentials) {
			creds = {
				region: credentials.region || '',
				profile: credentials.profile || '',
				role_arn: credentials.role_arn || ''
			};
		}
	});

	function markDirty() {
		isDirty = true;
	}

	async function handleSave() {
		loading = true;
		error = '';

		try {
			// Call UpdateCollection with the new signature
			// We need to pass the current pre (variables) as well to avoid wiping them
			// Since we don't have access to pre here easily, we might need it as a prop or
			// the backend should handle partial updates.
			// Current backend implementation: if pre != nil { coll.Vars.Pre = pre }
			// So if we pass null/undefined, it might keep existing (if Wails handles it)
			// Actually, let's look at collection.go:
			// if pre != nil { coll.Vars.Pre = pre }
			// if we pass an empty map or null, it might be tricky.

			// Better to pass everything or change backend to handle partial updates.
			// But for now, let's assume we can pass null for things we don't want to change.
			// Wait, the Go signature is pre map[string]interface{}.
			// If Wails sends null, it will be nil in Go.

			await UpdateCollection(collectionPath, '', '', variables, creds);

			onupdate({ ...creds });
			isDirty = false;
		} catch (e: any) {
			error = e.toString();
		} finally {
			loading = false;
		}
	}
</script>

<div class="flex h-full flex-col space-y-6">
	<div class="border-surface-200-800-token flex items-center justify-between border-b pb-4">
		<div class="flex items-center gap-2">
			<Lock size={20} class="text-primary-500" />
			<h3 class="h3 font-bold">AWS Credentials</h3>
		</div>
		<button
			class="btn flex items-center gap-2 preset-filled-primary-500"
			disabled={!isDirty || loading}
			onclick={handleSave}
		>
			{#if loading}
				<span class="animate-spin">⌛</span>
				Saving...
			{:else}
				<Save size={16} />
				Save Credentials
			{/if}
		</button>
	</div>

	{#if error}
		<div class="alert variant-filled-error" transition:slide>
			<p>{error}</p>
			<button class="ml-2 font-bold" onclick={() => (error = '')}>X</button>
		</div>
	{/if}

	<div class="bg-surface-50-950-token border-surface-200-800-token space-y-4 card border p-6">
		<div class="grid grid-cols-1 gap-4">
			<label class="label">
				<span class="text-sm font-bold opacity-80">AWS Region</span>
				<input
					class="mt-1 input"
					type="text"
					placeholder="e.g. us-east-1"
					bind:value={creds.region}
					oninput={markDirty}
				/>
				<p class="mt-1 text-xs opacity-60">The AWS region to use for this collection.</p>
			</label>

			<label class="label">
				<span class="text-sm font-bold opacity-80">AWS Profile</span>
				<input
					class="mt-1 input"
					type="text"
					placeholder="e.g. default"
					bind:value={creds.profile}
					oninput={markDirty}
				/>
				<p class="mt-1 text-xs opacity-60">
					The AWS CLI profile name from your ~/.aws/credentials file.
				</p>
			</label>

			<label class="label">
				<span class="text-sm font-bold opacity-80">Role ARN (Optional)</span>
				<input
					class="mt-1 input"
					type="text"
					placeholder="e.g. arn:aws:iam::123456789012:role/MyRole"
					bind:value={creds.role_arn}
					oninput={markDirty}
				/>
				<p class="mt-1 text-xs opacity-60">
					The IAM Role ARN to assume for actions in this collection.
				</p>
			</label>
		</div>
	</div>

	<div class="flex-1"></div>

	<div class="alert variant-soft-surface text-sm">
		<p>
			These credentials will apply to all actions within this collection unless overridden at the
			action level.
		</p>
	</div>
</div>
