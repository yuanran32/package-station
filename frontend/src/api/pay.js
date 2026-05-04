import request from './request'

export function createPayment(data) {
  return request.post('/api/pay/create', data)
}

export function payCallback(data) {
  return request.post('/api/pay/callback', data)
}

export function getBills() {
  return request.get('/api/pay/bill')
}
