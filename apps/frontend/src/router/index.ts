import { createRouter, createWebHistory, type RouteRecordRaw } from "vue-router";

// layouts
import DefaultLayout from "../layouts/DefaultLayout.vue";
import DashboardLayout from "../layouts/DashboardLayout.vue";
import BoardLayout from "../layouts/BoardLayout.vue";
import SearchLayout from "../layouts/SearchLayout.vue";
import PostLayout from "../layouts/PostLayout.vue";
import ThreadLayout from "../layouts/ThreadLayout.vue";

// pages
import HomeIndex from "../pages/HomeIndex.vue";
import ProfilePage from "../pages/ProfilePage.vue";
import BoardPage from "../pages/BoardPage.vue";
import SearchPage from "../pages/SearchPage.vue";
import PostPage from "../pages/PostPage.vue";
import ThreadPage from "../pages/ThreadPage.vue";

// dashboards (these are children inside DashboardLayout)
import OwnerDash from "../components/ownerDash.vue";
import CoOwnerDash from "../components/coOwnerDash.vue";
import AdminDash from "../components/adminDash.vue";
import ModeratorDash from "../components/moderatorDash.vue";

const routes: Array<RouteRecordRaw> = [
  {
    path: "/",
    component: DefaultLayout,
    children: [
      { path: "", name: "home", component: HomeIndex },
      {
        path: "profile/:id",
        name: "profile",
        component: ProfilePage,
        props: true,
      },
      { path: "search", name: "search", component: SearchPage },
    ],
  },

  {
    path: "/dash",
    component: DashboardLayout,
    children: [
      { path: "owner", name: "dash-owner", component: OwnerDash },
      { path: "coowner", name: "dash-coowner", component: CoOwnerDash },
      { path: "admin", name: "dash-admin", component: AdminDash },
      { path: "mod", name: "dash-mod", component: ModeratorDash },
    ],
  },

  {
    path: "/boards/:slug",
    component: BoardLayout,
    children: [
      { path: "", name: "board", component: BoardPage, props: true },
      // you can add /boards/:slug/threads, etc.
    ],
  },

  {
    path: "/posts/:id",
    component: PostLayout,
    children: [{ path: "", name: "post", component: PostPage, props: true }],
  },

  {
    path: "/threads/:id",
    component: ThreadLayout,
    children: [
      { path: "", name: "thread", component: ThreadPage, props: true },
    ],
  },

  // fallback
  { path: "/:catchAll(.*)", redirect: "/" },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
