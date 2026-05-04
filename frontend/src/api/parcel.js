import request from './request'

export function inboundParcel(data) {
  return request.post('/api/parcel/inbound', data)
}

export function outboundParcel(data) {
  return request.post('/api/parcel/outbound', data)
}

export function getParcelStatus(trackingNo) {
  return request.get('/api/parcel/status', {
    params: { tracking_no: trackingNo }
  })
}

export function getParcelList() {
  return request.get('/api/parcel/list')
}
