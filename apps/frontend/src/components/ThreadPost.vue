<template>
  <article class="card">
    <div class="board-meta mb-1">/{{ post.board?.slug || post.board_id }}</div>
    <div class="font-semibold text-lg">{{ post.title || "Untitled" }}</div>
    <div class="board-meta mb-2">
      by {{ post.user_id || "anon" }} • {{ post.created_at }}
    </div>
    <div class="markdown-body" v-html="post.content_html || post.content"></div>

    <div class="flex items-center gap-3 mt-3">
      <button class="px-2 py-1 rounded bg-ui/5" @click="$emit('like', post.id)">
        👍 {{ post.likes || 0 }}
      </button>
      <button
        class="px-2 py-1 rounded bg-ui/5"
        @click="$emit('dislike', post.id)"
      >
        👎 {{ post.dislikes || 0 }}
      </button>
      <router-link :to="`/posts/${post.id}`" class="ml-auto text-sm board-meta"
        >view post</router-link
      >
    </div>
  </article>
</template>

<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  props: { post: { type: Object, required: true } },
});
</script>
