import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'
import {
  resolve
} from 'path'
import typescript from '@rollup/plugin-typescript';

// https://vitejs.dev/config/
export default defineConfig({
  build: {
    lib: {
      // could also be a dictionary or array of multiple entry points
      entry: resolve(__dirname, 'src/index.ts'),
      formats: ["es"],
      name: 'mylib',
      // the proper extensions will be added
      filename: 'my-lib',
    },
  },
  plugins: [
    svelte(),
  ],
})
