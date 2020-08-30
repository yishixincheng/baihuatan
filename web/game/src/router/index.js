import Vue from 'vue'
import Router from 'vue-router'
import store from '../store'

Vue.use(Router)

const routes =[
    {
      path: '/',
      redirect: '/kpk'
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('../view/Login')
    },
    {
      path: '/kpk',
      name: 'kpk',
      component: () => import('../view/Kpk.vue')
    },
    {
      path: '/kpkroom',
      name: 'kpkroom',
      meta: {
        requireAuth: true,  // 添加该字段，表示进入这个路由是需要登录的
      },
      component: () => import('../view/Kpkroom.vue')
    }
]

const router = new Router({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

router.beforeEach((to,from,next)=>{

  if(to.matched.some(record=>!record.meta.requireAuth||record.meta.homePages)){
      next()
  }else{

    let token=localStorage.getItem("token");

    if(store.state.islogin==1||token ){

       next()

    }else{

      if(to.path=="/login"){

          next()

      }else{
          next();

          next({
              path:'/login'
          })
      }

    }

  }
})

export default router
