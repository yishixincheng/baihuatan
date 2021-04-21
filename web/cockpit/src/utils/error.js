import alertFunc from '@/utils/alert'
export class MsgError extends Error {
  alert(msg, type, sec) {
    msg = msg || this.message
    type = type || 'error'
    alertFunc(msg, type, sec)
  }
}
