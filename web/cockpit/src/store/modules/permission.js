import { asyncRoutes, constantRoutes } from '@/router'
import { equal, forEach, joinPath, trimChar, isArray } from '@/utils'
/**
 * Use meta.role to determine if the current user has permission
 * @param roles
 * @param route
 */
function hasPermission(permitRoutes, route) {
  if (route.path === '*' || route.path === '/') {
    return true
  }
  const priRoute = (route.meta && route.meta.priRoute) ? route.meta.priRoute : route.fullPath
  if (isArray(priRoute)) {
    for (let i = 0; i < priRoute.length; i++) {
      if (isInRoutes(priRoute[i], permitRoutes)) {
        return true
      }
    }
    return false
  }
  return isInRoutes(priRoute, permitRoutes)
}

/**
 * @param priRoute
 * @param permitRoutes
 * @returns {boolean}
 */
function isInRoutes(priRoute, permitRoutes) {
  if (!isArray(permitRoutes)) {
    return false
  }
  for (let i = 0; i < permitRoutes.length; i++) {
    if (equal(trimChar(priRoute, '/'), trimChar(permitRoutes[i], '/'))) {
      return true
    }
  }
  return false
}

/**
 * Filter asynchronous routing tables by recursion
 * @param routes asyncRoutes
 * @param permitRoutes
 */
export function filterAsyncRoutes(routes, permitRoutes) {
  const res = []
  routes.forEach(route => {
    const tmp = { ...route }
    if (tmp.children && tmp.children.length > 0) {
      tmp.children = filterAsyncRoutes(tmp.children, permitRoutes)
    }
    if (hasPermission(permitRoutes, route) || (tmp.children && tmp.children.length > 0)) {
      // 如果子类全部为隐身状态，则设置父类也为隐身状态
      let hiddenCount = 0
      forEach(tmp.children, v => {
        if (Object.prototype.hasOwnProperty.call(v, 'hidden') && v.hidden) {
          hiddenCount++
        }
      })
      if (hiddenCount > 0 && hiddenCount === tmp.children.length) {
        tmp.hidden = true
      }
      res.push(tmp)
    }
  })
  return res
}

/**
 * 补充路径
 * @param routes
 */
export function replenishPath(routes) {
  forEach(routes, route => {
    if (route.children && route.children.length > 0) {
      forEach(route.children, v => {
        v.parentPath = joinPath(route.parentPath, route.path)
      })
      replenishPath(route.children)
    }
    route.fullPath = joinPath(route.parentPath, route.path)
  })
}

const state = {
  routes: [],
  addRoutes: []
}

const mutations = {
  SET_ROUTES: (state, routes) => {
    state.addRoutes = routes
    state.routes = constantRoutes.concat(routes)
  }
}

const actions = {
  generateRoutes({ commit }, pri) {
    return new Promise(resolve => {
      let accessedRoutes
      replenishPath(asyncRoutes)
      if (pri.is_admin) {
        // 是超级用户
        accessedRoutes = asyncRoutes || []
      } else {
        accessedRoutes = filterAsyncRoutes(asyncRoutes, pri.auth_role_list)
      }
      let rootIndex = -1
      let redirectRoute = ''
      forEach(accessedRoutes, (v, i) => {
        if (v.path === '/') {
          rootIndex = i
          return true
        }
        if (v.path !== '*' && !v.hidden && !redirectRoute) {
          redirectRoute = v.path
          return rootIndex === -1
        }
      })
      if (rootIndex !== -1 && redirectRoute) {
        accessedRoutes[rootIndex].redirect = redirectRoute
      }

      commit('SET_ROUTES', accessedRoutes)
      resolve(accessedRoutes)
    })
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}
