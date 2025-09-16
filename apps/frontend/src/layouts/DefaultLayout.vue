<template>
  <div>
    <Navbar />
    <div class="container-grid">
      <aside>
        <!-- Left sidebar (boards quick list) -->
        <div class="card">
          <h4 class="font-semibold mb-2">Boards</h4>
          <div v-for="b in boards" :key="b.id" class="mb-2">
            <router-link
              :to="`/boards/${b.slug}`"
              class="block p-2 rounded hover:bg-ui/5"
            >
              <div class="font-semibold">{{ b.name }}</div>
              <div class="board-meta">/{{ b.slug }}</div>
            </router-link>
          </div>
        </div>
      </aside>

      <main>
        <!-- page content -->
        <router-view />
      </main>

      <aside>
        <!-- Right sidebar: quick threads / trending -->
        <div class="card">
          <h4 class="font-semibold mb-2">Trending</h4>
          <div v-for="t in threads" :key="t.id" class="mb-2">
            <a
              @click.prevent="gotoThread(t.id)"
              href="#"
              class="block p-2 rounded hover:bg-ui/5"
            >
              <div class="font-semibold text-sm">{{ t.title }}</div>
              <div class="board-meta text-xs">by {{ t.user_id || "anon" }}</div>
            </a>
          </div>
        </div>
      </aside>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, onMounted } from "vue";
import Navbar from "../components/Navbar.vue";
import api from "../services/api";
import { useRouter } from "vue-router";

export default defineComponent({
  components: { Navbar },
  setup() {
    const boards = ref<any[]>([]);
    const threads = ref<any[]>([]);
    const router = useRouter();

    async function load() {
      boards.value = (
        await api.get("/boards").catch(() => ({ data: [] }))
      ).data;
      if (boards.value.length) {
        threads.value = (
          await api
            .get(`/boards/${boards.value[0].slug}/threads`)
            .catch(() => ({ data: [] }))
        ).data.slice(0, 8);
      }
    }

    function gotoThread(id: string) {
      router.push(`/threads/${id}`);
    }

    onMounted(load);
    return { boards, threads, gotoThread };
  },
});
</script>

<style scoped>
main {
  min-height: 60vh;
}
</style>