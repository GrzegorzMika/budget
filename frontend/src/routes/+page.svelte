<script lang="ts">
	import { onMount } from 'svelte';
	import { fetchCategories } from '$lib';
	import { Circle3 } from 'svelte-loading-spinners';

	let { data } = $props();
	let { keycloak } = data;
	let token = $state('');

	let categories = $state([]);

	async function fetch() {
		categories = await fetchCategories(token);
	}

	async function getToken() {
		token = (await keycloak)?.token || '';
	}
	onMount(async () => {
		await getToken();
		await fetch();
	});
</script>

{#if !token}
	<div class="flex h-screen w-screen items-center justify-center">
		<Circle3 size="300" unit="px" duration="3s" />
	</div>
{:else}
	<div class="container">
		<h1 class="header-primary">Expense Entry</h1>
		<div class="expanses-container">
			<form method="POST">
				<div class="expenses-entry">
					<label for="timestamp">Date:</label>
					<input type="date" id="timestamp" name="timestamp" required />
				</div>
				<div class="expenses-entry">
					<label for="amount">Amount:</label>
					<input type="number" id="amount" name="amount" step="0.01" required />
				</div>
				<div class="expenses-entry">
					<label for="category">Category:</label>
					<select id="category" name="category" required>
						<option value="" disabled selected>Select Category</option>
						{#each categories as category}
							<option value={category}>{category}</option>
						{/each}
					</select>
				</div>
				<input type="hidden" name="jwt" value={token} />
				<button class="submit-button" type="submit">Submit Expense</button>
			</form>
		</div>
	</div>
{/if}
