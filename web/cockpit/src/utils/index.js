/**
 * Created by PanJiaChen on 16/11/18.
 */

/**
 * @param {string} path
 * @returns {Boolean}
 */
export function isExternal(path) {
  return /^(https?:|mailto:|tel:)/.test(path)
}

/**
 * @param {string} url
 * @returns {Boolean}
 */
export function validURL(url) {
  const reg = /^(https?|ftp):\/\/([a-zA-Z0-9.-]+(:[a-zA-Z0-9.&%$-]+)*@)*((25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9]?[0-9])){3}|([a-zA-Z0-9-]+\.)*[a-zA-Z0-9-]+\.(com|edu|gov|int|mil|net|org|biz|arpa|info|name|pro|aero|coop|museum|[a-zA-Z]{2}))(:[0-9]+)*(\/($|[a-zA-Z0-9.,?'\\+&%$#=~_-]+))*$/
  return reg.test(url)
}

/**
 * @param {string} str
 * @returns {Boolean}
 */
export function validLowerCase(str) {
  const reg = /^[a-z]+$/
  return reg.test(str)
}

/**
 * @param {string} str
 * @returns {Boolean}
 */
export function validUpperCase(str) {
  const reg = /^[A-Z]+$/
  return reg.test(str)
}

/**
 * @param {string} str
 * @returns {Boolean}
 */
export function validAlphabets(str) {
  const reg = /^[A-Za-z]+$/
  return reg.test(str)
}

/**
 * @param {string} email
 * @returns {Boolean}
 */
export function validEmail(email) {
  const reg = /^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/
  return reg.test(email)
}

export function isNumber(v) {
  if (isUndefined(v)) {
    return false
  }
  if (Number.isInteger(v)) {
    return true
  }
  return isString(v) && /^\s*\d+\s*$/.test(v)
}

export function isFloat(v) {
  if (typeIs(v) !== '[object Number]') {
    return false
  }
  return /^\s*\d+\.\d+\s*$/.test(v)
}

/**
 * @param {string} str
 * @returns {Boolean}
 */
export function isString(str) {
  return typeof str === 'string' || str instanceof String
}

export function isBoolean(bool) {
  return bool === true || bool === false
}

export function isBaseType(val) {
  return isNumber(val) || isString(val) || isBoolean(val) || isFloat(val)
}

/**
 * @param {Array} arg
 * @returns {Boolean}
 */
export function isArray(arg) {
  if (typeof Array.isArray === 'undefined') {
    return Object.prototype.toString.call(arg) === '[object Array]'
  }
  return Array.isArray(arg)
}

function typeIs(type) {
  return Object.prototype.toString.call(type)
}

/**
 * @param func
 * @returns {boolean}
 */
export function isFunc(func) {
  return typeIs(func) === '[object Function]'
}

/**
 * @param obj
 * @returns {boolean}
 */
export function isUndefined(obj) {
  return typeof obj === 'undefined'
}

/**
 * @param obj
 * @returns {boolean|boolean}
 */
export function isObj(obj) {
  return obj !== null && typeIs(obj) === '[object Object]'
}

/**
 * @param v
 * @param eZ 默认为 false, 代表如果时0，则返回真
 * @returns {boolean | boolean}
 */
export function isEmpty(v, eZ) {
  if (isUndefined(v)) {
    return true
  }
  if (isString(v)) {
    return /^\s*$/g.test(v) || (!eZ && v === '0')
  }
  if (!eZ && Number.isInteger(v)) {
    return v === 0
  }
  if (isArray(v)) {
    return v.length === 0
  }
  return v === null
}

/**
 * Parse the time to string
 * @param {(Object|string|number)} time
 * @param {string} cFormat
 * @returns {string | null}
 */
export function parseTime(time, cFormat) {
  if (arguments.length === 0 || !time) {
    return null
  }
  const format = cFormat || '{y}-{m}-{d} {h}:{i}:{s}'
  let date
  if (typeof time === 'object') {
    date = time
  } else {
    if ((typeof time === 'string')) {
      if ((/^[0-9]+$/.test(time))) {
        // support "1548221490638"
        time = parseInt(time)
      } else {
        // support safari
        // https://stackoverflow.com/questions/4310953/invalid-date-in-safari
        time = time.replace(new RegExp(/-/gm), '/')
      }
    }

    if ((typeof time === 'number') && (time.toString().length === 10)) {
      time = time * 1000
    }
    date = new Date(time)
  }
  const formatObj = {
    y: date.getFullYear(),
    m: date.getMonth() + 1,
    d: date.getDate(),
    h: date.getHours(),
    i: date.getMinutes(),
    s: date.getSeconds(),
    a: date.getDay()
  }
  const time_str = format.replace(/{([ymdhisa])+}/g, (result, key) => {
    const value = formatObj[key]
    // Note: getDay() returns 0 on Sunday
    if (key === 'a') { return ['日', '一', '二', '三', '四', '五', '六'][value ] }
    return value.toString().padStart(2, '0')
  })
  return time_str
}

/**
 * @param {number} time
 * @param {string} option
 * @returns {string}
 */
export function formatTime(time, option) {
  if (('' + time).length === 10) {
    time = parseInt(time) * 1000
  } else {
    time = +time
  }
  const d = new Date(time)
  const now = Date.now()

  const diff = (now - d) / 1000

  if (diff < 30) {
    return '刚刚'
  } else if (diff < 3600) {
    // less 1 hour
    return Math.ceil(diff / 60) + '分钟前'
  } else if (diff < 3600 * 24) {
    return Math.ceil(diff / 3600) + '小时前'
  } else if (diff < 3600 * 24 * 2) {
    return '1天前'
  }
  if (option) {
    return parseTime(time, option)
  } else {
    return (
      d.getMonth() +
      1 +
      '月' +
      d.getDate() +
      '日' +
      d.getHours() +
      '时' +
      d.getMinutes() +
      '分'
    )
  }
}

/**
 * @param {string} url
 * @returns {Object}
 */
export function getQueryObject(url) {
  url = url == null ? window.location.href : url
  const search = url.substring(url.lastIndexOf('?') + 1)
  const obj = {}
  const reg = /([^?&=]+)=([^?&=]*)/g
  search.replace(reg, (rs, $1, $2) => {
    const name = decodeURIComponent($1)
    let val = decodeURIComponent($2)
    val = String(val)
    obj[name] = val
    return rs
  })
  return obj
}

/**
 * @param {string} input value
 * @returns {number} output value
 */
export function byteLength(str) {
  // returns the byte length of an utf8 string
  let s = str.length
  for (var i = str.length - 1; i >= 0; i--) {
    const code = str.charCodeAt(i)
    if (code > 0x7f && code <= 0x7ff) s++
    else if (code > 0x7ff && code <= 0xffff) s += 2
    if (code >= 0xDC00 && code <= 0xDFFF) i--
  }
  return s
}

/**
 * @param {Array} actual
 * @returns {Array}
 */
export function cleanArray(actual) {
  const newArray = []
  for (let i = 0; i < actual.length; i++) {
    if (actual[i]) {
      newArray.push(actual[i])
    }
  }
  return newArray
}

/**
 * @param {Object} json
 * @returns {Array}
 */
export function param(json) {
  if (!json) return ''
  return cleanArray(
    Object.keys(json).map(key => {
      if (json[key] === undefined) return ''
      return encodeURIComponent(key) + '=' + encodeURIComponent(json[key])
    })
  ).join('&')
}

/**
 * @param {string} url
 * @returns {Object}
 */
export function param2Obj(url) {
  const search = decodeURIComponent(url.split('?')[1]).replace(/\+/g, ' ')
  if (!search) {
    return {}
  }
  const obj = {}
  const searchArr = search.split('&')
  searchArr.forEach(v => {
    const index = v.indexOf('=')
    if (index !== -1) {
      const name = v.substring(0, index)
      const val = v.substring(index + 1, v.length)
      obj[name] = val
    }
  })
  return obj
}

/**
 * @param {string} val
 * @returns {string}
 */
export function html2Text(val) {
  const div = document.createElement('div')
  div.innerHTML = val
  return div.textContent || div.innerText
}

/**
 * Merges two objects, giving the last one precedence
 * @param {Object} target
 * @param {(Object|Array)} source
 * @returns {Object}
 */
export function objectMerge(target, source) {
  if (typeof target !== 'object') {
    target = {}
  }
  if (Array.isArray(source)) {
    return source.slice()
  }
  Object.keys(source).forEach(property => {
    const sourceProperty = source[property]
    if (typeof sourceProperty === 'object') {
      target[property] = objectMerge(target[property], sourceProperty)
    } else {
      target[property] = sourceProperty
    }
  })
  return target
}

/**
 * @param {HTMLElement} element
 * @param {string} className
 */
export function toggleClass(element, className) {
  if (!element || !className) {
    return
  }
  let classString = element.className
  const nameIndex = classString.indexOf(className)
  if (nameIndex === -1) {
    classString += '' + className
  } else {
    classString =
      classString.substr(0, nameIndex) +
      classString.substr(nameIndex + className.length)
  }
  element.className = classString
}

/**
 * @param {string} type
 * @returns {Date}
 */
export function getTime(type) {
  if (type === 'start') {
    return new Date().getTime() - 3600 * 1000 * 24 * 90
  } else {
    return new Date(new Date().toDateString())
  }
}

/**
 * @param {Function} func
 * @param {number} wait
 * @param {boolean} immediate
 * @return {*}
 */
export function debounce(func, wait, immediate) {
  let timeout, args, context, timestamp, result

  const later = function() {
    // 据上一次触发时间间隔
    const last = +new Date() - timestamp

    // 上次被包装函数被调用时间间隔 last 小于设定时间间隔 wait
    if (last < wait && last > 0) {
      timeout = setTimeout(later, wait - last)
    } else {
      timeout = null
      // 如果设定为immediate===true，因为开始边界已经调用过了此处无需调用
      if (!immediate) {
        result = func.apply(context, args)
        if (!timeout) context = args = null
      }
    }
  }

  return function(...args) {
    context = this
    timestamp = +new Date()
    const callNow = immediate && !timeout
    // 如果延时不存在，重新设定延时
    if (!timeout) timeout = setTimeout(later, wait)
    if (callNow) {
      result = func.apply(context, args)
      context = args = null
    }

    return result
  }
}

/**
 * This is just a simple version of deep copy
 * Has a lot of edge cases bug
 * If you want to use a perfect deep copy, use lodash's _.cloneDeep
 * @param {Object} source
 * @returns {Object}
 */
export function deepClone(source) {
  if (!source && typeof source !== 'object') {
    throw new Error('error arguments', 'deepClone')
  }
  const targetObj = source.constructor === Array ? [] : {}
  Object.keys(source).forEach(keys => {
    if (source[keys] && typeof source[keys] === 'object') {
      targetObj[keys] = deepClone(source[keys])
    } else {
      targetObj[keys] = source[keys]
    }
  })
  return targetObj
}

/**
 * @param {Array} arr
 * @returns {Array}
 */
export function uniqueArr(arr) {
  return Array.from(new Set(arr))
}

/**
 * @returns {string}
 */
export function createUniqueString() {
  const timestamp = +new Date() + ''
  const randomNum = parseInt((1 + Math.random()) * 65536) + ''
  return (+(randomNum + timestamp)).toString(32)
}

/**
 * Check if an element has a class
 * @param {HTMLElement} elm
 * @param {string} cls
 * @returns {boolean}
 */
export function hasClass(ele, cls) {
  return !!ele.className.match(new RegExp('(\\s|^)' + cls + '(\\s|$)'))
}

/**
 * Add class to element
 * @param {HTMLElement} elm
 * @param {string} cls
 */
export function addClass(ele, cls) {
  if (!hasClass(ele, cls)) ele.className += ' ' + cls
}

/**
 * Remove class from element
 * @param {HTMLElement} elm
 * @param {string} cls
 */
export function removeClass(ele, cls) {
  if (hasClass(ele, cls)) {
    const reg = new RegExp('(\\s|^)' + cls + '(\\s|$)')
    ele.className = ele.className.replace(reg, ' ')
  }
}

/**
 * @param list
 * @param func
 */
export function forEach(list, func) {
  if (!list && typeof list !== 'object') {
    return
  }
  for (const i in list) {
    if (!Object.prototype.hasOwnProperty.call(list, i)) {
      continue
    }
    const ret = func(list[i], i)
    if (ret === false) {
      break
    }
  }
}

/**
 * @param v1
 * @param v2
 * @returns {boolean}
 */
export function equal(v1, v2) {
  if (v1 === v2) {
    return true
  }
  if (isBaseType(v1) && isBaseType(v2)) {
    return v1.toString() === v2.toString()
  }
  if (isArray(v1) && isArray(v2) && v1.length === v2.length) {
    for (const i in v1) {
      if (!equal(v1[i], v2[i])) {
        return false
      }
    }
    return true
  }
  if (isObj(v1) && isObj(v2) && Object.keys(v1).length === Object.keys(v2).length) {
    let isEqual = true
    forEach(v1, (v, k) => {
      if (!equal(v, v2[k])) {
        isEqual = false
        return false // 跳出循环
      }
    })
    return isEqual
  }
  return false
}

/**
 * 摘取对象部分属性
 * @param base
 * @param keys
 */
export function pick(base, keys) {
  const newObj = {}
  forEach(keys, key => {
    const pos = key.indexOf(':')
    let dataType = null
    if (pos !== -1) {
      key = key.substr(0, pos)
      dataType = key.substr(0, pos + 1)
    }
    const match = key.match(/^(.+?)\s+[Aa][Ss]\s+(.+)$/)
    let newKey = ''
    if (match) {
      newKey = match[2]
      key = match[1]
    } else {
      newKey = key
    }
    if (hasKey(base, key)) {
      switch (dataType) {
        case 'Number':
          newObj[newKey] = Number(base[key])
          break
        case 'String':
          newObj[newKey] = String(base[key])
          break
        case 'Boolean':
          newObj[newKey] = Boolean(base[key])
          break
        case null:
          newObj[newKey] = base[key]
          break
      }
    }
  })
  return newObj
}

export function findInListPos(find, arr, key) {
  let pos = -1
  forEach(arr, (v, idx) => {
    if ((key && equal(v[key], find)) || (isUndefined(key) && equal(v, find))) {
      pos = idx
      return false
    }
  })
  return pos
}

export function inArray(val, arr, key) {
  return findInListPos(val, arr, key) !== -1
}

// 获取数组元素中含有关键字的列表
export function fetchListHasKeyword(list, keyword, ...keys) {
  const reg = new RegExp(keyword)
  const afterList = []
  forEach(list, v => {
    let isFind = false
    forEach(keys, k => {
      let matchStr = v
      if (isObj(v) || isArray(v)) {
        matchStr = v[k]
      }
      if (reg.test(matchStr)) {
        isFind = true
        return false
      }
    })
    if (isFind) {
      afterList.push(v)
    }
  })
  return afterList
}
// 过滤数组中掉包含value的元素
export function filterValueFromList(list, value, ...keys) {
  const afterList = []
  if (keys.length === 0) {
    keys = [undefined]
  }
  forEach(list, v => {
    let isFind = false
    forEach(keys, k => {
      let matchStr = v
      if (!isUndefined(k) && (isObj(v) || isArray(v))) {
        matchStr = v[k]
      }
      if (isObj(value) || isArray(value)) {
        if (findInListPos(matchStr, value) !== -1) {
          isFind = true
          return false
        }
      } else if (equal(value, matchStr)) {
        isFind = true
        return false
      }
    })
    if (!isFind) {
      afterList.push(v)
    }
  })
  return afterList
}

// 摘取数组元素中key的元素，组成新的数组
export function arrayCloumn(arr, key, dealFunc) {
  const tmp = []
  forEach(arr, v => {
    if (!isObj(v) || isUndefined(v[key])) {
      return true
    }
    tmp.push(isFunc(dealFunc) ? dealFunc(v[key]) : v[key])
  })
  return tmp
}
export function trimChar(str, char, type) {
  if (isUndefined(str)) {
    str = ''
  }
  if (char) {
    if (type === 'left') {
      return str.replace(new RegExp('^' + char + '+', 'g'), '')
    } else if (type === 'right') {
      return str.replace(new RegExp(char + '+$', 'g'), '')
    }
    return str.replace(new RegExp('^' + char + '+|' + char + '+$', 'g'), '')
  }
  return str.replace(/^\s+|\s+$/g, '')
}

/**
 * 连接路径
 * @param path1
 * @param path2
 * @returns {string}
 */
export function joinPath(path1, path2) {
  const p1 = trimChar(path1 || '', '/')
  const p2 = trimChar(path2 || '', '/')
  return p1 ? (p1 + (p2 ? ('/' + p2) : '')) : p2
}

/**
 * @param str
 * @param suffix
 * @returns {boolean|boolean}
 */
export function hasSuffix(str, suffix) {
  if (isUndefined(suffix) || isUndefined(str)) {
    return false
  }
  if (!isString(str)) {
    str = str.toString()
  }
  if (!isString(suffix)) {
    suffix = suffix.toString()
  }
  return str.length >= suffix.length && equal(str.substr(str.length - suffix.length), suffix)
}

/**
 * @param str
 * @param prefix
 * @returns {boolean|boolean}
 */
export function hasPrefix(str, prefix) {
  if (isUndefined(prefix) || isUndefined(str)) {
    return false
  }
  if (!isString(str)) {
    str = str.toString()
  }
  if (!isString(prefix)) {
    prefix = prefix.toString()
  }
  return str.length >= prefix.length && equal(str.substr(0, prefix.length), prefix)
}

/**
 * @param obj
 * @param key
 * @returns {boolean}
 */
export function hasKey(obj, key) {
  return Object.prototype.hasOwnProperty.call(obj, key)
}

/**
 * 遍历树节点
 * @param tree
 * @param eachFunc
 * @returns {boolean}
 */
export function eachTree(tree, eachFunc) {
  let isBreak = false
  forEach(tree, (node, index) => {
    if (eachFunc(node, index) === false) {
      // 跳出循环
      isBreak = true
      return false
    }
    if (node.children && node.children.length > 0) {
      if (eachTree(node.children, eachFunc) === true) {
        return false
      }
    }
  })
  return isBreak
}

/**
 * @param url
 * @returns
 */
export function parseUrl(url) {
  url = url.trim()
  const m = url.match(/^(https?):\/\/([^:]+)(:(\d+))?(\/([^?#]+))(\?([^#]+))?(#(.+))?$/)
  if (!m) {
    return { url: '', path: url, value: url }
  }
  return { url, protocol: m[1], host: m[2], port: m[4], path: m[5], value: m[6], query: m[8], hash: m[10] }
}

/**
 * @param list
 * @param perCount
 * @returns {[]|*[]}
 */
export function arrayGroup(list, perCount) {
  if (perCount <= 0) {
    return [list]
  }
  const len = list.length
  if (len === 0) {
    return []
  }
  return Array.from(Array(Math.ceil(len / perCount)))
    .map((_, i) => list.slice(i * perCount, (i + 1) * perCount))
}

/**
 * 转化成逆序索引
 * @param index
 * @param page
 * @param limit
 * @param total
 */
export function toROrder(index, page, limit, total) {
  return total - (page - 1) * limit - index
}

export function b64EncodeUnicode(str) {
  return btoa(encodeURIComponent(str).replace(/%([0-9A-F]{2})/g, function(match, p1) { return String.fromCharCode('0x' + p1) }))
}

/**
 * 多值返回
 */
export class MulReturn {
  constructor(exec) {
    this.exec = exec
  }
  then(f) {
    this.exec((...rets) => {
      f(...rets)
    })
  }
}
