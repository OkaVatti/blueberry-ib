<template>
  <article class="card flex gap-3">
    <div>
      <div class="avatar" :title="post.user_id">
        <!-- Placeholder avatar initial -->
        <div
          style="
            display: flex;
            align-items: center;
            justify-content: center;
            height: 100%;
            color: var(--color-text);
            font-weight: 700;
          "
        >
          {{ initials }}
        </div>
      </div>
    </div>

    <div class="flex-1">
      <div class="flex items-center justify-between">
        <div>
          <div class="font-semibold">
            {{ post.title || post.user_id || "anon" }}
          </div>
          <div class="board-meta">{{ post.created_at }}</div>
        </div>
        <div class="flex gap-2">
          <button
            class="px-2 py-1 rounded text-sm"
            @click="$emit('like', post.id)"
          >
            👍 {{ post.likes || 0 }}
          </button>
          <button
            class="px-2 py-1 rounded text-sm"
            @click="$emit('dislike', post.id)"
          >
            👎 {{ post.dislikes || 0 }}
          </button>
        </div>
      </div>

      <div
        class="markdown-body mt-2"
        v-html="post.content_html || post.content"
      ></div>
    </div>
  </article>
</template>

<script lang="ts">
import { defineComponent, computed } from "vue";

export default defineComponent({
  props: { post: { type: Object, required: true } },
  setup(props) {
    const initials = computed(() => {
      const id = props.post.user_id || props.post.userId || "A";
      return String((id || "").toString().slice(0, 2)).toUpperCase();
    });
    return { initials };
  },
});
</script>
