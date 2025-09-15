<template>
  <div class="p-3 border rounded mb-3">
    <input v-model="title" placeholder="Thread title" class="input mb-2" />
    <input
      v-model="postTitle"
      placeholder="Initial post title (optional)"
      class="input mb-2"
    />
    <textarea
      v-model="content"
      placeholder="Content (markdown supported)"
      class="input mb-2"
      rows="4"
    ></textarea>
    <div class="flex gap-2">
      <button @click="submit" class="btn">Create Thread</button>
      <div v-if="err" class="text-rose-600 text-sm">{{ err }}</div>
    </div>
  </div>
</template>

<script lang="ts">
import { ref } from "vue";
import api from "../services/api";
import { useUserStore } from "../stores/user";

export default {
  props: {
    boardSlug: { type: String, required: true },
  },
  emits: ["created"],
  setup(props, { emit }) {
    const title = ref("");
    const postTitle = ref("");
    const content = ref("");
    const err = ref<string | null>(null);
    const userStore = useUserStore();

    async function submit() {
      err.value = null;
      if (!userStore.authenticated) {
        err.value = "Login to create a thread";
        return;
      }
      try {
        const res = await api.post(`/boards/${props.boardSlug}/threads`, {
          title: title.value,
          post_title: postTitle.value,
          content: content.value,
        });
        title.value = "";
        postTitle.value = "";
        content.value = "";
        emit("created", res.data.thread);
      } catch (e: any) {
        err.value = e?.response?.data?.error || "failed to create thread";
      }
    }

    return { title, postTitle, content, submit, err };
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
  background: #0ea5e9;
  color: white;
  border-radius: 6px;
}
</style>
