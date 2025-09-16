<template>
  <div>
    <div class="card mb-4">
      <input
        v-model="q"
        @keyup.enter="doSearch"
        placeholder="Search posts, comments, users"
        class="input"
      />
      <div class="mt-2 flex gap-2">
        <button @click="doSearch" class="btn-primary">Search</button>
      </div>
    </div>

    <div v-if="loading" class="card">Searching…</div>

    <div v-if="results.posts?.length" class="card mb-3">
      <h3 class="font-semibold mb-2">Posts</h3>
      <div v-for="p in results.posts" :key="p.id" class="mb-2">
        <div class="font-semibold">{{ p.title || "untitled" }}</div>
        <div class="board-meta">{{ p.user_id || "anon" }}</div>
      </div>
    </div>

    <div v-if="results.comments?.length" class="card mb-3">
      <h3 class="font-semibold mb-2">Comments</h3>
      <div v-for="c in results.comments" :key="c.id" class="mb-2">
        <div v-html="c.content_html || c.content"></div>
        <div class="board-meta">on post {{ c.post_id }}</div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from "vue";
import api from "../services/api";

export default defineComponent({
  setup() {
    const q = ref("");
    const loading = ref(false);
    const results = ref({ posts: [], comments: [] });

    async function doSearch() {
      loading.value = true;
      results.value = { posts: [], comments: [] };
      // naive client-side search: fetch boards and filter posts/comments
      const boards = (await api.get("/boards").catch(() => ({ data: [] })))
        .data;
      const matchesPosts: any[] = [];
      const matchesComments: any[] = [];
      for (const b of boards) {
        const posts = (
          await api.get(`/boards/${b.slug}/posts`).catch(() => ({ data: [] }))
        ).data;
        for (const p of posts) {
          if (
            (p.title || "").toLowerCase().includes(q.value.toLowerCase()) ||
            (p.content || "").toLowerCase().includes(q.value.toLowerCase())
          ) {
            matchesPosts.push(p);
          }
          // comments
          const comments = (
            await api.get(`/posts/${p.id}/comments`).catch(() => ({ data: [] }))
          ).data;
          for (const c of comments) {
            if (
              (c.content || "").toLowerCase().includes(q.value.toLowerCase())
            ) {
              matchesComments.push(c);
            }
          }
        }
      }
      results.value.posts = matchesPosts;
      results.value.comments = matchesComments;
      loading.value = false;
    }

    return { q, doSearch, loading, results };
  },
});
</script>

<style scoped>
.input {
  padding: 8px;
  border-radius: 8px;
  width: 100%;
}
</style>
