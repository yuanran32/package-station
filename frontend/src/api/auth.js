import request from './request'

export function registerUser(data) {
  return request.post('/api/user/register', data)
}

export function loginUser(data) {
  return request.post('/api/user/login', data)
}
