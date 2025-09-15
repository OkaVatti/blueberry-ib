<template>
  <div>
    <div class="p-3 border rounded mb-3">
      <h1 class="text-xl font-semibold">{{ post?.title }}</h1>
      <div class="text-sm mt-2 whitespace-pre-wrap">{{ post?.content }}</div>
      <div class="text-xs text-slate-500 mt-3">
        Likes: {{ post?.likes || 0 }} • Dislikes: {{ post?.dislikes || 0 }}
      </div>
      <div class="flex gap-2 mt-2">
        <button @click="vote(1)" class="px-2 py-1 bg-sky-100 rounded">
          Like
        </button>
        <button @click="vote(-1)" class="px-2 py-1 bg-rose-100 rounded">
          Dislike
        </button>
      </div>
    </div>

    <div class="mb-3">
      <h3 class="font-semibold">Comments</h3>
      <textarea
        v-model="commentText"
        class="input mb-2"
        rows="3"
        placeholder="Add a comment"
      ></textarea>
      <div>
        <button @click="submitComment" class="btn">Comment</button>
        <div v-if="commentErr" class="text-rose-600 text-sm">
          {{ commentErr }}
        </div>
      </div>
      <div class="mt-3">
        <div v-for="c in comments" :key="c.id" class="p-2 border rounded mb-2">
          <div class="text-sm whitespace-pre-wrap">{{ c.content }}</div>
          <div class="text-xs text-slate-500">
            by {{ c.user_id || "anon" }} • {{ c.created_at }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { ref, onMounted } from "vue";
import api from "../services/api";
import { useRoute } from "vue-router";
import { useUserStore } from "../stores/user";
import { wsClient } from "../services/ws";

export default {
  setup() {
    const route = useRoute();
    const id = route.params.id as string;
    const post = ref<any | null>(null);
    const comments = ref<any[]>([]);
    const commentText = ref("");
    const commentErr = ref<string | null>(null);
    const userStore = useUserStore();

    async function fetchPost() {
      const res = await api
        .get(`/posts/${id}/comments`)
        .catch(() => ({ data: [] }));
      comments.value = res.data;
      // fetch post list by searching boards? instead call backend /boards/:slug/posts earlier
      // but simple approach: try to fetch from posts list endpoints (not provided), so rely on comment res for now
      // In our backend we don't have GET /posts/:id, so fetch all boards/posts is needed
      // quick attempt: fetch boards and search posts
      const boards = await api.get("/boards").catch(() => ({ data: [] }));
      for (const b of boards.data) {
        const ps = await api
          .get(`/boards/${b.slug}/posts`)
          .catch(() => ({ data: [] }));
        const found = ps.data.find((p: any) => p.id === id);
        if (found) {
          post.value = found;
          break;
        }
      }
    }

    async function submitComment() {
      commentErr.value = null;
      if (!userStore.authenticated) {
        commentErr.value = "Login to comment";
        return;
      }
      try {
        const res = await api.post(`/posts/${id}/comments`, {
          content: commentText.value,
        });
        commentText.value = "";
        comments.value.push(res.data);
      } catch (e: any) {
        commentErr.value = e?.response?.data?.error || "failed to comment";
      }
    }

    async function vote(v: number) {
      if (!userStore.authenticated) {
        alert("login to vote");
        return;
      }
      await api
        .post(`/posts/${id}/vote`, { value: v })
        .then((r) => {
          if (post.value) {
            post.value.likes = r.data.likes;
            post.value.dislikes = r.data.dislikes;
          }
        })
        .catch(() => {});
    }

    function wsHandler(msg: any) {
      if (msg.type === "comment.created") {
        const c = msg.data;
        if (c.post_id === id || c.postId === id) {
          comments.value.push(c);
        }
      } else if (msg.type === "post.created") {
        const p = msg.data;
        if (p.id === id) {
          post.value = p;
        }
      }
    }

    onMounted(async () => {
      await fetchPost();
      wsClient.addListener(wsHandler);
    });

    return {
      post,
      comments,
      commentText,
      submitComment,
      commentErr,
      vote,
      userStore,
    };
  },
};
</script>

<style scoped>
.input {
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 6px;
  width: 100%;
}
.btn {
  padding: 8px 12px;
  background: #10b981;
  color: white;
  border-radius: 6px;
  margin-top: 6px;
}
</style>
