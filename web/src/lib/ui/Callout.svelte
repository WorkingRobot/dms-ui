<script>
	import Info from '@lucide/svelte/icons/info';
	import TriangleAlert from '@lucide/svelte/icons/triangle-alert';
	// tone: 'info' (default) | 'warning'
	let { tone = 'info', children } = $props();
	const Icon = $derived(tone === 'warning' ? TriangleAlert : Info);
</script>

<div
	class="flex items-start gap-3 rounded-[var(--radius)] border px-4 py-3 text-sm {tone ===
	'warning'
		? 'border-warning/40 bg-warning/10'
		: 'border-accent/30 bg-accent/[0.07]'}"
>
	<Icon
		size={16}
		class="mt-0.5 shrink-0 {tone === 'warning' ? 'text-warning' : 'text-accent'}"
	/>
	<div class="min-w-0 leading-relaxed [&_a]:text-accent [&_a]:underline">
		{@render children?.()}
	</div>
</div>
