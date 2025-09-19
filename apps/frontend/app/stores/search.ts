// /apps/frontend/stores/search.ts
import { defineStore } from "pinia";
import { ref, watch } from "vue";

export const useSearchStore = defineStore("search", () => {
  const config = useRuntimeConfig();
  const apiBase = config.public.apiBase ?? "/api";

  const query = ref("");
  const results = ref<any[]>([]);
  const loading = ref(false);
  const error = ref<string | null>(null);

  let debounceTimer: number | null = null;

  watch(query, (q) => {
    if (debounceTimer) window.clearTimeout(debounceTimer);
    debounceTimer = window.setTimeout(() => {
      run(q);
    }, 260);
  });

  async function run(q?: string) {
    const term = q ?? query.value;
    if (!term || term.trim().length === 0) {
      results.value = [];
      return { ok: true, data: [] };
    }
    loading.value = true;
    try {
      const res = await $fetch(
        `${apiBase}/search?q=${encodeURIComponent(term)}`
      );
      results.value = res as any[];
      return { ok: true, data: res };
    } catch (err: any) {
      error.value = err?.message ?? "search failed";
      return { ok: false, error: error.value };
    } finally {
      loading.value = false;
    }
  }

  return {
    query,
    results,
    loading,
    error,
    run,
  };
});
