<script>
	import ExternalLink from '@lucide/svelte/icons/external-link';
	let { title = '', description = '', doc = '', children } = $props();
</script>

<div
	class="mb-6 flex flex-col gap-4 md:flex-row md:items-end md:justify-between"
>
	<div class="min-w-0">
		<h1 class="text-xl font-semibold tracking-tight">{title}</h1>
		{#if description}
			<p class="mt-1 text-sm text-muted-foreground">{description}</p>
		{/if}
		{#if doc}
			<a
				href={doc}
				target="_blank"
				rel="noreferrer"
				class="mt-1.5 inline-flex items-center gap-1 text-xs text-muted-foreground transition hover:text-accent"
			>
				<ExternalLink size={12} />docker-mailserver docs
			</a>
		{/if}
	</div>
	{#if children}<div class="flex flex-wrap items-center gap-2">{@render children()}</div>{/if}
</div>
