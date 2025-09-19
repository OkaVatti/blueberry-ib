// /apps/frontend/stores/moderation.ts
import { defineStore } from "pinia";
import { ref } from "vue";
import type { Post, Comment } from "~/types/post";

export const useModerationStore = defineStore("moderation", () => {
  const config = useRuntimeConfig();
  const apiBase = config.public.apiBase ?? "/api";

  const reports = ref<any[]>([]);
  const loading = ref(false);
  const error = ref<string | null>(null);

  async function fetchReports(params?: {
    type?: string;
    limit?: number;
    cursor?: string;
  }) {
    loading.value = true;
    try {
      const q = new URLSearchParams();
      if (params?.type) q.set("type", params.type);
      if (params?.limit) q.set("limit", `${params.limit}`);
      if (params?.cursor) q.set("cursor", params.cursor);
      const list = await $fetch(
        `${apiBase}/moderation/reports?${q.toString()}`
      );
      reports.value = list;
      return { ok: true, data: list };
    } catch (err: any) {
      return { ok: false, error: err?.message ?? "fetch reports failed" };
    } finally {
      loading.value = false;
    }
  }

  async function deletePost(postId: string, reason?: string) {
    try {
      await $fetch(
        `${apiBase}/moderation/posts/${encodeURIComponent(postId)}`,
        { method: "DELETE", body: { reason } }
      );
      return { ok: true };
    } catch (err: any) {
      return { ok: false, error: err?.message ?? "delete post failed" };
    }
  }

  async function deleteComment(commentId: string) {
    try {
      await $fetch(
        `${apiBase}/moderation/comments/${encodeURIComponent(commentId)}`,
        { method: "DELETE" }
      );
      return { ok: true };
    } catch (err: any) {
      return { ok: false, error: err?.message ?? "delete comment failed" };
    }
  }

  async function banUser(
    userId: string,
    payload: { reason?: string; durationSeconds?: number | null }
  ) {
    try {
      const res = await $fetch(
        `${apiBase}/moderation/users/${encodeURIComponent(userId)}/ban`,
        { method: "POST", body: payload }
      );
      return { ok: true, data: res };
    } catch (err: any) {
      return { ok: false, error: err?.message ?? "ban failed" };
    }
  }

  async function liftBan(userId: string) {
    try {
      await $fetch(
        `${apiBase}/moderation/users/${encodeURIComponent(userId)}/ban`,
        { method: "DELETE" }
      );
      return { ok: true };
    } catch (err: any) {
      return { ok: false, error: err?.message ?? "lift ban failed" };
    }
  }

  return {
    reports,
    loading,
    error,
    fetchReports,
    deletePost,
    deleteComment,
    banUser,
    liftBan,
  };
});
