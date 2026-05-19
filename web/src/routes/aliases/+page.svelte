<script>
	import { onMount } from 'svelte';
	import { api, enc } from '$lib/api.js';
	import { toast } from '$lib/toast.svelte.js';
	import { Card, PageHeader, Button, Modal, Field } from '$lib/ui';
	import Plus from '@lucide/svelte/icons/plus';
	import Trash2 from '@lucide/svelte/icons/trash-2';

	let aliases = $state([]);
	let loading = $state(true);
	let busy = $state(false);
	let addOpen = $state(false);
	let src = $state('');
	let tgt = $state('');

	async function load() {
		loading = true;
		try {
			aliases = (await api.get('/aliases')) ?? [];
		} catch (e) {
			toast.error(e.message);
		}
		loading = false;
	}
	onMount(load);

	async function add() {
		busy = true;
		try {
			await api.post('/aliases', { source: src, target: tgt });
			toast.success('Alias added');
			addOpen = false;
			src = tgt = '';
			await load();
		} catch (e) {
			toast.error(e.message);
		}
		busy = false;
	}

	async function del(a) {
		try {
			await api.del(
				`/aliases?source=${enc(a.source)}&target=${enc(a.target)}`
			);
			toast.success('Alias removed');
			await load();
		} catch (e) {
			toast.error(e.message);
		}
	}
</script>

<PageHeader
	title="Aliases"
	description="Virtual address forwarding (postfix-virtual.cf)"
	doc="https://docker-mailserver.github.io/docker-mailserver/latest/config/account-management/overview/#aliases"
>
	<Button onclick={() => (addOpen = true)}><Plus size={16} />Add alias</Button>
</PageHeader>

<Card title="{aliases.length} alias(es)">
	{#if loading}
		<p class="text-sm text-muted-foreground">Loading…</p>
	{:else if aliases.length === 0}
		<p class="text-sm text-muted-foreground">No aliases configured.</p>
	{:else}
		<div class="overflow-x-auto rounded-md border border-border">
			<table class="w-full min-w-[480px] text-sm">
				<thead class="bg-muted/60 text-xs text-muted-foreground">
					<tr>
						<th class="px-4 py-2.5 text-left font-medium">Source</th>
						<th class="px-4 py-2.5 text-left font-medium">Target</th>
						<th class="px-4 py-2.5"></th>
					</tr>
				</thead>
				<tbody>
					{#each aliases as a}
						<tr class="border-t border-border">
							<td class="px-4 py-3 font-medium">{a.source}</td>
							<td class="px-4 py-3 text-muted-foreground">→ {a.target}</td>
							<td class="px-4 py-3 text-right">
								<Button size="sm" variant="ghost" onclick={() => del(a)}>
									<Trash2 size={14} class="text-danger" />
								</Button>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
</Card>

<Modal bind:open={addOpen} title="Add alias">
	<div class="space-y-4">
		<Field label="Source address" bind:value={src} placeholder="alias@domain.tld" required />
		<Field label="Target" bind:value={tgt} placeholder="real@domain.tld" required />
	</div>
	{#snippet footer()}
		<Button variant="outline" onclick={() => (addOpen = false)}>Cancel</Button>
		<Button loading={busy} onclick={add}>Add</Button>
	{/snippet}
</Modal>
