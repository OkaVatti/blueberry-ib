// /apps/frontend/stores/ui.ts
import { defineStore } from "pinia";
import { ref, computed } from "vue";

type ThemeName = "fern" | "kitten" | "cafe" | "portland" | "north" | "default";

export const useUIStore = defineStore("ui", () => {
  const sidebarOpen = ref(false);
  const theme = ref<ThemeName>("default");
  const compactMode = ref(false);

  const isDark = computed(() => {
    if (theme.value === "default") {
      if (import.meta.client && window.matchMedia)
        return window.matchMedia("(prefers-color-scheme: dark)").matches;
      return false;
    }
    return theme.value === "north"; // example map
  });

  function toggleSidebar() {
    sidebarOpen.value = !sidebarOpen.value;
  }
  function openSidebar() {
    sidebarOpen.value = true;
  }
  function closeSidebar() {
    sidebarOpen.value = false;
  }
  function setTheme(t: ThemeName) {
    theme.value = t;
  }
  function toggleCompact() {
    compactMode.value = !compactMode.value;
  }

  return {
    sidebarOpen,
    theme,
    compactMode,
    isDark,
    toggleSidebar,
    openSidebar,
    closeSidebar,
    setTheme,
    toggleCompact,
  };
});
