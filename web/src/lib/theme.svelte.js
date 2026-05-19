// Light/dark theme, persisted to localStorage and reflected on <html>.

export const theme = $state({ dark: false });

export function initTheme() {
	if (typeof document === 'undefined') return;
	theme.dark = document.documentElement.classList.contains('dark');
}

export function toggleTheme() {
	theme.dark = !theme.dark;
	document.documentElement.classList.toggle('dark', theme.dark);
	try {
		localStorage.setItem('dms-theme', theme.dark ? 'dark' : 'light');
	} catch {}
}
