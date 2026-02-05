<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import Sidebar from '$lib/components/sidebar/Sidebar.svelte';
	import AppBar from '$lib/components/AppBar.svelte';

	let { children } = $props();
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>

<div class="flex h-screen w-screen flex-col overflow-hidden">
	<AppBar />
	<div class="flex flex-1 overflow-hidden">
		<Sidebar />
		<main class="bg-surface-100-800-token flex-1 overflow-auto">
			{@render children()}
		</main>
	</div>
</div>
