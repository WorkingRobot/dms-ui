<script>
	import { onMount } from 'svelte';
	import { api } from '$lib/api.js';
	import { toast, ok } from '$lib/toast.svelte.js';
	import { Card, PageHeader, Button, Badge } from '$lib/ui';

	let jails = $state([]);
	let loading = $state(true);
	let banIP = $state('');

	async function load() {
		loading = true;
		try {
			jails = (await api.get('/fail2ban')) ?? [];
		} catch (e) {
			toast.error(e.message);
		}
		loading = false;
	}
	onMount(load);

	async function ban() {
		if (!banIP) return;
		try {
			ok(await api.post('/fail2ban/ban', { ip: banIP }), `Banned ${banIP}`);
			banIP = '';
			await load();
		} catch (e) {
			toast.error(e.message);
		}
	}
	async function unban(ip) {
		try {
			ok(await api.post('/fail2ban/unban', { ip }), `Unbanned ${ip}`);
			await load();
		} catch (e) {
			toast.error(e.message);
		}
	}
</script>

<PageHeader
	title="Fail2ban"
	description="Intrusion banning across jails"
	doc="https://docker-mailserver.github.io/docker-mailserver/latest/config/security/fail2ban/"
>
	<input
		bind:value={banIP}
		placeholder="1.2.3.4"
		class="h-9 w-full sm:w-40 rounded-[calc(var(--radius)-0.25rem)] border border-input bg-background px-3 text-sm"
	/>
	<Button variant="danger" onclick={ban}>Ban IP</Button>
</PageHeader>

{#if loading}
	<p class="text-sm text-muted-foreground">Loading…</p>
{:else}
	<div class="grid gap-4">
		{#each jails as j}
			<Card title="Jail: {j.name}">
				{#snippet actions()}
					<Badge tone={j.currentlyBanned > 0 ? 'danger' : 'success'}>
						{j.currentlyBanned} banned
					</Badge>
				{/snippet}
				<div class="mb-3 flex gap-6 text-xs text-muted-foreground">
					<span>Currently failed: {j.currentlyFailed}</span>
					<span>Total failed: {j.totalFailed}</span>
					<span>Total banned: {j.totalBanned}</span>
				</div>
				{#if j.bannedIPs?.length}
					<div class="flex flex-wrap gap-1.5">
						{#each j.bannedIPs as ip}
							<button
								class="group inline-flex items-center gap-1 rounded-full bg-muted px-2 py-0.5 text-xs hover:bg-danger/15"
								onclick={() => unban(ip)}
								title="Click to unban"
							>
								{ip}
								<span class="text-muted-foreground group-hover:text-danger"
									>✕</span
								>
							</button>
						{/each}
					</div>
				{:else}
					<p class="text-xs text-muted-foreground">No banned IPs.</p>
				{/if}
			</Card>
		{/each}
	</div>
{/if}
