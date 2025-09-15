<template>
  <div
    class="min-h-screen bg-slate-50 dark:bg-slate-900 text-slate-900 dark:text-slate-100"
  >
    <Navbar />
    <main class="max-w-5xl mx-auto p-4">
      <router-view />
    </main>
  </div>
</template>

<script lang="ts">
import Navbar from "./components/Navbar.vue";
export default { components: { Navbar } };
</script>
