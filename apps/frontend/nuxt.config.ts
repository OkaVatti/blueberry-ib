import tailwindcss from "@tailwindcss/vite";

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: "2025-07-15",
  devtools: { enabled: true },
  builder: "vite",
  buildDir: "blueberry-ib-build",
  dev: true,

  devServer: {
    host: "0.0.0.0",
    port: 9876,
    https: {
      key: "ssl/lilith.key",
      cert: "ssl/blueberry.lan.crt",
    },
  },

  dir: {
    assets: "app/assets",
    layouts: "app/layouts",
    pages: "app/pages",
    public: "public",
  },

  modules: [
    "@nuxtjs/sitemap",
    "@nuxtjs/robots",
    "@nuxt/eslint",
    "@nuxt/image",
    "@nuxt/fonts",
    "@nuxtjs/color-mode",
    "nuxt-feather-icons",
    "@pinia/nuxt",
  ],

  imports: {
    dirs: ['stores']
  },

  // Vite-specific server config (covers vite dev server)
  vite: {
    plugins: [tailwindcss()],
  },

  typescript: {
    strict: true,
  },
});