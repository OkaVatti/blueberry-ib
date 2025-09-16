<template>
  <div class="card">
    <div class="flex items-center justify-between mb-3">
      <h2 class="text-xl font-bold">Admin Dashboard</h2>
      <div class="board-meta">manage users & moderation</div>
    </div>

    <div class="grid grid-cols-2 gap-4">
      <section>
        <h3 class="font-semibold mb-2">Users</h3>
        <div class="space-y-2">
          <div
            v-for="u in users"
            :key="u.id"
            class="p-3 border rounded flex items-center justify-between"
          >
            <div>
              <div class="font-semibold">{{ u.username }}</div>
              <div class="board-meta">{{ u.email }} • role: {{ u.role }}</div>
            </div>
            <div class="flex gap-2">
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
      </section>

      <section>
        <h3 class="font-semibold mb-2">Audit Logs</h3>
        <div class="space-y-2 text-xs">
          <div v-for="l in logs" :key="l.id" class="p-2 border rounded">
            <div class="font-semibold">{{ l.action }}</div>
            <div class="board-meta">
              actor {{ l.actor_id || "system" }} • {{ l.created_at }}
            </div>
            <div class="text-sm mt-1">{{ l.details }}</div>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<script lang="ts">
import { ref, onMounted } from "vue";
import api from "../services/api";

export default {
  setup() {
    const users = ref<any[]>([]);
    const logs = ref<any[]>([]);
    async function refresh() {
      users.value = (
        await api.get("/admin/users").catch(() => ({ data: [] }))
      ).data;
      logs.value = (
        await api.get("/admin/logs").catch(() => ({ data: [] }))
      ).data;
    }
    async function ban(id: string) {
      await api
        .post(`/admin/users/${id}/ban`, {
          duration_minutes: 60,
          reason: "admin action",
        })
        .catch(() => {});
      refresh();
    }
    async function unban(id: string) {
      await api.post(`/admin/users/${id}/ban`, { unban: true }).catch(() => {});
      refresh();
    }
    onMounted(refresh);
    return { users, logs, refresh, ban, unban };
  },
};
</script>
