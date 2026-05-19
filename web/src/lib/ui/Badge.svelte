<script>
	let { tone = 'muted', children } = $props();
	const tones = {
		muted: 'bg-muted text-muted-foreground',
		success: 'bg-success/15 text-success',
		danger: 'bg-danger/15 text-danger',
		warning: 'bg-warning/15 text-warning',
		accent: 'bg-accent/15 text-accent'
	};
</script>

<span
	class="inline-flex items-center rounded-full px-2 py-0.5 text-[11px] font-medium {tones[
		tone
	]}"
>
	{@render children?.()}
</span>
