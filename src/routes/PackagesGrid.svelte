<script lang="ts">
	import { request, gql } from 'graphql-request';
	import { createQuery } from '@tanstack/svelte-query';
	import PackageCard from './PackageCard.svelte';
	import { page } from '$app/stores';

	const endpoint = 'https://indexer.portal-loop.gno.testnet.teritori.com/graphql/query';

	const query = createQuery({
		queryKey: ['packages'],
		staleTime: Infinity,
		queryFn: async () => {
			const res = await request(
				endpoint,
				gql`
					query {
						transactions(filter: { message: { type_url: add_package }, to_block_height: 1 }) {
							index
							hash
							success
							block_height
							gas_wanted
							gas_used
							memo
							messages {
								typeUrl
								route
								value {
									... on MsgAddPackage {
										creator
										package {
											path
										}
										deposit
									}
								}
							}
							response {
								data
								info
								log
							}
						}
					}
				`
			);

			// TODO: extract and type
			const txs = (res as any).transactions;
			const packages: { [key: string]: { tx: any; msg: any } } = {};
			for (const tx of txs) {
				let i = 0;
				for (const msg of tx.messages) {
					const path = msg.value.package.path;
					if (
						!path ||
						(packages[path] && packages[path].tx.response.log.includes(`msg:${i},success:true,`))
					) {
						i++;
						continue;
					}
					packages[path] = {
						tx,
						msg
					};
					i++;
				}
			}

			return { map: packages, list: Object.entries(packages).reverse() };
		}
	});

	let search: string = '';
</script>

<div style="display: flex; justify-content: center; width: 100%;">
	{#if $query.isLoading}
		<p>Loading...</p>
	{:else if $query.isError}
		<p>Error: {$query.error.message}</p>
	{:else if $query.isSuccess}
		{#if $page.url.hash === ''}
			<div style="width: 100%; display: flex; flex-direction: column; gap: 5px;">
				<div style="display: flex; gap: 5px">
					<a href="/#" style="display: flex; align-items: center;">
						<img src="/favicon.png" alt="logo" style="height: 58.5px;" />
					</a>
					<div
						style="font-size: 1.2rem; color: black; background-color: white; padding-left: 20px; padding-right: 20px; flex-grow: 1; border-radius: 5px; display: flex; align-items: center; border: 1px solid black;"
					>
						🔍
						<!-- svelte-ignore a11y-autofocus -->
						<input
							autofocus
							style="all: unset; margin-left: 10px; flex-grow: 1;"
							name="search"
							bind:value={search}
						/>
					</div>
				</div>
				<div
					style="display: grid; grid-template-columns: repeat(auto-fit, minmax(320px, 1fr)); gap: 5px; width: 100%"
				>
					{#each $query.data.list.filter(([name]) => !search || name
								.toLowerCase()
								.includes(search.toLowerCase())) as [name, elem]}
						<PackageCard minify {name} {elem} />
					{/each}
				</div>
			</div>
		{:else}
			<div style="width: 100%; display: flex; flex-direction: column; gap: 5px;">
				<div style="display: flex; gap: 5px;">
					<a href="/#" style="display: flex; align-items: center;">
						<img src="/favicon.png" alt="logo" style="height: 58.5px;" />
					</a>
					<div style="flex-grow: 1;">
						<PackageCard
							name={$page.url.hash && $page.url.hash.length > 1 ? $page.url.hash.substring(1) : ''}
							elem={$query.data.map[
								$page.url.hash && $page.url.hash.length > 1 ? $page.url.hash.substring(1) : ''
							]}
						/>
					</div>
				</div>
				<p
					style="margin: 0; padding: 0; padding: 5px; padding-left: 9px; padding-right: 9px; white-space: pre-wrap; text-align: left; color: white; background-color: black; border-radius: 5px"
				>
					{$query.data.map[
						$page.url.hash && $page.url.hash.length > 1 ? $page.url.hash.substring(1) : ''
					].tx.response.log}
				</p>
			</div>
		{/if}
	{/if}
</div>
