import { sveltekit } from '@sveltejs/kit/vite';
import tailwindcss from '@tailwindcss/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [tailwindcss(), sveltekit()],
	server: {
		// `npm run dev` proxies API calls to a locally-running dms-ui binary.
		proxy: {
			'/api': 'http://127.0.0.1:8099'
		}
	}
});
