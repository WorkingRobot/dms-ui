// Thin fetch wrapper around the dms-ui JSON API. Every handler returns either
// a value or { error } with the appropriate status; we normalise that here.

async function req(method, path, body) {
	const opts = { method, headers: {} };
	if (body !== undefined) {
		opts.headers['Content-Type'] = 'application/json';
		opts.body = JSON.stringify(body);
	}
	const res = await fetch(`/api${path}`, opts);
	const text = await res.text();
	let data = null;
	if (text) {
		try {
			data = JSON.parse(text);
		} catch {
			data = { error: text };
		}
	}
	if (!res.ok) {
		throw new Error((data && data.error) || `HTTP ${res.status}`);
	}
	return data;
}

export const api = {
	get: (p) => req('GET', p),
	post: (p, b) => req('POST', p, b ?? {}),
	put: (p, b) => req('PUT', p, b ?? {}),
	del: (p) => req('DELETE', p)
};

export const enc = encodeURIComponent;
