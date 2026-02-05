<script lang="ts">
	/**
	 * Myna - The Postman for AWS serverless.
	 *
	 * Project Details:
	 * - Name: Myna
	 * - Description: Bridges the gap between API-first tools and modern serverless systems by providing
	 *   Git-first workflows, repeatable execution, and a unified interface for Lambda, S3, SQS, SNS, and EventBridge.
	 */

	function handleAddCollection() {
		console.log('Add collections handler triggered');
	}

	function handleOpenCollection() {
		console.log('Open collection handler triggered');
	}
</script>

<div
	class="mx-auto flex max-w-3xl flex-col items-center justify-center space-y-12 px-4 py-16 text-center"
>
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
