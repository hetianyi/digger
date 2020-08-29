import axios from '@/libs/api.request'

export const nodeList = () => {
  return axios.request({
    url: 'api/v1/nodes',
    method: 'get'
  })
}

