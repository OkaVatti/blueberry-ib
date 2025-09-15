<template>
  <nav
    class="flex items-center justify-between p-3 rounded-md mb-4 bg-white/5 glass"
  >
    <div class="flex items-center gap-3">
      <router-link
        to="/"
        class="text-lg font-bold"
        style="color: var(--color-primary)"
        >Blueberry</router-link
      >
      <span class="board-meta">a new take on imageboards + microblogs</span>
    </div>

    <div class="flex items-center gap-3">
      <div v-if="userStore.authenticated" class="flex items-center gap-3">
        <span class="text-sm">{{ userStore.user?.username }}</span>
      </div>
      <button
        @click="toggleTheme"
        class="px-3 py-1 rounded border"
        :title="themeLabel"
      >
        <svg
          v-if="isDark"
          xmlns="http://www.w3.org/2000/svg"
          class="h-5 w-5"
          viewBox="0 0 20 20"
          fill="currentColor"
        >
          <path
            d="M17.293 13.293a8 8 0 01-10.586-10.586A8 8 0 1017.293 13.293z"
          />
        </svg>
        <svg
          v-else
          xmlns="http://www.w3.org/2000/svg"
          class="h-5 w-5"
          viewBox="0 0 20 20"
          fill="currentColor"
        >
          <path
            d="M10 2a1 1 0 011 1v1a1 1 0 11-2 0V3a1 1 0 011-1zM4.22 4.22a1 1 0 011.415 0L6.64 5.22a1 1 0 11-1.415 1.415L4.22 5.636a1 1 0 010-1.415zM2 10a1 1 0 011-1h1a1 1 0 110 2H3a1 1 0 01-1-1zm8 7a1 1 0 011-1v-1a1 1 0 10-2 0v1a1 1 0 011 1zm5.78-2.78a1 1 0 010-1.415l1.005-1.005a1 1 0 111.415 1.415l-1.005 1.005a1 1 0 01-1.415 0zM17 9a1 1 0 100 2h1a1 1 0 100-2h-1z"
          />
          <path d="M10 5.5a4.5 4.5 0 100 9 4.5 4.5 0 000-9z" />
        </svg>
      </button>

      <router-link
        to="/dash/admin"
        v-if="canAdmin"
        class="px-3 py-1 rounded border"
        >Admin</router-link
      >
      <router-link to="/login" v-else class="px-3 py-1 rounded border"
        >Login</router-link
      >
    </div>
  </nav>
</template>

<script lang="ts">
import { defineComponent, computed } from "vue";
import { useThemeStore } from "../stores/theme";
import { useUserStore } from "../stores/user";

export default defineComponent({
  setup() {
    const theme = useThemeStore();
    const userStore = useUserStore();
    const isDark = computed(() => theme.mode === "dark");
    function toggleTheme() {
      theme.toggle();
    }
    const canAdmin = computed(
      () =>
        userStore.authenticated &&
        ["owner", "coowner", "admin", "moderator"].includes(
          userStore.user?.role || ""
        )
    );
    const themeLabel = computed(() =>
      theme.mode === "dark" ? "Dark (default)" : "Light"
    );
    return { toggleTheme, isDark, userStore, canAdmin, themeLabel };
  },
});
</script>

<style scoped>
nav {
  backdrop-filter: blur(6px);
}
</style>
