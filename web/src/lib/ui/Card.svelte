<script>
	import ExternalLink from '@lucide/svelte/icons/external-link';
	let { title = '', subtitle = '', doc = '', children, actions } = $props();
</script>

<section
	class="overflow-hidden rounded-[var(--radius)] border border-border bg-card shadow-sm"
>
	{#if title || actions}
		<header
			class="flex flex-wrap items-center justify-between gap-x-4 gap-y-2 border-b border-border px-5 py-4"
		>
			<div class="min-w-0">
				<h2 class="flex items-center gap-1.5 text-sm font-semibold tracking-tight">
					{title}
					{#if doc}
						<a
							href={doc}
							target="_blank"
							rel="noreferrer"
							class="text-muted-foreground transition hover:text-accent"
							title="docker-mailserver documentation"
						>
							<ExternalLink size={13} />
						</a>
					{/if}
				</h2>
				{#if subtitle}
					<p class="mt-0.5 text-xs text-muted-foreground">{subtitle}</p>
				{/if}
			</div>
			{#if actions}<div class="flex flex-wrap items-center gap-2">{@render actions()}</div>{/if}
		</header>
	{/if}
	<div class="min-w-0 p-5">
		{@render children?.()}
	</div>
</section>
