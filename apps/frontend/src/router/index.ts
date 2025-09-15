import { createRouter, createWebHistory } from "vue-router";
import BoardsList from "../components/BoardsList.vue";
import BoardView from "../components/BoardView.vue";
import PostView from "../components/PostView.vue";
import Login from "../components/Login.vue";

const routes = [
  { path: "/", name: "home", component: BoardsList },
  { path: "/login", name: "login", component: Login },
  { path: "/boards/:slug", name: "board", component: BoardView, props: true },
  { path: "/posts/:id", name: "post", component: PostView, props: true },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
