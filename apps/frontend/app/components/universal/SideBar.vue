<script setup lang="ts">
import { ref } from "vue";

const user = ref({
  name: "Lily",
  username: "lily_dev",
  avatar: "L",
  ink: 420,
  roles: ["admin", "moderator"], // example roles
});

const boards = ref([
  { id: "g", name: "General" },
  { id: "art", name: "Art" },
  { id: "tech", name: "Technology" },
  { id: "weird", name: "Weird" },
]);

function createBoard() {
  /* open modal / navigate */ console.log("create board");
}
function openBoard(id: string) {
  console.log("open", id);
}
</script>

<template>
  <div class="flex flex-col gap-4">
    <div class="flex items-center gap-3">
      <div class="avatar">{{ user.avatar }}</div>
      <div>
        <div class="font-semibold">{{ user.name }}</div>
        <div class="text-xs muted">@{{ user.username }}</div>
        <div class="text-xs mt-1 kv">Ink: {{ user.ink }}</div>
      </div>
    </div>

    <div class="nb-panel--thin">
      <div class="flex items-center justify-between">
        <div class="text-sm font-bold">Your Boards</div>
        <button class="nb-btn" @click="createBoard">New</button>
      </div>

      <div class="mt-3 space-y-2">
        <div
          v-for="b in boards"
          :key="b.id"
          class="board-item"
          @click="() => openBoard(b.id)"
        >
          <div
            class="kv w-9 h-9 flex-center"
            style="background: var(--muted); border-radius: 6px"
          >
            {{ b.name[0] }}
          </div>
          <div>
            <div class="font-semibold">{{ b.name }}</div>
            <div class="text-xs muted">Public</div>
          </div>
        </div>
      </div>
    </div>

    <div class="nb-panel">
      <div class="font-semibold">Settings</div>
      <ul class="mt-2 space-y-2 text-sm">
        <li class="nb-navitem">Profile</li>
        <li class="nb-navitem">Appearance</li>
        <li class="nb-navitem">Notifications</li>
      </ul>
    </div>

    <div class="mt-auto text-xs muted">
      <div>Version: alpha</div>
      <div class="mt-2">Built with ♥ for weird internet culture</div>
    </div>
  </div>
</template>
