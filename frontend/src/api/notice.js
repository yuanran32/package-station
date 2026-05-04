import request from './request'

export function sendNotice(data) {
  return request.post('/api/notice/send', data)
}

export function createNotifySocket(token) {
  const base = import.meta.env.VITE_WS_BASE_URL || 'ws://localhost:8080'
  return new WebSocket(`${base}/ws/notify?token=${encodeURIComponent(token)}`)
}
