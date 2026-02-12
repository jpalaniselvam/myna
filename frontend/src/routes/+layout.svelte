<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import Sidebar from '$lib/components/sidebar/Sidebar.svelte';
	import AppBar from '$lib/components/AppBar.svelte';
	import { layoutStore } from '$lib/stores/layout.svelte';
	import { slide } from 'svelte/transition';
	import { cubicInOut } from 'svelte/easing';

	let { children } = $props();
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>

<div class="flex h-screen w-screen flex-col overflow-hidden">
	<AppBar />
	<div class="flex flex-1 overflow-auto">
		{#if layoutStore.isSidebarOpen}
			<div transition:slide={{ axis: 'x', duration: 400, easing: cubicInOut }}>
				<Sidebar />
			</div>
		{/if}
		<main class="bg-surface-100-800-token flex-1 overflow-auto">
			{@render children()}
		</main>
	</div>
</div>
