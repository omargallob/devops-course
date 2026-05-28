// Astro config: enables MDX for lesson content, Tailwind CSS v4 via Vite
// plugin, and proxies /api requests to the Go backend (http://localhost:8080).
import { defineConfig } from 'astro/config';
import mdx from '@astrojs/mdx';
import tailwindcss from '@tailwindcss/vite';

export default defineConfig({
  integrations: [mdx()],
  vite: {
    plugins: [tailwindcss()],
    server: {
      proxy: {
        '/api': 'http://localhost:8080',
      },
    },
  },
});
