import { createApp } from "vue";
import { createPinia } from "pinia";
import router from "./router";
import App from "./App.vue";
import { setAuthToken } from "./services/api";
import { useUserStore } from "./stores/user";
// in main.ts after creating app and pinia:
import { useThemeStore } from "./stores/theme";
import "./style.css";
import "./assets/css/main.css"; // if you use main.css; optional


import "./style.css";

const app = createApp(App);
const pinia = createPinia();
app.use(pinia);
app.use(router);

// hydrate auth token from localStorage
const token = localStorage.getItem("blueberry_token");
if (token) {
  setAuthToken(token);
  // lazy-load user into store if available
  const userStore = useUserStore();
  const raw = localStorage.getItem("blueberry_user");
  if (raw) {
    try {
      userStore.user = JSON.parse(raw);
      userStore.token = token;
      userStore.authenticated = true;
    } catch {}
  }
}

// in main.ts after creating app and pinia:
const theme = useThemeStore();
theme.applyMode(theme.mode);


app.mount("#app");
