<template>
  <div class="card">
    <h2 class="text-xl font-bold">{{ user?.username || "Profile" }}</h2>
    <div class="board-meta">id: {{ user?.id }}</div>
    <div class="mt-3">
      <p>Email: {{ user?.email }}</p>
      <p>Role: {{ user?.role }}</p>
      <p>Ink: {{ user?.ink }}</p>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, onMounted } from "vue";
import api from "../services/api";
import { useRoute } from "vue-router";

export default defineComponent({
  setup() {
    const route = useRoute();
    const id = route.params.id as string | undefined;
    const user = ref<any | null>(null);

    async function load() {
      if (!id) return;
      // backend doesn't have GET /users/:id yet — fetch via admin list and filter (fallback)
      const res = await api.get("/admin/users").catch(() => ({ data: [] }));
      user.value = res.data.find((u: any) => u.id === id) || null;
    }

    onMounted(load);
    return { user };
  },
});
</script>
