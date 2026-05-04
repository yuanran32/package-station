import request from './request'

export function receiveCoupon(data) {
  return request.post('/api/coupon/receive', data)
}

export function getMyCoupons() {
  return request.get('/api/coupon/my')
}

export function useCoupon(data) {
  return request.post('/api/coupon/use', data)
}
