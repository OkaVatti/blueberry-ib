<template>
  <div class="container-grid">
    <!-- left: boards -->
    <aside>
      <div class="card">
        <h3 class="font-semibold mb-2">Boards</h3>
        <div
          v-for="b in boards"
          :key="b.id"
          class="mb-2 p-2 rounded hover:bg-ui/5 cursor-pointer"
          @click="openBoard(b.slug)"
        >
          <div class="flex items-center justify-between">
            <div>
              <div class="font-semibold">{{ b.name }}</div>
              <div class="board-meta">/{{ b.slug }} • {{ b.description }}</div>
            </div>
            <div class="text-xs text-muted">{{ formatDate(b.created_at) }}</div>
          </div>
        </div>
      </div>
    </aside>

    <!-- center: feed -->
    <section>
      <div class="card mb-4">
        <div class="flex items-center justify-between">
          <h2 class="text-xl font-bold">Home</h2>
          <div class="board-meta">mix of boards & microfeed</div>
        </div>
      </div>

      <div class="space-y-4">
        <MicroPost v-for="post in microposts" :key="post.id" :post="post" />
        <ThreadPost v-for="p in threadPosts" :key="p.id" :post="p" />
      </div>
    </section>

    <!-- right: trending / threads -->
    <aside>
      <div class="card">
        <h3 class="font-semibold mb-2">Trending threads</h3>
        <div
          v-for="t in threads"
          :key="t.id"
          class="mb-2 p-2 rounded hover:bg-ui/5 cursor-pointer"
          @click="openThread(t.id)"
        >
          <div class="font-semibold">{{ t.title }}</div>
          <div class="board-meta">by {{ t.user_id || "anon" }}</div>
        </div>
      </div>
    </aside>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, onMounted } from "vue";
import api from "../services/api";
import MicroPost from "../components/MicroPost.vue";
import ThreadPost from "../components/ThreadPost.vue";
import { useRouter } from "vue-router";

export default defineComponent({
  components: { MicroPost, ThreadPost },
  setup() {
    const boards = ref<any[]>([]);
    const microposts = ref<any[]>([]);
    const threadPosts = ref<any[]>([]);
    const threads = ref<any[]>([]);
    const router = useRouter();

    async function fetch() {
      const b = await api.get("/boards").catch(() => ({ data: [] }));
      boards.value = b.data;
      // microposts: gather recent posts across boards limit 20
      const mp: any[] = [];
      for (const bd of boards.value) {
        const res = await api
          .get(`/boards/${bd.slug}/posts`)
          .catch(() => ({ data: [] }));
        for (const p of res.data.slice(0, 5)) mp.push(p);
      }
      microposts.value = mp.slice(0, 40);
      // thread posts: fetch recent thread posts (quick approach)
      const th = await api
        .get(`/boards/${boards.value[0]?.slug || "b"}/threads`)
        .catch(() => ({ data: [] }));
      threads.value = th.data.slice(0, 10);
      // fetch first thread posts as examples
      if (threads.value.length) {
        const tposts = await api
          .get(`/threads/${threads.value[0].id}/posts`)
          .catch(() => ({ data: [] }));
        threadPosts.value = tposts.data.slice(0, 10);
      }
    }

    function openBoard(slug: string) {
      router.push(`/boards/${slug}`);
    }
    function openThread(id: string) {
      router.push(`/threads/${id}`);
    }
    function formatDate(_created_at?:
            /// <reference types="../../node_modules/.vue-global-types/vue_3.5_0.d.ts" />
            any) {
      return "";
    }

    onMounted(fetch);

    return {
      boards,
      microposts,
      threadPosts,
      threads,
      openBoard,
      openThread,
      formatDate,
    };
  },
});
</script>
