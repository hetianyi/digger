import axios from '@/libs/api.request'

export const taskList = (projectId, page, pageSize, status) => {
  return axios.request({
    url: 'api/v1/tasks',
    params: {
      page: page,
      pageSize: pageSize,
      projectId: projectId,
      status: status,
    },
    method: 'get'
  })
}


export const pauseTask = (taskId) => {
  return axios.request({
    url: 'api/v1/tasks/'+ taskId +'/pause',
    method: 'put'
  })
}

export const continueTask = (taskId) => {
  return axios.request({
    url: 'api/v1/tasks/'+ taskId +'/continue',
    method: 'put'
  })
}

export const stopTask = (taskId) => {
  return axios.request({
    url: 'api/v1/tasks/'+ taskId +'/stop',
    method: 'put'
  })
}

export const deleteTask = (taskId) => {
  return axios.request({
    url: 'api/v1/tasks/'+ taskId,
    method: 'delete'
  })
}
