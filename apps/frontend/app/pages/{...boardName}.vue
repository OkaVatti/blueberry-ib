<!-- apps/frontend/app/pages/{...boardName}.vue -->
<script setup lang="ts">
import { ref, onMounted, computed } from "vue";
import { useRoute } from "vue-router";
import { useBoardStore } from "~/stores/boards";
import { usePostStore } from "~/stores/posts";
import { useThreadStore } from "~/stores/threads";
import { useAuthStore } from "~/stores/auth";

const route = useRoute();
const boardStore = useBoardStore();
const postStore = usePostStore();
const threadStore = useThreadStore();
const authStore = useAuthStore();

const boardSlug = computed(() => route.params.boardName as string);
const board = ref<any>(null);
const posts = ref<any[]>([]);
const threads = ref<any[]>([]);
const activeTab = ref<"posts" | "threads">("posts");
const loading = ref(true);
const showCreateModal = ref(false);
const newPost = ref({ title: "", content: "", isThread: false });

onMounted(async () => {
  await loadBoard();
});

async function loadBoard() {
  loading.value = true;
  try {
    const result = await boardStore.getBoard(boardSlug.value);
    if (result.ok) {
      board.value = boardStore.current;
      await loadContent();
    }
  } catch (error) {
    console.error("Failed to load board:", error);
  } finally {
    loading.value = false;
  }
}

async function loadContent() {
  if (activeTab.value === "posts") {
    await loadPosts();
  } else {
    await loadThreads();
  }
}

async function loadPosts() {
  try {
    const response = await $fetch(`/api/v1/boards/${boardSlug.value}/posts`);
    posts.value = response;
  } catch (error) {
    console.error("Failed to load posts:", error);
  }
}

async function loadThreads() {
  try {
    const result = await threadStore.listBoardThreads(board.value.id);
    if (result.ok) {
      threads.value = threadStore.summaries;
    }
  } catch (error) {
    console.error("Failed to load threads:", error);
  }
}

function switchTab(tab: "posts" | "threads") {
  activeTab.value = tab;
  loadContent();
}

function openCreateModal() {
  if (!authStore.user) {
    navigateTo("/auth/login");
    return;
  }
  showCreateModal.value = true;
}

async function createContent() {
  if (!newPost.value.content.trim()) return;

  try {
    if (newPost.value.isThread) {
      await threadStore.createThread({
        boardId: board.value.id,
        title: newPost.value.title,
        body: newPost.value.content,
      });
    } else {
      await postStore.createPost({
        boardId: board.value.id,
        title: newPost.value.title,
        body: newPost.value.content,
      });
    }

    showCreateModal.value = false;
    newPost.value = { title: "", content: "", isThread: false };
    await loadContent();
  } catch (error) {
    console.error("Failed to create content:", error);
  }
}

function navigateToPost(postId: string) {
  navigateTo(`/p/${postId}`);
}

function navigateToThread(threadId: string) {
  navigateTo(`/th/${threadId}`);
}

async function likePost(postId: string) {
  if (!authStore.user) {
    navigateTo("/auth/login");
    return;
  }
  await postStore.like(postId);
  await loadPosts();
}

async function repostPost(postId: string) {
  if (!authStore.user) {
    navigateTo("/auth/login");
    return;
  }
  await postStore.repost(postId);
}
</script>

<template>
  <div class="board-container">
    <!-- Board Header -->
    <div v-if="board" class="board-header nb-panel mb-6">
      <div class="flex items-start justify-between">
        <div>
          <h1 class="text-3xl font-black mb-2">
            /{{ board.slug }}/ - {{ board.name }}
          </h1>
          <p class="text-lg muted">{{ board.description }}</p>
          <div class="flex gap-4 mt-3 text-sm kv">
            <span><HashIcon/> {{ board.stats?.posts || 0 }} posts</span>
            <span><LayersIcon /> {{ board.stats?.threads || 0 }} threads</span>
            <span><UserIcon /> {{ board.stats?.members || 0 }} members</span>
          </div>
        </div>
        <div class="flex gap-2">
          <button class="nb-btn" @click="boardStore.getBoard(boardSlug)">
            <RefreshIcon />
          </button>
          <button class="nb-btn nb-btn--primary" @click="openCreateModal">
            <PlusIcon />
          </button>
        </div>
      </div>
    </div>

    <!-- Content Tabs -->
    <div class="tabs mb-4">
      <button
        class="tab-btn"
        :class="{ active: activeTab === 'posts' }"
        @click="switchTab('posts')"
      >
        <HashIcon /> Posts
      </button>
      <button
        class="tab-btn"
        :class="{ active: activeTab === 'threads' }"
        @click="switchTab('threads')"
      >
        <LayersIcon /> Threads
      </button>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="text-center py-12">
      <div class="kv text-lg">Loading content...</div>
    </div>

    <!-- Posts View -->
    <div v-else-if="activeTab === 'posts'" class="posts-grid">
      <div v-if="posts.length === 0" class="nb-panel text-center py-8">
        <p class="muted">No posts yet. Be the first to post!</p>
      </div>
      <div v-else class="space-y-4">
        <div
          v-for="post in posts"
          :key="post.id"
          class="post cursor-pointer"
          @click="navigateToPost(post.id)"
        >
          <div class="post-meta">
            <div class="flex items-center gap-3">
              <div class="avatar">
                {{ (post.user?.displayName || post.user?.username || "A")[0] }}
              </div>
              <div>
                <div class="font-semibold">
                  {{
                    post.user?.displayName || post.user?.username || "Anonymous"
                  }}
                </div>
                <div class="text-xs muted">
                  @{{ post.user?.username || "anon" }} ·
                  {{ new Date(post.createdAt).toLocaleDateString() }}
                </div>
              </div>
            </div>
          </div>
          <div class="post-content">
            <h3 v-if="post.title" class="font-bold mb-1">{{ post.title }}</h3>
            <div v-html="post.contentHTML || post.content"></div>
          </div>
          <div class="post-actions" @click.stop>
            <button class="nb-btn nb-btn--small" @click="likePost(post.id)">
              👍 {{ post.likes }}
            </button>
            <button class="nb-btn nb-btn--small">👎 {{ post.dislikes }}</button>
            <button class="nb-btn nb-btn--small">
              💬 {{ post.commentCount || 0 }}
            </button>
            <button class="nb-btn nb-btn--small" @click="repostPost(post.id)">
              🔁 Repost
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Threads View -->
    <div v-else-if="activeTab === 'threads'" class="threads-list">
      <div v-if="threads.length === 0" class="nb-panel text-center py-8">
        <p class="muted">No threads yet. Start a discussion!</p>
      </div>
      <div v-else class="space-y-4">
        <div
          v-for="thread in threads"
          :key="thread.id"
          class="thread-card nb-panel nb-panel--thin cursor-pointer"
          @click="navigateToThread(thread.id)"
        >
          <div class="flex items-center justify-between mb-2">
            <div class="flex items-center gap-2">
              <span v-if="thread.pinned" class="tag tag--accent"
                >📌 Pinned</span
              >
              <span v-if="thread.closed" class="tag">🔒 Locked</span>
              <h3 class="font-bold">{{ thread.title || "Untitled Thread" }}</h3>
            </div>
            <div class="text-sm muted">{{ thread.postCount || 0 }} posts</div>
          </div>
          <div class="text-sm muted">
            Started by @{{ thread.user?.username || "anon" }} · Last activity
            {{
              new Date(
                thread.updatedAt || thread.createdAt
              ).toLocaleDateString()
            }}
          </div>
        </div>
      </div>
    </div>

    <!-- Create Post/Thread Modal -->
    <div
      v-if="showCreateModal"
      class="modal-overlay"
      @click="showCreateModal = false"
    >
      <div class="modal nb-panel" @click.stop>
        <h2 class="text-xl font-bold mb-4">
          Create New {{ newPost.isThread ? "Thread" : "Post" }}
        </h2>

        <div class="mb-4">
          <label class="flex items-center gap-2 mb-3">
            <input v-model="newPost.isThread" type="checkbox" />
            <span class="font-semibold">Create as thread</span>
          </label>
        </div>

        <div class="mb-4">
          <input
            v-model="newPost.title"
            type="text"
            placeholder="Title (optional)"
            class="input-field"
          />
        </div>

        <div class="mb-4">
          <textarea
            v-model="newPost.content"
            placeholder="What's on your mind?"
            rows="6"
            class="input-field"
          ></textarea>
        </div>

        <div class="flex gap-2 justify-end">
          <button class="nb-btn" @click="showCreateModal = false">
            Cancel
          </button>
          <button class="nb-btn nb-btn--primary" @click="createContent">
            Create
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.tabs {
  display: flex;
  gap: 0.5rem;
  border-bottom: var(--border-heavy) solid var(--muted);
  padding-bottom: 0;
}

.tab-btn {
  padding: 0.75rem 1.5rem;
  font-weight: 700;
  background: transparent;
  border: none;
  cursor: pointer;
  position: relative;
  color: var(--ink-muted);
  transition: color 0.2s;
}

.tab-btn.active {
  color: var(--ink);
}

.tab-btn.active::after {
  content: "";
  position: absolute;
  bottom: -4px;
  left: 0;
  right: 0;
  height: 4px;
  background: var(--accent);
}

.thread-card:hover {
  transform: translateX(2px);
  box-shadow: 8px 8px 0 0 var(--muted);
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal {
  max-width: 600px;
  width: 90%;
  max-height: 90vh;
  overflow-y: auto;
}

.input-field {
  width: 100%;
  padding: 0.75rem;
  border: var(--border-light) solid var(--muted);
  border-radius: 6px;
  font-family: inherit;
  background: var(--panel);
  color: var(--ink);
}

.input-field:focus {
  outline: none;
  border-color: var(--accent);
}
</style>
