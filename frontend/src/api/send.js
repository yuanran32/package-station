import request from './request'

export function createSendOrder(data) {
  return request.post('/api/send/order', data)
}

export function getAdminSendOrders(status) {
  return request.get('/api/admin/send/orders', {
    params: status ? { status } : {}
  })
}

export function processSendOrder(data) {
  return request.post('/api/admin/send/process', data)
}
