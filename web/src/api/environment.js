import request from '../utils/request'

export function fetchList(query) {
  return request({
    url: '/v1/environment/list',
    method: 'get',
    params: query
  })
}

export function createRow(projectId, data) {
  return request({
    url: '/v1/environment/' + projectId + '/create',
    method: 'post',
    data
  })
}

export function updateRow(id, data) {
  return request({
    url: '/v1/environment/' + id,
    method: 'post',
    data
  })
}

export function deleteRow(id) {
  return request({
    url: '/v1/environment/' + id + '/delete',
    method: 'post'
  })
}

export function copyRow(id, data) {
  return request({
    url: '/v1/environment/' + id + '/copy',
    method: 'post',
    data
  })
}
