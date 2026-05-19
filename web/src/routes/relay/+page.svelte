<script>
	import { onMount } from 'svelte';
	import { api } from '$lib/api.js';
	import { toast } from '$lib/toast.svelte.js';
	import { Card, PageHeader, Button, Modal, Field } from '$lib/ui';

	let info = $state(null);
	let loading = $state(true);
	let busy = $state(false);

	let domOpen = $state(false);
	let dDomain = $state('');
	let dHost = $state('');
	let dPort = $state('');

	let authOpen = $state(false);
	let aDomain = $state('');
	let aUser = $state('');
	let aPass = $state('');

	let exclDomain = $state('');

	async function load() {
		loading = true;
		try {
			info = await api.get('/relay');
		} catch (e) {
			toast.error(e.message);
		}
		loading = false;
	}
	onMount(load);

	async function addDomain() {
		busy = true;
		try {
			await api.post('/relay/domain', {
				Domain: dDomain,
				Host: dHost,
				Port: dPort
			});
			toast.success('Relay domain added');
			domOpen = false;
			await load();
		} catch (e) {
			toast.error(e.message);
		}
		busy = false;
	}

	async function addAuth() {
		busy = true;
		try {
			await api.post('/relay/auth', {
				Domain: aDomain,
				User: aUser,
				Password: aPass
			});
			toast.success('Relay credentials saved');
			authOpen = false;
			aPass = '';
			await load();
		} catch (e) {
			toast.error(e.message);
		}
		busy = false;
	}

	async function exclude() {
		try {
			await api.post('/relay/exclude', { Domain: exclDomain });
			toast.success(`${exclDomain} excluded from relay`);
			exclDomain = '';
			await load();
		} catch (e) {
			toast.error(e.message);
		}
	}
</script>

<PageHeader
	title="Relay"
	description="Outbound relay host configuration"
	doc="https://docker-mailserver.github.io/docker-mailserver/latest/config/advanced/mail-forwarding/relay-hosts/"
>
	<Button variant="outline" onclick={() => (authOpen = true)}>Add credentials</Button>
	<Button onclick={() => (domOpen = true)}>Add relay domain</Button>
</PageHeader>

{#if loading}
	<p class="text-sm text-muted-foreground">Loading…</p>
{:else}
	<div class="grid gap-6 lg:grid-cols-2">
		<Card title="Relay map" subtitle="postfix-relaymap.cf">
			<pre
				class="max-h-64 overflow-auto whitespace-pre-wrap break-words rounded bg-muted px-3 py-2 text-xs">{info?.relayMap?.trim() || 'empty'}</pre>
			<div class="mt-4 flex gap-2">
				<input
					bind:value={exclDomain}
					placeholder="domain to exclude"
					class="h-9 flex-1 rounded-[calc(var(--radius)-0.25rem)] border border-input bg-background px-3 text-sm"
				/>
				<Button variant="outline" onclick={exclude}>Exclude</Button>
			</div>
		</Card>
		<Card title="SASL credentials" subtitle="postfix-sasl-password.cf (masked)">
			<pre
				class="max-h-64 overflow-auto whitespace-pre-wrap break-words rounded bg-muted px-3 py-2 text-xs">{info?.saslPassword?.trim() || 'empty'}</pre>
		</Card>
	</div>
{/if}

<Modal bind:open={domOpen} title="Add relay domain">
	<div class="space-y-4">
		<Field label="Domain" bind:value={dDomain} required />
		<Field label="Relay host" bind:value={dHost} required />
		<Field label="Port" bind:value={dPort} placeholder="587" />
	</div>
	{#snippet footer()}
		<Button variant="outline" onclick={() => (domOpen = false)}>Cancel</Button>
		<Button loading={busy} onclick={addDomain}>Add</Button>
	{/snippet}
</Modal>

<Modal bind:open={authOpen} title="Add relay credentials">
	<div class="space-y-4">
		<Field label="Domain" bind:value={aDomain} required />
		<Field label="Username" bind:value={aUser} required />
		<Field label="Password" type="password" bind:value={aPass} required />
	</div>
	{#snippet footer()}
		<Button variant="outline" onclick={() => (authOpen = false)}>Cancel</Button>
		<Button loading={busy} onclick={addAuth}>Save</Button>
	{/snippet}
</Modal>
