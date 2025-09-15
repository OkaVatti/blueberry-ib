<template>
  <nav
    class="bg-white/80 dark:bg-slate-800/80 backdrop-blur p-3 rounded-md mb-4 flex items-center justify-between"
  >
    <div>
      <router-link to="/" class="font-bold text-lg">Blueberry</router-link>
      <span class="text-sm text-slate-500 ml-2">prototype</span>
    </div>

    <div class="flex items-center gap-3">
      <div v-if="userStore.authenticated" class="flex items-center gap-3">
        <span class="text-sm">👤 {{ userStore.user?.username }}</span>
        <button
          @click="logout"
          class="px-3 py-1 text-sm rounded bg-red-100 dark:bg-red-700"
        >
          Logout
        </button>
      </div>
      <div v-else>
        <router-link
          to="/login"
          class="px-3 py-1 rounded bg-sky-100 dark:bg-sky-700"
          >Login / Register</router-link
        >
      </div>
    </div>
  </nav>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import { useUserStore } from "../stores/user";

export default defineComponent({
  setup() {
    const userStore = useUserStore();
    function logout() {
      userStore.logout();
    }
    return { userStore, logout };
  },
});
</script>
