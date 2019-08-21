import request from '@/utils/request'

export function fetchList(query) {
  return request({
    url: '/v1/project/list',
    method: 'get',
    params: query
  })
}

export function createRow(data) {
  return request({
    url: '/v1/project/create',
    method: 'post',
    data
  })
}

export function updateRow(id, data) {
  return request({
    url: '/v1/project/' + id,
    method: 'post',
    data
  })
}

export function deleteRow(id) {
  return request({
    url: '/v1/project/' + id + '/delete',
    method: 'post'
  })
}
