<!-- apps/frontend/src/pages/HomeIndex.vue -->
<template>
  <MainLayout />
</template>

<script lang="ts">
import { defineComponent } from "vue";
import MainLayout from "../components/MainLayout.vue";

export default defineComponent({
  name: "HomeIndex",
  components: { MainLayout }
});
</script>