import { login, logout, getInfo, getProjects } from '@/api/user'
import { getToken, setToken, removeToken } from '@/utils/auth'
import { cookieGet, cookieSet, cookieRemove } from '@/utils/cookie'
import router, { resetRouter } from '@/router'
import store from '../index'

const state = {
  token: getToken(),
  projects: {},
  selectProject: {},
  selectEnvironments: [],
  selectEnvironment: {},
  selectServers: [],
  name: '',
  avatar: '',
  introduction: '',
  roles: []
}

const mutations = {
  SET_TOKEN: (state, token) => {
    state.token = token
  },
  SET_PROJECTS: (state, projects) => {
    state.projects = projects
  },
  SET_SELECT_PROJECT: (state, project) => {
    state.selectProject = project
  },
  SET_SELECT_ENVIRONMENTS: (state, environments) => {
    state.selectEnvironments = environments
  },
  SET_SELECT_ENVIRONMENT: (state, environment) => {
    state.selectEnvironment = environment
  },
  SET_SELECT_SERVERS: (state, servers) => {
    state.selectServers = servers
  },
  SET_INTRODUCTION: (state, introduction) => {
    state.introduction = introduction
  },
  SET_NAME: (state, name) => {
    state.name = name
  },
  SET_AVATAR: (state, avatar) => {
    state.avatar = avatar
  },
  SET_ROLES: (state, roles) => {
    state.roles = roles
  }
}

const actions = {
  // user login
  login({ commit }, userInfo) {
    const { username, password } = userInfo
    return new Promise((resolve, reject) => {
      login({ username: username.trim(), password: password }).then(response => {
        const { data } = response
        commit('SET_TOKEN', data.token)
        setToken(data.token)
        resolve()
      }).catch(error => {
        reject(error)
      })
    })
  },

  // get user info
  getInfo({ commit, state }) {
    return new Promise((resolve, reject) => {
      getInfo().then(response => {
        const { data } = response

        if (!data) {
          reject('Verification failed, please Login again.')
        }

        const { roles, name, avatar, introduction } = data

        // roles must be a non-empty array
        if (!roles || roles.length <= 0) {
          reject('getInfo: roles must be a non-null array!')
        }

        commit('SET_ROLES', roles)
        commit('SET_NAME', name)
        commit('SET_AVATAR', avatar)
        commit('SET_INTRODUCTION', introduction)

        resolve()
      }).catch(error => {
        reject(error)
      }).then(function() {
        // router.push('/dashboard')
      })
    })
  },
  // render user projects
  renderProjects({ commit, state }) {
    return new Promise((resolve, reject) => {
      getProjects().then(response => {
        const projects = response.data

        if (!projects || projects.length <= 0) {
          reject('renderProjects: projects must be a non-null array!')
        }

        commit('SET_PROJECTS', projects)

        resolve()
      }).catch(error => {
        reject(error)
      }).then(function() {
        store.dispatch('user/changeProject', cookieGet('projectId'))
      })
    })
  },
  changeProject({ commit, dispatch }, id) {
    return new Promise(async resolve => {
      let project

      for (let i = 0; i < state.projects.length; i++) {
        if (state.projects[i].id === id) {
          project = state.projects[i]
          break
        }
      }

      if (project === undefined) {
        project = state.projects[0]
      }

      commit('SET_SELECT_PROJECT', project)
      commit('SET_SELECT_ENVIRONMENTS', project.envs)
      commit('SET_SELECT_SERVERS', project.servers)
      cookieSet('projectId', project.id)
      console.log('SET_SELECT_PROJECT to ', project.name)

      await dispatch('changeEnvironment', cookieGet('envId'))

      resolve()
    })
  },
  // change environment
  changeEnvironment({ commit, dispatch }, id) {
    return new Promise(async resolve => {
      var env = undefined

      for (let i = 0; i < state.selectProject.envs.length; i++) {
        if (state.selectProject.envs[i].id === id) {
          env = state.selectProject.envs[i]
          break
        }
      }

      if (env === undefined) {
        env = state.selectProject.envs[0]
      }

      cookieSet('envId', env.id)
      commit('SET_SELECT_ENVIRONMENT', env)
      console.log('SET_SELECT_ENVIRONMENT to ', env.name)

      await dispatch('renderMenu')

      resolve()
    })
  },
  // render menu
  renderMenu({ commit, dispatch }) {
    return new Promise(async resolve => {
      resetRouter()

      // generate accessible routes map based on roles
      const accessRoutes = await dispatch('permission/generateRoutes', state.roles, { root: true })

      // dynamically add accessible routes
      router.addRoutes(accessRoutes)

      await store.dispatch('tagsView/delAllViews')

      resolve()
    })
  },
  // user logout
  logout({ commit, state }) {
    return new Promise((resolve, reject) => {
      logout(state.token).then(() => {
        commit('SET_TOKEN', '')
        commit('SET_ROLES', [])
        removeToken()
        resetRouter()
        resolve()
      }).catch(error => {
        reject(error)
      })
    })
  },

  // remove token
  resetToken({ commit }) {
    return new Promise(resolve => {
      commit('SET_TOKEN', '')
      commit('SET_ROLES', [])
      removeToken()
      resolve()
    })
  },

  // Dynamically modify permissions
  changeRoles({ commit, dispatch }, role) {
    return new Promise(async resolve => {
      const token = role + '-token'

      commit('SET_TOKEN', token)
      setToken(token)

      const { roles } = await dispatch('getInfo')

      resetRouter()

      // generate accessible routes map based on roles
      const accessRoutes = await dispatch('permission/generateRoutes', roles, { root: true })

      // dynamically add accessible routes
      router.addRoutes(accessRoutes)

      resolve()
    })
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}
