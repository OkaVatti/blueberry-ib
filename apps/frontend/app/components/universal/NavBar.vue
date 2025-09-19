<script setup lang="ts">
import { defineEmits } from "vue";
const emit = defineEmits<{
  (e: "toggle-sidebar"): void;
}>();

function onToggle() {
  emit("toggle-sidebar");
}
</script>

<template>
  <div class="w-full flex items-center justify-between">
    <div class="flex items-center gap-4">
      <button class="nb-btn" @click="onToggle">☰</button>
      <div>
        <div class="kv text-lg">Blueberry</div>
        <div class="text-xs muted">neobrutalist imageboard · microblog</div>
      </div>
    </div>

    <div class="flex items-center gap-3">
      <slot />
      <button class="nb-btn">Explore</button>
      <button class="nb-btn nb-btn--primary">New Post</button>
      <div class="ml-2">
        <div class="avatar">L</div>
      </div>
    </div>
  </div>
</template>
