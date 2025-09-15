import { createRouter, createWebHistory } from "vue-router";
import BoardsList from "../components/BoardsList.vue";
import BoardView from "..//components/BoardView.vue";
import PostView from "..//components/PostView.vue";
import Login from "../components/Login.vue";
import ThreadView from "../components/ThreadView.vue";

// add imports
import OwnerDash from "../components/ownerDash.vue";
import CoOwnerDash from "../components/coOwnerDash.vue";
import AdminDash from "../components/adminDash.vue";
import ModeratorDash from "../components/moderatorDash.vue";

const routes = [
  { path: "/", name: "home", component: BoardsList },
  { path: "/login", name: "login", component: Login },
  { path: "/boards/:slug", name: "board", component: BoardView, props: true },
  { path: "/posts/:id", name: "post", component: PostView, props: true },
  { path: "/threads/:id", name: "thread", component: ThreadView, props: true },
  // Add routes
  {
    path: "/dash/owner",
    name: "ownerDash",
    component: OwnerDash,
    meta: { role: "owner" },
  },
  {
    path: "/dash/coowner",
    name: "coownerDash",
    component: CoOwnerDash,
    meta: { role: "coowner" },
  },
  {
    path: "/dash/admin",
    name: "adminDash",
    component: AdminDash,
    meta: { role: "admin" },
  },
  {
    path: "/dash/mod",
    name: "modDash",
    component: ModeratorDash,
    meta: { role: "moderator" },
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
