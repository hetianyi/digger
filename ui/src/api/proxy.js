import axios from '@/libs/api.request'

export const proxyList = (key, page, pageSize) => {
  return axios.request({
    url: 'api/v1/proxies',
    params: {
      page: page,
      pageSize: pageSize,
      key: key,
    },
    method: 'get'
  })
}

export const deleteProxies = (idList) => {
  return axios.request({
    url: 'api/v1/proxies',
    data: idList,
    method: 'delete'
  })
}

export const saveProxy = (proxy) => {
  return axios.request({
    url: 'api/v1/proxies',
    data: proxy,
    method: 'post'
  })
}
