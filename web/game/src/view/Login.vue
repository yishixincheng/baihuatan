<template>
    <mu-container class="login">
        <mu-flex>
            <h3 class="login-title">登录百花潭</h3>
        </mu-flex>

        <mu-form ref="form" :model="form" class="mu-demo-form" label-position="top">
            <mu-form-item prop="username" label="用户名" :rules="usernameRules">
                <mu-text-field v-model="form.username"></mu-text-field>
            </mu-form-item>
            <mu-form-item prop="password" label="密码" :rules="passwordRules">
                <mu-text-field v-model="form.password"></mu-text-field>
            </mu-form-item>

            <mu-flex justify-content="center" align-items="center">
                <mu-button full-width color="primary" @click="onSubmit">登录</mu-button>
            </mu-flex> 

        </mu-form>

    </mu-container>
  
</template>

<style>

.mu-demo-form {
  width: 100%;
}
.login-title {
    width: 100%;
    text-align: center;
    color: #81c784;
    padding: 30px 0;
}

</style>

<script>
export default {
    name : "login",
    data() {
        return {
            usernameRules: [
                { validate: (val) => !!val, message: '必须填写用户名'},
                { validate: (val) => val.length >= 2, message: '用户名长度大于2'}
            ],
            passwordRules: [
                { validate: (val) => !!val, message:'密码必填'},
            ],
            form: {
                username:'',
                password:'',
            }
        }
    },
    mounted() {

    },
    methods: {
        onSubmit() {
            this.$refs.form.validate().then((result) => {
                console.log('form valid: ', result)
            });

            let data = {
                username: this.form.username,
                password: this.form.password,
            }

            this.$axios.post('/oauth/token?grant_type=password',data, {headers:{'Authorization':'Basic '+ this.$xl.b64EncodeUnicode('bht_user_clientId:123456')}}).then(response => {

                    response = response.data;
                    if (response.access_token) {
                        // 登录成功
                        this.$toast.success('登录成功')
                        localStorage.setItem('token', response.access_token.TokenValue);
                        localStorage.setItem('refreshToken', response.access_token.RefreshToken.TokenValue)
                        this.$axios.defaults.headers.common['Authorization'] = response.access_token.TokenValue;
                        this.$store.commit("login_succ",{token:response.access_token.RefreshToken.TokenValue});
                        this.$router.go(-1)
                    } else {
                        this.$toast.error("用户名或密码错误")
                    }

            }).catch(function (error) { // 请求失败处理
                    console.log(error);
            });

        }
    }
}
</script>