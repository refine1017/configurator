import request from '@/utils/request'

export function fetchList(query) {
  return request({
    url: '/v1/table/list',
    method: 'get',
    params: query,
  })
}

export function createRow(envId, collect, data) {
  return request({
    url: '/v1/table/' + envId + '/' + collect,
    method: 'post',
    data
  })
}

export function updateRow(envId, collect, id, data) {
  return request({
    url: '/v1/table/' + envId + '/' + collect + '/' + id,
    method: 'post',
    data
  })
}

export function deleteRow(envId, collect, id) {
  return request({
    url: '/v1/table/' + envId + '/' + collect + '/' + id + "/delete",
    method: 'post',
  })
}

export function uploadData(envId, collect, data) {
  return request({
    url: '/v1/table/' + envId + '/' + collect + '/upload',
    method: 'post',
    data
  })
}
