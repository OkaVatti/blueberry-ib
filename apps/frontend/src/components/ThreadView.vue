<template>
  <div>
    <div v-if="thread">
      <h1 class="text-2xl font-bold mb-2">{{ thread.title }}</h1>
      <div class="mb-4">
        <button @click="like" class="px-2 py-1 bg-sky-100 rounded">Like</button>
        <button @click="dislike" class="px-2 py-1 bg-rose-100 rounded">
          Dislike
        </button>
      </div>
      <div class="mb-4">
        <h3 class="font-semibold">Posts in thread</h3>
        <div v-for="p in posts" :key="p.id" class="p-3 border rounded mb-2">
          <div class="font-semibold">{{ p.title }}</div>
          <div v-html="p.content_html || p.content" class="mt-1"></div>
          <div class="text-xs text-slate-500 mt-2">
            Likes: {{ p.likes || 0 }} • Dislikes: {{ p.dislikes || 0 }}
          </div>
        </div>
      </div>

      <div class="mb-4">
        <h4 class="font-semibold">Add a reply</h4>
        <textarea v-model="replyContent" class="input mb-2" rows="4"></textarea>
        <div><button @click="submitReply" class="btn">Reply</button></div>
        <div v-if="err" class="text-rose-600 text-sm mt-2">{{ err }}</div>
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
    const thread = ref<any | null>(null);
    const posts = ref<any[]>([]);
    const replyContent = ref("");
    const err = ref<string | null>(null);
    const userStore = useUserStore();

    async function fetchThread() {
      const t = await api.get(`/threads/${id}`);
      thread.value = t.data;
      const res = await api.get(`/threads/${id}/posts`);
      posts.value = res.data;
    }

    async function submitReply() {
      err.value = null;
      if (!userStore.authenticated) {
        err.value = "Login to reply";
        return;
      }
      try {
        const res = await api.post(`/threads/${id}/posts`, {
          title: "",
          content: replyContent.value,
        });
        posts.value.push(res.data);
        replyContent.value = "";
      } catch (e: any) {
        err.value = e?.response?.data?.error || "failed to reply";
      }
    }

    async function like() {
      if (!userStore.authenticated) {
        alert("login to vote");
        return;
      }
      await api.post(`/threads/${id}/vote`, { value: 1 }).catch(() => {});
    }
    async function dislike() {
      if (!userStore.authenticated) {
        alert("login to vote");
        return;
      }
      await api.post(`/threads/${id}/vote`, { value: -1 }).catch(() => {});
    }

    onMounted(() => {
      fetchThread();
    });

    return { thread, posts, replyContent, submitReply, err, like, dislike };
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
