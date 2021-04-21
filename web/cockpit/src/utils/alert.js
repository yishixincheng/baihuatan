import { Message } from 'element-ui'

let errPreAlert = null

export default function(msg, type, sec) {
  if (!msg) {
    return
  }
  type = type || 'info'
  sec = sec || 3
  if (type === 'error' && errPreAlert) {
    errPreAlert.close()
  }
  const _ = Message({
    message: msg,
    type: type,
    duration: sec * 1000
  })
  if (type === 'error') {
    errPreAlert = _
  }
}
