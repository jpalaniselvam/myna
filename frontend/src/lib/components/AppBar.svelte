<script lang="ts">
	import { MenuIcon } from '@lucide/svelte';
	import { AppBar } from '@skeletonlabs/skeleton-svelte';
	import { layoutStore } from '$lib/stores/layout.svelte';
</script>

<AppBar>
	<AppBar.Toolbar class="grid-cols-[auto_1fr_auto]">
		<AppBar.Lead>
			<button
				type="button"
				class="btn-icon btn-icon-lg hover:preset-tonal"
				onclick={() => layoutStore.toggleSidebar()}
			>
				<MenuIcon />
			</button>
		</AppBar.Lead>
		<AppBar.Headline>
			<p class="text-2xl">Myna</p>
		</AppBar.Headline>
	</AppBar.Toolbar>
</AppBar>
