<!-- apps/frontend/src/App.vue -->
<template>
  <div id="app">
    <router-view />
  </div>
</template>

<script lang="ts">
import { defineComponent, onMounted } from "vue";
import { useThemeStore } from "./stores/theme";
import { useUserStore } from "./stores/user";
import { setAuthToken } from "./services/api";

export default defineComponent({
  name: "App",
  setup() {
    const themeStore = useThemeStore();
    const userStore = useUserStore();

    onMounted(() => {
      // Apply saved theme
      themeStore.applyMode(themeStore.mode);
      
      // Restore user session if exists
      const savedToken = localStorage.getItem("blueberry_token");
      const savedUser = localStorage.getItem("blueberry_user");
      
      if (savedToken && savedUser) {
        try {
          const user = JSON.parse(savedUser);
          userStore.setSession(savedToken, user);
        } catch (err) {
          // Clear invalid data
          localStorage.removeItem("blueberry_token");
          localStorage.removeItem("blueberry_user");
        }
      }
    });

    return {
      themeStore,
      userStore
    };
  }
});
</script>

<style>
/* Global styles are now in main.css */
#app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}
</style>