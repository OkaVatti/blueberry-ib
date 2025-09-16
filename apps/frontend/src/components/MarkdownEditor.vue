<template>
  <div class="card">
    <div class="flex items-center justify-between mb-2">
      <h3 class="font-semibold">Compose (Markdown)</h3>
      <div class="flex items-center gap-2">
        <button @click="togglePreview" class="px-2 py-1 rounded border text-sm">
          {{ showPreview ? "Edit" : "Preview" }}
        </button>
        <button
          @click="$emit('submit', content)"
          class="px-3 py-1 rounded btn-primary"
        >
          Post
        </button>
      </div>
    </div>

    <div v-if="!showPreview">
      <input
        v-model="title"
        placeholder="Title (optional)"
        class="input mb-2"
      />
      <textarea
        v-model="content"
        rows="8"
        class="input w-full"
        placeholder="Write Markdown..."
      ></textarea>
    </div>

    <div v-else>
      <div
        class="markdown-body"
        v-html="previewHtml || '<em>Loading preview...</em>'"
      ></div>
    </div>

    <div class="text-xs text-muted mt-2">
      Server-rendered preview (sanitized)
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, watch } from "vue";
import api from "../services/api";
import { debounce } from "lodash-es";

export default defineComponent({
  emits: ["submit"],
  setup() {
    const content = ref("");
    const title = ref("");
    const showPreview = ref(false);
    const previewHtml = ref<string | null>(null);

    const fetchPreview = debounce(async (md: string) => {
      if (!md) {
        previewHtml.value = "<em>Nothing to preview</em>";
        return;
      }
      try {
        const res = await api.post("/render", { markdown: md });
        previewHtml.value = res.data.html;
      } catch {
        previewHtml.value = "<em>Preview failed</em>";
      }
    }, 300);

    watch(content, (v) => {
      if (showPreview.value) fetchPreview(v);
    });

    function togglePreview() {
      showPreview.value = !showPreview.value;
      if (showPreview.value) {
        fetchPreview(content.value);
      }
    }

    return { content, title, showPreview, togglePreview, previewHtml };
  },
});
</script>

<style scoped>
.input {
  padding: 8px;
  border-radius: 8px;
  border: 1px solid rgba(255, 255, 255, 0.04);
  width: 100%;
}
</style>
