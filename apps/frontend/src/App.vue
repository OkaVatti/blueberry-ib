<template>
  <div style="font-family: Inter, system-ui, -apple-system; padding: 24px">
    <h1>Blueberry — Prototype</h1>

    <div style="display: flex; gap: 20px">
      <div style="width: 320px">
        <Login @login="onLogin" />
      </div>

      <div style="flex: 1">
        <BoardsList :token="token" />
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from "vue";
import Login from "./components/Login.vue";
import BoardsList from "./components/BoardsList.vue";

export default defineComponent({
  components: { Login, BoardsList },
  setup() {
    const token = ref<string | null>(null);
    function onLogin(t: string) {
      token.value = t;
    }
    return { token, onLogin };
  },
});
</script>
