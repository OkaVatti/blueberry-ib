// /apps/frontend/stores/boards.ts
import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Board, BoardSummary } from '~/types/board'

export const useBoardStore = defineStore('boards', () => {
  const config = useRuntimeConfig()
  const apiBase = config.public.apiBase ?? '/api'

  const boards = ref<BoardSummary[]>([])
  const current = ref<Board | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function listBoards() {
    loading.value = true
    try {
      const data = await $fetch(`${apiBase}/boards`)
      boards.value = data as BoardSummary[]
      return { ok: true, data }
    } catch (err: any) {
      error.value = err?.message ?? 'failed loading boards'
      return { ok: false, error: error.value }
    } finally {
      loading.value = false
    }
  }

  async function getBoard(slugOrId: string) {
    loading.value = true
    try {
      const b = await $fetch(`${apiBase}/boards/${encodeURIComponent(slugOrId)}`)
      current.value = b as Board
      return { ok: true, data: b }
    } catch (err: any) {
      error.value = err?.message ?? 'failed loading board'
      return { ok: false, error: error.value }
    } finally {
      loading.value = false
    }
  }

  async function createBoard(payload: Partial<Board>) {
    loading.value = true
    try {
      const b = await $fetch(`${apiBase}/boards`, { method: 'POST', body: payload })
      // refresh list
      await listBoards()
      return { ok: true, data: b }
    } catch (err: any) {
      return { ok: false, error: err?.data?.message ?? err?.message }
    } finally {
      loading.value = false
    }
  }

  async function updateBoard(id: string, payload: Partial<Board>) {
    try {
      const b = await $fetch(`${apiBase}/boards/${encodeURIComponent(id)}`, { method: 'PUT', body: payload })
      // optimistic refresh
      await getBoard(id)
      return { ok: true, data: b }
    } catch (err: any) {
      return { ok: false, error: err?.message ?? 'update failed' }
    }
  }

  async function deleteBoard(id: string) {
    try {
      await $fetch(`${apiBase}/boards/${encodeURIComponent(id)}`, { method: 'DELETE' })
      await listBoards()
      return { ok: true }
    } catch (err: any) {
      return { ok: false, error: err?.message ?? 'delete failed' }
    }
  }

  return {
    boards,
    current,
    loading,
    error,
    listBoards,
    getBoard,
    createBoard,
    updateBoard,
    deleteBoard
  }
})
