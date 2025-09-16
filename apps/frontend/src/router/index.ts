import { createRouter, createWebHistory } from "vue-router";
import MainLayout from "../components/MainLayout.vue";
import BoardsList from "../components/BoardsList.vue";
import BoardView from "..//components/BoardView.vue";
import PostView from "..//components/PostView.vue";
import Login from "../components/Login.vue";
import ThreadView from "../components/ThreadView.vue";
import OwnerDash from "../components/ownerDash.vue";
import CoOwnerDash from "../components/coOwnerDash.vue";
import AdminDash from "../components/adminDash.vue";
import ModeratorDash from "../components/moderatorDash.vue";

const routes = [
  { path: "/", name: "home", component: MainLayout },
  { path: "/boards", name: "boards", component: BoardsList },
  { path: "/login", name: "login", component: Login },
  { path: "/boards/:slug", name: "board", component: BoardView, props: true },
  { path: "/posts/:id", name: "post", component: PostView, props: true },
  { path: "/threads/:id", name: "thread", component: ThreadView, props: true },
  { path: "/dash/owner", name: "ownerDash", component: OwnerDash },
  { path: "/dash/coowner", name: "coownerDash", component: CoOwnerDash },
  { path: "/dash/admin", name: "adminDash", component: AdminDash },
  { path: "/dash/mod", name: "modDash", component: ModeratorDash },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
