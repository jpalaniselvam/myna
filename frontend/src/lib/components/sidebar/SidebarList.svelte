<script lang="ts">
	import { TreeView, createTreeViewCollection } from '@skeletonlabs/skeleton-svelte';
	import { Folder, FileCode, Database, ChevronRight, ChevronDown } from '@lucide/svelte';
	import { workspaceStore } from '$lib/stores/workspace.svelte';
	import type { collection } from '../../../../wailsjs/go/models';

	interface Node {
		id: string;
		name: string;
		kind: 'root' | 'collection' | 'folder' | 'action';
		collectionPath?: string;
		subPath?: string;
		fileName?: string;
		children?: Node[];
	}

	let {
		collections,
		activeCollection,
		collectionDetails = {},
		onselect,
		onselectAction
	} = $props<{
		collections: string[];
		activeCollection: string | null;
		collectionDetails?: Record<string, collection.CollectionResponse>;
		onselect: (path: string) => void;
		onselectAction: (collectionPath: string, subPath: string, fileName: string) => void;
	}>();

	function transformActions(
		actionsMap: any,
		collectionPath: string,
		parentId: string,
		subPath: string = ''
	): Node[] {
		if (!actionsMap) return [];
		return Object.keys(actionsMap)
			.sort((a, b) => {
				const aIsDir = typeof actionsMap[a] === 'object';
				const bIsDir = typeof actionsMap[b] === 'object';
				if (aIsDir && !bIsDir) return -1;
				if (!aIsDir && bIsDir) return 1;
				return a.localeCompare(b);
			})
			.map((key) => {
				const item = actionsMap[key];
				const id = `${parentId}/${key}`;
				if (typeof item === 'object' && item !== null) {
					const currentSubPath = subPath ? `${subPath}/${key}` : key;
					return {
						id,
						name: key,
						kind: 'folder',
						children: transformActions(item, collectionPath, id, currentSubPath)
					};
				} else {
					return {
						id,
						name: key.replace('.toml', ''),
						kind: 'action',
						collectionPath,
						subPath,
						fileName: key
					};
				}
			});
	}

	let treeData = $derived(
		collections.map((path: string) => {
			const normalized = path.replace(/\\/g, '/');
			const name = path.split(/[\\/]/).pop() || 'Collection';
			const details = collectionDetails[normalized];
			return {
				id: normalized,
				name,
				kind: 'collection',
				collectionPath: path,
				children: details ? transformActions(details.actions, path, normalized) : []
			};
		})
	);

	let treeCollection = $derived(
		createTreeViewCollection<Node>({
			nodeToValue: (node) => node.id,
			nodeToString: (node) => node.name,
			rootNode: {
				id: 'root',
				name: '',
				kind: 'root',
				children: treeData
			}
		})
	);

	function handleNodeClick(node: Node) {
		if (node.kind === 'collection') {
			onselect(node.collectionPath!);
		} else if (node.kind === 'action') {
			onselectAction(node.collectionPath!, node.subPath!, node.fileName!);
		}
	}

	let activeTabId = $derived(workspaceStore.activeTabId);
</script>

<div class="flex-1 overflow-y-auto p-2">
	<TreeView collection={treeCollection} selectionMode="single">
		<TreeView.Tree>
			{#each treeCollection.rootNode.children || [] as node, index (node.id)}
				{@render treeNode(node, [index])}
			{/each}
		</TreeView.Tree>
	</TreeView>
</div>

{#snippet treeNode(node: Node, indexPath: number[])}
	<TreeView.NodeProvider value={{ node, indexPath }}>
		{#if node.children && node.children.length > 0}
			<TreeView.Branch>
				<TreeView.BranchControl
					onclick={() => handleNodeClick(node)}
					class="rounded-md px-2 py-1 hover:bg-surface-500/20 {activeTabId === node.id
						? 'preset-filled-primary-500'
						: ''}"
				>
					<TreeView.BranchIndicator>
						<ChevronRight size={14} class="block in-data-[state=open]:hidden" />
						<ChevronDown size={14} class="hidden in-data-[state=open]:block" />
					</TreeView.BranchIndicator>
					<TreeView.BranchText class="flex items-center gap-2">
						{#if node.kind === 'collection'}
							<Database size={16} class="text-primary-500" />
						{:else}
							<Folder size={16} class="text-secondary-500" />
						{/if}
						<span class={node.kind === 'collection' ? 'font-bold' : 'font-medium'}>{node.name}</span
						>
					</TreeView.BranchText>
				</TreeView.BranchControl>
				<TreeView.BranchContent>
					<TreeView.BranchIndentGuide />
					{#each node.children as childNode, childIndex (childNode.id)}
						{@render treeNode(childNode, [...indexPath, childIndex])}
					{/each}
				</TreeView.BranchContent>
			</TreeView.Branch>
		{:else if node.kind === 'folder'}
			<TreeView.Branch>
				<TreeView.BranchControl
					onclick={() => handleNodeClick(node)}
					class="rounded-md px-2 py-1 hover:bg-surface-500/20"
				>
					<TreeView.BranchIndicator>
						<ChevronRight size={14} class="block in-data-[state=open]:hidden" />
						<ChevronDown size={14} class="hidden in-data-[state=open]:block" />
					</TreeView.BranchIndicator>
					<TreeView.BranchText class="flex items-center gap-2">
						<Folder size={16} class="text-secondary-500" />
						<span class="font-medium">{node.name}</span>
					</TreeView.BranchText>
				</TreeView.BranchControl>
				<TreeView.BranchContent>
					<TreeView.BranchIndentGuide />
					<div class="p-2 text-xs text-surface-400 italic">Empty</div>
				</TreeView.BranchContent>
			</TreeView.Branch>
		{:else}
			<TreeView.Item
				onclick={() => handleNodeClick(node)}
				class="ml-6 flex items-center gap-2 rounded-md px-2 py-1 hover:bg-surface-500/20 {activeTabId ===
				node.id
					? 'preset-filled-primary-500'
					: ''}"
			>
				{#if node.kind === 'collection'}
					<Database size={16} class="text-primary-500" />
				{:else}
					<FileCode size={16} class="text-tertiary-500" />
				{/if}
				<span>{node.name}</span>
			</TreeView.Item>
		{/if}
	</TreeView.NodeProvider>
{/snippet}

<style>
</style>
