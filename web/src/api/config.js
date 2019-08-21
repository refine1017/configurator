import request from '../utils/request'

export function fetchList(query) {
  return request({
    url: '/v1/config/list',
    method: 'get',
    params: query
  })
}

export function createRow(envId, data) {
  return request({
    url: '/v1/config/' + envId + '/create',
    method: 'post',
    data
  })
}

export function updateRow(envId, index, data) {
  return request({
    url: '/v1/config/' + envId + '/' + index,
    method: 'post',
    data
  })
}

export function deleteRow(envId, index) {
  return request({
    url: '/v1/config/' + envId + '/' + index + '/delete',
    method: 'post'
  })
}

export function pushConfig(envId, server, config) {
  return request({
    url: '/v1/config/' + envId + '/push/' + config + '/' + server,
    method: 'post'
  })
}

export function mergeConfigInfo(envId, config) {
  return request({
    url: '/v1/config/' + envId + '/merge/' + config + '/info',
    method: 'post'
  })
}

export function mergeConfig(envId, config) {
  return request({
    url: '/v1/config/' + envId + '/merge/' + config,
    method: 'post'
  })
}
