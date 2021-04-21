import request from '@/utils/request'
import { b64EncodeUnicode } from '@/utils'

export function login(data) {
  return request({
    url: '@oauth/token?grant_type=password',
    // url: '/vue-element-admin/user/login',
    method: 'post',
    headers: {
      Authorization: 'Basic ' + b64EncodeUnicode('admin_client_id:123456') },
    data
  })
}

export function getInfo(token) {
  return request({
    url: '/vue-element-admin/user/info',
    method: 'post',
    params: { token }
  })
}

export function logout() {
  return request({
    url: '/vue-element-admin/user/logout',
    method: 'post'
  })
}

// 刷新TOKEN
export function refreshToken(rToken) {
  return request({
    url: '/vue-element-admin/user/refresh-token',
    method: 'post',
    data: { refresh_token: rToken }
  })
}
