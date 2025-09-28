<!-- apps/frontend/app/pages/p/{...postViewPage}.vue -->
<script setup lang="ts">
import { ref, onMounted, computed } from "vue";
import { useRoute } from "vue-router";
import { usePostStore } from "~/stores/posts";
import { useAuthStore } from "~/stores/auth";

const route = useRoute();
const postStore = usePostStore();
const authStore = useAuthStore();

const postId = computed(() => route.params.postViewPage as string);
const post = ref<any>(null);
const comments = ref<any[]>([]);
const loading = ref(true);
const newComment = ref("");
const showReplyTo = ref<string | null>(null);

onMounted(async () => {
  await loadPost();
  await loadComments();
});

async function loadPost() {
  loading.value = true;
  try {
    const result = await postStore.fetchPost(postId.value);
    if (result.ok) {
      post.value = result.data;
    }
  } catch (error) {
    console.error("Failed to load post:", error);
  } finally {
    loading.value = false;
  }
}

async function loadComments() {
  try {
    const result = await postStore.fetchComments(postId.value);
    if (result.ok) {
      comments.value = result.data;
    }
  } catch (error) {
    console.error("Failed to load comments:", error);
  }
}

async function submitComment() {
  if (!newComment.value.trim()) return;

  if (!authStore.user) {
    navigateTo("/auth/login");
    return;
  }

  try {
    const result = await postStore.addComment(postId.value, {
      body: newComment.value,
    });

    if (result.ok) {
      newComment.value = "";
      await loadComments();
    }
  } catch (error) {
    console.error("Failed to post comment:", error);
  }
}

async function likePost() {
  if (!authStore.user) {
    navigateTo("/auth/login");
    return;
  }

  await postStore.like(postId.value);
  await loadPost();
}

async function dislikePost() {
  if (!authStore.user) {
    navigateTo("/auth/login");
    return;
  }

  await postStore.dislike(postId.value);
  await loadPost();
}

async function repost() {
  if (!authStore.user) {
    navigateTo("/auth/login");
    return;
  }

  await postStore.repost(postId.value);
}

function formatDate(date: string) {
  const d = new Date(date);
  const now = new Date();
  const diff = (now.getTime() - d.getTime()) / 1000;

  if (diff < 60) return "just now";
  if (diff < 3600) return `${Math.floor(diff / 60)}m ago`;
  if (diff < 86400) return `${Math.floor(diff / 3600)}h ago`;
  if (diff < 604800) return `${Math.floor(diff / 86400)}d ago`;

  return d.toLocaleDateString();
}

function navigateToUser(username: string) {
  navigateTo(`/u/${username}`);
}

function navigateToBoard(slug: string) {
  navigateTo(`/${slug}`);
}
</script>

<template>
  <div class="post-view-container">
    <!-- Back Navigation -->
    <div class="mb-4">
      <button class="nb-btn" @click="$router.back()">← Back</button>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="text-center py-12">
      <div class="kv text-lg">Loading post...</div>
    </div>

    <!-- Post Content -->
    <article v-else-if="post" class="post-full nb-panel mb-6">
      <!-- Post Header -->
      <div class="post-header">
        <div class="flex items-start justify-between">
          <div class="flex items-center gap-3">
            <div
              class="avatar avatar--large"
              @click="navigateToUser(post.user?.username)"
            >
              {{ (post.user?.displayName || post.user?.username || "A")[0] }}
            </div>
            <div>
              <div class="font-bold text-lg">
                {{
                  post.user?.displayName || post.user?.username || "Anonymous"
                }}
              </div>
              <div class="text-sm muted">
                @{{ post.user?.username || "anon" }} ·
                <span class="link" @click="navigateToBoard(post.board?.slug)">
                  /{{ post.board?.slug }}/
                </span>
                ·
                {{ formatDate(post.createdAt) }}
              </div>
            </div>
          </div>

          <div v-if="authStore.user?.id === post.user?.id" class="dropdown">
            <button class="nb-btn">⋮</button>
          </div>
        </div>
      </div>

      <!-- Post Body -->
      <div class="post-body">
        <h1 v-if="post.title" class="text-2xl font-bold mb-3">
          {{ post.title }}
        </h1>
        <div class="content" v-html="post.contentHTML || post.content"></div>

        <!-- Media Attachments -->
        <div v-if="post.mediaAttachments?.length" class="media-grid mt-4">
          <img
            v-for="(media, idx) in post.mediaAttachments"
            :key="idx"
            :src="media.url"
            :alt="media.altText"
            class="media-item"
          />
        </div>
      </div>
    </article>
  </div>
</template>
