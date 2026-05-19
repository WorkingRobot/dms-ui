<script>
	import '../app.css';
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import { api } from '$lib/api.js';
	import { initTheme, theme, toggleTheme } from '$lib/theme.svelte.js';
	import { Toaster } from '$lib/ui';
	import Gauge from '@lucide/svelte/icons/gauge';
	import Users from '@lucide/svelte/icons/users';
	import Forward from '@lucide/svelte/icons/forward';
	import KeyRound from '@lucide/svelte/icons/key-round';
	import ShieldCheck from '@lucide/svelte/icons/shield-check';
	import Send from '@lucide/svelte/icons/send';
	import Ban from '@lucide/svelte/icons/ban';
	import ScrollText from '@lucide/svelte/icons/scroll-text';
	import SlidersHorizontal from '@lucide/svelte/icons/sliders-horizontal';
	import Mail from '@lucide/svelte/icons/mail';
	import Sun from '@lucide/svelte/icons/sun';
	import Moon from '@lucide/svelte/icons/moon';
	import Menu from '@lucide/svelte/icons/menu';
	import X from '@lucide/svelte/icons/x';

	let { children } = $props();
	let me = $state({ email: '', container: '' });
	let sidebarOpen = $state(false);

	// Close the mobile drawer whenever the route changes.
	$effect(() => {
		$page.url.pathname;
		sidebarOpen = false;
	});

	const nav = [
		{ href: '/', label: 'Dashboard', icon: Gauge },
		{ href: '/accounts', label: 'Accounts', icon: Users },
		{ href: '/aliases', label: 'Aliases', icon: Forward },
		{ href: '/masters', label: 'Dovecot Masters', icon: KeyRound },
		{ href: '/dkim', label: 'DKIM', icon: ShieldCheck },
		{ href: '/relay', label: 'Relay', icon: Send },
		{ href: '/fail2ban', label: 'Fail2ban', icon: Ban },
		{ href: '/logs', label: 'Logs', icon: ScrollText },
		{ href: '/env', label: 'Environment', icon: SlidersHorizontal }
	];

	const active = (href) =>
		href === '/'
			? $page.url.pathname === '/'
			: $page.url.pathname.startsWith(href);

	onMount(async () => {
		initTheme();
		try {
			me = await api.get('/whoami');
		} catch {}
	});
</script>

<div class="flex min-h-screen">
	<!-- Mobile top bar -->
	<header
		class="fixed inset-x-0 top-0 z-30 flex h-14 items-center gap-3 border-b border-border bg-sidebar px-4 md:hidden"
	>
		<button
			class="rounded-md p-1.5 text-muted-foreground hover:bg-muted"
			onclick={() => (sidebarOpen = true)}
			aria-label="Open menu"
		>
			<Menu size={20} />
		</button>
		<div
			class="flex h-7 w-7 items-center justify-center rounded-lg bg-accent text-accent-foreground"
		>
			<Mail size={15} />
		</div>
		<span class="text-sm font-semibold">DMS Admin</span>
	</header>

	<!-- Backdrop (mobile, when drawer open) -->
	{#if sidebarOpen}
		<div
			class="fixed inset-0 z-40 bg-black/50 md:hidden"
			onclick={() => (sidebarOpen = false)}
			role="presentation"
		></div>
	{/if}

	<aside
		class="fixed inset-y-0 left-0 z-50 flex w-60 flex-col border-r border-border bg-sidebar transition-transform duration-200 md:static md:z-0 md:translate-x-0 {sidebarOpen
			? 'translate-x-0'
			: '-translate-x-full'}"
	>
		<div class="flex items-center gap-2.5 px-5 py-5">
			<button
				class="absolute right-3 top-4 rounded-md p-1.5 text-muted-foreground hover:bg-muted md:hidden"
				onclick={() => (sidebarOpen = false)}
				aria-label="Close menu"
			>
				<X size={18} />
			</button>
			<div
				class="flex h-8 w-8 items-center justify-center rounded-lg bg-accent text-accent-foreground"
			>
				<Mail size={17} />
			</div>
			<div>
				<p class="text-sm font-semibold leading-tight">DMS Admin</p>
				<p class="text-[11px] text-muted-foreground">docker-mailserver</p>
			</div>
		</div>

		<nav class="flex-1 space-y-0.5 px-3 py-2">
			{#each nav as item}
				{@const Icon = item.icon}
				<a
					href={item.href}
					class="flex items-center gap-3 rounded-[calc(var(--radius)-0.25rem)] px-3 py-2 text-sm transition {active(
						item.href
					)
						? 'bg-accent/12 font-medium text-accent'
						: 'text-muted-foreground hover:bg-muted hover:text-foreground'}"
				>
					<Icon size={17} />
					{item.label}
				</a>
			{/each}
		</nav>

		<div class="mt-auto border-t border-border px-3 py-3">
			<button
				onclick={toggleTheme}
				class="mb-2 flex w-full items-center gap-3 rounded-[calc(var(--radius)-0.25rem)] px-3 py-2 text-sm text-muted-foreground transition hover:bg-muted hover:text-foreground"
			>
				{#if theme.dark}<Sun size={17} />Light mode{:else}<Moon
						size={17}
					/>Dark mode{/if}
			</button>
			<div class="rounded-[calc(var(--radius)-0.25rem)] bg-muted px-3 py-2">
				<p class="truncate text-xs font-medium" title={me.email}>
					{me.email || 'unknown admin'}
				</p>
				<p class="truncate text-[11px] text-muted-foreground">
					{me.container}
				</p>
			</div>
		</div>
	</aside>

	<main class="mt-14 min-w-0 flex-1 px-4 py-6 md:mt-0 md:px-8 md:py-8">
		{@render children?.()}
	</main>
</div>

<Toaster />
