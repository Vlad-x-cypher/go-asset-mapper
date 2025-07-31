import tailwindcss from "@tailwindcss/vite";
import { defineConfig } from "vite";

export default defineConfig (() => ({
    base: '/', 
    build: {
        outDir: './public',
        assetsDir: 'assets',
        assetsInlineLimit: 0,
        manifest: true,
        rollupOptions: {
            input: {
                app: 'src/app.js',
            }
        },
    },
    plugins: [
        tailwindcss(),
    ],
}))
