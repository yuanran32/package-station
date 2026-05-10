import request from './request'

export function createAdminCoupon(data) {
  return request.post('/api/admin/coupon/create', data)
}

export function getAdminCouponList(params) {
  return request.get('/api/admin/coupon/list', { params })
}

export function receiveCoupon(data) {
  return request.post('/api/coupon/receive', data)
}

export function getMyCoupons() {
  return request.get('/api/coupon/my')
}

export function useCoupon(data) {
  return request.post('/api/coupon/use', data)
}
