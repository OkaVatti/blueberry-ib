// /apps/frontend/stores/threads.ts
import { defineStore } from "pinia";
import { ref } from "vue";
import type { Thread, ThreadSummary } from "~/types/thread";

export const useThreadStore = defineStore("threads", () => {
  const config = useRuntimeConfig();
  const apiBase = config.public.apiBase ?? "/api";

  const threads = ref<Record<string, Thread>>({});
  const summaries = ref<ThreadSummary[]>([]);
  const loading = ref(false);
  const error = ref<string | null>(null);

  async function fetchThread(threadId: string) {
    loading.value = true;
    try {
      const t = await $fetch(
        `${apiBase}/threads/${encodeURIComponent(threadId)}`
      );
      threads.value[threadId] = t as Thread;
      return { ok: true, data: t as Thread };
    } catch (err: any) {
      return { ok: false, error: err?.message ?? "fetch thread failed" };
    } finally {
      loading.value = false;
    }
  }

  async function listBoardThreads(
    boardId: string,
    params?: { limit?: number; cursor?: string }
  ) {
    loading.value = true;
    try {
      const q = new URLSearchParams();
      if (params?.limit) q.set("limit", `${params.limit}`);
      if (params?.cursor) q.set("cursor", params.cursor);
      const list = await $fetch(
        `${apiBase}/boards/${encodeURIComponent(
          boardId
        )}/threads?${q.toString()}`
      );
      summaries.value = list as ThreadSummary[];
      return { ok: true, data: list };
    } catch (err: any) {
      return { ok: false, error: err?.message ?? "list threads failed" };
    } finally {
      loading.value = false;
    }
  }

  async function createThread(payload: Partial<Thread>) {
    try {
      const t = await $fetch(`${apiBase}/threads`, {
        method: "POST",
        body: payload,
      });
      threads.value[(t as Thread).id] = t as Thread;
      return { ok: true, data: t };
    } catch (err: any) {
      return { ok: false, error: err?.message ?? "create thread failed" };
    }
  }

  async function reply(
    threadId: string,
    payload: { body: string; media?: any[] }
  ) {
    try {
      const r = await $fetch(
        `${apiBase}/threads/${encodeURIComponent(threadId)}/posts`,
        { method: "POST", body: payload }
      );
      // optionally refresh thread
      await fetchThread(threadId);
      return { ok: true, data: r };
    } catch (err: any) {
      return { ok: false, error: err?.message ?? "reply failed" };
    }
  }

  async function pinThread(threadId: string) {
    try {
      await $fetch(`${apiBase}/threads/${encodeURIComponent(threadId)}/pin`, {
        method: "POST",
      });
      return { ok: true };
    } catch (err: any) {
      return { ok: false, error: err?.message };
    }
  }

  async function closeThread(threadId: string) {
    try {
      await $fetch(`${apiBase}/threads/${encodeURIComponent(threadId)}/close`, {
        method: "POST",
      });
      return { ok: true };
    } catch (err: any) {
      return { ok: false, error: err?.message };
    }
  }

  return {
    threads,
    summaries,
    loading,
    error,
    fetchThread,
    listBoardThreads,
    createThread,
    reply,
    pinThread,
    closeThread,
  };
});
