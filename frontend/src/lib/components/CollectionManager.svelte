<script lang="ts">
	import {
		CreateCollection,
		UpdateCollection,
		DeleteCollection,
		GetCollection,
		CreateAction,
		UpdateAction,
		DeleteAction,
		GetAction
	} from '../../../wailsjs/go/main/App';
	import { Tabs } from '@skeletonlabs/skeleton-svelte';

	let activeTab = $state('create');
	let result = $state('');
	let error = $state('');

	// Create State
	let createBaseDir = $state('');
	let createName = $state('');
	let createDesc = $state('');

	// Update State
	let updatePath = $state('');
	let updateName = $state('');
	let updateDesc = $state('');
	let updatePreJson = $state('{}');

	// Delete State
	let deletePath = $state('');

	// Get State
	let getPath = $state('');
	let collectionData: any = $state(null);

	// Action State
	let actionOp = $state('create');
	let actionCollectionPath = $state('');
	let actionSubPath = $state('');
	let actionFileName = $state('');
	let actionData = $state('{}');
	let actionResult: any = $state(null);

	async function handleCreate() {
		try {
			error = '';
			result = '';
			if (!createBaseDir || !createName) {
				error = 'Base Directory and Name are required.';
				return;
			}
			await CreateCollection(createBaseDir, createName, createDesc);
			result = `Collection '${createName}' created successfully at ${createBaseDir}`;
		} catch (e: any) {
			error = String(e);
		}
	}

	async function handleUpdate() {
		try {
			error = '';
			result = '';
			const pre = JSON.parse(updatePreJson);
			await UpdateCollection(updatePath, updateName, updateDesc, pre, null);
			result = `Collection updated successfully.`;
		} catch (e: any) {
			error = String(e);
		}
	}

	async function handleDelete() {
		try {
			error = '';
			result = '';
			await DeleteCollection(deletePath);
			result = `Collection deleted successfully.`;
		} catch (e: any) {
			error = String(e);
		}
	}

	async function handleGet() {
		try {
			error = '';
			result = '';
			collectionData = null;
			collectionData = await GetCollection(getPath);
			result = `Collection loaded.`;
		} catch (e: any) {
			error = String(e);
		}
	}

	async function handleAction() {
		try {
			error = '';
			result = '';
			actionResult = null;

			let dataObj = null;
			if (actionOp === 'create' || actionOp === 'update') {
				try {
					dataObj = JSON.parse(actionData);
				} catch (e) {
					throw new Error('Invalid JSON Data');
				}
			}

			if (actionOp === 'create') {
				await CreateAction({
					collection_path: actionCollectionPath,
					sub_path: actionSubPath,
					file_name: actionFileName,
					data: dataObj
				});
				result = 'Action Created';
			} else if (actionOp === 'update') {
				await UpdateAction({
					collection_path: actionCollectionPath,
					sub_path: actionSubPath,
					file_name: actionFileName,
					data: dataObj
				});
				result = 'Action Updated';
			} else if (actionOp === 'delete') {
				await DeleteAction({
					collection_path: actionCollectionPath,
					sub_path: actionSubPath,
					file_name: actionFileName
				});
				result = 'Action Deleted';
			} else if (actionOp === 'get') {
				const res = await GetAction({
					collection_path: actionCollectionPath,
					sub_path: actionSubPath,
					file_name: actionFileName
				});
				actionResult = res;
				result = 'Action Retrieved';
			}
		} catch (e: any) {
			error = String(e);
		}
	}
</script>

<div class="mx-auto w-full max-w-4xl space-y-8 p-4">
	<h2 class="h2 font-bold">Collection Manager</h2>

	<Tabs value={activeTab} onValueChange={(e) => (activeTab = e.value)}>
		<Tabs.List class="grid w-full grid-cols-5 gap-1">
			<Tabs.Trigger value="create">Create</Tabs.Trigger>
			<Tabs.Trigger value="update">Update</Tabs.Trigger>
			<Tabs.Trigger value="delete">Delete</Tabs.Trigger>
			<Tabs.Trigger value="get">Get</Tabs.Trigger>
			<Tabs.Trigger value="actions">Actions</Tabs.Trigger>
		</Tabs.List>

		<!-- CREATE -->
		<Tabs.Content value="create" class="space-y-4 pt-4">
			<div class="space-y-4 card p-6">
				<label class="label">
					<span>Base Directory</span>
					<input
						class="input"
						type="text"
						bind:value={createBaseDir}
						placeholder="/path/to/collections"
					/>
				</label>
				<label class="label">
					<span>Name</span>
					<input class="input" type="text" bind:value={createName} placeholder="my-collection" />
				</label>
				<label class="label">
					<span>Description</span>
					<input
						class="input"
						type="text"
						bind:value={createDesc}
						placeholder="My awesome collection"
					/>
				</label>
				<div class="flex justify-end">
					<button class="variant-filled-primary btn" onclick={handleCreate}>
						Create Collection
					</button>
				</div>
			</div>
		</Tabs.Content>

		<!-- UPDATE -->
		<Tabs.Content value="update" class="space-y-4 pt-4">
			<div class="space-y-4 card p-6">
				<label class="label">
					<span>Collection Path</span>
					<input
						class="input"
						type="text"
						bind:value={updatePath}
						placeholder="/path/to/collections/my-collection"
					/>
				</label>
				<label class="label">
					<span>New Name (Optional)</span>
					<input
						class="input"
						type="text"
						bind:value={updateName}
						placeholder="Leave empty to keep current name"
					/>
				</label>
				<label class="label">
					<span>Description</span>
					<input class="input" type="text" bind:value={updateDesc} />
				</label>
				<label class="label">
					<span>Pre-variables (JSON)</span>
					<textarea class="textarea" bind:value={updatePreJson} rows="4"></textarea>
				</label>
				<div class="flex justify-end">
					<button class="variant-filled-warning btn" onclick={handleUpdate}>
						Update Collection
					</button>
				</div>
			</div>
		</Tabs.Content>

		<!-- DELETE -->
		<Tabs.Content value="delete" class="space-y-4 pt-4">
			<div class="space-y-4 card p-6">
				<label class="label">
					<span>Collection Path</span>
					<input
						class="input"
						type="text"
						bind:value={deletePath}
						placeholder="/path/to/collections/my-collection"
					/>
				</label>
				<div class="flex justify-end">
					<button class="variant-filled-error btn" onclick={handleDelete}>
						Delete Collection
					</button>
				</div>
			</div>
		</Tabs.Content>

		<!-- GET -->
		<Tabs.Content value="get" class="space-y-4 pt-4">
			<div class="space-y-4 card p-6">
				<label class="label">
					<span>Collection Path</span>
					<input
						class="input"
						type="text"
						bind:value={getPath}
						placeholder="/path/to/collections/my-collection"
					/>
				</label>
				<div class="flex justify-end">
					<button class="variant-filled-success btn" onclick={handleGet}>
						Get Collection Details
					</button>
				</div>

				{#if collectionData}
					<div class="variant-soft space-y-4 card p-4">
						<h3 class="h3">Details</h3>

						<div class="space-y-2">
							<span class="font-bold">Environments:</span>
							<ul class="list-inside list-disc">
								{#each collectionData.environments as env}
									<li>{env}</li>
								{/each}
							</ul>
						</div>

						<div class="space-y-2">
							<span class="font-bold">Actions:</span>
							<pre class="pre">{JSON.stringify(collectionData.actions, null, 2)}</pre>
						</div>

						<div class="space-y-2">
							<span class="font-bold">Pre-variables:</span>
							<pre class="pre">{JSON.stringify(collectionData.pre, null, 2)}</pre>
						</div>
					</div>
				{/if}
			</div>
		</Tabs.Content>

		<!-- ACTIONS -->
		<Tabs.Content value="actions" class="space-y-4 pt-4">
			<div class="space-y-4 card p-6">
				<h3 class="mb-2 h3">Manage Actions</h3>
				<!-- Operation Selector -->
				<div class="mb-4 flex space-x-4">
					<label class="flex items-center space-x-2">
						<input class="radio" type="radio" value="create" bind:group={actionOp} />
						<span>Create</span>
					</label>
					<label class="flex items-center space-x-2">
						<input class="radio" type="radio" value="update" bind:group={actionOp} />
						<span>Update</span>
					</label>
					<label class="flex items-center space-x-2">
						<input class="radio" type="radio" value="delete" bind:group={actionOp} />
						<span>Delete</span>
					</label>
					<label class="flex items-center space-x-2">
						<input class="radio" type="radio" value="get" bind:group={actionOp} />
						<span>Get</span>
					</label>
				</div>

				<!-- Fields -->
				<label class="label">
					<span>Collection Path</span>
					<input
						class="input"
						type="text"
						bind:value={actionCollectionPath}
						placeholder="/path/to/collections/my-collection"
					/>
				</label>

				<label class="label">
					<span>Sub-Path (Optional)</span>
					<input
						class="input"
						type="text"
						bind:value={actionSubPath}
						placeholder="workflows/daily"
					/>
				</label>

				<label class="label">
					<span>File Name</span>
					<input
						class="input"
						type="text"
						bind:value={actionFileName}
						placeholder="my-action.toml"
					/>
				</label>

				{#if actionOp === 'create' || actionOp === 'update'}
					<label class="label">
						<span>Action JSON Data</span>
						<textarea
							class="textarea"
							bind:value={actionData}
							rows="6"
							placeholder="Enter valid JSON action definition..."
						></textarea>
					</label>
				{/if}

				<div class="flex justify-end pt-2">
					<button class="variant-filled-secondary btn" onclick={handleAction}>
						Execute {actionOp.toUpperCase()}
					</button>
				</div>

				{#if actionResult}
					<div class="variant-soft mt-4 space-y-2 card p-4">
						<span class="font-bold">Action Content:</span>
						<pre class="max-h-[300px] overflow-auto pre">{JSON.stringify(
								actionResult,
								null,
								2
							)}</pre>
					</div>
				{/if}
			</div>
		</Tabs.Content>
	</Tabs>

	<!-- Global Feedback -->
	{#if error}
		<aside class="alert variant-filled-error">
			<div class="alert-message">
				<h3 class="h3">Error</h3>
				<p>{error}</p>
			</div>
		</aside>
	{/if}
	{#if result}
		<aside class="alert variant-filled-success">
			<div class="alert-message">
				<h3 class="h3">Success</h3>
				<p>{result}</p>
			</div>
		</aside>
	{/if}
</div>
