<script>
	import PackagesGrid from './PackagesGrid.svelte';
</script>

<svelte:head>
	<title>Loop Doctor</title>
	<meta name="description" content="Svelte demo app" />
</svelte:head>

<section>
	<PackagesGrid />
</section>

<style>
	section {
		display: flex;
		flex-direction: column;
		justify-content: flex-start;
		align-items: center;
		flex: 0.6;
	}
</style>
