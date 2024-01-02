import { defineConfig } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";
import copy from "rollup-plugin-copy";
import path from "path";
import { visualizer } from "rollup-plugin-visualizer";

export default defineConfig({
  build: {
    minify: true,
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
