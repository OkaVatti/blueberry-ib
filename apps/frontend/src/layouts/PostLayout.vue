<template>
  <div>
    <Navbar />
    <div class="container-grid">
      <aside>
        <div class="card">
          <h4 class="font-semibold mb-2">Post</h4>
        </div>
      </aside>

      <main>
        <router-view />
      </main>

      <aside>
        <div class="card">
          <h4 class="font-semibold mb-2">Post Tools</h4>
          <div class="board-meta">share / report</div>
        </div>
      </aside>
    </div>
  </div>
</template>

<script lang="ts">
import Navbar from "../components/Navbar.vue";
export default { components: { Navbar } };
</script>
