<template>
  <div>
    <div class="flex items-center justify-between mb-3">
      <div>
        <h1 class="text-2xl font-bold">{{ board?.name }}</h1>
        <p class="text-sm text-slate-500">{{ board?.description }}</p>
      </div>
    </div>

    <NewPostForm
      v-if="userStore.authenticated"
      :boardSlug="slug"
      @created="onPostCreated"
    />

    <div class="mt-4 grid gap-3">
      <div v-for="p in posts" :key="p.id" class="p-3 border rounded">
        <router-link :to="`/posts/${p.id}`" class="font-semibold">{{
          p.title
        }}</router-link>
        <div class="text-sm mt-1 whitespace-pre-wrap">{{ p.content }}</div>
        <div class="text-xs text-slate-500 mt-2">
          Likes: {{ p.likes || 0 }} • Dislikes: {{ p.dislikes || 0 }}
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { ref, onMounted } from "vue";
import api from "../services/api";
import { useRoute } from "vue-router";
import NewPostForm from "./NewPostForm.vue";
import { useUserStore } from "../stores/user";
import { wsClient } from "../services/ws";

export default {
  components: { NewPostForm },
  props: { slug: { type: String, required: false } },
  setup(props) {
    const route = useRoute();
    const slug = props.slug || (route.params.slug as string) || "";
    const board = ref<any | null>(null);
    const posts = ref<any[]>([]);
    const userStore = useUserStore();

    async function fetchBoard() {
      const boardsRes = await api.get("/boards");
      const found = boardsRes.data.find((b: any) => b.slug === slug);
      board.value = found || null;
      if (!found) return;
      const res = await api.get(`/boards/${slug}/posts`);
      posts.value = res.data;
    }

    function onPostCreated(post: any) {
      // add to top
      posts.value.unshift(post);
    }

    function wsHandler(msg: any) {
      if (msg.type === "post.created") {
        const post = msg.data;
        if (
          post.board_id === board.value?.id ||
          post.boardId === board.value?.id
        ) {
          posts.value.unshift(post);
        }
      }
    }

    onMounted(async () => {
      await fetchBoard();
      wsClient.addListener(wsHandler);
    });

    return { board, posts, slug, onPostCreated, userStore };
  },
};
</script>
