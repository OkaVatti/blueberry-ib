<template>
  <div>
    <Navbar />
    <div class="container-grid">
      <aside>
        <div class="card">
          <h4 class="font-semibold mb-2">Search</h4>
          <div class="board-meta">filter by board, user, type</div>
        </div>
      </aside>

      <main>
        <router-view />
      </main>

      <aside>
        <div class="card">
          <h4 class="font-semibold mb-2">Saved searches</h4>
        </div>
      </aside>
    </div>
  </div>
</template>

<script lang="ts">
import Navbar from "../components/Navbar.vue";
export default { components: { Navbar } };
</script>
