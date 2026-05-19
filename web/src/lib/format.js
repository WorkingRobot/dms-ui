export function bytes(n) {
	if (n == null || n < 0) return 'n/a';
	const u = ['B', 'KiB', 'MiB', 'GiB', 'TiB'];
	let i = 0;
	while (n >= 1024 && i < u.length - 1) {
		n /= 1024;
		i++;
	}
	return `${n.toFixed(i ? 1 : 0)} ${u[i]}`;
}

export function num(n) {
	if (n == null || n < 0) return 'n/a';
	return n.toLocaleString();
}
