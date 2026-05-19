<script>
	import { onMount } from 'svelte';
	import { api, enc } from '$lib/api.js';
	import { bytes, num } from '$lib/format.js';
	import { toast } from '$lib/toast.svelte.js';
	import {
		Card,
		PageHeader,
		Button,
		Modal,
		Field,
		Badge,
		PasswordField,
		Confirm
	} from '$lib/ui';
	import Plus from '@lucide/svelte/icons/plus';
	import Trash2 from '@lucide/svelte/icons/trash-2';
	import KeyRound from '@lucide/svelte/icons/key-round';
	import HardDrive from '@lucide/svelte/icons/hard-drive';
	import ChevronDown from '@lucide/svelte/icons/chevron-down';
	import Search from '@lucide/svelte/icons/search';

	let accounts = $state([]);
	let query = $state('');
	let loading = $state(true);
	const filtered = $derived(
		query
			? accounts.filter(
					(a) =>
						a.email.toLowerCase().includes(query.toLowerCase()) ||
						(a.aliases ?? []).some((al) =>
							al.toLowerCase().includes(query.toLowerCase())
						)
				)
			: accounts
	);
	let busy = $state(false);
	let expanded = $state(null);
	let detail = $state(null);

	let addOpen = $state(false);
	let newEmail = $state('');
	let newPass = $state('');

	let pwOpen = $state(false);
	let pwTarget = $state('');
	let pwValue = $state('');

	let quotaOpen = $state(false);
	let quotaTarget = $state('');
	let quotaValue = $state('');

	let delOpen = $state(false);
	let delTarget = $state('');
	let delPurge = $state(false);

	let restrictOpen = $state(false);
	let rTarget = $state('');
	let rSend = $state([]);
	let rRecv = $state([]);

	async function load() {
		loading = true;
		try {
			accounts = await api.get('/accounts');
		} catch (e) {
			toast.error(e.message);
		}
		loading = false;
	}
	onMount(load);

	async function toggleDetail(email) {
		if (expanded === email) {
			expanded = null;
			return;
		}
		expanded = email;
		detail = null;
		try {
			detail = await api.get(`/accounts/${enc(email)}`);
		} catch (e) {
			toast.error(e.message);
		}
	}

	async function addAccount() {
		busy = true;
		try {
			await api.post('/accounts', { email: newEmail, password: newPass });
			toast.success(`Created ${newEmail}`);
			addOpen = false;
			newEmail = newPass = '';
			await load();
		} catch (e) {
			toast.error(e.message);
		}
		busy = false;
	}

	async function changePw() {
		busy = true;
		try {
			await api.put(`/accounts/${enc(pwTarget)}/password`, {
				password: pwValue
			});
			toast.success(`Password updated for ${pwTarget}`);
			pwOpen = false;
			pwValue = '';
		} catch (e) {
			toast.error(e.message);
		}
		busy = false;
	}

	async function saveQuota() {
		busy = true;
		try {
			await api.put(`/accounts/${enc(quotaTarget)}/quota`, {
				quota: quotaValue
			});
			toast.success(`Quota set for ${quotaTarget}`);
			quotaOpen = false;
			await load();
		} catch (e) {
			toast.error(e.message);
		}
		busy = false;
	}

	async function doDelete() {
		busy = true;
		try {
			await api.del(
				`/accounts/${enc(delTarget)}${delPurge ? '?purge=true' : ''}`
			);
			toast.success(`Deleted ${delTarget}`);
			delOpen = false;
			await load();
		} catch (e) {
			toast.error(e.message);
		}
		busy = false;
	}

	async function openRestrict(email) {
		rTarget = email;
		restrictOpen = true;
		rSend = rRecv = [];
		try {
			rSend = (await api.get('/restrictions/send')) ?? [];
			rRecv = (await api.get('/restrictions/receive')) ?? [];
		} catch (e) {
			toast.error(e.message);
		}
	}
	async function setRestrict(dir, on) {
		try {
			if (on) await api.post(`/restrictions/${dir}`, { email: rTarget });
			else await api.del(`/restrictions/${dir}?email=${enc(rTarget)}`);
			toast.success(`Updated ${dir} restriction`);
			rSend = (await api.get('/restrictions/send')) ?? [];
			rRecv = (await api.get('/restrictions/receive')) ?? [];
		} catch (e) {
			toast.error(e.message);
		}
	}
</script>

<PageHeader
	title="Accounts"
	description="Mailboxes and their storage usage"
	doc="https://docker-mailserver.github.io/docker-mailserver/latest/config/account-management/overview/"
>
	<Button onclick={() => (addOpen = true)}>
		<Plus size={16} />Add account
	</Button>
</PageHeader>

<Card title="{accounts.length} mailbox(es)">
	{#snippet actions()}
		<div class="relative">
			<Search
				size={14}
				class="absolute left-2.5 top-1/2 -translate-y-1/2 text-muted-foreground"
			/>
			<input
				bind:value={query}
				placeholder="Search…"
				class="h-8 w-full rounded-[calc(var(--radius)-0.25rem)] border border-input bg-background pl-8 pr-3 text-xs sm:w-56"
			/>
		</div>
	{/snippet}
	{#if loading}
		<p class="text-sm text-muted-foreground">Loading…</p>
	{:else if filtered.length === 0}
		<p class="py-8 text-center text-sm text-muted-foreground">
			{accounts.length === 0
				? 'No mailboxes yet. Click “Add account” to create one.'
				: `No mailboxes match “${query}”.`}
		</p>
	{:else}
		<div class="overflow-x-auto rounded-md border border-border">
			<table class="w-full min-w-[640px] text-sm">
				<thead class="bg-muted/60 text-xs text-muted-foreground">
					<tr>
						<th class="px-4 py-2.5 text-left font-medium">Email</th>
						<th class="px-4 py-2.5 text-left font-medium">Usage</th>
						<th class="px-4 py-2.5 text-left font-medium">Quota</th>
						<th class="px-4 py-2.5 text-left font-medium">Aliases</th>
						<th class="px-4 py-2.5 text-right font-medium">Actions</th>
					</tr>
				</thead>
				<tbody>
					{#each filtered as a (a.email)}
						<tr class="border-t border-border align-middle">
							<td class="px-4 py-3">
								<button
									class="flex items-center gap-1.5 font-medium hover:text-accent"
									onclick={() => toggleDetail(a.email)}
								>
									<ChevronDown
										size={14}
										class="transition {expanded === a.email
											? 'rotate-180'
											: ''}"
									/>
									{a.email}
								</button>
							</td>
							<td class="px-4 py-3 text-muted-foreground">{a.used}</td>
							<td class="px-4 py-3">
								{#if a.quota === '~'}
									<Badge>unlimited</Badge>
								{:else}
									<Badge tone="accent">{a.quota}</Badge>
								{/if}
							</td>
							<td class="px-4 py-3 text-muted-foreground">
								{a.aliases?.length ?? 0}
							</td>
							<td class="px-4 py-3">
								<div class="flex justify-end gap-1.5">
									<Button
										size="sm"
										variant="ghost"
										onclick={() => {
											pwTarget = a.email;
											pwOpen = true;
										}}><KeyRound size={14} /></Button
									>
									<Button
										size="sm"
										variant="ghost"
										onclick={() => {
											quotaTarget = a.email;
											quotaValue = '';
											quotaOpen = true;
										}}><HardDrive size={14} /></Button
									>
									<Button
										size="sm"
										variant="ghost"
										onclick={() => openRestrict(a.email)}>Restrict</Button
									>
									<Button
										size="sm"
										variant="ghost"
										onclick={() => {
											delTarget = a.email;
											delPurge = false;
											delOpen = true;
										}}
									>
										<Trash2 size={14} class="text-danger" />
									</Button>
								</div>
							</td>
						</tr>
						{#if expanded === a.email}
							<tr class="border-t border-border bg-muted/30">
								<td colspan="5" class="px-6 py-4">
									{#if !detail}
										<p class="text-xs text-muted-foreground">Loading…</p>
									{:else}
										<div class="mb-3 flex flex-wrap gap-x-6 gap-y-1 text-xs">
											<span
												>Storage: <b>{bytes(detail.usedBytes)}</b></span
											>
											<span>Messages: <b>{num(detail.messages)}</b></span>
										</div>
										<div class="overflow-x-auto">
										<table class="w-full min-w-[420px] text-xs">
											<thead class="text-muted-foreground">
												<tr>
													<th class="py-1 text-left">Folder</th>
													<th class="py-1 text-right">Messages</th>
													<th class="py-1 text-right">Size</th>
												</tr>
											</thead>
											<tbody>
												{#each detail.folders ?? [] as f}
													<tr class="border-t border-border/60">
														<td class="py-1">{f.name}</td>
														<td class="py-1 text-right">{num(f.messages)}</td>
														<td class="py-1 text-right">{bytes(f.vsize)}</td>
													</tr>
												{/each}
											</tbody>
										</table>
										</div>
									{/if}
								</td>
							</tr>
						{/if}
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
</Card>

<Modal bind:open={addOpen} title="Add account">
	<div class="space-y-4">
		<Field
			label="Email"
			bind:value={newEmail}
			placeholder="user@domain.tld"
			required
		/>
		<PasswordField label="Password" bind:value={newPass} required />
	</div>
	{#snippet footer()}
		<Button variant="outline" onclick={() => (addOpen = false)}>Cancel</Button>
		<Button loading={busy} onclick={addAccount}>Create</Button>
	{/snippet}
</Modal>

<Modal bind:open={pwOpen} title="Change password">
	<p class="mb-3 text-xs text-muted-foreground">{pwTarget}</p>
	<PasswordField label="New password" bind:value={pwValue} required />
	{#snippet footer()}
		<Button variant="outline" onclick={() => (pwOpen = false)}>Cancel</Button>
		<Button loading={busy} onclick={changePw}>Update</Button>
	{/snippet}
</Modal>

<Modal bind:open={quotaOpen} title="Set quota">
	<p class="mb-3 text-xs text-muted-foreground">{quotaTarget}</p>
	<Field
		label="Quota"
		bind:value={quotaValue}
		placeholder="5G, 512M, or 0 for unlimited"
		hint="Requires ENABLE_QUOTAS=1 in the mail container."
		required
	/>
	{#snippet footer()}
		<Button variant="outline" onclick={() => (quotaOpen = false)}>Cancel</Button
		>
		<Button loading={busy} onclick={saveQuota}>Save</Button>
	{/snippet}
</Modal>

<Confirm
	bind:open={delOpen}
	title="Delete account"
	phrase={delTarget}
	message="Deleting {delTarget} is permanent and cannot be undone."
	confirmLabel="Delete account"
	loading={busy}
	onconfirm={doDelete}
>
	<label class="flex items-center gap-2 text-sm">
		<input type="checkbox" bind:checked={delPurge} />
		Also delete the mailbox data on disk
	</label>
</Confirm>

<Modal bind:open={restrictOpen} title="Send / receive restrictions">
	<p class="mb-4 text-xs text-muted-foreground">{rTarget}</p>
	<div class="space-y-3">
		<label class="flex items-center justify-between text-sm">
			Block sending
			<input
				type="checkbox"
				checked={rSend.includes(rTarget)}
				onchange={(e) => setRestrict('send', e.currentTarget.checked)}
			/>
		</label>
		<label class="flex items-center justify-between text-sm">
			Block receiving
			<input
				type="checkbox"
				checked={rRecv.includes(rTarget)}
				onchange={(e) => setRestrict('receive', e.currentTarget.checked)}
			/>
		</label>
	</div>
	{#snippet footer()}
		<Button variant="outline" onclick={() => (restrictOpen = false)}
			>Done</Button
		>
	{/snippet}
</Modal>
