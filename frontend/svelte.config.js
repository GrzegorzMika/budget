import adapter from '@sveltejs/adapter-auto';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	// Consult https://svelte.dev/docs/kit/integrations
	// for more information about preprocessors
	preprocess: vitePreprocess(),

	kit: {
		// adapter-auto only supports some environments, see https://svelte.dev/docs/kit/adapter-auto for a list.
		// If your environment is not supported, or you settled on a specific environment, switch out the adapter.
		// See https://svelte.dev/docs/kit/adapters for more information about adapters.
		adapter: adapter(),
		csp: {
			directives: {
				'base-uri': ['self'],
				'default-src': ['self'],
				'script-src': ['self'],
				'style-src': ['self'],
				'img-src': ['self'],
				'media-src': ['none'],
				'font-src': ['self'],
				'connect-src': [
					'https://api.budget.gregdev.dev',
					'https://idp.budget.gregdev.dev',
					'http://localhost:3000',
					'http://127.0.0.1:8080',
					'self',
				],
				'form-action': ['https://budget.gregdev.dev', 'http://localhost:5173'],
				'frame-src': ['none'],
				'worker-src': ['none']
			}
		}
	}
};

export default config;
