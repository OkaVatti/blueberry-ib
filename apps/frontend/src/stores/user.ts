import { defineStore } from "pinia";
import api, { setAuthToken } from "../services/api";
import type { User } from "../types";

export const useUserStore = defineStore("user", {
  state: () => ({
    token: null as string | null,
    user: null as User | null,
    authenticated: false,
  }),
  actions: {
    async login(usernameOrEmail: string, password: string) {
      const res = await api.post("/auth/login", { usernameOrEmail, password });
      const token = res.data.token;
      const usr = res.data.user;
      this.setSession(token, usr);
      return res.data;
    },
    async register(username: string, email: string, password: string) {
      const res = await api.post("/auth/register", {
        username,
        email,
        password,
      });
      const token = res.data.token;
      const usr = res.data.user;
      this.setSession(token, usr);
      return res.data;
    },
    logout() {
      this.token = null;
      this.user = null;
      this.authenticated = false;
      setAuthToken(null);
      localStorage.removeItem("blueberry_token");
      localStorage.removeItem("blueberry_user");
    },
    setSession(token: string, user: any) {
      this.token = token;
      this.user = user;
      this.authenticated = true;
      setAuthToken(token);
      localStorage.setItem("blueberry_token", token);
      localStorage.setItem("blueberry_user", JSON.stringify(user));
    },
  },
});
