import axios from '@/libs/api.request'

export const configList = () => {
  return axios.request({
    url: 'api/v1/configs',
    method: 'get'
  })
}

export const updateList = (data) => {
  return axios.request({
    url: 'api/v1/configs',
    data,
    method: 'put'
  })
}
