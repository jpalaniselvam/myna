<script lang="ts">
	import { collectionStore } from '$lib/stores/collections.svelte';
	import { workspaceStore } from '$lib/stores/workspace.svelte';
	import { CreateCollection, SelectDirectory, GetCollection } from '../../../wailsjs/go/main/App';
	import ErrorBanner from './sidebar/ErrorBanner.svelte';
	import CreateCollectionDialog from '$lib/components/dialogs/CreateCollectionDialog.svelte';

	let showCreateDialog = $state(false);
	let error = $state('');

	function handleAddCollection() {
		showCreateDialog = true;
	}

	async function handleOpenCollection() {
		try {
			error = '';
			const path = await SelectDirectory();
			if (!path) return;

			// Verify it's a valid collection
			await GetCollection(path);

			const name = path.split(/[\\/]/).pop() || 'Collection';
			collectionStore.add(path);
			workspaceStore.openTab('collection', path, name, { path });
			collectionStore.setActive(path);
		} catch (e: any) {
			error = 'This directory does not appear to be a valid Myna collection.';
			console.error(e);
		}
	}

	async function onCreateSubmit(data: { name: string; baseDir: string }) {
		try {
			error = '';
			await CreateCollection(data.baseDir, data.name, '');

			// Success
			const separator = data.baseDir.includes('\\') ? '\\' : '/';
			const cleanBase = data.baseDir.endsWith(separator) ? data.baseDir.slice(0, -1) : data.baseDir;
			const fullPath = `${cleanBase}${separator}${data.name}`;

			collectionStore.add(fullPath);
			workspaceStore.openTab('collection', fullPath, data.name, { path: fullPath });
			collectionStore.setActive(fullPath);
			showCreateDialog = false;
		} catch (e: any) {
			error = e.toString() || 'Failed to create collection';
			console.error(error);
		}
	}
</script>

<CreateCollectionDialog bind:open={showCreateDialog} onsubmit={onCreateSubmit} />

<div
	class="mx-auto flex max-w-3xl flex-col items-center justify-center space-y-12 px-4 py-16 text-center"
>
	{#if error}
		<div class="fixed top-4 right-4 z-50 w-full max-w-xs transition-all">
			<ErrorBanner message={error} ondismiss={() => (error = '')} />
		</div>
	{/if}
	<!-- Hero Section -->
	<header class="space-y-6">
		<h1 class="h1 text-5xl font-bold text-primary-500 md:text-6xl">Myna</h1>
		<p class="text-surface-600-300-token mx-auto max-w-2xl text-xl md:text-2xl">
			The Postman for AWS serverless. Bridge the gap between API-first tools and modern serverless
			systems.
		</p>
	</header>

	<!-- Collections Action Section -->
	<section
		class="bg-surface-100-800-token w-full max-w-2xl space-y-8 card rounded-xl border border-surface-500/10 p-10 shadow-2xl"
	>
		<div class="space-y-3">
			<h2 class="h2 font-bold">Collections</h2>
			<p class="text-surface-600-300-token">
				Manage your AWS requests in Git-friendly TOML collections. Store and replay events across
				different environments.
			</p>
		</div>

		<div class="flex flex-col justify-center gap-4 pt-4 sm:flex-row">
			<button
				class="btn preset-filled-primary-500 btn-lg shadow-lg transition-all hover:brightness-110"
				onclick={handleAddCollection}
			>
				Add collections
			</button>
			<button
				class="preset-outline-surface-500 btn btn-lg transition-all hover:preset-tonal"
				onclick={handleOpenCollection}
			>
				Open collection
			</button>
		</div>
	</section>

	<!-- Features Grid -->
	<div class="grid w-full grid-cols-1 gap-6 text-left md:grid-cols-2">
		<article
			class="card border border-surface-500/10 p-6 transition-colors hover:border-primary-500/30"
		>
			<h3 class="mb-3 h3 font-bold text-primary-500">Unified Interface</h3>
			<p class="text-surface-600-300-token leading-relaxed">
				Use a consistent TOML-based configuration for Lambda, S3, SQS, SNS, and EventBridge. No more
				navigating complex AWS Console UIs.
			</p>
		</article>
		<article
			class="card border border-surface-500/10 p-6 transition-colors hover:border-primary-500/30"
		>
			<h3 class="mb-3 h3 font-bold text-primary-500">Git-First Workflow</h3>
			<p class="text-surface-600-300-token leading-relaxed">
				Keep your test cases in version control next to your code. Design for repeatable execution
			</p>
		</article>
	</div>

	<!-- Footer Links -->
	<footer class="mt-8 flex w-full justify-center gap-8 border-t border-surface-500/10 pt-8">
		<a
			href="https://github.com/jpalaniselvam/myna"
			target="_blank"
			rel="noreferrer"
			class="text-surface-600-300-token flex items-center gap-2 anchor transition-colors hover:text-primary-500"
		>
			<span>GitHub</span>
		</a>
		<a
			href="https://jpalaniselvam.github.io/myna"
			target="_blank"
			rel="noreferrer"
			class="text-surface-600-300-token flex items-center gap-2 anchor transition-colors hover:text-primary-500"
		>
			<span>Documentation</span>
		</a>
	</footer>
</div>
