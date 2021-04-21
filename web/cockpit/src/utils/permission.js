import store from '@/store'
import { eachTree, equal, MulReturn, trimChar } from '@/utils/index'

/**
 * @param {Array} value
 * @returns {Boolean}
 * @example see @/views/permission/directive.vue
 */
export default function checkPermission(value) {
  if (value && value instanceof Array && value.length > 0) {
    const roles = store.getters && store.getters.roles
    const permissionRoles = value

    const hasPermission = roles.some(role => {
      return permissionRoles.includes(role)
    })
    return hasPermission
  } else {
    console.error(`need roles! Like v-permission="['admin','editor']"`)
    return false
  }
}

/**
 * 是否包含路由
 * @param routes
 * @returns {MulReturn}
 */
export function hasRoute(...routes) {
  return new MulReturn((resolve) => {
    const isHas = []
    for (let i = 0; i < routes.length; i++) {
      isHas[i] = false
    }
    eachTree(store.getters.permission_routes, node => {
      if (node.fullPath) {
        const nodeFullPath = trimChar(node.fullPath, '/')
        for (let i = 0; i < routes.length; i++) {
          if (equal(nodeFullPath, trimChar(routes[i], '/'))) {
            isHas[i] = true
          }
        }
      }
    })
    resolve(...isHas)
  })
}

