<script>
	// Tooltip that works on touch as well as pointer devices: hover/focus on
	// desktop, tap-to-toggle on mobile (tap elsewhere or Escape to close).
	let { text, children } = $props();

	let open = $state(false);
	let root;
	const id = `tt-${Math.random().toString(36).slice(2, 8)}`;

	function show() {
		open = true;
	}
	function hide() {
		open = false;
	}
	function toggle() {
		open = !open;
	}

	// On touch there is no hover; a tap toggles. Close when the next pointer
	// event lands outside the trigger, or on Escape.
	function onWindowPointerDown(e) {
		if (open && root && !root.contains(e.target)) open = false;
	}
	function onKeydown(e) {
		if (e.key === 'Escape') open = false;
	}
</script>

<svelte:window onpointerdown={onWindowPointerDown} onkeydown={onKeydown} />

<span class="relative inline-block" bind:this={root}>
	<button
		type="button"
		class="inline-flex cursor-help items-center"
		aria-describedby={open ? id : undefined}
		aria-expanded={open}
		onclick={toggle}
		onmouseenter={show}
		onmouseleave={hide}
		onfocus={show}
		onblur={hide}
	>
		{@render children?.()}
	</button>
	{#if open}
		<span
			{id}
			role="tooltip"
			class="absolute left-1/2 top-full z-50 mt-1.5 w-max max-w-[16rem] -translate-x-1/2 rounded-[calc(var(--radius)-0.25rem)] border border-border bg-card px-2.5 py-1.5 text-[11px] leading-snug text-foreground shadow-md"
		>
			{text}
		</span>
	{/if}
</span>
