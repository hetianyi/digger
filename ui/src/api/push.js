import axios from '@/libs/api.request'

export const pushList = (key, page, pageSize) => {
  return axios.request({
    url: 'api/v1/pushes',
    params: {
      page: page,
      pageSize: pageSize,
      key: key,
    },
    method: 'get'
  })
}

export const deletePushes = (idList) => {
  return axios.request({
    url: 'api/v1/pushes',
    data: idList,
    method: 'delete'
  })
}

export const savePush = (data) => {
  return axios.request({
    url: 'api/v1/pushes',
    data: data,
    method: 'post'
  })
}
