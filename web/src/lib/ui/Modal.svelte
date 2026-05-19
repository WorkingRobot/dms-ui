<script>
	import X from '@lucide/svelte/icons/x';
	let { open = $bindable(false), title = '', children, footer } = $props();
</script>

{#if open}
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-2 backdrop-blur-sm sm:p-4"
		onclick={(e) => e.target === e.currentTarget && (open = false)}
		onkeydown={(e) => e.key === 'Escape' && (open = false)}
		role="presentation"
	>
		<div
			class="flex max-h-[90vh] w-full max-w-md flex-col overflow-hidden rounded-[var(--radius)] border border-border bg-card shadow-xl"
		>
			<header
				class="flex items-center justify-between border-b border-border px-5 py-4"
			>
				<h3 class="text-sm font-semibold">{title}</h3>
				<button
					class="rounded p-1 text-muted-foreground hover:bg-muted"
					onclick={() => (open = false)}
					aria-label="Close"
				>
					<X size={16} />
				</button>
			</header>
			<div class="flex-1 overflow-y-auto px-5 py-4">{@render children?.()}</div>
			{#if footer}
				<footer
					class="flex flex-wrap justify-end gap-2 border-t border-border px-5 py-3"
				>
					{@render footer()}
				</footer>
			{/if}
		</div>
	</div>
{/if}
