<template>
    <div>
        fdfd
    </div>    
</template>

<style scoped>

</style>

<script>
export default {
    "name" : "kpkroom",
    data() {
        return {
            type: 1,
            path: 'ws://192.168.43.94:9090/gamekpk/ws',
            socket : null,
        }
    },
    mounted() {
        // 匹配房间类型
        this.type = this.$route.query.type || 1;
        this.init();
    },
    methods: {
        init() {
            if(typeof(WebSocket) === "undefined"){
                alert("您的浏览器不支持socket")
            }else{
                // 实例化socket
                this.socket = new WebSocket(this.path + '?joinroomtype='+ this.type, [this.$store.state.token])
                // 监听socket连接
                this.socket.onopen = this.open
                // 监听socket错误信息
                this.socket.onerror = this.error
                // 监听socket消息
                this.socket.onmessage = this.onMessage
                this.socket.onclose = this.close
            }
        },
        open() {
            console.log("已链接上")
        },
        error(e) {
            console.log(e);
        },
        onMessage(msg) {
            console.log(msg)
        },
        send: function () {
            this.socket.send(params)
        },
        close: function () {
            console.log("socket已经关闭")
        }
    },
    destroyed () {
        // 销毁监听
        this.socket.close()
    }
}
</script>