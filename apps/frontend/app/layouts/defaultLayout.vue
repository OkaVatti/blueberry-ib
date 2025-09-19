<script setup lang="ts">
import { ref } from "vue";
import NavBar from "../components/universal/NavBar.vue";
import SideBar from "../components/universal/SideBar.vue";

const sidebarOpen = ref(false);
function toggleSidebar() {
  sidebarOpen.value = !sidebarOpen.value;
}
</script>

<template>
  <div class="min-h-screen flex flex-col">
    <header class="nb-topbar">
      <NavBar @toggle-sidebar="toggleSidebar">
        <!-- Search slot optional -->
      </NavBar>
    </header>

    <div class="flex-1 flex overflow-hidden">
      <aside class="nb-sidebar hidden md:block">
        <SideBar />
      </aside>

      <!-- mobile slide-in sidebar -->
      <transition name="slide">
        <aside
          v-if="sidebarOpen"
          class="nb-sidebar fixed left-0 top-16 z-50 md:hidden"
        >
          <div class="mb-4 flex justify-end">
            <button class="nb-btn nb-btn--primary" @click="toggleSidebar">
              Close
            </button>
          </div>
          <SideBar />
        </aside>
      </transition>

      <main class="flex-1 overflow-auto p-6">
        <slot />
      </main>
    </div>

    <footer
      class="p-3 text-sm text-gray-500 flex justify-center"
      style="border-top: var(--border-light) dashed var(--muted)"
    >
      © {{ new Date().getFullYear() }} Blueberry — built brutal & beautiful
    </footer>
  </div>
</template>

<style scoped>

.slide-enter-active,
.slide-leave-active {
  transition: transform 0.18s ease;
}

.slide-enter-from {
  transform: translateX(-120%);
}

.slide-enter-to {
  transform: translateX(0);
}

.slide-leave-from {
  transform: translateX(0);
}

.slide-leave-to {
  transform: translateX(-120%);
}
</style>
