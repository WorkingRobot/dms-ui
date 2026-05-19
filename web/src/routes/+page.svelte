<script>
	import { onMount } from 'svelte';
	import { api } from '$lib/api.js';
	import { num } from '$lib/format.js';
	import { Card, PageHeader, Badge, Tooltip } from '$lib/ui';
	import { serviceContext } from '$lib/dmsmeta.js';
	import Users from '@lucide/svelte/icons/users';
	import Forward from '@lucide/svelte/icons/forward';
	import Inbox from '@lucide/svelte/icons/inbox';
	import ShieldCheck from '@lucide/svelte/icons/shield-check';
	import Info from '@lucide/svelte/icons/info';

	let d = $state(null);
	let envMap = $state({});
	let err = $state('');
	let loading = $state(true);

	async function refresh() {
		try {
			d = await api.get('/dashboard');
			err = '';
		} catch (e) {
			err = e.message;
		}
		loading = false;
	}

	// Env is fetched once so each service can be shown with context (e.g.
	// "disabled because ENABLE_CLAMAV is not set") instead of a bare "stopped".
	async function loadEnv() {
		try {
			const s = await api.get('/settings');
			envMap = Object.fromEntries((s.env ?? []).map((e) => [e.key, e.value]));
		} catch {}
	}

	const services = $derived(
		(d?.services ?? [])
			.map((svc) => ({ ...svc, ctx: serviceContext(svc, envMap) }))
			.sort((a, b) => {
				const rank = { danger: 0, success: 1, muted: 2 };
				return rank[a.ctx.tone] - rank[b.ctx.tone] || a.name.localeCompare(b.name);
			})
	);

	onMount(() => {
		refresh();
		loadEnv();
		const t = setInterval(refresh, 15000); // live-ish dashboard
		return () => clearInterval(t);
	});
	// Rspamd is optional; when it isn't running its stats are unavailable, so
	// the widgets that depend on it are gated rather than showing a bare 0.
	const rspamdLive = $derived(!!d?.rspamd);
	const rspamdTip =
		'Rspamd is not enabled (ENABLE_RSPAMD is off), so spam statistics are unavailable.';

	const stats = $derived(
		d
			? [
					{ label: 'Accounts', value: num(d.accountCount), icon: Users },
					{ label: 'Aliases', value: num(d.aliasCount), icon: Forward },
					{ label: 'Mail queue', value: num(d.queueLength), icon: Inbox },
					{
						label: 'Spam blocked',
						value: rspamdLive ? num(d.rspamd.spam) : 'n/a',
						icon: ShieldCheck,
						gated: !rspamdLive
					}
				]
			: []
	);
</script>

<PageHeader
	title="Dashboard"
	description="docker-mailserver {d?.version ?? ''} · {d?.container ?? ''}"
	doc="https://docker-mailserver.github.io/docker-mailserver/latest/"
/>

{#if loading}
	<p class="text-sm text-muted-foreground">Loading…</p>
{:else if err}
	<Card title="Error"><p class="text-sm text-danger">{err}</p></Card>
{:else}
	<div class="mb-6 grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
		{#each stats as s}
			{@const Icon = s.icon}
			<div
				class="rounded-[var(--radius)] border border-border bg-card p-5 shadow-sm {s.gated
					? 'opacity-60'
					: ''}"
			>
				<div class="flex items-center justify-between">
					<span class="text-xs text-muted-foreground">{s.label}</span>
					{#if s.gated}
						<Tooltip text={rspamdTip}>
							<Info size={16} class="text-muted-foreground" />
						</Tooltip>
					{:else}
						<Icon size={16} class="text-accent" />
					{/if}
				</div>
				<p class="mt-2 text-2xl font-semibold">{s.value}</p>
			</div>
		{/each}
	</div>

	<div class="grid gap-6 lg:grid-cols-2">
		<Card
			title="Services"
			subtitle="supervisord-managed processes"
			doc="https://docker-mailserver.github.io/docker-mailserver/latest/config/environment/"
		>
			<div class="grid grid-cols-1 gap-2">
				{#each services as svc (svc.name)}
					<div
						class="flex flex-col gap-1 rounded-md border border-border px-3 py-2 sm:flex-row sm:items-start sm:justify-between sm:gap-4"
					>
						<div class="min-w-0">
							<p class="text-xs font-medium">{svc.name}</p>
							<p class="mt-0.5 text-[11px] leading-snug text-muted-foreground">
								{svc.ctx.note}
							</p>
						</div>
						<div class="shrink-0">
							<Badge tone={svc.ctx.tone}>{svc.ctx.label}</Badge>
						</div>
					</div>
				{/each}
			</div>
			<p class="mt-3 text-[11px] leading-snug text-muted-foreground">
				<span class="text-warning">disabled</span> = intentionally off via
				its environment variable (expected).
				<span class="text-danger">stopped</span> = a service that should be
				running isn't; investigate.
				<span class="text-muted-foreground">unknown</span> = not a
				recognised docker-mailserver service.
			</p>
		</Card>

		<Card title="Rspamd" subtitle="message filtering statistics">
			{#if d.rspamd}
				<dl class="space-y-2 text-sm">
					{#each [['Scanned', d.rspamd.scanned], ['Ham', d.rspamd.ham], ['Spam', d.rspamd.spam], ['Learned', d.rspamd.learned], ['Avg scan (s)', d.rspamd.avgScanMs]] as [k, v]}
						<div class="flex justify-between">
							<dt class="text-muted-foreground">{k}</dt>
							<dd class="font-medium">{num(v)}</dd>
						</div>
					{/each}
				</dl>
			{:else}
				<div class="flex items-center gap-2 text-sm text-muted-foreground">
					<span>Rspamd not available</span>
					<Tooltip text={rspamdTip}>
						<Info size={15} class="text-muted-foreground" />
					</Tooltip>
				</div>
			{/if}
		</Card>
	</div>
{/if}
