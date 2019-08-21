import request from '@/utils/request'

export function fetchList(query) {
  return request({
    url: '/v1/admin/list',
    method: 'get',
    params: query
  })
}

export function createRow(data) {
  return request({
    url: '/v1/admin/create',
    method: 'post',
    data
  })
}

export function updateRow(id, data) {
  return request({
    url: '/v1/admin/' + id,
    method: 'post',
    data
  })
}

export function deleteRow(id) {
  return request({
    url: '/v1/admin/' + id + '/delete',
    method: 'post'
  })
}

export function fetchLogList(query) {
  return request({
    url: '/v1/admin/logs',
    method: 'get',
    params: query
  })
}
