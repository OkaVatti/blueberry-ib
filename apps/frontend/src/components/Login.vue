<template>
  <div class="max-w-md mx-auto bg-white p-4 rounded shadow">
    <h2 class="text-xl font-semibold mb-2">Login / Register</h2>

    <div class="grid gap-2">
      <input
        v-model="usernameOrEmail"
        placeholder="username or email"
        class="input"
      />
      <input
        v-model="password"
        placeholder="password"
        type="password"
        class="input"
      />
      <div class="flex gap-2">
        <button @click="doLogin" class="btn">Login</button>
        <button @click="doRegister" class="btn-ghost">Register</button>
      </div>
      <div v-if="err" class="text-sm text-rose-600">{{ err }}</div>
    </div>
  </div>
</template>

<script lang="ts">
import { ref } from "vue";
import { useUserStore } from "../stores/user";
import { useRouter } from "vue-router";

export default {
  setup() {
    const usernameOrEmail = ref("");
    const password = ref("");
    const err = ref<string | null>(null);
    const userStore = useUserStore();
    const router = useRouter();

    async function doLogin() {
      err.value = null;
      try {
        await userStore.login(usernameOrEmail.value, password.value);
        router.push("/");
      } catch (e: any) {
        err.value = e?.response?.data?.error || "login failed";
      }
    }

    async function doRegister() {
      err.value = null;
      try {
        // use same value for email if user just types username
        await userStore.register(
          usernameOrEmail.value,
          usernameOrEmail.value,
          password.value
        );
        router.push("/");
      } catch (e: any) {
        err.value = e?.response?.data?.error || "register failed";
      }
    }

    return { usernameOrEmail, password, doLogin, doRegister, err };
  },
};
</script>

<style scoped>
.input {
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 6px;
}
.btn {
  padding: 8px 12px;
  background: #0ea5e9;
  color: white;
  border-radius: 6px;
}
.btn-ghost {
  padding: 8px 12px;
  border-radius: 6px;
  background: #f1f5f9;
}
</style>
