<!-- apps/frontend/app/pages/auth/login.vue -->
<script setup lang="ts">
import { ref } from "vue";
import { useAuthStore } from "~/stores/auth";

const authStore = useAuthStore();
const mode = ref<"login" | "register">("login");
const error = ref("");
const loading = ref(false);

const loginForm = ref({
  usernameOrEmail: "",
  password: "",
});

const registerForm = ref({
  username: "",
  email: "",
  password: "",
  confirmPassword: "",
  displayName: "",
});

async function handleLogin() {
  error.value = "";
  loading.value = true;

  try {
    const result = await authStore.login(
      loginForm.value.usernameOrEmail,
      loginForm.value.password
    );

    if (result.ok) {
      navigateTo("/");
    } else {
      error.value = result.error?.message || "Login failed";
    }
  } catch (e) {
    error.value = "An unexpected error occurred";
  } finally {
    loading.value = false;
  }
}

async function handleRegister() {
  error.value = "";

  if (registerForm.value.password !== registerForm.value.confirmPassword) {
    error.value = "Passwords do not match";
    return;
  }

  if (registerForm.value.password.length < 6) {
    error.value = "Password must be at least 6 characters";
    return;
  }

  loading.value = true;

  try {
    const result = await authStore.signup({
      username: registerForm.value.username,
      email: registerForm.value.email,
      password: registerForm.value.password,
      displayName:
        registerForm.value.displayName || registerForm.value.username,
    });

    if (result.ok) {
      navigateTo("/");
    } else {
      error.value = result.error?.message || "Registration failed";
    }
  } catch (e) {
    error.value = "An unexpected error occurred";
  } finally {
    loading.value = false;
  }
}

function switchMode() {
  mode.value = mode.value === "login" ? "register" : "login";
  error.value = "";
}
</script>

<template>
  <div class="auth-container">
    <div class="auth-box nb-panel">
      <div class="text-center mb-6">
        <h1 class="text-3xl font-black mb-2">🫐 Blueberry</h1>
        <p class="muted">
          {{ mode === "login" ? "Welcome back!" : "Join the community" }}
        </p>
      </div>

      <!-- Error Message -->
      <div v-if="error" class="error-box mb-4">
        {{ error }}
      </div>

      <!-- Login Form -->
      <form v-if="mode === 'login'" @submit.prevent="handleLogin">
        <div class="mb-4">
          <label class="label">Username or Email</label>
          <input
            v-model="loginForm.usernameOrEmail"
            type="text"
            class="input-field"
            required
            :disabled="loading"
          />
        </div>

        <div class="mb-4">
          <label class="label">Password</label>
          <input
            v-model="loginForm.password"
            type="password"
            class="input-field"
            required
            :disabled="loading"
          />
        </div>

        <button
          type="submit"
          class="nb-btn nb-btn--primary w-full mb-4"
          :disabled="loading"
        >
          {{ loading ? "Logging in..." : "Login" }}
        </button>

        <p class="text-center text-sm">
          Don't have an account?
          <button type="button" class="link" @click="switchMode">
            Register here
          </button>
        </p>
      </form>

      <!-- Register Form -->
      <form v-else @submit.prevent="handleRegister">
        <div class="mb-4">
          <label class="label">Username</label>
          <input
            v-model="registerForm.username"
            type="text"
            class="input-field"
            pattern="[a-zA-Z0-9_]{3,20}"
            title="3-20 characters, letters, numbers and underscores only"
            required
            :disabled="loading"
          />
          <div class="text-xs muted mt-1">
            3-20 characters, alphanumeric and underscores
          </div>
        </div>

        <div class="mb-4">
          <label class="label">Display Name (optional)</label>
          <input
            v-model="registerForm.displayName"
            type="text"
            class="input-field"
            :disabled="loading"
          />
        </div>

        <div class="mb-4">
          <label class="label">Email</label>
          <input
            v-model="registerForm.email"
            type="email"
            class="input-field"
            required
            :disabled="loading"
          />
        </div>

        <div class="mb-4">
          <label class="label">Password</label>
          <input
            v-model="registerForm.password"
            type="password"
            class="input-field"
            minlength="6"
            required
            :disabled="loading"
          />
          <div class="text-xs muted mt-1">Minimum 6 characters</div>
        </div>

        <div class="mb-4">
          <label class="label">Confirm Password</label>
          <input
            v-model="registerForm.confirmPassword"
            type="password"
            class="input-field"
            required
            :disabled="loading"
          />
        </div>

        <button
          type="submit"
          class="nb-btn nb-btn--primary w-full mb-4"
          :disabled="loading"
        >
          {{ loading ? "Creating account..." : "Register" }}
        </button>

        <p class="text-center text-sm">
          Already have an account?
          <button type="button" class="link" @click="switchMode">
            Login here
          </button>
        </p>
      </form>
    </div>

    <div class="text-center mt-6 text-sm muted">
      <p>By using Blueberry, you agree to our community guidelines.</p>
      <p class="mt-2">Be kind, be weird, be yourself. 💜</p>
    </div>
  </div>
</template>

<style scoped>
.auth-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 2rem;
  background: linear-gradient(
    135deg,
    var(--bg) 0%,
    rgba(30, 144, 255, 0.05) 100%
  );
}

.auth-box {
  width: 100%;
  max-width: 420px;
}

.label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 600;
  font-size: 0.875rem;
}

.input-field {
  width: 100%;
  padding: 0.75rem;
  border: var(--border-light) solid var(--muted);
  border-radius: 6px;
  font-family: inherit;
  background: var(--panel);
  color: var(--ink);
  transition: border-color 0.2s;
}

.input-field:focus {
  outline: none;
  border-color: var(--accent);
  box-shadow: 0 0 0 3px rgba(30, 144, 255, 0.1);
}

.input-field:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.w-full {
  width: 100%;
}

.error-box {
  padding: 0.75rem;
  background: rgba(239, 68, 68, 0.1);
  border: var(--border-light) solid var(--danger);
  border-radius: 6px;
  color: var(--danger);
  font-weight: 600;
  font-size: 0.875rem;
}

.link {
  color: var(--accent);
  font-weight: 600;
  text-decoration: underline;
  background: none;
  border: none;
  cursor: pointer;
}

.link:hover {
  opacity: 0.8;
}
</style>
