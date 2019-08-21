import request from '@/utils/request'

export function getSetting() {
  return request({
    url: '/v1/app/setting',
    method: 'get',
  })
}
