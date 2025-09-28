<!-- apps/frontend/app/pages/indexHome.vue -->
<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useBoardStore } from "~/stores/boards";
import { usePostStore } from "~/stores/posts";
import { useAuthStore } from "~/stores/auth";

const boardStore = useBoardStore();
const postStore = usePostStore();
const authStore = useAuthStore();

const featuredBoards = ref<any[]>([]);
const popularPosts = ref<any[]>([]);
const loading = ref(true);

onMounted(async () => {
  loading.value = true;
  try {
    // Fetch boards
    await boardStore.listBoards();
    featuredBoards.value = boardStore.boards.slice(0, 6);

    // Fetch popular posts (would need backend endpoint)
    // For now, just use placeholder
    popularPosts.value = [
      {
        id: "1",
        title: "Welcome to Blueberry!",
        content: "A neobrutalist imageboard and microblogging platform.",
        author: { username: "admin", displayName: "Administrator" },
        board: { name: "Random", slug: "b" },
        likes: 42,
        comments: 7,
        createdAt: new Date().toISOString(),
      },
    ];
  } finally {
    loading.value = false;
  }
});

function navigateToBoard(slug: string) {
  navigateTo(`/${slug}`);
}

function createPost() {
  if (!authStore.user) {
    navigateTo("/auth/login");
  } else {
    // Open post creation modal
    console.log("Open post modal");
  }
}
</script>

<template>
  <div class="home-container">
    <!-- Hero Section -->
    <div class="hero nb-panel mb-6">
      <h1 class="text-4xl font-black mb-2">Welcome to Blueberry 🫐</h1>
      <p class="text-lg muted mb-4">
        An experimental neobrutalist imageboard + microblogging platform
      </p>
      <div class="flex gap-3">
        <button class="nb-btn nb-btn--primary" @click="createPost">
          Create Post
        </button>
        <button class="nb-btn" @click="navigateTo('/boards')">
          Explore Boards
        </button>
      </div>
    </div>

    <!-- Featured Boards -->
    <section class="mb-8">
      <h2 class="text-2xl font-bold mb-4">Featured Boards</h2>
      <div v-if="loading" class="text-center py-8">
        <div class="kv">Loading boards...</div>
      </div>
      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <div
          v-for="board in featuredBoards"
          :key="board.id"
          class="board-card nb-panel nb-panel--thin cursor-pointer hover:transform hover:translate-x-1"
          @click="navigateToBoard(board.slug)"
        >
          <div class="flex items-center gap-3 mb-2">
            <div class="board-icon">
              {{ board.name[0] }}
            </div>
            <div>
              <h3 class="font-bold">{{ board.name }}</h3>
              <div class="text-sm muted">/{{ board.slug }}/</div>
            </div>
          </div>
          <p class="text-sm muted mb-2">{{ board.description }}</p>
          <div class="flex gap-4 text-xs kv">
            <span>{{ board.postCount || 0 }} posts</span>
            <span>{{ board.memberCount || 0 }} members</span>
          </div>
        </div>
      </div>
    </section>

    <!-- Popular Posts -->
    <section class="mb-8">
      <h2 class="text-2xl font-bold mb-4">Popular Posts</h2>
      <div class="space-y-4">
        <div v-for="post in popularPosts" :key="post.id" class="post">
          <div class="post-meta">
            <div class="flex items-center gap-3">
              <div class="avatar">
                {{ post.author.displayName[0] }}
              </div>
              <div>
                <div class="font-semibold">{{ post.author.displayName }}</div>
                <div class="text-xs muted">
                  @{{ post.author.username }} · /{{ post.board.slug }}/ ·
                  {{ new Date(post.createdAt).toLocaleDateString() }}
                </div>
              </div>
            </div>
          </div>
          <div class="post-content">
            <h3 v-if="post.title" class="font-bold mb-1">{{ post.title }}</h3>
            <p>{{ post.content }}</p>
          </div>
          <div class="post-actions">
            <button class="nb-btn nb-btn--small">👍 {{ post.likes }}</button>
            <button class="nb-btn nb-btn--small">💬 {{ post.comments }}</button>
            <button class="nb-btn nb-btn--small">🔁 Repost</button>
          </div>
        </div>
      </div>
    </section>

    <!-- Announcements -->
    <section class="mb-8">
      <div class="nb-panel" style="background: var(--accent); color: white">
        <h3 class="font-bold mb-2">📢 Site Announcement</h3>
        <p>
          Welcome to the Blueberry alpha! We're still building features and
          squashing bugs. Report issues and join the discussion in /tech/.
        </p>
      </div>
    </section>
  </div>
</template>

<style scoped>
.hero {
  background: linear-gradient(
    135deg,
    var(--panel) 0%,
    rgba(30, 144, 255, 0.05) 100%
  );
}

.board-icon {
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--accent);
  color: white;
  border-radius: 8px;
  font-weight: 900;
  font-size: 1.25rem;
  box-shadow: 4px 4px 0 0 var(--muted);
}

.board-card:hover .board-icon {
  box-shadow: 6px 6px 0 0 var(--muted);
}

.nb-btn--small {
  padding: 0.25rem 0.5rem;
  font-size: 0.875rem;
}
</style>
