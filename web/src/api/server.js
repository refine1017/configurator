import request from '../utils/request'

export function fetchList(query) {
  return request({
    url: '/v1/server/list',
    method: 'get',
    params: query
  })
}

export function createRow(projectId, data) {
  return request({
    url: '/v1/server/' + projectId + '/create',
    method: 'post',
    data
  })
}

export function updateRow(id, data) {
  return request({
    url: '/v1/server/' + id,
    method: 'post',
    data
  })
}

export function deleteRow(id) {
  return request({
    url: '/v1/server/' + id + '/delete',
    method: 'post'
  })
}
