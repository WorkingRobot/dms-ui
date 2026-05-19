// Minimal toast store (Svelte 5 runes); avoids pulling an extra dependency
// while matching the unobtrusive corner-toast style of Pocket ID.

let seq = 0;
export const toasts = $state([]);

function push(kind, message) {
	const id = ++seq;
	toasts.push({ id, kind, message });
	setTimeout(() => dismiss(id), 4500);
}

export function dismiss(id) {
	const i = toasts.findIndex((t) => t.id === id);
	if (i >= 0) toasts.splice(i, 1);
}

export const toast = {
	success: (m) => push('success', m),
	error: (m) => push('error', m),
	info: (m) => push('info', m)
};

// ok shows the `setup` command output returned by the API when it's
// meaningful, otherwise a sensible fallback. Keeps long multi-line output
// readable by collapsing whitespace.
export function ok(res, fallback) {
	const out = (res && res.output ? String(res.output) : '').trim();
	toast.success(out && out.length <= 240 ? out : fallback);
}
