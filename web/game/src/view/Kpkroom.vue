<template>
    <mu-container class="container">
        <!-- 房间视图，用户加入房间等待开始 -->

        <div class="room-view" v-if="showRoomView">

            <!-- 背景图片加用户列表 -->

            <mu-row class="join-user-list">

                <mu-col v-for="(item,i) in userList" :key="i" :class="{'join-room-ower':room.ownerUID}">
                    <div class="join-user-avatar">
                       <img :src="item.avatar" />
                    </div>
                    <div class="join-user-name">
                        {{item.userName}}
                    </div>
                    <div class="join-user-rank">
                        {{item.rank}}
                    </div>    
                </mu-col>    

            </mu-row>    

        </div>    

        <!-- 答题视图-->
        <div class="kpk-view" v-if="showKpkView">
            答题视图
        </div>    

    </mu-container>    
</template>

<style scoped>

.container {
    width: 100vw;
    height: 100vh;
    padding: 0;
    margin: 0;
    overflow: hidden;
    position: relative;
}

.room-view {
    width: 100vw;
    height: 100vh;
    position: relative;
    background: #e8e8e8;
}
.join-user-list {
    position: absolute;
    height: 120px;
    bottom: 80px;
    left: 20px;
}
/**.col */
.join-user-list .col {
    width: 80px;
    text-align: center;
}

/**头像 */
.join-user-avatar {
    width: 64px;
    height: 64px;
    padding: 2px;
    background: #fff;
    overflow: hidden;
    border-radius: 50%;
    margin: auto;
}
.join-user-avatar img {
    width: 60px;
    height: 60px;
    border-radius: 50%;
}
.join-user-name {
    font-size: 18px;
}
.join-user-rank {
    font-size: 12px;
}
.join-room-ower .join-user-avatar {
    background: #fbc02d;
}



</style>

<script>
export default {
    "name" : "kpkroom",
    data() {
        return {
            type: 1,
            path: 'ws://192.168.43.94:9090/gamekpk/ws',
            socket : null,
            showRoomView: false,
            showKpkView: false,
            userList: [],  // 用户列表
            room: {}
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
        onMessage(event) {

            try {
                let data = JSON.parse(event.data)
                if (!data) {
                    console.log("非json格式")
                    console.log(event.data)
                    return;
                }
                console.log(data)
                switch(data.method) {
                    case "userjoin":
                        // 用户加入
                        this.respUserJoin(data)
                        break;
                    case "start":
                        // 开始游戏
                        this.respStart(data)
                        break;
                    case "newquestion":
                        // 新的题目
                        this.respNewQuestion(data)
                        break;
                    case "fightdynamic":
                        // 动态
                        this.respFightDynamic(data)
                        break;
                    case "gameover":
                        // 游戏结束
                        this.respGameOver(data)
                        break;
                    case "answerresult":
                        // 答题结果
                        this.respAnswerResult(data)    
                        break;
                    case "error":
                        this.respError(data)
                        break;                            
                }

            } catch (err) {
                console.log(err)
            }
           

        },
        send: function () {
            this.socket.send(params)
        },
        close: function () {
            console.log("socket已经关闭")
        },

        // 用户加入响应
        respUserJoin(data) {
            console.log("用户加入")
            if (data.status > 1) {
                // 已进入pk模式
                return;
            }
            if (!this.showRoomView) {
                this.showRoomView = true
            }
            this.userList = data.userList||[]
            console.log(this.userList)
            this.room.ownerUID = data.ownerUID
            this.room.roomID = data.roomID

        },
        /**
         * 开始游戏响应
         */
        respStart(data) {


        },
        /**
         * 返回新的题目
         */
        respNewQuestion(data) {

        },
        /**
         * 战斗状态响应
         */
        respFightDynamic(data) {

        },
        /**
         * 答题结果
         */
        respAnswerResult(data) {

        },
        /**
         * 游戏结束响应
         */
        respGameOver(data) {

        },
        /**
         * 错误消息响应
         */
        respError(data) {
            this.$toast.error(data.msg||"错误")
        }
    },
    destroyed () {
        // 销毁监听
        this.socket.close()
    }
}
</script>