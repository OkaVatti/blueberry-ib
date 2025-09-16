<template>
  <div>
    <BoardView :slug="slug" />
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import BoardView from "../components/BoardView.vue";
import { useRoute } from "vue-router";

export default defineComponent({
  components: { BoardView },
  setup() {
    const route = useRoute();
    return { slug: route.params.slug as string };
  },
});
</script>
