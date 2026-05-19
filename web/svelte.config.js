import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
export default {
	preprocess: vitePreprocess(),
	kit: {
		// SPA: a single index.html fallback, served by the Go binary's
		// spaHandler. Output goes straight into the Go embed directory.
		adapter: adapter({
			pages: '../internal/web/dist',
			assets: '../internal/web/dist',
			fallback: 'index.html',
			precompress: false,
			strict: false
		})
	}
};
