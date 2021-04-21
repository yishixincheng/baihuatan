/**
 * 注入函数
 */
import alertFunc from '@/utils/alert'
import { MsgError } from '@/utils/error'
import { isBoolean, isFunc, isObj, isString, isUndefined, equal, hasKey, pick } from '@/utils'

const ErrLocked = -10
class Response {
  constructor(respData) {
    this.respData = respData || {}
  }
  getData() {
    if (!isUndefined(this.respData.result)) {
      return this.respData.result
    }
    if (!isUndefined(this.respData.data)) {
      return this.respData.data
    }
    return null
  }
  get() {
    return this.respData
  }
  alert(msg, type, sec) {
    msg = msg || this.respData.message || this.respData.msg
    if (!msg) {
      return
    }
    alertFunc(msg, type, sec)
  }
}
class ErrResponse extends Response {
  alert(msg, type, sec) {
    super.alert(msg, type || 'error', sec)
  }
}
class SuccResponse extends Response {
  alert(msg, type, sec) {
    super.alert(msg, type || 'success', sec)
  }
}
function errorTip(msg) {
  alertFunc(msg, 'error')
}

function isLocked(proxy) {
  const key = proxy.__lastLockName || 'default'
  if (!isUndefined(proxy.__locks)) {
    if (isUndefined(proxy.__locks[key])) {
      return false
    }
    if (proxy.__locks[key] === 1) {
      proxy.__locks[key] = 2
      return false
    }
    return true
  }
  return false
}
function unLock(proxy) {
  if (!isUndefined(proxy.__locks)) {
    delete proxy.__locks[proxy.__lastLockName || 'default']
  }
}
function unLoading(proxy) {
  const key = proxy.__lastIsLoadingKey || 'isLoading'
  if (isUndefined(key)) {
    return
  }
  if (proxy[key]) {
    proxy[key] = false
  }
}

async function request(url, params, callBack, failCallback) {
  const isAsync = isFunc(callBack)
  if (isLocked(this)) {
    return isAsync ? null : Promise.reject(new ErrResponse({ text: '上锁中', code: ErrLocked }))
  }
  let promise
  if (isFunc(url)) {
    promise = url(params)
  } else {
    promise = this.$store.dispatch(url, params)
  }
  if (isAsync) {
    promise.then(resp => {
      callBack(new SuccResponse(resp))
    }).catch(error => {
      // 捕获异常
      if (error instanceof MsgError) {
        errorTip(error.message)
      }
      if (isFunc(failCallback)) {
        failCallback(new ErrResponse(error))
      }
      console.log(error, 'async')
    }).finally(_ => {
      unLoading(this)
      unLock(this)
    })
    return
  }
  let _error = null
  const response = await promise.catch(error => {
    if (error instanceof MsgError) {
      errorTip(error.message)
    }
    if (isFunc(failCallback)) {
      failCallback(new ErrResponse(error))
    }
    _error = error
    console.log(error, 'sync')
  }).finally(_ => {
    unLoading(this)
    unLock(this)
  })
  if (isUndefined(response)) {
    return Promise.reject(new ErrResponse(_error))
  }
  return new SuccResponse(response)
}

function lock(name) {
  if (isUndefined(this.__locks)) {
    this.__locks = {}
  }
  name = name || 'default'
  if (this.__locks[name] !== 2) {
    this.__locks[name] = 1 // 1 准备锁; 2 上锁
  }
  this.__lastLockName = name
  return this
}
function loading(key) {
  key = key || 'isLoading'
  if (isUndefined(key)) {
    return this
  }
  this[key] = true
  this.__lastIsLoadingKey = key
  return this
}

/**
 * 验证器
 */
function validator(message, regExp, isAlert) {
  return (rule, value, callback) => {
    let bool
    if (regExp instanceof RegExp) {
      bool = !regExp.test(value)
    } else if (isString(regExp) && regExp) {
      bool = !(new RegExp(regExp)).test(value)
    } else if (isFunc(regExp)) {
      const res = regExp(value)
      if (isBoolean(res)) {
        bool = !res
      } else if (isObj(res)) {
        if (hasKey(res, 'message')) {
          message = res.message
        }
        bool = !res.valid
      }
    } else {
      bool = equal(value, '')
    }
    if (bool) {
      if (isAlert === true) {
        this.$message({
          message: message,
          type: 'warning',
          duration: 3000
        })
      }
      callback(new Error(message))
    } else {
      callback()
    }
  }
}

/**
 * 跳转页面
 * @param url
 */
function go(url) {
  if (isUndefined(url)) {
    return
  }
  if (equal(url, -1)) {
    if (this.$store.getters.preRouter) {
      this.$router.push(pick(this.$store.getters.preRouter, ['path', 'query', 'params']))
    } else {
      this.$router.go(-1)
    }
    return
  }
  this.$router.push(url)
}

export default {
  install(vue) {
    vue.prototype.request = request
    vue.prototype.lock = lock
    vue.prototype.loading = loading
    vue.prototype.validator = validator
    vue.prototype.go = go
  }
}
