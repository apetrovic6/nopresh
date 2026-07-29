import { defineConfig } from 'vite'
import { devtools } from '@tanstack/devtools-vite'

import { tanstackStart } from '@tanstack/react-start/plugin/vite'
import { nitro } from 'nitro/vite'

import viteReact from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

const config = defineConfig({
  resolve: { tsconfigPaths: true },
  // nitro() targets a self-hostable Node server: `vite build` emits
  // .output/server/index.mjs, run with `node .output/server/index.mjs`.
  plugins: [devtools(), tailwindcss(), tanstackStart(), nitro({ preset: "bun" }), viteReact()],
})

export default config
