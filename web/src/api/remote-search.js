import request from '@/utils/request'

export function searchUser(name) {
  return request({
    url: '/dev/search/user',
    method: 'get',
    params: { name }
  })
}

export function transactionList(query) {
  return request({
    url: '/dev/transaction/list',
    method: 'get',
    params: query
  })
}
