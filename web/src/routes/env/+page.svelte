<script>
	import { onMount } from 'svelte';
	import { api } from '$lib/api.js';
	import { toast } from '$lib/toast.svelte.js';
	import { PageHeader, Badge, Callout } from '$lib/ui';
	import { ENV_META, ENV_GROUPS } from '$lib/dmsmeta.js';

	let data = $state(null);
	let loading = $state(true);
	let q = $state('');

	onMount(async () => {
		try {
			data = await api.get('/settings');
		} catch (e) {
			toast.error(e.message);
		}
		loading = false;
	});

	const ql = $derived(q.trim().toLowerCase());

	const groups = $derived.by(() => {
		// Union of every known DMS variable with the live container env. Known
		// vars that aren't set still show (blank); they're part of the config
		// surface even when unset. Unknown live vars fall into "Other".
		const live = new Map((data?.env ?? []).map((e) => [e.key, e.value]));
		const keys = new Set([...Object.keys(ENV_META), ...live.keys()]);
		const env = [...keys]
			.map((key) => {
				const m = ENV_META[key];
				return {
					key,
					value: live.get(key) ?? '',
					group: m?.group ?? 'Other',
					desc: m?.desc ?? '',
					def: m?.default ?? ''
				};
			})
			.filter(
				(e) =>
					!ql ||
					e.key.toLowerCase().includes(ql) ||
					e.value.toLowerCase().includes(ql) ||
					e.desc.toLowerCase().includes(ql)
			);
		return ENV_GROUPS.map((name) => ({
			name,
			items: env
				.filter((e) => e.group === name)
				.sort((a, b) => a.key.localeCompare(b.key))
		})).filter((g) => g.items.length);
	});

	const total = $derived(groups.reduce((n, g) => n + g.items.length, 0));
</script>

<PageHeader
	title="Environment"
	description="Effective docker-mailserver environment {data
		? `· ${data.version} · ${data.container}`
		: ''}"
	doc="https://docker-mailserver.github.io/docker-mailserver/latest/config/environment/"
>
	<input
		bind:value={q}
		placeholder="Filter…"
		class="h-8 w-full rounded-[calc(var(--radius)-0.25rem)] border border-input bg-background px-3 text-xs sm:w-56"
	/>
	{#if !loading}<Badge>{total} vars</Badge>{/if}
</PageHeader>

<div class="mb-6">
	<Callout tone="warning">
		These are read by docker-mailserver at container start, so changing them
		means editing <code class="text-xs">mail.env</code> /
		<code class="text-xs">compose.yaml</code> on the host and
		<b>restarting the mail container</b>; not something this panel does, to
		avoid disrupting live mail. Variables under <b>Other</b> are ones this
		panel has no description for (custom or newly added); see the
		<a
			href="https://docker-mailserver.github.io/docker-mailserver/latest/config/environment/"
			target="_blank"
			rel="noreferrer">environment variable reference</a
		>.
	</Callout>
</div>

{#if loading}
	<p class="text-sm text-muted-foreground">Loading…</p>
{:else if total === 0}
	<p class="py-10 text-center text-sm text-muted-foreground">
		{data ? `No variables match “${q}”.` : 'Could not read environment.'}
	</p>
{:else}
	<div class="space-y-7">
		{#each groups as g (g.name)}
			<section>
				<div class="mb-2 flex items-center gap-2">
					<h2 class="text-xs font-semibold uppercase tracking-wide text-muted-foreground">
						{g.name}
					</h2>
					<span class="text-[11px] text-muted-foreground">·</span>
					<span class="text-[11px] text-muted-foreground">{g.items.length}</span>
				</div>
				<div
					class="overflow-hidden rounded-[var(--radius)] border border-border bg-card"
				>
					{#each g.items as e (e.key)}
						<div
							class="flex flex-col gap-1 border-b border-border px-4 py-3 last:border-b-0 sm:flex-row sm:items-start sm:gap-4"
						>
							<div class="min-w-0 sm:w-64 sm:shrink-0">
								<p class="font-mono text-xs font-medium">{e.key}</p>
								{#if e.desc}
									<p class="mt-0.5 text-[11px] leading-snug text-muted-foreground">
										{e.desc}
									</p>
								{/if}
							</div>
							<div class="min-w-0 flex-1">
								{#if e.value === ''}
									<p class="font-mono text-xs italic text-muted-foreground">
										(blank)
									</p>
								{:else}
									<p class="break-all font-mono text-xs">{e.value}</p>
								{/if}
								{#if e.def && String(e.def) !== e.value}
									<p class="mt-0.5 text-[11px] text-muted-foreground">
										default: <span class="font-mono">{e.def}</span>
									</p>
								{/if}
							</div>
						</div>
					{/each}
				</div>
			</section>
		{/each}
	</div>
{/if}
