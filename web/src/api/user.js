import request from '@/utils/request'

export function login(data) {
  return request({
    url: '/v1/user/login',
    method: 'post',
    data
  })
}

export function getInfo() {
  return request({
    url: '/v1/user/info',
    method: 'get',
  })
}

export function getProjects() {
  return request({
    url: '/v1/user/projects',
    method: 'get'
  })
}

export function logout() {
  return request({
    url: '/v1/user/logout',
    method: 'post'
  })
}
