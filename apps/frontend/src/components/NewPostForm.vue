<template>
  <div class="p-3 border rounded">
    <input v-model="title" placeholder="Title" class="input mb-2" />
    <textarea
      v-model="content"
      placeholder="Content"
      class="input mb-2"
      rows="4"
    ></textarea>
    <div class="flex gap-2">
      <button @click="submit" class="btn">Create Post</button>
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
    const content = ref("");
    const err = ref<string | null>(null);
    const userStore = useUserStore();

    async function submit() {
      err.value = null;
      if (!userStore.authenticated) {
        err.value = "You must be logged in to post";
        return;
      }
      try {
        const res = await api.post(`/boards/${props.boardSlug}/posts`, {
          title: title.value,
          content: content.value,
        });
        title.value = "";
        content.value = "";
        emit("created", res.data);
      } catch (e: any) {
        err.value = e?.response?.data?.error || "failed to create post";
      }
    }

    return { title, content, submit, err };
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
