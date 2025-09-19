<script setup lang="ts">
// Root app wrapper: import the default layout and global CSS
import DefaultLayout from "./layouts/defaultLayout.vue";

// nothing else required here; page content is injected with <NuxtPage/>
</script>

<template>
  <DefaultLayout>
    <NuxtPage />
  </DefaultLayout>
</template>

<!-- Import vanilla CSS (theme + styles). Use the explicit path so pnpm/nuxt finds them in your app/ folder -->
<style src="~/app/assets/css/theme.css"></style>
<style src="~/app/assets/css/style.css"></style>
