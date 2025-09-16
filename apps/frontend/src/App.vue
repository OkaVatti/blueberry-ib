<template>
  <div>
    <Navbar />
    <main>
      <router-view />
    </main>
  </div>
</template>

<script lang="ts">
import Navbar from ".//components/Navbar.vue";
export default { components: { Navbar } };
</script>

<style>
main {
  min-height: calc(100vh - 64px);
}
</style>
