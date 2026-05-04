import { defineStore } from 'pinia'
import { loginUser } from '../api/auth'
import { getProfile } from '../api/user'

const saved = JSON.parse(localStorage.getItem('package_station_auth') || '{}')

function normalizeLogin(payload, account) {
  const data = payload || {}
  const user = data.user || data.profile || {}
  const token = data.token || data.access_token || data.jwt || ''
  const role = data.role || user.role || (account === 'admin' ? 'admin' : 'user')
  return {
    token,
    role,
    user: {
      ...user,
      username: user.username || account,
      role
    }
  }
}

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: saved.token || '',
    role: saved.role || '',
    user: saved.user || null
  }),
  getters: {
    isLoggedIn: (state) => Boolean(state.token),
    isAdmin: (state) => state.role === 'admin'
  },
  actions: {
    persist() {
      localStorage.setItem(
        'package_station_auth',
        JSON.stringify({
          token: this.token,
          role: this.role,
          user: this.user
        })
      )
    },
    async login(form) {
      const payload = await loginUser(form)
      const auth = normalizeLogin(payload, form.account)
      this.token = auth.token
      this.role = auth.role
      this.user = auth.user
      if (this.token) {
        await this.refreshProfile().catch(() => {})
      }
      this.persist()
      return this.role
    },
    async refreshProfile() {
      const profile = await getProfile()
      if (profile) {
        this.user = { ...this.user, ...profile }
        this.role = profile.role || this.role
        this.persist()
      }
      return this.user
    },
    logout() {
      this.token = ''
      this.role = ''
      this.user = null
      localStorage.removeItem('package_station_auth')
    }
  }
})
