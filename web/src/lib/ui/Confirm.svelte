<script>
	import Modal from './Modal.svelte';
	import Button from './Button.svelte';
	import AlertTriangle from '@lucide/svelte/icons/triangle-alert';

	// Destructive confirmation that requires typing the exact `phrase`
	// (e.g. the account address) before the action button enables.
	let {
		open = $bindable(false),
		title = 'Confirm',
		phrase = '',
		message = '',
		confirmLabel = 'Delete',
		loading = false,
		onconfirm,
		children
	} = $props();

	let typed = $state('');
	$effect(() => {
		if (open) typed = '';
	});
</script>

<Modal bind:open {title}>
	<div class="space-y-3 text-sm">
		<div class="flex items-start gap-2 text-danger">
			<AlertTriangle size={16} class="mt-0.5 shrink-0" />
			<p>{message}</p>
		</div>
		{@render children?.()}
		<p class="text-xs text-muted-foreground">
			Type <code class="rounded bg-muted px-1 py-0.5 font-mono">{phrase}</code>
			to confirm:
		</p>
		<input
			bind:value={typed}
			autocomplete="off"
			class="h-9 w-full rounded-[calc(var(--radius)-0.25rem)] border border-input bg-background px-3 text-sm focus:outline-none focus-visible:ring-2 focus-visible:ring-ring"
		/>
	</div>
	{#snippet footer()}
		<Button variant="outline" onclick={() => (open = false)}>Cancel</Button>
		<Button
			variant="danger"
			{loading}
			disabled={typed !== phrase}
			onclick={onconfirm}>{confirmLabel}</Button
		>
	{/snippet}
</Modal>
