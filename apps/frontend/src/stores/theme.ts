// apps/frontend/src/stores/theme.ts
import { defineStore } from "pinia";

export const useThemeStore = defineStore("theme", {
  state: () => ({
    mode: (localStorage.getItem("blueberry_theme") as "dark" | "light" | null) || "dark"
  }),
  actions: {
    applyMode(mode: "dark" | "light") {
      this.mode = mode;
      localStorage.setItem("blueberry_theme", mode);
      if (mode === "dark") {
        document.body.classList.remove("light");
        document.body.classList.add("dark");
      } else {
        document.body.classList.remove("dark");
        document.body.classList.add("light");
      }
    },
    toggle() {
      this.applyMode(this.mode === "dark" ? "light" : "dark");
    }
  }
});
