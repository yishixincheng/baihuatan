// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import router from './router'
import Xl from './lib/xl.js'
import axios from 'axios'
import store from './store'
import MuseUI from 'muse-ui'
import Toast from 'muse-ui-toast'
import 'muse-ui/dist/muse-ui.css'
import theme from 'muse-ui/lib/theme'
import 'typeface-roboto'
import less from 'less'

Vue.config.productionTip = false
axios.defaults.headers.post["Content-type"] = "application/json"
axios.defaults.baseURL="/api"

Vue.prototype.$axios = axios
Vue.prototype.$xl = Xl

Vue.use(less)
Vue.use(MuseUI)
Vue.use(Toast)

theme.add("bhtstyle", {
  primary: '#43a047',
  secondary: '#66bb6a',
  success: '#81c784',
  warning: '#f9a825',
  info: '#2196f3',
  error: '#f44336',
  track: '#bdbdbd',
  text: {
    primary: 'rgba(0, 0, 0, 0.87)',
    secondary: 'gba(0, 0, 0, 0.54)',
    alternate: '#fff',
    disabled: 'rgba(0, 0, 0, 0.38)',
    hint: 'rgba(0, 0, 0, 0.38)' // 提示文字颜色
  },
  divider: 'rgba(0, 0, 0, 0.12)',
  background: {
    paper: '#fff',
    chip: '#e0e0e0',
    default: '#fafafa'
  }
})
theme.use("bhtstyle")

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  store,
  components: { App },
  template: '<App/>',
  created() {
    let token = localStorage.getItem("token")
    this.$axios.defaults.headers.common['Authorization'] = token
    
    //在页面加载时读取sessionStorage里的状态信息
    if (sessionStorage.getItem("store") ) {
        this.$store.replaceState(Object.assign({}, this.$store.state,JSON.parse(sessionStorage.getItem("store"))))
    }

    //在页面刷新时将vuex里的信息保存到sessionStorage里
    window.addEventListener("beforeunload",()=>{
        sessionStorage.setItem("store",JSON.stringify(this.$store.state))
    })

    console.log("app created")
  }
})
