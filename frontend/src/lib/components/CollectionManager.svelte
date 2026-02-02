<script lang="ts">
    import { CreateCollection, UpdateCollection, DeleteCollection, GetCollection } from '../../../wailsjs/go/main/App';
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

    async function handleCreate() {
        try {
            error = '';
            result = '';
            if (!createBaseDir || !createName) {
                error = "Base Directory and Name are required.";
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
            await UpdateCollection(updatePath, updateName, updateDesc, pre);
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
</script>

<div class="w-full max-w-4xl mx-auto space-y-8 p-4">
    <h2 class="h2 font-bold">Collection Manager</h2>
    
    <Tabs value={activeTab} onValueChange={(e) => activeTab = e.value}>
        <Tabs.List class="grid grid-cols-4 w-full gap-1">
            <Tabs.Trigger value="create">Create</Tabs.Trigger>
            <Tabs.Trigger value="update">Update</Tabs.Trigger>
            <Tabs.Trigger value="delete">Delete</Tabs.Trigger>
            <Tabs.Trigger value="get">Get</Tabs.Trigger>
        </Tabs.List>

        <!-- CREATE -->
        <Tabs.Content value="create" class="space-y-4 pt-4">
            <div class="card p-6 space-y-4">
                <label class="label">
                    <span>Base Directory</span>
                    <input class="input" type="text" bind:value={createBaseDir} placeholder="/path/to/collections" />
                </label>
                <label class="label">
                    <span>Name</span>
                    <input class="input" type="text" bind:value={createName} placeholder="my-collection" />
                </label>
                <label class="label">
                    <span>Description</span>
                    <input class="input" type="text" bind:value={createDesc} placeholder="My awesome collection" />
                </label>
                <div class="flex justify-end">
                    <button class="btn variant-filled-primary" onclick={handleCreate}>
                        Create Collection
                    </button>
                </div>
            </div>
        </Tabs.Content>

        <!-- UPDATE -->
        <Tabs.Content value="update" class="space-y-4 pt-4">
            <div class="card p-6 space-y-4">
                <label class="label">
                    <span>Collection Path</span>
                    <input class="input" type="text" bind:value={updatePath} placeholder="/path/to/collections/my-collection" />
                </label>
                <label class="label">
                    <span>New Name (Optional)</span>
                    <input class="input" type="text" bind:value={updateName} placeholder="Leave empty to keep current name" />
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
                    <button class="btn variant-filled-warning" onclick={handleUpdate}>
                        Update Collection
                    </button>
                </div>
            </div>
        </Tabs.Content>

        <!-- DELETE -->
        <Tabs.Content value="delete" class="space-y-4 pt-4">
            <div class="card p-6 space-y-4">
                <label class="label">
                    <span>Collection Path</span>
                    <input class="input" type="text" bind:value={deletePath} placeholder="/path/to/collections/my-collection" />
                </label>
                <div class="flex justify-end">
                    <button class="btn variant-filled-error" onclick={handleDelete}>
                        Delete Collection
                    </button>
                </div>
            </div>
        </Tabs.Content>

        <!-- GET -->
        <Tabs.Content value="get" class="space-y-4 pt-4">
            <div class="card p-6 space-y-4">
                <label class="label">
                    <span>Collection Path</span>
                    <input class="input" type="text" bind:value={getPath} placeholder="/path/to/collections/my-collection" />
                </label>
                <div class="flex justify-end">
                    <button class="btn variant-filled-success" onclick={handleGet}>
                        Get Collection Details
                    </button>
                </div>

                {#if collectionData}
                    <div class="card p-4 variant-soft space-y-4">
                        <h3 class="h3">Details</h3>
                        
                        <div class="space-y-2">
                            <span class="font-bold">Environments:</span>
                            <ul class="list-disc list-inside">
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
