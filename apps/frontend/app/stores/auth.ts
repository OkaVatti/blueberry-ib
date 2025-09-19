/* eslint-disable @typescript-eslint/no-explicit-any */
// /apps/frontend/stores/auth.ts
import { defineStore } from "pinia";
import { ref } from "vue";
import type { User, SignupDTO, AuthTokens } from "~/types/user";
import type { ApiResult } from "~/types/api";
import { useRuntimeConfig } from "#imports";

const STORAGE_KEY = "blueberry:auth:v1";

export const useAuthStore = defineStore("auth", () => {
  const config = useRuntimeConfig();
  const apiBase = config.public.apiBase ?? "";

  const user = ref<User | null>(null);
  const tokens = ref<AuthTokens | null>(null);
  const loading = ref(false);
  const lastError = ref<string | null>(null);

  // hydrate on client
  if (import.meta.client) {
    try {
      const raw = localStorage.getItem(STORAGE_KEY);
      if (raw) {
        const parsed = JSON.parse(raw) as { user?: User; tokens?: AuthTokens };
        user.value = parsed.user ?? null;
        tokens.value = parsed.tokens ?? null;
      }
    } catch {
      // ignore
    }
  }

  async function persist() {
    if (!import.meta.client) return;
    try {
      localStorage.setItem(
        STORAGE_KEY,
        JSON.stringify({ user: user.value, tokens: tokens.value })
      );
    } catch {
      // swallow
    }
  }

  function clearPersist() {
    if (!import.meta.client) return;
    try {
      localStorage.removeItem(STORAGE_KEY);
    } catch { /* empty */ }
  }

  async function login(
    principal: string,
    password: string
  ): Promise<ApiResult<{ user: User; tokens: AuthTokens }>> {
    loading.value = true;
    lastError.value = null;
    try {
      const res = await $fetch<{ user: User; tokens: AuthTokens }>(
        `${apiBase}/auth/login`,
        {
          method: "POST",
          body: { principal, password },
        }
      );
      user.value = res.user;
      tokens.value = res.tokens;
      await persist();
      return { ok: true, data: res };
    } catch (err: unknown) {
      const e = err as any;
      const message = e?.data?.message ?? e?.message ?? "login failed";
      lastError.value = message;
      return {
        ok: false,
        error: { message, status: e?.response?.status ?? e?.status },
      };
    } finally {
      loading.value = false;
    }
  }

  async function signup(
    payload: SignupDTO
  ): Promise<ApiResult<{ user: User; tokens: AuthTokens }>> {
    loading.value = true;
    lastError.value = null;
    try {
      const res = await $fetch<{ user: User; tokens: AuthTokens }>(
        `${apiBase}/auth/signup`,
        {
          method: "POST",
          body: payload,
        }
      );
      user.value = res.user;
      tokens.value = res.tokens;
      await persist();
      return { ok: true, data: res };
    } catch (err: unknown) {
      const e = err as any;
      const message = e?.data?.message ?? e?.message ?? "signup failed";
      lastError.value = message;
      return {
        ok: false,
        error: { message, status: e?.response?.status ?? e?.status },
      };
    } finally {
      loading.value = false;
    }
  }

  async function logout(): Promise<ApiResult<null>> {
    // attempt best-effort server-side invalidate, but always clear local state
    try {
      if (tokens.value?.refreshToken) {
        await $fetch(`${apiBase}/auth/logout`, {
          method: "POST",
          body: { refreshToken: tokens.value.refreshToken },
        }).catch(() => {});
      }
    } finally {
      user.value = null;
      tokens.value = null;
      clearPersist();
    }
    return { ok: true, data: null };
  }

  async function refresh(): Promise<ApiResult<{ tokens: AuthTokens }>> {
    // Call to server to refresh tokens using refresh token
    if (!tokens.value?.refreshToken)
      return { ok: false, error: { message: "no refresh token" } };
    try {
      const res = await $fetch<{ tokens: AuthTokens }>(
        `${apiBase}/auth/refresh`,
        { method: "POST", body: { refreshToken: tokens.value.refreshToken } }
      );
      tokens.value = res.tokens;
      await persist();
      return { ok: true, data: res };
    } catch (err: unknown) {
      // clear auth if refresh fails
      user.value = null;
      tokens.value = null;
      clearPersist();
      const e = err as any;
      return {
        ok: false,
        error: {
          message: e?.data?.message ?? e?.message ?? "refresh failed",
          status: e?.response?.status ?? e?.status,
        },
      };
    }
  }

  async function fetchProfile(): Promise<ApiResult<User>> {
    try {
      const res = await $fetch<User>(`${apiBase}/me`);
      user.value = res;
      await persist();
      return { ok: true, data: res };
    } catch (err: unknown) {
      const e = err as any;
      return {
        ok: false,
        error: {
          message: e?.data?.message ?? e?.message ?? "fetch profile failed",
          status: e?.response?.status ?? e?.status,
        },
      };
    }
  }

  return {
    user,
    tokens,
    loading,
    lastError,
    login,
    signup,
    logout,
    refresh,
    fetchProfile,
    persist,
  };
});
