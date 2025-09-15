<template>
  <div>
    <h2>Boards</h2>
    <ul>
      <li v-for="b in boards" :key="b.id" style="margin-bottom: 8px">
        <strong>{{ b.name }}</strong> /{{ b.slug }} — {{ b.description }}
        <div>
          <button @click="openBoard(b.slug)">Open</button>
        </div>
      </li>
    </ul>

    <div v-if="selectedBoard">
      <hr />
      <h3>Board: {{ selectedBoard.name }}</h3>
      <div style="display: flex; gap: 8px; margin-bottom: 8px">
        <input v-model="newTitle" placeholder="title" />
        <input v-model="newContent" placeholder="content" />
        <button @click="createPost">Post</button>
      </div>

      <div>
        <h4>Posts</h4>
        <div
          v-for="p in posts"
          :key="p.id"
          style="border: 1px solid #ddd; padding: 8px; margin-bottom: 8px"
        >
          <div>
            <strong>{{ p.title }}</strong>
          </div>
          <div style="white-space: pre-wrap">{{ p.content }}</div>
          <div style="margin-top: 6px">
            Likes: {{ p.likes }} Dislikes: {{ p.dislikes }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, onMounted } from "vue";
import api, { setAuthToken } from "../services/api";

export default defineComponent({
  props: {
    token: { type: String, required: false },
  },
  setup(props) {
    const boards = ref<any[]>([]);
    const selectedBoard = ref<any | null>(null);
    const posts = ref<any[]>([]);
    const newTitle = ref("");
    const newContent = ref("");
    const ws = ref<WebSocket | null>(null);

    async function fetchBoards() {
      const res = await api.get("/boards");
      boards.value = res.data;
    }

    async function openBoard(slug: string) {
      const b = boards.value.find((x: any) => x.slug === slug);
      if (!b) return;
      selectedBoard.value = b;
      const res = await api.get(`/boards/${slug}/posts`);
      posts.value = res.data;
    }

    async function createPost() {
      if (!selectedBoard.value) return;
      try {
        await api.post(`/boards/${selectedBoard.value.slug}/posts`, {
          title: newTitle.value,
          content: newContent.value,
        });
        newTitle.value = "";
        newContent.value = "";
        // optimistic fetch
        const res = await api.get(`/boards/${selectedBoard.value.slug}/posts`);
        posts.value = res.data;
      } catch (e) {
        alert("need to be logged in to post");
      }
    }

    onMounted(() => {
      if (props.token) setAuthToken(props.token);
      fetchBoards();

      ws.value = new WebSocket(`ws://${location.hostname}:8080/ws`);
      ws.value.onmessage = (ev) => {
        try {
          const msg = JSON.parse(ev.data);
          if (msg.type === "post.created") {
            // if board is open and matches, add to posts
            const post = msg.data;
            if (
              selectedBoard.value &&
              post.board_id === selectedBoard.value.id
            ) {
              posts.value.unshift(post);
            }
          }
        } catch (err) {}
      };
    });

    return {
      boards,
      openBoard,
      selectedBoard,
      posts,
      newTitle,
      newContent,
      createPost,
    };
  },
});
</script>
