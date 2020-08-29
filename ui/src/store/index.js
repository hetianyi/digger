import Vue from 'vue'
import Vuex from 'vuex'

import user from './module/user'
import app from './module/app'

Vue.use(Vuex)

let oldPid = window.localStorage.getItem('currentTaskPageSelectedProjectId')
if (oldPid == null) {
  oldPid = 0
}
oldPid = Number(oldPid)
console.log('恢复的currentTaskPageSelectedProjectId=' + oldPid)

let oldStatus = window.localStorage.getItem('currentTaskPageSelectedStatus')
if (oldStatus == null) {
  oldStatus = -1
}
oldStatus = Number(oldStatus)
if (isNaN(oldStatus)) {
  oldStatus = -1
}
console.log('恢复的currentTaskPageSelectedStatus=' + oldStatus)

export default new Vuex.Store({
  state: {
    //
    currentTaskPageSelectedProjectId: oldPid,
    currentTaskPageSelectedStatus: oldStatus,
  },
  mutations: {
    selectProjectChange (state, payload) {
      state.currentTaskPageSelectedProjectId  = payload
      console.log('store currentTaskPageSelectedProjectId ...' + payload)
      window.localStorage.setItem('currentTaskPageSelectedProjectId', payload)
    },
    selectStatusChange (state, payload) {
      state.currentTaskPageSelectedStatus  = payload
      console.log('store currentTaskPageSelectedStatus ...' + payload)
      window.localStorage.setItem('currentTaskPageSelectedStatus', payload)
    },
  },
  actions: {
    //
  },
  modules: {
    user,
    app
  }
})
