<template>
  <div>
    <Navbar />
    <div class="container-grid">
      <aside>
        <div class="card">
          <h4 class="font-semibold mb-2">Dashboards</h4>
          <router-link to="/dash/owner" class="block p-2 rounded hover:bg-ui/5"
            >Owner</router-link
          >
          <router-link
            to="/dash/coowner"
            class="block p-2 rounded hover:bg-ui/5"
            >Co-Owner</router-link
          >
          <router-link to="/dash/admin" class="block p-2 rounded hover:bg-ui/5"
            >Admin</router-link
          >
          <router-link to="/dash/mod" class="block p-2 rounded hover:bg-ui/5"
            >Moderator</router-link
          >
        </div>
      </aside>

      <main>
        <div class="card mb-4">
          <h2 class="text-xl font-bold">Management</h2>
          <div class="board-meta">tools for moderators / admins / owners</div>
        </div>

        <router-view />
      </main>

      <aside>
        <div class="card">
          <h4 class="font-semibold mb-2">Quick Actions</h4>
          <div class="board-meta">moderation tools, reports, etc.</div>
        </div>
      </aside>
    </div>
  </div>
</template>

<script lang="ts">
import Navbar from "../components/Navbar.vue";
export default { components: { Navbar } };
</script>

<style scoped></style>
