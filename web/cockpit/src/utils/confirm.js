import { MessageBox } from 'element-ui'

export default function(message, title, confirmText, cancelText, type) {
  return MessageBox.confirm(message, title || '', {
    confirmButtonText: confirmText || '确定',
    cancelButtonText: cancelText || '取消',
    type: type || 'warning'
  })
}
