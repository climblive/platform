import { svelte } from "@sveltejs/vite-plugin-svelte";
import path from "path";
import copy from "rollup-plugin-copy";
import { defineConfig } from "vite";

export default defineConfig({
  build: {
    minify: true,
  },
  resolve: {
    alias: {
      "@": path.resolve(__dirname, 'src'),
    }
  },
  plugins: [
    svelte(),
    copy({
      copyOnce: true,
      targets: [
        {
          src: path.resolve(
            __dirname,
            "node_modules/@shoelace-style/shoelace/dist/assets"
          ),
          dest: path.resolve(__dirname, "public/shoelace"),
        },
      ],
    }),
  ],
});
