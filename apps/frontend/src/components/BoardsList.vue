<template>
  <div>
    <h1 class="text-2xl font-bold mb-3">Boards</h1>

    <div class="grid gap-3">
      <div
        v-for="b in boards"
        :key="b.id"
        class="p-3 border rounded flex items-start justify-between"
      >
        <div>
          <router-link
            :to="`/boards/${b.slug}`"
            class="font-semibold text-lg"
            >{{ b.name }}</router-link
          >
          <div class="text-sm text-slate-500">
            /{{ b.slug }} • {{ b.description }}
          </div>
        </div>
        <div>
          <button
            @click="open(b.slug)"
            class="px-2 py-1 text-sm rounded bg-slate-100"
          >
            Open
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { ref, onMounted } from "vue";
import api from "../services/api";
import { useRouter } from "vue-router";

export default {
  setup() {
    const boards = ref<any[]>([]);
    const router = useRouter();

    async function fetchBoards() {
      const res = await api.get("/boards");
      boards.value = res.data;
    }

    function open(slug: string) {
      router.push(`/boards/${slug}`);
    }

    onMounted(() => {
      fetchBoards();
    });

    return { boards, open };
  },
};
</script>
