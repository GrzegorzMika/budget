<script lang="ts">
	import '../app.css';
	import { onMount } from 'svelte';

	let { data, children } = $props();
	let { keycloak } = data;

	onMount(() => {
		if (!keycloak) {
			window.location.reload();
		}
	});
</script>

{@render children()}
