<script>
	import { onMount } from 'svelte';
	import { api } from '$lib/api.js';
	import { toast } from '$lib/toast.svelte.js';
	import { Card, PageHeader, Button, Modal, Field, Badge } from '$lib/ui';
	import Plus from '@lucide/svelte/icons/plus';
	import Copy from '@lucide/svelte/icons/copy';

	let keys = $state([]);
	let loading = $state(true);
	let busy = $state(false);
	let genOpen = $state(false);
	let domain = $state('');
	let selector = $state('mail');
	let keyType = $state('rsa');
	let keySize = $state('2048');
	let force = $state(false);
	let result = $state('');

	async function load() {
		loading = true;
		try {
			keys = (await api.get('/dkim')) ?? [];
		} catch (e) {
			toast.error(e.message);
		}
		loading = false;
	}
	onMount(load);

	async function generate() {
		busy = true;
		result = '';
		try {
			const r = await api.post('/dkim', {
				domain,
				selector,
				keyType,
				keySize,
				force
			});
			result = r.output ?? '';
			toast.success(`DKIM key generated for ${domain}`);
			await load();
		} catch (e) {
			toast.error(e.message);
		}
		busy = false;
	}

	function copy(txt) {
		navigator.clipboard?.writeText(txt);
		toast.info('Copied to clipboard');
	}

	// Suggested companion records. SPF/DMARC are conservative defaults an
	// admin can tighten; derived purely from the domain so the DKIM page
	// shows the full publish-this-trio in one place.
	const spfFor = (d) => `v=spf1 mx ~all`;
	const dmarcFor = (d) =>
		`v=DMARC1; p=quarantine; rua=mailto:postmaster@${d}; adkim=s; aspf=s`;

	// Chunk a long TXT value into 255-char quoted strings (BIND zonefile rule).
	function txtChunks(v) {
		const parts = [];
		for (let i = 0; i < v.length; i += 255)
			parts.push(`"${v.slice(i, i + 255)}"`);
		return parts.join(' ');
	}

	function zonefile(k) {
		return [
			`; --- ${k.domain} mail authentication ---`,
			`${k.domain}.\t\t\tIN\tTXT\t"${spfFor(k.domain)}"`,
			`_dmarc.${k.domain}.\t\tIN\tTXT\t"${dmarcFor(k.domain)}"`,
			`${k.recordKey}.\tIN\tTXT\t${txtChunks(k.txtValue)}`,
			''
		].join('\n');
	}
</script>

<PageHeader
	title="DKIM"
	description="rspamd DKIM signing keys and DNS records"
	doc="https://docker-mailserver.github.io/docker-mailserver/latest/config/best-practices/dkim_dmarc_spf/"
>
	<Button onclick={() => { genOpen = true; result = ''; }}>
		<Plus size={16} />Generate key
	</Button>
</PageHeader>

{#if loading}
	<p class="text-sm text-muted-foreground">Loading…</p>
{:else if keys.length === 0}
	<Card title="No keys">
		<p class="text-sm text-muted-foreground">
			No DKIM keys found. Generate one to start signing outbound mail.
		</p>
	</Card>
{:else}
	<div class="grid gap-4">
		{#each keys as k}
			<Card title={k.domain} subtitle="selector: {k.selector} · {k.keyType}">
				{#snippet actions()}
					<button
						class="flex items-center gap-1 text-xs text-accent hover:underline"
						onclick={() => copy(zonefile(k))}
					>
						<Copy size={12} />Copy zonefile
					</button>
					<Badge tone="accent">{k.keyType}</Badge>
				{/snippet}
				<div class="space-y-3 text-sm">
					<div>
						<p class="mb-1 text-xs text-muted-foreground">Record name</p>
						<code
							class="block break-all rounded bg-muted px-3 py-2 text-xs"
							>{k.recordKey}</code
						>
					</div>
					<div>
						<div class="mb-1 flex items-center justify-between">
							<p class="text-xs text-muted-foreground">TXT value</p>
							<button
								class="flex items-center gap-1 text-xs text-accent hover:underline"
								onclick={() => copy(k.txtValue)}
							>
								<Copy size={12} />Copy
							</button>
						</div>
						<code
							class="block max-h-32 overflow-auto whitespace-pre-wrap break-all rounded bg-muted px-3 py-2 text-xs"
							>{k.txtValue}</code
						>
					</div>

					<div class="grid gap-3 border-t border-border pt-3 sm:grid-cols-2">
						<div>
							<div class="mb-1 flex items-center justify-between">
								<p class="text-xs text-muted-foreground">
									SPF · <span class="font-mono">{k.domain}</span>
								</p>
								<button
									class="flex items-center gap-1 text-xs text-accent hover:underline"
									onclick={() => copy(spfFor(k.domain))}
								>
									<Copy size={12} />Copy
								</button>
							</div>
							<code
								class="block break-all rounded bg-muted px-3 py-2 text-xs"
								>{spfFor(k.domain)}</code
							>
						</div>
						<div>
							<div class="mb-1 flex items-center justify-between">
								<p class="text-xs text-muted-foreground">
									DMARC · <span class="font-mono">_dmarc.{k.domain}</span>
								</p>
								<button
									class="flex items-center gap-1 text-xs text-accent hover:underline"
									onclick={() => copy(dmarcFor(k.domain))}
								>
									<Copy size={12} />Copy
								</button>
							</div>
							<code
								class="block break-all rounded bg-muted px-3 py-2 text-xs"
								>{dmarcFor(k.domain)}</code
							>
						</div>
					</div>
					<p class="text-[11px] text-muted-foreground">
						Suggested defaults; tighten <code>p=</code> / SPF once mail flow
						is verified. “Copy zonefile” yields all three records as a
						paste-ready BIND block.
					</p>
				</div>
			</Card>
		{/each}
	</div>
{/if}

<Modal bind:open={genOpen} title="Generate DKIM key">
	<div class="space-y-4">
		<Field label="Domain" bind:value={domain} placeholder="example.com" required />
		<Field label="Selector" bind:value={selector} />
		<label class="block">
			<span class="mb-1.5 block text-xs font-medium text-muted-foreground"
				>Key type</span
			>
			<select
				bind:value={keyType}
				class="h-9 w-full rounded-[calc(var(--radius)-0.25rem)] border border-input bg-background px-3 text-sm"
			>
				<option value="rsa">RSA</option>
				<option value="ed25519">ED25519</option>
			</select>
		</label>
		{#if keyType === 'rsa'}
			<label class="block">
				<span class="mb-1.5 block text-xs font-medium text-muted-foreground"
					>Key size</span
				>
				<select
					bind:value={keySize}
					class="h-9 w-full rounded-[calc(var(--radius)-0.25rem)] border border-input bg-background px-3 text-sm"
				>
					<option>1024</option>
					<option>2048</option>
					<option>4096</option>
				</select>
			</label>
		{/if}
		<label class="flex items-center gap-2 text-sm">
			<input type="checkbox" bind:checked={force} />
			Force / rotate existing key
		</label>
		{#if result}
			<pre
				class="max-h-48 overflow-auto rounded bg-muted px-3 py-2 text-[11px]">{result}</pre>
		{/if}
	</div>
	{#snippet footer()}
		<Button variant="outline" onclick={() => (genOpen = false)}>Close</Button>
		<Button loading={busy} onclick={generate}>Generate</Button>
	{/snippet}
</Modal>
