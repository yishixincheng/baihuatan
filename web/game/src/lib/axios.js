import axios from 'axios'

export default function newAxios(params) {

    params = params || {}

    let baseURL = params.baseURL || "/api"
    let timeOut = params.timeOut || 3000
    let headers = params.headers || {"Content-Type":"application/json;chartset=uft-8"}

　　let instance = axios.create({
　　　　baseURL,
　　　　timeOut,
       headers,
　　});

    //拦截request请求
    instance.interceptors.request.use(
        config=>{
            let token = localStorage.getItem("token")
            if (token && !config.headers['Authorization']) {
                config.headers['Authorization'] = token
            }
            config.headers['Baihuatan'] = 'v1'
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

            }
            return response.data || [];

        },
        error => {
            console.log(error)
            return Promise.reject(error.response.status) // 返回接口返回的错误信息
        }
    )

    return instance

}