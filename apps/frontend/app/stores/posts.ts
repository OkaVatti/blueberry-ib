// /apps/frontend/stores/posts.ts
import { defineStore } from "pinia";
import { ref } from "vue";
import type { Post, Comment } from "~/types/post";

export const usePostStore = defineStore("posts", () => {
  const config = useRuntimeConfig();
  const apiBase = config.public.apiBase ?? "/api";

  const posts = ref<Record<string, Post>>({});
  const loading = ref(false);
  const error = ref<string | null>(null);

  async function fetchPost(postId: string) {
    loading.value = true;
    try {
      const p = await $fetch(`${apiBase}/posts/${encodeURIComponent(postId)}`);
      posts.value[p.id] = p;
      return { ok: true, data: p as Post };
    } catch (err: any) {
      return { ok: false, error: err?.message ?? "fetch post failed" };
    } finally {
      loading.value = false;
    }
  }

  async function createPost(payload: Partial<Post>) {
    try {
      const p = await $fetch(`${apiBase}/posts`, {
        method: "POST",
        body: payload,
      });
      posts.value[(p as Post).id] = p as Post;
      return { ok: true, data: p };
    } catch (err: any) {
      return { ok: false, error: err?.message ?? "create failed" };
    }
  }

  async function deletePost(postId: string) {
    try {
      await $fetch(`${apiBase}/posts/${encodeURIComponent(postId)}`, {
        method: "DELETE",
      });
      delete posts.value[postId];
      return { ok: true };
    } catch (err: any) {
      return { ok: false, error: err?.message ?? "delete failed" };
    }
  }

  async function like(postId: string) {
    try {
      const res = await $fetch(
        `${apiBase}/posts/${encodeURIComponent(postId)}/like`,
        { method: "POST" }
      );
      if (posts.value[postId])
        posts.value[postId].reactionCounts.likes =
          (posts.value[postId].reactionCounts.likes ?? 0) + 1;
      return { ok: true, data: res };
    } catch (err: any) {
      return { ok: false, error: err?.message };
    }
  }

  async function dislike(postId: string) {
    try {
      const res = await $fetch(
        `${apiBase}/posts/${encodeURIComponent(postId)}/dislike`,
        { method: "POST" }
      );
      if (posts.value[postId])
        posts.value[postId].reactionCounts.dislikes =
          (posts.value[postId].reactionCounts.dislikes ?? 0) + 1;
      return { ok: true, data: res };
    } catch (err: any) {
      return { ok: false, error: err?.message };
    }
  }

  async function repost(postId: string) {
    try {
      const res = await $fetch(
        `${apiBase}/posts/${encodeURIComponent(postId)}/repost`,
        { method: "POST" }
      );
      return { ok: true, data: res };
    } catch (err: any) {
      return { ok: false, error: err?.message };
    }
  }

  async function fetchComments(
    postId: string,
    params?: { limit?: number; cursor?: string }
  ) {
    try {
      const q = new URLSearchParams();
      if (params?.limit) q.set("limit", `${params.limit}`);
      if (params?.cursor) q.set("cursor", params.cursor);
      const list = await $fetch(
        `${apiBase}/posts/${encodeURIComponent(
          postId
        )}/comments?${q.toString()}`
      );
      return { ok: true, data: list as Comment[] };
    } catch (err: any) {
      return { ok: false, error: err?.message ?? "fetch comments failed" };
    }
  }

  async function addComment(postId: string, payload: Partial<Comment>) {
    try {
      const c = await $fetch(
        `${apiBase}/posts/${encodeURIComponent(postId)}/comments`,
        { method: "POST", body: payload }
      );
      return { ok: true, data: c as Comment };
    } catch (err: any) {
      return { ok: false, error: err?.message ?? "comment failed" };
    }
  }

  return {
    posts,
    loading,
    error,
    fetchPost,
    createPost,
    deletePost,
    like,
    dislike,
    repost,
    fetchComments,
    addComment,
  };
});
