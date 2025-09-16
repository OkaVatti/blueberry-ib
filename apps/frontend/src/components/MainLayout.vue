<template>
  <div class="container-grid">
    <!-- left -->
    <aside>
      <div class="card">
        <h3 class="font-semibold mb-2">Boards</h3>
        <div
          v-for="b in boards"
          :key="b.id"
          class="mb-2 p-2 rounded hover:bg-ui/5 cursor-pointer"
          @click="openBoard(b.slug)"
        >
          <div class="font-semibold">{{ b.name }}</div>
          <div class="board-meta">/{{ b.slug }} • {{ b.description }}</div>
        </div>
      </div>
    </aside>

    <!-- center -->
    <section>
      <div class="card mb-4">
        <div class="flex items-center justify-between">
          <h2 class="text-xl font-bold">Home</h2>
          <div class="board-meta">mix of boards & microfeed</div>
        </div>
      </div>

      <MarkdownEditor @submit="handleSubmit" />

      <div class="space-y-4 mt-4">
        <MicroPost v-for="post in feed" :key="post.id" :post="post" />
      </div>
    </section>

    <!-- right -->
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
import MarkdownEditor from "./MarkdownEditor.vue";
import MicroPost from ".//MicroPost.vue";
import { useUserStore } from "../stores/user";
import { useRouter } from "vue-router";

export default defineComponent({
  components: { MarkdownEditor, MicroPost },
  setup() {
    const boards = ref<any[]>([]);
    const feed = ref<any[]>([]);
    const threads = ref<any[]>([]);
    const userStore = useUserStore();
    const router = useRouter();

    async function fetch() {
      const b = await api.get("/boards").catch(() => ({ data: [] }));
      boards.value = b.data;
      // quick feed: recent posts from first few boards
      const f: any[] = [];
      for (const bd of boards.value.slice(0, 5)) {
        const res = await api
          .get(`/boards/${bd.slug}/posts`)
          .catch(() => ({ data: [] }));
        for (const p of res.data.slice(0, 10)) f.push(p);
      }
      // sort by created_at desc
      f.sort(
        (a: any, b: any) =>
          new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
      );
      feed.value = f.slice(0, 40);
      // trending threads
      if (boards.value.length) {
        const th = await api
          .get(`/boards/${boards.value[0].slug}/threads`)
          .catch(() => ({ data: [] }));
        threads.value = th.data.slice(0, 8);
      }
    }

    function openBoard(slug: string) {
      router.push(`/boards/${slug}`);
    }
    function openThread(id: string) {
      router.push(`/threads/${id}`);
    }

    onMounted(fetch);

    async function handleSubmit(markdown: string) {
      if (!userStore.authenticated) {
        alert("login to post");
        return;
      }
      // create post on default board (choose first) — ideally let user pick board
      if (!boards.value.length) {
        alert("no board");
        return;
      }
      const slug = boards.value[0].slug;
      try {
        await api.post(`/boards/${slug}/posts`, {
          title: "",
          content: markdown,
        });
        // refresh feed
        await fetch();
      } catch (err) {
        alert("post failed");
      }
    }

    return { boards, feed, threads, openBoard, openThread, handleSubmit };
  },
});
</script>
