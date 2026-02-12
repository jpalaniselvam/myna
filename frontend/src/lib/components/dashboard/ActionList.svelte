<script lang="ts">
	import { Folder, FileCode, Play, ChevronRight, ChevronDown } from '@lucide/svelte';
	import { TreeView, createTreeViewCollection } from '@skeletonlabs/skeleton-svelte';

	interface Node {
		id: string;
		name: string;
		path?: string;
		children?: Node[];
	}

	let { actions = {}, collectionPath } = $props<{
		actions: any;
		collectionPath: string;
	}>();

	function transformActions(actionsMap: any, parentId: string = 'root'): Node[] {
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
					return {
						id,
						name: key,
						children: transformActions(item, id)
					};
				} else {
					return {
						id,
						name: key.replace('.toml', ''),
						path: item
					};
				}
			});
	}

	let collection = $derived(
		createTreeViewCollection<Node>({
			nodeToValue: (node) => node.id,
			nodeToString: (node) => node.name,
			rootNode: {
				id: 'root',
				name: '',
				children: transformActions(actions)
			}
		})
	);
</script>

<div class="w-full">
	{#if Object.keys(actions).length === 0}
		<div class="p-4 text-surface-400 italic">No actions found in this collection.</div>
	{:else}
		<TreeView {collection} selectionMode="single">
			<TreeView.Tree>
				{#each collection.rootNode.children || [] as node, index (node.id)}
					{@render treeNode(node, [index])}
				{/each}
			</TreeView.Tree>
		</TreeView>
	{/if}
</div>

{#snippet treeNode(node: Node, indexPath: number[])}
	<TreeView.NodeProvider value={{ node, indexPath }}>
		{#if node.children}
			<TreeView.Branch>
				<TreeView.BranchControl>
					<TreeView.BranchIndicator>
						<ChevronRight size={14} class="block in-data-[state=open]:hidden" />
						<ChevronDown size={14} class="hidden in-data-[state=open]:block" />
					</TreeView.BranchIndicator>
					<TreeView.BranchText>
						<Folder size={16} class="text-secondary-500" />
						<span class="font-medium">{node.name}</span>
					</TreeView.BranchText>
				</TreeView.BranchControl>
				<TreeView.BranchContent>
					<TreeView.BranchIndentGuide />
					{#each node.children as childNode, childIndex (childNode.id)}
						{@render treeNode(childNode, [...indexPath, childIndex])}
					{/each}
				</TreeView.BranchContent>
			</TreeView.Branch>
		{:else}
			<TreeView.Item class="group flex items-center justify-between">
				<div class="flex items-center gap-2">
					<FileCode size={16} class="text-primary-500" />
					<span>{node.name}</span>
				</div>
				<div class="opacity-0 transition-opacity group-hover:opacity-100">
					<button
						class="variant-soft-primary btn-icon btn-icon-sm"
						title="Run Action"
						onclick={(e) => {
							e.stopPropagation();
							console.log('Run action:', node.path);
						}}
					>
						<Play size={12} />
					</button>
				</div>
			</TreeView.Item>
		{/if}
	</TreeView.NodeProvider>
{/snippet}
