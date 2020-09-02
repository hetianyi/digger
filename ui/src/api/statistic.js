import axios from '@/libs/api.request'

export const getStatistic = (params) => {
  return axios.request({
    url: 'api/v1/statistics',
    params,
    method: 'get'
  })
}

