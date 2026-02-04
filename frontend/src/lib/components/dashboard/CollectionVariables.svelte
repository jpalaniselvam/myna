<script lang="ts">
	import { UpdateCollection } from 'wailsjs/go/main/App';
	import { slide } from 'svelte/transition';

	let { collectionPath, variables, onupdate } = $props<{
		collectionPath: string;
		variables: Record<string, any>;
		onupdate: (newVars: Record<string, any>) => void;
	}>();

	// Convert object to array of entries for editing
	let pairs = $state<{ key: string; value: string }[]>([]);
	let loading = $state(false);
	let error = $state('');
	let isDirty = $state(false);

	$effect(() => {
		// Initialize pairs from variables prop
		if (variables) {
			pairs = Object.entries(variables).map(([k, v]) => ({
				key: k,
				value: String(v) // Force string for simplicity in UI for now
			}));
		} else {
			pairs = [];
		}
	});

	function addPair() {
		pairs.push({ key: '', value: '' });
		isDirty = true;
	}

	function removePair(index: number) {
		pairs.splice(index, 1);
		isDirty = true;
	}

	function markDirty() {
		isDirty = true;
	}

	async function handleSave() {
		loading = true;
		error = '';

		try {
			// Reconstruct object
			const pre: Record<string, any> = {};
			for (const p of pairs) {
				if (p.key.trim()) {
					pre[p.key.trim()] = p.value;
				}
			}

			// Call UpdateCollection
			// Name: "", Desc: "" (to skip updating them? Backend checks != "")
			// BEWARE: Backend implementation for Desc: if desc != "" update it.
			// If we pass "", it keeps existing.
			// But we must ensure we don't accidentally wipe description if API requires specific separate calls?
			// collection.go Update:
			// if desc != "" { coll.Metadata.Description = desc }
			// if pre != nil { coll.Vars.Pre = pre }
			// So passing "" for desc is safe.

			await UpdateCollection(collectionPath, '', '', pre);

			onupdate(pre);
			isDirty = false;
		} catch (e: any) {
			error = e.toString();
		} finally {
			loading = false;
		}
	}
</script>

<div class="space-y-4">
	<div class="flex items-center justify-between">
		<h3 class="h3">Collection Variables</h3>
		<button class="variant-filled-primary btn" disabled={!isDirty || loading} onclick={handleSave}>
			{loading ? 'Saving...' : 'Save Changes'}
		</button>
	</div>

	{#if error}
		<div class="alert variant-filled-error" transition:slide>
			{error}
			<button class="ml-2 font-bold" onclick={() => (error = '')}>X</button>
		</div>
	{/if}

	<div class="table-container">
		<table class="table-hover table">
			<thead>
				<tr>
					<th class="w-1/3">Key</th>
					<th>Value</th>
					<th class="w-16">Action</th>
				</tr>
			</thead>
			<tbody>
				{#each pairs as pair, i}
					<tr>
						<td>
							<input
								class="input"
								type="text"
								placeholder="Key"
								bind:value={pair.key}
								oninput={markDirty}
							/>
						</td>
						<td>
							<input
								class="input"
								type="text"
								placeholder="Value"
								bind:value={pair.value}
								oninput={markDirty}
							/>
						</td>
						<td>
							<button
								class="variant-filled-error btn-icon btn-icon-sm"
								onclick={() => removePair(i)}
							>
								-
							</button>
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>

	<button class="variant-ghost-secondary btn w-full" onclick={addPair}> + Add Variable </button>
</div>
