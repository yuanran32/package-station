import { defineStore } from 'pinia'

const STORAGE_KEY = 'package_station_notices'

function loadSavedNotices() {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    const parsed = JSON.parse(raw || '[]')
    return Array.isArray(parsed) ? parsed : []
  } catch {
    return []
  }
}

export const useNoticeStore = defineStore('notices', {
  state: () => ({
    notices: loadSavedNotices()
  }),
  actions: {
    persist() {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(this.notices))
    },
    addNotice(content, time = new Date().toLocaleString()) {
      this.notices.unshift({
        id: Date.now() + Math.random(),
        time,
        content
      })
      if (this.notices.length > 200) {
        this.notices = this.notices.slice(0, 200)
      }
      this.persist()
    },
    clearNotices() {
      this.notices = []
      this.persist()
    }
  }
})

