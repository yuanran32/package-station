import request from './request'

export function getProfile() {
  return request.get('/api/user/profile')
}

export function updateProfile(data) {
  return request.put('/api/user/profile', data)
}

export function getPickupHistory() {
  return request.get('/api/user/pickup-history')
}

export function getIdentityCode() {
  return request.get('/api/user/qrcode')
}
