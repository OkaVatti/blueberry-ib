<script setup lang="ts">
import { ref, watch, onUnmounted, defineEmits } from "vue";
const emit = defineEmits<{
  (e: "search", q: string): void;
}>();

const q = ref("");
let t: ReturnType<typeof setTimeout> | null = null;

function triggerSearch() {
  emit("search", q.value);
}

/* simple debounce */
watch(q, (val) => {
  if (t) clearTimeout(t);
  t = setTimeout(() => {
    emit("search", val);
    t = null;
  }, 260);
});

onUnmounted(() => {
  if (t) clearTimeout(t);
});
function clearQ() {
  q.value = "";
  emit("search", "");
}
</script>

<template>
  <div class="nb-search relative">
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width="14"
      height="14"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      stroke-width="2"
      class="opacity-80"
    >
      <circle cx="11" cy="11" r="7" />
      <line x1="21" y1="21" x2="16.65" y2="16.65" />
    </svg>

    <input
      v-model="q"
      placeholder="Search posts, boards, tags..."
      class="ml-2 outline-none bg-transparent"
      @keyup.enter="triggerSearch"
    />

    <button v-if="q" class="kv text-xs ml-2" @click="clearQ">×</button>
  </div>
</template>

<style scoped>
.nb-search {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.35rem 0.6rem;
  border: var(--border-light) solid var(--muted);
  border-radius: 8px;
  background: var(--panel);
}
.nb-search input {
  min-width: 220px;
  font-weight: 600;
}
@media (max-width: 640px) {
  .nb-search input {
    min-width: 120px;
  }
}
</style>
