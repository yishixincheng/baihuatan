import axios from 'axios'
import { MessageBox } from 'element-ui'
import alertFunc from '@/utils/alert'
import store from '@/store'
import { getToken } from '@/utils/auth'
import { trimChar, hasPrefix, isEmpty } from '@/utils'

let reTryCount = 0
// create an axios instance
const service = axios.create({
  baseURL: '', // url = base url + request url
  // withCredentials: true, // send cookies when cross-domain requests
  timeout: 5000 // request timeout
})

// request interceptor
service.interceptors.request.use(
  config => {
    // do something before request is sent
    if (store.getters.token && !config.headers['Authorization']) {
      // let each request carry token
      // ['X-Token'] is a custom headers key
      // please modify it according to the actual situation
      config.headers['Authorization'] = getToken()
    }
    config.headers['Baihuatan'] = 'v1'
    if (isEmpty(config.oriUrl)) {
      config.oriUrl = config.url
    }
    if (hasPrefix(config.oriUrl, '@')) {
      config.url = '/' + trimChar(config.oriUrl, '@', 'left')
    } else if (hasPrefix(config.oriUrl, '/vue-element-admin')) {
      config.url = process.env.VUE_APP_MOCK_API + config.oriUrl
    } else if (!hasPrefix(config.oriUrl, process.env.VUE_APP_BASE_API)) {
      config.url = process.env.VUE_APP_BASE_API + config.oriUrl
    }
    return config
  },
  error => {
    // do something with request error
    console.log(error) // for debug
    return Promise.reject(error)
  }
)

// response interceptor
service.interceptors.response.use(
  /**
   * If you want to get http information such as headers or status
   * Please return  response => response
  */

  /**
   * Determine the request status by custom code
   * Here is just an example
   * You can also judge the status by HTTP Status Code
   */
  response => {
    const config = response.config
    const res = response.data
    // if the custom code is not 20000, it is judged as an error.
    res.code = Number(res.code) || 0
    if (res.code !== 20000 && res.code !== 200 && res.code !== 0) {
      // 401 token过期， 402 刷新token过期
      if (res.code === 401 && reTryCount === 0) {
        // token过期，需要获取刷新token，并重试请求
        reTryCount++
        return store.dispatch('user/refreshToken').then(() => {
          // 成功刷新token, 则重新发送请求
          reTryCount = 0
          config.headers['Authorization'] = getToken()
          return service({
            url: config.oriUrl,
            method: config.method,
            data: config.data,
            params: config.params,
            headers: config.headers
          })
        })
      }
      // 获取刷新token过期
      if (res.code === 402 || reTryCount > 0) {
        // to re-login
        reTryCount = 0
        MessageBox.confirm('您需要重新登录', '登录提示', {
          confirmButtonText: '重新登录',
          cancelButtonText: '取消',
          type: 'warning'
        }).then(() => {
          store.dispatch('user/resetToken').then(() => {
            location.reload()
          })
        })
        return
      }
      // 错误提示
      alertFunc(res.message || res.msg || 'Error', 'error')
      return Promise.reject(res)
    } else {
      return res
    }
  },
  error => {
    console.log('err' + error) // for debug
    alertFunc(error.message || 'Error', 'error')
    return Promise.reject(error)
  }
)

export default service
