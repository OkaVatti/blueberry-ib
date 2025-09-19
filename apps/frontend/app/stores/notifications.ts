// /apps/frontend/stores/notifications.ts
import { defineStore } from "pinia";
import { ref } from "vue";

export interface NotificationItem {
  id: string;
  type: "info" | "success" | "warning" | "error";
  title?: string;
  body: string;
  createdAt: string;
  read?: boolean;
}

export const useNotificationsStore = defineStore("notifications", () => {
  const items = ref<NotificationItem[]>([]);

  function push(
    n: Omit<NotificationItem, "id" | "createdAt"> & { id?: string }
  ) {
    const item: NotificationItem = {
      id: n.id ?? `${Date.now()}-${Math.random().toString(36).slice(2, 7)}`,
      title: n.title,
      body: n.body,
      type: n.type,
      createdAt: new Date().toISOString(),
      read: false,
    };
    items.value.unshift(item);
    // auto-expire optional
    setTimeout(() => {
      items.value = items.value.filter((i) => i.id !== item.id);
    }, 12_000);
    return item;
  }

  function markRead(id: string) {
    const it = items.value.find((i) => i.id === id);
    if (it) it.read = true;
  }
  function clearAll() {
    items.value = [];
  }

  return { items, push, markRead, clearAll };
});
