<script>
	let {
		variant = 'default',
		size = 'md',
		type = 'button',
		disabled = false,
		loading = false,
		onclick,
		children,
		...rest
	} = $props();

	const variants = {
		default: 'bg-accent text-accent-foreground hover:opacity-90',
		outline:
			'border border-border bg-transparent hover:bg-muted text-foreground',
		ghost: 'bg-transparent hover:bg-muted text-foreground',
		danger: 'bg-danger text-white hover:opacity-90'
	};
	const sizes = {
		sm: 'h-8 px-3 text-xs',
		md: 'h-9 px-4 text-sm',
		icon: 'h-9 w-9'
	};
</script>

<button
	{type}
	disabled={disabled || loading}
	{onclick}
	class="inline-flex items-center justify-center gap-2 rounded-[calc(var(--radius)-0.25rem)] font-medium transition focus:outline-none focus-visible:ring-2 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50 {variants[
		variant
	]} {sizes[size]}"
	{...rest}
>
	{#if loading}
		<span
			class="h-3.5 w-3.5 animate-spin rounded-full border-2 border-current border-t-transparent"
		></span>
	{/if}
	{@render children?.()}
</button>
