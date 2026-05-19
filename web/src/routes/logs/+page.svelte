<script>
	import { api } from '$lib/api.js';
	import { toast } from '$lib/toast.svelte.js';
	import { Card, PageHeader, Button } from '$lib/ui';
	import RefreshCw from '@lucide/svelte/icons/refresh-cw';

	let tab = $state('mail');
	let lines = $state(500);
	let content = $state('');
	let queue = $state([]);
	let loading = $state(false);

	async function load() {
		loading = true;
		try {
			if (tab === 'mail') {
				content = (await api.get(`/logs?lines=${lines}`)).log;
			} else if (tab === 'fail2ban') {
				content = (await api.get('/fail2ban/log')).log;
			} else {
				queue = (await api.get('/queue')) ?? [];
			}
		} catch (e) {
			toast.error(e.message);
		}
		loading = false;
	}
	// Re-loads on mount and whenever the tab changes.
	$effect(() => {
		tab;
		load();
	});
</script>

<PageHeader
	title="Logs"
	description="Mail log, fail2ban log and the Postfix queue"
	doc="https://docker-mailserver.github.io/docker-mailserver/latest/config/debugging/"
>
	<Button variant="outline" loading={loading} onclick={load}>
		<RefreshCw size={14} />Refresh
	</Button>
</PageHeader>

<div class="mb-4 flex gap-1 rounded-[var(--radius)] border border-border bg-card p-1 text-xs sm:text-sm">
	{#each [['mail', 'Mail log'], ['fail2ban', 'Fail2ban log'], ['queue', 'Mail queue']] as [id, label]}
		<button
			class="flex-1 rounded-[calc(var(--radius)-0.25rem)] px-2 py-1.5 transition sm:px-3 {tab ===
			id
				? 'bg-accent/12 font-medium text-accent'
				: 'text-muted-foreground hover:bg-muted'}"
			onclick={() => (tab = id)}>{label}</button
		>
	{/each}
</div>

{#if tab === 'queue'}
	<Card title="{queue.length} queued message(s)">
		{#if queue.length === 0}
			<p class="text-sm text-muted-foreground">Mail queue is empty.</p>
		{:else}
			<div class="overflow-x-auto">
			<table class="w-full min-w-[560px] text-xs">
				<thead class="text-muted-foreground">
					<tr>
						<th class="py-1 text-left">ID</th>
						<th class="py-1 text-left">Size</th>
						<th class="py-1 text-left">Sender</th>
						<th class="py-1 text-left">Recipient</th>
						<th class="py-1 text-left">Reason</th>
					</tr>
				</thead>
				<tbody>
					{#each queue as q}
						<tr class="border-t border-border/60">
							<td class="py-1 font-mono">{q.id}</td>
							<td class="py-1">{q.size}</td>
							<td class="py-1">{q.sender}</td>
							<td class="py-1">{q.recipient}</td>
							<td class="py-1 text-danger">{q.reason}</td>
						</tr>
					{/each}
				</tbody>
			</table>
			</div>
		{/if}
	</Card>
{:else}
	<Card title={tab === 'mail' ? 'mail.log' : 'fail2ban.log'}>
		{#snippet actions()}
			{#if tab === 'mail'}
				<select
					bind:value={lines}
					onchange={load}
					class="h-8 rounded border border-input bg-background px-2 text-xs"
				>
					<option value={200}>200 lines</option>
					<option value={500}>500 lines</option>
					<option value={1000}>1000 lines</option>
					<option value={2000}>2000 lines</option>
				</select>
			{/if}
		{/snippet}
		<pre
			class="max-h-[65vh] overflow-auto rounded bg-muted px-3 py-2 text-[11px] leading-relaxed">{content ||
				'(empty)'}</pre>
	</Card>
{/if}
