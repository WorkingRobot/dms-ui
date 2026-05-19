<script>
	import { onMount } from 'svelte';
	import { api, enc } from '$lib/api.js';
	import { toast } from '$lib/toast.svelte.js';
	import { Card, PageHeader, Button, Modal, Field, PasswordField } from '$lib/ui';
	import Plus from '@lucide/svelte/icons/plus';
	import Trash2 from '@lucide/svelte/icons/trash-2';

	let users = $state([]);
	let loading = $state(true);
	let busy = $state(false);
	let addOpen = $state(false);
	let u = $state('');
	let p = $state('');

	async function load() {
		loading = true;
		try {
			users = (await api.get('/masters')) ?? [];
		} catch (e) {
			toast.error(e.message);
		}
		loading = false;
	}
	onMount(load);

	async function add() {
		busy = true;
		try {
			await api.post('/masters', { user: u, password: p });
			toast.success('Master user added');
			addOpen = false;
			u = p = '';
			await load();
		} catch (e) {
			toast.error(e.message);
		}
		busy = false;
	}

	async function del(user) {
		try {
			await api.del(`/masters/${enc(user)}`);
			toast.success('Removed');
			await load();
		} catch (e) {
			toast.error(e.message);
		}
	}
</script>

<PageHeader
	title="Dovecot Masters"
	description="SASL master users with full-mailbox access (migration / sync)"
	doc="https://docker-mailserver.github.io/docker-mailserver/latest/config/account-management/supplementary/master-accounts/"
>
	<Button onclick={() => (addOpen = true)}><Plus size={16} />Add master</Button>
</PageHeader>

<Card title="{users.length} master user(s)">
	{#if loading}
		<p class="text-sm text-muted-foreground">Loading…</p>
	{:else if users.length === 0}
		<p class="text-sm text-muted-foreground">No master users.</p>
	{:else}
		<ul class="divide-y divide-border rounded-md border border-border">
			{#each users as user}
				<li class="flex items-center justify-between px-4 py-3 text-sm">
					<span class="font-medium">{user}</span>
					<Button size="sm" variant="ghost" onclick={() => del(user)}>
						<Trash2 size={14} class="text-danger" />
					</Button>
				</li>
			{/each}
		</ul>
	{/if}
</Card>

<Modal bind:open={addOpen} title="Add Dovecot master user">
	<div class="space-y-4">
		<Field label="Username" bind:value={u} required />
		<PasswordField label="Password" bind:value={p} required />
	</div>
	{#snippet footer()}
		<Button variant="outline" onclick={() => (addOpen = false)}>Cancel</Button>
		<Button loading={busy} onclick={add}>Add</Button>
	{/snippet}
</Modal>
