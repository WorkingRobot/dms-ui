<script>
	let {
		label = '',
		value = $bindable(''),
		type = 'text',
		placeholder = '',
		hint = '',
		required = false
	} = $props();
</script>

<label class="block">
	<span class="mb-1.5 block text-xs font-medium text-muted-foreground"
		>{label}{#if required}<span class="text-danger"> *</span>{/if}</span
	>
	{#if type === 'password'}
		<input
			type="password"
			bind:value
			{placeholder}
			class="h-9 w-full rounded-[calc(var(--radius)-0.25rem)] border border-input bg-background px-3 text-sm focus:outline-none focus-visible:ring-2 focus-visible:ring-ring"
		/>
	{:else}
		<input
			{type}
			bind:value
			{placeholder}
			class="h-9 w-full rounded-[calc(var(--radius)-0.25rem)] border border-input bg-background px-3 text-sm focus:outline-none focus-visible:ring-2 focus-visible:ring-ring"
		/>
	{/if}
	{#if hint}<span class="mt-1 block text-xs text-muted-foreground">{hint}</span
		>{/if}
</label>
