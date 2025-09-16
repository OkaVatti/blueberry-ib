<template>
  <div>
    <Navbar />
    <div class="container-grid">
      <aside>
        <div class="card">
          <h4 class="font-semibold mb-2">Board</h4>
          <div class="board-meta">/{{ slug }}</div>
        </div>
      </aside>

      <main>
        <router-view />
      </main>

      <aside>
        <div class="card">
          <h4 class="font-semibold mb-2">Board tools</h4>
          <div class="board-meta">filters, sort</div>
        </div>
      </aside>
    </div>
  </div>
</template>

<script lang="ts">
import Navbar from "../components/Navbar.vue";
import { defineComponent } from "vue";
export default defineComponent({
  components: { Navbar },
  computed: {
    slug() {
      // pull slug from route param if available
      return (this.$route.params.slug as string) || "";
    },
  },
});
</script>
