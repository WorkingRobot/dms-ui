<script>
	import { generatePassword, strength } from '$lib/pw.js';
	import { toast } from '$lib/toast.svelte.js';
	import Wand from '@lucide/svelte/icons/wand-sparkles';
	import Eye from '@lucide/svelte/icons/eye';
	import EyeOff from '@lucide/svelte/icons/eye-off';
	import Copy from '@lucide/svelte/icons/copy';

	let { label = 'Password', value = $bindable(''), required = false } = $props();
	let show = $state(false);
	const s = $derived(strength(value));
	const bars = ['bg-danger', 'bg-danger', 'bg-warning', 'bg-success', 'bg-success'];

	function gen() {
		value = generatePassword(20);
		show = true;
	}
	function copy() {
		if (!value) return;
		navigator.clipboard?.writeText(value);
		toast.info('Password copied');
	}
</script>

<div>
	<div class="mb-1.5 flex items-center justify-between">
		<span class="text-xs font-medium text-muted-foreground"
			>{label}{#if required}<span class="text-danger"> *</span>{/if}</span
		>
		<div class="flex items-center gap-2 text-xs">
			<button
				type="button"
				class="flex items-center gap-1 text-accent hover:underline"
				onclick={gen}><Wand size={12} />Generate</button
			>
			{#if value}
				<button
					type="button"
					class="flex items-center gap-1 text-muted-foreground hover:text-foreground"
					onclick={copy}><Copy size={12} />Copy</button
				>
			{/if}
		</div>
	</div>
	<div class="relative">
		{#if show}
			<input
				type="text"
				bind:value
				class="h-9 w-full rounded-[calc(var(--radius)-0.25rem)] border border-input bg-background px-3 pr-9 text-sm focus:outline-none focus-visible:ring-2 focus-visible:ring-ring"
			/>
		{:else}
			<input
				type="password"
				bind:value
				class="h-9 w-full rounded-[calc(var(--radius)-0.25rem)] border border-input bg-background px-3 pr-9 text-sm focus:outline-none focus-visible:ring-2 focus-visible:ring-ring"
			/>
		{/if}
		<button
			type="button"
			class="absolute right-2 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
			onclick={() => (show = !show)}
			aria-label={show ? 'Hide' : 'Show'}
		>
			{#if show}<EyeOff size={15} />{:else}<Eye size={15} />{/if}
		</button>
	</div>
	{#if value}
		<div class="mt-1.5 flex gap-1">
			{#each [0, 1, 2, 3] as i}
				<span
					class="h-1 flex-1 rounded-full {i < s
						? bars[s]
						: 'bg-border'}"
				></span>
			{/each}
		</div>
	{/if}
</div>
