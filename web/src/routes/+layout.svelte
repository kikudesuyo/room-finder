<script lang="ts">
	import './layout.css';

	let { children } = $props();
</script>

<svelte:head>
	<title>Room Finder</title>
	<meta name="description" content="賃貸物件を探すための個人用Agent" />
</svelte:head>

<main class="bg-surface text-on-surface min-h-screen">
	{@render children()}
</main>
