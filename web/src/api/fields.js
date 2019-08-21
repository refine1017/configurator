import request from '../utils/request'

export function fetchList(query) {
  return request({
    url: '/v1/fields/list',
    method: 'get',
    params: query,
  })
}

export function createRow(envId, configIndex, data) {
  return request({
    url: '/v1/fields/' + envId + '/' + configIndex + '/create',
    method: 'post',
    data
  })
}

export function updateRow(envId, configIndex, index, data) {
  return request({
    url: '/v1/fields/' + envId + '/' + configIndex + '/' + index,
    method: 'post',
    data
  })
}

export function deleteRow(envId, configIndex, index) {
  return request({
    url: '/v1/fields/' + envId + '/' + configIndex + '/' + index + '/delete',
    method: 'post'
  })
}
