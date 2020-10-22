import axios from '@/libs/api.request'

export const getProject = (projectId) => {
  return axios.request({
    url: 'api/v1/projects/' + projectId,
    method: 'get'
  })
}

export const getProjectSnapshotConfig = (taskId) => {
  return axios.request({
    url: 'api/v1/snapshot',
    params: {
      taskId: taskId,
    },
    method: 'get'
  })
}

export const listProjects = (page, pageSize, order) => {
  const params = {
    page,
    pageSize,
    order,
  }
  return axios.request({
    url: 'api/v1/projects',
    params,
    method: 'get'
  })
}

export const deleteProject = (projectId) => {
  return axios.request({
    url: `api/v1/projects/` + projectId,
    method: 'delete'
  })
}


export const parseProjectConfigFile = (data) => {
  return axios.request({
    url: `api/v2/play/parse`,
    data,
    method: 'post'
  })
}

export const debugStage = (data) => {
  return axios.request({
    url: `api/v2/play`,
    data,
    method: 'post'
  })
}


export const updateProjectBaseInfo = (data) => {
  return axios.request({
    url: 'api/v1/projects/' + data.id,
    data,
    method: 'put'
  })
}

export const saveProjectConfig = (data) => {
  return axios.request({
    url: 'api/v1/projects/' + data.id + '/config',
    data,
    method: 'put'
  })
}

export const createProject = (data) => {
  return axios.request({
    url: 'api/v1/projects',
    data,
    method: 'post'
  })
}


export const savePlugins = (data) => {
  return axios.request({
    url: 'api/v1/projects/' + data.projectId + '/plugins',
    data,
    method: 'put'
  })
}

export const startNewTask = (projectId) => {
  return axios.request({
    url: 'api/v1/tasks',
    params: {
      projectId: projectId
    },
    method: 'post'
  })
}


export const updateProjectProxies = (projectId, proxyIds) => {
  return axios.request({
    url: 'api/v1/projects/'+ projectId +'/proxies',
    data: proxyIds,
    method: 'post'
  })
}
