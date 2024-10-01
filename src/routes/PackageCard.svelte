<script lang="ts">
	export let name: string;
	export let elem;
	export let minify: boolean = false;

	const prettyName = (name: string): string => {
		if (name.startsWith('gno.land/')) {
			name = name.substring('gno.land/'.length);
		}

		if (minify) {
			const parts = name.split('/');
			for (let i = 0; i < parts.length; i++) {
				let part = parts[i];
				if (part.length > 17) {
					part = part.substring(0, 7) + '...' + part.substring(part.length - 7);
				}
				parts[i] = part;
			}
			name = parts.join('/');
		}

		return name;
	};
</script>

<div
	style="padding: 0; padding: 5px; padding-left: 9px; padding-right: 9px; box-sizing: border-box; overflow: hidden; border: 1px solid black; border-radius: 5px; background-color: {elem
		.tx.success
		? '#C8FFAE'
		: '#F46F76'}"
>
	<div>
		<p
			style="margin: 0; color: black; text-align: left; direction: rtl; overflow: hidden; text-overflow: ellipsis;  white-space: nowrap;"
		>
			{prettyName(name)}
		</p>
	</div>
	<div>
		<a style="color: black;" href="/#{name}" on:click={() => window.scrollTo(0, 0)}>📜 log</a>
		{#if elem.tx.success}
			<a style="color: black;" href="https://{name}" target="_blank">𖣐 gnoweb</a>
			{#if elem.tx.block_height > 0}
				<a
					style="color: black;"
					href="https://gnoscan.io/realms/details?chainId=portal-loop&path={name}"
					target="_blank">🩻 gnoscan</a
				>
			{/if}
		{/if}
	</div>
</div>
