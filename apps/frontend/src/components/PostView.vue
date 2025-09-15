<template>
  <div>
    <div class="p-3 border rounded mb-3">
      <h1 class="text-xl font-semibold">{{ post?.title }}</h1>
      <div
        class="text-sm mt-2"
        v-html="post?.content_html || post?.content"
      ></div>
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

    <div>
      <h3 class="font-semibold">Comments</h3>
      <div class="mb-3">
        <textarea
          v-model="commentText"
          class="input mb-2"
          rows="3"
          placeholder="Add a comment"
        ></textarea>
        <div><button @click="submitComment" class="btn">Comment</button></div>
        <div v-if="commentErr" class="text-rose-600 text-sm mt-2">
          {{ commentErr }}
        </div>
      </div>
      <div class="mt-3">
        <div v-for="c in comments" :key="c.id" class="p-2 border rounded mb-2">
          <div v-html="c.content_html || c.content" class="text-sm"></div>
          <div class="text-xs text-slate-500">
            by {{ c.user_id || "anon" }} • {{ c.created_at }}
          </div>
          <div class="flex gap-2 mt-2">
            <button
              @click="voteComment(c.id, 1)"
              class="px-2 py-1 bg-sky-100 rounded text-xs"
            >
              Like
            </button>
            <button
              @click="voteComment(c.id, -1)"
              class="px-2 py-1 bg-rose-100 rounded text-xs"
            >
              Dislike
            </button>
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

export default {
  setup() {
    const route = useRoute();
    const id = route.params.id as string;
    const post = ref<any | null>(null);
    const comments = ref<any[]>([]);
    const commentText = ref("");
    const commentErr = ref<string | null>(null);
    const userStore = useUserStore();

    async function fetchPostAndComments() {
      // load comments
      const res = await api.get(`/posts/${id}/comments`);
      comments.value = res.data;
      // load the post by scanning boards (or add GET /posts/:id server-side for faster fetch)
      const boards = await api.get("/boards");
      for (const b of boards.data) {
        const ps = await api.get(`/boards/${b.slug}/posts`);
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

    async function voteComment(cid: string, v: number) {
      if (!userStore.authenticated) {
        alert("login to vote");
        return;
      }
      await api
        .post(`/comments/${cid}/vote`, { value: v })
        .then((_r) => {
          // optimistic update: find comment and set counts (requires backend to return counts — currently it returns likes/dislikes)
        })
        .catch(() => {});
    }

    onMounted(async () => {
      await fetchPostAndComments();
    });

    return {
      post,
      comments,
      commentText,
      submitComment,
      commentErr,
      vote,
      voteComment,
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
}
</style>
