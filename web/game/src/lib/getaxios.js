import axios from 'axios'

export default function getAxios() {

　　let instance = axios.create({
　　　　baseURL:"/api",
　　　　timeOut:3000,
       headers:{
           post:{"Content-Type":"application/json"},
       }
　　});

    //拦截request请求
    instance.interceptors.request.use(
        config=>{
            let token = localStorage.getItem("token")
            if (token) {
                config.headers['Authorization'] = token
            }
            return config;
        },
        err => {
            return Promise.reject(err)
        }
    )

    //拦截response请求
    instance.interceptors.response.use(
        response=>{

            if(response.status==200){

                const data=response.data;

                // if (data.code === 400){
                //     //登录过期,权限不足
                //     console.warn("登陆过期");
                //     //清除token
                //     store.commit('setToken','')
                //     window.localStorage.removeItem('token')
                //     //跳转登录
                //     router.replace({
                //         path:"/login"
                //     })
                // }

                return data;

            }

            return response;

        },
        error => {
            return Promise.reject(error.response.status) // 返回接口返回的错误信息
        }
    )

    return instance

}