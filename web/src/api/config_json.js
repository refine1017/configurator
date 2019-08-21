import request from '@/utils/request'

export function fetchList(query) {
  return request({
    url: '/v1/json/list',
    method: 'get',
    params: query
  })
}

export function createRow(envId, collect, data) {
  return request({
    url: '/v1/json/' + envId + '/' + collect,
    method: 'post',
    data
  })
}

export function updateRow(envId, collect, id, data) {
  return request({
    url: '/v1/json/' + envId + '/' + collect + '/' + id,
    method: 'post',
    data
  })
}

export function deleteRow(envId, collect, id) {
  return request({
    url: '/v1/json/' + envId + '/' + collect + '/' + id + '/delete',
    method: 'post'
  })
}

export function getData(query) {
  return request({
    url: '/v1/json/' + query.envId + '/' + query.collect + '/' + query.id + '/data',
    method: 'get'
  })
}

export function setData(query, data) {
  return request({
    url: '/v1/json/' + query.envId + '/' + query.collect + '/' + query.id + '/data',
    method: 'post',
    data
  })
}

export function uploadData(envId, collect, data) {
  return request({
    url: '/v1/json/' + envId + '/' + collect + '/upload',
    method: 'post',
    data
  })
}
