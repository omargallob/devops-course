// Astro config: enables MDX for lesson content, Tailwind CSS v4 via Vite
// plugin, and proxies /api requests to the Go backend. The proxy target is
// configurable via the API_URL env var (defaults to http://localhost:8080,
// overridden to http://api:8080 inside Docker Compose).
import { defineConfig } from 'astro/config';
import mdx from '@astrojs/mdx';
import tailwindcss from '@tailwindcss/vite';

const apiURL = process.env.API_URL || 'http://localhost:8080';

export default defineConfig({
  integrations: [mdx()],
  vite: {
    plugins: [tailwindcss()],
    server: {
      proxy: {
        '/api': apiURL,
      },
    },
  },
});
