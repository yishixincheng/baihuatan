import { login, logout, getInfo, refreshToken } from '@/api/user'
import { getToken, setToken, removeToken, setRefreshToken, getRefreshToken } from '@/utils/auth'
import router, { resetRouter } from '@/router'
import { isEmpty } from '@/utils'

const state = {
  token: getToken(),
  refreshToken: getRefreshToken(),
  name: '',
  avatar: '',
  roles: [],
  userInfo: null,
  isGatedUser: false
}

const mutations = {
  SET_TOKEN: (state, token) => {
    state.token = token
  },
  SET_REFRESH_TOKEN: (state, refreshToken) => {
    // 刷新令牌
    state.refreshToken = refreshToken
  },
  SET_USERINFO: (state, userInfo) => {
    userInfo.avatar = userInfo.avatar || 'https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif'
    userInfo.name = userInfo.name || userInfo.staffName
    state.userInfo = userInfo
  },
  SET_AUTHROUTELIST: (state, authRoleList) => {
    // 授权的路由列表
    state.authRoleList = authRoleList
  },
  SET_ROLES: (state, roles) => {
    state.roles = roles
  },
  SET_ISADMIN: (state, isAdmin) => {
    state.isAdmin = isAdmin || false
  },
  SET_ISGATEDUSER: (state, isGatedRoute) => {
    state.isGatedUser = isGatedRoute
  }
}

const actions = {
  // user login
  login({ commit }, userInfo) {
    const { username, password } = userInfo
    return new Promise((resolve, reject) => {
      login({ username: username.trim(), password: password }).then(response => {
        if (!isEmpty(response.error)) {
          // 登录失败
          reject(response.error)
          return
        }
        console.log(response)
        const { access_token } = response
        const refreshToken = access_token.RefreshToken ? access_token.RefreshToken.TokenValue : ''
        commit('SET_TOKEN', access_token.TokenValue)
        commit('SET_REFRESH_TOKEN', refreshToken)
        setToken(access_token.TokenValue)
        setRefreshToken(refreshToken)
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },

  // get user info
  getInfo({ commit, state }) {
    return new Promise((resolve, reject) => {
      commit('SET_ISGATEDUSER', true)
      getInfo(state.token).then(response => {
        const { data } = response
        if (!data) {
          reject('验证失败，请重新登录！')
        }
        const { is_admin, user_info, auth_route_list } = data
        commit('SET_ISADMIN', is_admin)
        commit('SET_AUTHROUTELIST', auth_route_list)
        commit('SET_USERINFO', user_info)
        resolve({ is_admin, auth_route_list })
      }).catch(error => {
        reject(error)
      })
    })
  },

  // user logout
  logout({ commit, state, dispatch }) {
    return new Promise((resolve, reject) => {
      logout(state.token).then(() => {
        commit('SET_TOKEN', '')
        commit('SET_ROLES', [])
        commit('SET_ISGATEDUSER', false)
        removeToken()
        resetRouter()

        // reset visited views and cached views
        // to fixed https://github.com/PanJiaChen/vue-element-admin/issues/2485
        dispatch('tagsView/delAllViews', null, { root: true })

        resolve()
      }).catch(error => {
        reject(error)
      })
    })
  },

  // 刷新令牌
  refreshToken({ commit, state }) {
    return new Promise((resolve, reject) => {
      refreshToken(state.refreshToken).then((response) => {
        const { result } = response
        commit('SET_TOKEN', result.token)
        commit('SET_REFRESH_TOKEN', result.refreshToken)
        setToken(result.token)
        setRefreshToken(result.refreshToken)
        resolve()
      })
    })
  },

  // remove token
  resetToken({ commit }) {
    return new Promise(resolve => {
      commit('SET_TOKEN', '')
      commit('SET_REFRESH_TOKEN', '')
      commit('SET_ROLES', [])
      removeToken()
      resolve()
    })
  },

  // dynamically modify permissions
  async changeRoles({ commit, dispatch }, role) {
    const token = role + '-token'

    commit('SET_TOKEN', token)
    setToken(token)

    const { roles } = await dispatch('getInfo')

    resetRouter()

    // generate accessible routes map based on roles
    const accessRoutes = await dispatch('permission/generateRoutes', roles, { root: true })
    // dynamically add accessible routes
    router.addRoutes(accessRoutes)

    // reset visited views and cached views
    dispatch('tagsView/delAllViews', null, { root: true })
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}
