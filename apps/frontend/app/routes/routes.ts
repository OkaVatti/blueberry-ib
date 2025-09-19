// /apps/frontend/routes/routes.ts
import type { ID } from "~/types/board"; // path hint for IDEs; adjust if needed
import type { Post, Comment } from "~/types/post";

// A minimal typed route descriptor for runtime / docs
export type LayoutName =
  | "default"
  | "boardLayout"
  | "postLayout"
  | "profileLayout"
  | "threadLayout"
  | "adminLayout"
  | "searchLayout";

export interface RouteDescriptor {
  path: string; // route path (nuxt route)
  name?: string; // optional route name
  pageFile: string; // relative path to file under app/pages
  layout: LayoutName;
  componentsUsed?: string[]; // component file paths (components referenced by this page)
  params?: Record<string, string>; // param name -> description/type
  authRequired?: boolean;
  adminOnly?: boolean;
}

/**
 * NOTE:
 * Dynamic page paths in your tree:
 *  - /p/{...postViewPage}.vue  -> route /p/:postId? or /p/:postSlug
 *  - /th/{...threadViewPage}.vue -> /th/:threadId?
 *  - /u/{...userProfile}.vue -> /u/:username
 *  - /{...boardName}.vue -> /:boardName
 *  - /admin/{...dashBoard}.vue -> /admin/:dashboard?
 *
 * Adjust param names to match your page's use of definePageMeta or useRoute() calls.
 */

export const ROUTES: RouteDescriptor[] = [
  {
    path: "/",
    name: "home",
    pageFile: "app/pages/indexHome.vue",
    layout: "default",
    componentsUsed: [
      "app/components/universal/NavBar.vue",
      "app/components/universal/SideBar.vue",
      "app/components/universal/SearchBar.vue",
      "app/components/home/homeView.vue",
      "app/components/home/popularPostsView.vue",
    ],
    authRequired: false,
  },

  // search
  {
    path: "/search",
    name: "search",
    pageFile: "app/pages/searchPage.vue",
    layout: "searchLayout",
    componentsUsed: [
      "app/components/universal/SearchBar.vue",
      "app/components/dashboards/utils/siteSearch.vue",
    ],
    authRequired: false,
  },

  // settings
  {
    path: "/settings",
    name: "settings",
    pageFile: "app/pages/settingsPage.vue",
    layout: "default",
    componentsUsed: ["app/components/profiles/userProfileDisplay.vue"],
    authRequired: true,
  },

  // boards list / directory
  {
    path: "/boards",
    name: "board-directory",
    pageFile: "app/components/home/boardDirectory.vue",
    layout: "default",
    componentsUsed: [
      "app/components/boards/boardInfoDisplay.vue",
      "app/components/boards/boardContentDisplay.vue",
    ],
    authRequired: false,
  },

  // dynamic board pages: /:boardName
  {
    path: "/:boardName",
    name: "board",
    pageFile: "app/pages/{...boardName}.vue",
    layout: "boardLayout",
    componentsUsed: [
      "app/components/boards/boardThreadDisplay.vue",
      "app/components/boards/boardPostDisplay.vue",
      "app/components/boards/boardInfoDisplay.vue",
    ],
    params: { boardName: "slug (string) of the board" },
    authRequired: false,
  },

  // post view pages: /p/:postId
  {
    path: "/p/:postId",
    name: "post-view",
    pageFile: "app/pages/p/{...postViewPage}.vue",
    layout: "postLayout",
    componentsUsed: [
      "app/components/posts/postView.vue",
      "app/components/posts/commentCard.vue",
      "app/components/universal/NavBar.vue",
    ],
    params: { postId: "post id or slug" },
    authRequired: false,
  },

  // thread view: /th/:threadId
  {
    path: "/th/:threadId",
    name: "thread-view",
    pageFile: "app/pages/th/{...threadViewPage}.vue",
    layout: "threadLayout",
    componentsUsed: [
      "app/components/posts/threadView.vue",
      "app/components/posts/threadCard.vue",
      "app/components/posts/postCard.vue",
    ],
    params: { threadId: "thread id or slug" },
    authRequired: false,
  },

  // user profile: /u/:username
  {
    path: "/u/:username",
    name: "user-profile",
    pageFile: "app/pages/u/{...userProfile}.vue",
    layout: "profileLayout",
    componentsUsed: [
      "app/components/profiles/profileView.vue",
      "app/components/profiles/userPostsView.vue",
      "app/components/profiles/userLikesDisplay.vue",
      "app/components/profiles/userCommentsView.vue",
    ],
    params: { username: "username string" },
    authRequired: false,
  },

  // admin dashboard(s)
  {
    path: "/admin/:dashboard?",
    name: "admin",
    pageFile: "app/pages/admin/{...dashBoard}.vue",
    layout: "adminLayout",
    componentsUsed: [
      "app/components/dashboards/administratorDash.vue",
      "app/components/dashboards/moderatorDash.vue",
      "app/components/dashboards/ownerDash.vue",
      "app/components/dashboards/utils/reportedPostsDisplay.vue",
    ],
    authRequired: true,
    adminOnly: true,
  },

  // catch-all fallback: map to home (optional)
  {
    path: "/:catchAll(.*)*",
    name: "not-found",
    pageFile: "app/pages/index.vue",
    layout: "default",
    componentsUsed: ["app/components/universal/NavBar.vue"],
    authRequired: false,
  },
];
