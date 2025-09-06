import { sentryVitePlugin } from "@sentry/vite-plugin";
import { svelte } from "@sveltejs/vite-plugin-svelte";
import path from "path";
import { defineConfig } from "vite";

export default defineConfig({
  build: {
    minify: true,
    sourcemap: true,
    cssCodeSplit: true,
  },
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "src"),
    },
  },
  plugins: [
    svelte(),
    sentryVitePlugin({
      org: "climblive",
      project: "app",
      authToken: process.env.SENTRY_AUTH_TOKEN,
    }),
  ],
});
