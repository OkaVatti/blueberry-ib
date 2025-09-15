<template>
  <div>
    <h2 class="text-xl font-bold mb-3">Moderator Dashboard</h2>

    <div>
      <h3 class="font-semibold">Recent Posts (last 200)</h3>
      <button @click="refresh" class="btn mb-2">Refresh</button>
      <div v-for="p in posts" :key="p.id" class="p-2 border rounded mb-1">
        <div class="font-semibold">{{ p.title }}</div>
        <div class="text-sm">{{ p.content }}</div>
        <div class="flex gap-2 mt-2">
          <button
            @click="delPost(p.id)"
            class="px-2 py-1 bg-rose-100 rounded text-xs"
          >
            Delete
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { ref, onMounted } from "vue";
import api from "../services/api";
import { useUserStore } from "../stores/user";

export default {
  setup() {
    const posts = ref<any[]>([]);
    const userStore = useUserStore();
    async function refresh() {
      // quick: fetch boards and first 50 posts across them
      const boardsRes = await api.get("/boards").catch(() => ({ data: [] }));
      const list: any[] = [];
      for (const b of boardsRes.data) {
        try {
          const ps = await api
            .get(`/boards/${b.slug}/posts`)
            .catch(() => ({ data: [] }));
          for (const p of ps.data.slice(0, 50)) list.push(p);
        } catch {}
      }
      posts.value = list.slice(0, 200);
    }
    async function delPost(id: string) {
      await api.delete(`/posts/${id}`).catch(() => {});
      refresh();
    }

    onMounted(() => {
      if (userStore.authenticated) refresh();
    });

    return { posts, refresh, delPost };
  },
};
</script>

<style scoped>
.btn {
  padding: 6px 10px;
  background: #0ea5e9;
  color: white;
  border-radius: 6px;
  margin-bottom: 8px;
}
</style>
