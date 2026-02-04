<script lang="ts">
	import { slide } from 'svelte/transition';

	let { message, ondismiss } = $props<{ message: string; ondismiss: () => void }>();
</script>

{#if message}
	<aside class="alert variant-filled-error m-2" transition:slide>
		<div class="alert-message flex w-full items-center justify-between">
			<p class="text-sm">{message}</p>
			<button class="variant-filled-error btn-icon btn-icon-sm" onclick={ondismiss}>X</button>
		</div>
	</aside>
{/if}
