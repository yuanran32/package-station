import request from './request'

export function generatePickupCode(data) {
  return request.post('/api/pickup/code', data)
}

export function recordPickup(data) {
  return request.post('/api/pickup/record', data)
}

export function recordDelivery(data) {
  return request.post('/api/delivery/record', data)
}
