<template>
  <div>
    <h2 class="text-xl font-bold mb-3">Owner Dashboard</h2>
    <div class="grid grid-cols-2 gap-4">
      <div>
        <h3 class="font-semibold">Users</h3>
        <button @click="refreshUsers" class="btn mb-2">Refresh</button>
        <div v-for="u in users" :key="u.id" class="p-2 border rounded mb-1">
          <div>
            <strong>{{ u.username }}</strong> ({{ u.role }})
          </div>
          <div class="text-sm">Email: {{ u.email }}</div>
          <div class="flex gap-2 mt-2">
            <button
              @click="banUser(u.id)"
              class="px-2 py-1 bg-rose-100 rounded text-xs"
            >
              Ban
            </button>
            <button
              @click="unbanUser(u.id)"
              class="px-2 py-1 bg-sky-100 rounded text-xs"
            >
              Unban
            </button>
            <button
              @click="deleteUser(u.id)"
              class="px-2 py-1 bg-gray-200 rounded text-xs"
            >
              Delete
            </button>
          </div>
        </div>
      </div>

      <div>
        <h3 class="font-semibold">Audit Logs</h3>
        <button @click="refreshLogs" class="btn mb-2">Refresh</button>
        <div
          v-for="l in logs"
          :key="l.id"
          class="p-2 border rounded mb-1 text-xs"
        >
          <div>
            <strong>{{ l.action }}</strong> by {{ l.actor_id || "system" }}
          </div>
          <div>
            {{ l.created_at }} • target: {{ l.target_type }} {{ l.target_id }}
          </div>
          <div class="text-sm">{{ l.details }}</div>
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
    const userStore = useUserStore();
    const users = ref<any[]>([]);
    const logs = ref<any[]>([]);

    async function refreshUsers() {
      const res = await api.get("/admin/users").catch(() => ({ data: [] }));
      users.value = res.data;
    }
    async function refreshLogs() {
      const res = await api.get("/admin/logs").catch(() => ({ data: [] }));
      logs.value = res.data;
    }
    async function banUser(id: string) {
      await api
        .post(`/admin/users/${id}/ban`, {
          duration_minutes: 0,
          reason: "owner action",
        })
        .catch(() => {});
      await refreshUsers();
      await refreshLogs();
    }
    async function unbanUser(id: string) {
      await api.post(`/admin/users/${id}/ban`, { unban: true }).catch(() => {});
      await refreshUsers();
      await refreshLogs();
    }
    async function deleteUser(id: string) {
      // reuse admin delete post? here we simply set banned + audit; if you want full delete endpoint, add server route
      await api
        .post(`/admin/users/${id}/ban`, {
          duration_minutes: 0,
          reason: "owner delete",
        })
        .catch(() => {});
      await refreshUsers();
      await refreshLogs();
    }

    onMounted(() => {
      if (!userStore.authenticated) return;
      refreshUsers();
      refreshLogs();
    });

    return {
      users,
      logs,
      refreshUsers,
      refreshLogs,
      banUser,
      unbanUser,
      deleteUser,
    };
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
