import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getProfile } from '../api'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const userInfo = ref(null)

  const setToken = (t) => {
    token.value = t
    localStorage.setItem('token', t)
  }

  const fetchProfile = async () => {
    try {
      const data = await getProfile()
      userInfo.value = data
    } catch {
      userInfo.value = null
    }
  }

  const logout = () => {
    token.value = ''
    userInfo.value = null
    localStorage.removeItem('token')
  }

  return { token, userInfo, setToken, fetchProfile, logout }
})
