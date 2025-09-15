import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import tailwindcss from '@tailwindcss/vite'

export default defineConfig({
  plugins: [vue(), tailwindcss()],
  server: {
    port: 9876,
    host: true,
    proxy: {
      "/api": {
        target: "http://127.0.0.1:8080",
        changeOrigin: true,
        secure: false,
        ws: true,
        timeout: 60000, // 60s
        configure: (proxy, _options) => {
          proxy.on("error", (err: any, _req: any, res: any) => {
            console.error("[vite-proxy] proxy error:", err && err.message);
            if (res && !res.headersSent) {
              res.writeHead(502, { "Content-Type": "text/plain" });
            }
            try {
              res.end("Bad gateway (proxy)");
            } catch {}
          });
        },
      },
    },
  },
});
