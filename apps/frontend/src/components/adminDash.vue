<template>
  <div>
    <h2 class="text-xl font-bold mb-3">Admin Dashboard</h2>

    <div>
      <h3 class="font-semibold">Users</h3>
      <button @click="refresh" class="btn mb-2">Refresh</button>
      <div v-for="u in users" :key="u.id" class="p-2 border rounded mb-1">
        <div>
          <strong>{{ u.username }}</strong> ({{ u.role }})
        </div>
        <div class="flex gap-2 mt-1">
          <button
            @click="ban(u.id)"
            class="px-2 py-1 bg-rose-100 rounded text-xs"
          >
            Ban
          </button>
          <button
            @click="unban(u.id)"
            class="px-2 py-1 bg-sky-100 rounded text-xs"
          >
            Unban
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { ref, onMounted } from "vue";
import api from "../services/api";
import { useUserStore } from "../stores/user";

export default {
  setup() {
    const users = ref<any[]>([]);
    const userStore = useUserStore();
    async function refresh() {
      const res = await api.get("/admin/users").catch(() => ({ data: [] }));
      users.value = res.data;
    }
    async function ban(id: string) {
      await api
        .post(`/admin/users/${id}/ban`, {
          duration_minutes: 60,
          reason: "admin",
        })
        .catch(() => {});
      refresh();
    }
    async function unban(id: string) {
      await api.post(`/admin/users/${id}/ban`, { unban: true }).catch(() => {});
      refresh();
    }

    onMounted(() => {
      if (userStore.authenticated) refresh();
    });

    return { users, refresh, ban, unban };
  },
};
</script>

<style scoped>
.btn {
  padding: 6px 10px;
  background: #0ea5e9;
  color: white;
  border-radius: 6px;
  margin-bottom: 8px;
}
</style>
