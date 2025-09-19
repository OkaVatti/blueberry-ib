// /apps/frontend/stores/users.ts
import { defineStore } from "pinia";
import { ref } from "vue";
import type { User, UserProfile } from "~/types/user";

export const useUserStore = defineStore("users", () => {
  const config = useRuntimeConfig();
  const apiBase = config.public.apiBase ?? "/api";

  const byId = ref<Record<string, User>>({});
  const loading = ref(false);
  const error = ref<string | null>(null);

  async function fetchUser(usernameOrId: string) {
    loading.value = true;
    try {
      const u = await $fetch(
        `${apiBase}/users/${encodeURIComponent(usernameOrId)}`
      );
      byId.value[(u as User).id] = u as User;
      return { ok: true, data: u as User };
    } catch (err: any) {
      return { ok: false, error: err?.message ?? "fetch user failed" };
    } finally {
      loading.value = false;
    }
  }

  async function updateProfile(id: string, payload: Partial<UserProfile>) {
    try {
      const u = await $fetch(`${apiBase}/users/${encodeURIComponent(id)}`, {
        method: "PUT",
        body: payload,
      });
      byId.value[id] = u as User;
      return { ok: true, data: u };
    } catch (err: any) {
      return { ok: false, error: err?.message ?? "update failed" };
    }
  }

  async function fetchUserPosts(
    usernameOrId: string,
    params?: { limit?: number; cursor?: string }
  ) {
    try {
      const q = new URLSearchParams();
      if (params?.limit) q.set("limit", `${params.limit}`);
      if (params?.cursor) q.set("cursor", params.cursor);
      const list = await $fetch(
        `${apiBase}/users/${encodeURIComponent(
          usernameOrId
        )}/posts?${q.toString()}`
      );
      return { ok: true, data: list };
    } catch (err: any) {
      return { ok: false, error: err?.message ?? "fetch posts failed" };
    }
  }

  return {
    byId,
    loading,
    error,
    fetchUser,
    updateProfile,
    fetchUserPosts,
  };
});
