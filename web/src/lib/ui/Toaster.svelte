<script>
	import { toasts, dismiss } from '$lib/toast.svelte.js';
	import Check from '@lucide/svelte/icons/circle-check';
	import Alert from '@lucide/svelte/icons/circle-alert';
	import Info from '@lucide/svelte/icons/info';
	const icon = { success: Check, error: Alert, info: Info };
	const tone = {
		success: 'border-success/40 text-success',
		error: 'border-danger/40 text-danger',
		info: 'border-accent/40 text-accent'
	};
</script>

<div class="fixed bottom-4 right-4 z-[60] flex w-80 flex-col gap-2">
	{#each toasts as t (t.id)}
		{@const Icon = icon[t.kind]}
		<div
			class="flex items-start gap-2 rounded-[var(--radius)] border bg-card px-4 py-3 text-sm shadow-lg {tone[
				t.kind
			]}"
		>
			<Icon size={16} class="mt-0.5 shrink-0" />
			<span class="flex-1 text-foreground">{t.message}</span>
			<button
				class="text-muted-foreground hover:text-foreground"
				onclick={() => dismiss(t.id)}>✕</button
			>
		</div>
	{/each}
</div>
