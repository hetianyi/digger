import axios from '@/libs/api.request'

export const resultList = (taskId, page, pageSize) => {
  return axios.request({
    url: 'api/v1/results',
    params: {
      page: page,
      pageSize: pageSize,
      taskId: taskId,
    },
    method: 'get'
  })
}
