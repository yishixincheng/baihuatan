<template>
    <mu-container class="container">
        <!-- 房间视图，用户加入房间等待开始 -->

        <div class="room-view" v-if="showRoomView">

            <!-- 背景图片加用户列表 -->

            <mu-row class="join-user-list">

                <mu-col v-for="(item,i) in userList" :key="i" :class="room.ownerUID == item.userID ? 'join-room-ower' : ''">
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

            <div class="join-user-start-tip">

                <span v-if="gameStatus==0">
                    等待其他玩家加入，请稍后...
                </span>
                <span v-else-if="gameStatus==1">
                    准备开始 {{gameCountDown}}
                </span>    

            </div>        

        </div>    

        <!-- 答题视图-->
        <div class="kpk-view" v-if="showKpkView">
            
            <div class="kpk-user-process">

                 <div class="kpk-user-road" v-for="(item,i) in userList" :key="i">
                     <div class="_user_profile">
                         <div class="_avatar">
                             <img :src="item.avatar" />
                         </div>
                         <div class="_username">
                             {{item.userName}}
                         </div>    
                     </div>
                     <!-- 跑道 -->
                     <div class="_road_track">

                         <div class="_ruler">
                             <div class="_scale_line" v-for="(idx,i) in [1,2,3,4,5,6,7,8,9,10,11]" :key="i">
                             </div>    
                         </div>    
                         <div :class="'_road' + ' ' + '_road_' + item.roadID">
                         </div>    
                         <div :class="'_pet' + ' ' + '_pet_' + item.petID"  :id="'Pet_'+item.userID">
                         </div>

                     </div>

                 </div>    

            </div> 

            <section class="kpk-question">

                  <header class="_question_title">
                      <span>{{number}}. </span>
                      <span>{{question.title}}</span>
                  </header>

                  <div class="_options-list">

                      <div @click="selectOption" data-option="A" :class="['_option', resultOptionStyle.A]" wx:if="question.option_1">
                          <i></i>
                          <span>A.</span>
                          <span>{{question.option_1}}</span>
                      </div>

                      <div @click="selectOption" data-option="B" :class="['_option', resultOptionStyle.B]" wx:if="question.option_2">
                          <i></i>
                          <span>B.</span>
                          <span>{{question.option_2}}</span>
                      </div>

                      <div @click="selectOption" data-option="C" :class="['_option', resultOptionStyle.C]" wx:if="question.option_3">
                          <i></i>
                          <span>C.</span>
                          <span>{{question.option_3}}</span>
                      </div>

                      <div @click="selectOption" data-option="D" :class="['_option', resultOptionStyle.D]" wx:if="question.option_4">
                          <i></i>
                          <span>D.</span>
                          <span>{{question.option_4}}</span>
                      </div>

                  </div>

                  <!-- 答题结果对话框 -->
                  <div class="_answer_body" v-if="openAnswerResult">

                      <div class="_answer_result _right" v-if="resultIsRight">
                          回答正确
                      </div>
                      <div class="_answer_result _error" v-else>
                          回答错误 正确答案：{{resultRightChoice}}
                      </div>    

                      <div class="_answer_annotation" v-if="resultAnnotation">
                          {{resultAnnotation}}
                      </div>

                      <div class="_next_btn" @click="nextOption">
                          <a>下一题</a>
                      </div>
                  </div> 


            </section>

            <!-- 答题结束对话框 -->
            <transition name="fade">
                <div class="kpk-gameover-dlg" v-if="openOverGameResult">

                    <div class="_gd_tip">游戏结束</div>
                    <div class="_gd_score">
                        <mu-row>
                            <mu-col>名次</mu-col>
                            <mu-col>玩家</mu-col>
                            <mu-col>得分</mu-col>
                        </mu-row>
                        <mu-row v-for="(item,i) in userRank" :key="i">
                             <mu-col>{{i+1}}</mu-col>
                             <mu-col>{{item.userName}}</mu-col>
                             <mu-col>{{item.score}}</mu-col>
                        </mu-row>        
                    </div>
                    <div class="_gd_back">
                        <span>{{gameCountDown}}秒后系统自动返回</span>
                    </div>
                </div> 
            </transition>   

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

.join-user-start-tip {
    position: absolute;
    bottom: 10px;
    text-align: center;
    width: 100vw;
}

.kpk-gameover-dlg {
    width: 300px;
    height: 220px;
    border-radius: 10px;
    box-shadow: 1px 1px 3px #f8f8f8;
    position: absolute;
    z-index: 20;
    left: calc((100% - 300px)/2);
    top: calc((100% - 220px)/2);
    background: linear-gradient(45deg,  #ffecb3, #ffe0b2);
}

/**
 * pk视图
 */
 .kpk-view {
     display: flex;
     flex-direction: column;
     height: 100vh;
     width: 100vw;
 }
 .kpk-user-road ._user_profile ._avatar {
     width: 32px;
     height: 32px;
     overflow: hidden;
     padding: 2px;
     background: #81c784;
     border-radius: 50%;
 }
 ._avatar img{
     width: 28px;
     height: 28px;
     border-radius: 50%;
 }

 .kpk-user-road {
     height: 60px;
     position: relative;
     padding: 5px;
     display: flex;
     flex-direction: row;
 }
 .kpk-user-road ._user_profile {
     width: 40px;
 }
 .kpk-user-road ._road_track {
     flex: 1;
     position: relative;
 }
 ._road_track ._road {
     height: 40px;
     position: absolute;
     left: 0;
     top: 5px;
     width: 100%;
     opacity: 0.9;
     z-index: 1;
     border-radius: 20px;
 }
 ._road_0 {
     background-size: auto 40px;
     background-image: url('../assets/road_0.png');
 }

 ._road_track ._ruler {
     height: 60px;
     flex: 1;
     padding: 0 20px;
     display: flex; 
     justify-content: space-between;
     align-items: center;
     z-index: 0;
 }

.kpk-user-road:nth-child(even){
    background: #f9fbe7;
 }
 .kpk-user-road:nth-child(odd){
    background: #e8f5e9;
 }

 ._ruler ._scale_line {
     height: 60px;
     width: 1px;
     background: #fbc02d;
     position: relative;
     top: -5px;
 }

 .kpk-user-road:nth-child(even)  ._ruler ._scale_line{
     background: #81c784;
 }

 ._road_track ._pet {
     height: 32px;
     width: 32px;
     position: absolute;
     left: 0;
     top: 8px;
     z-index: 2;
 }
 ._pet_0 {
     background-size: 32px 32px;
     background-image: url('../assets/pet_0.png');
 }
 ._pet_1 {
     background-size: 32px 32px;
     background-image: url('../assets/pet_1.png');
 }

 .kpk-question {
     padding: 20px 10px;
     position: relative;
     flex: 1;
 }
 ._question_title {
     font-size: 20px;
 }
 ._question_title span:nth-child(1) {
     font-size: 21px;
     color: #666;
 }
 ._question_title span:nth-child(2) {
     font-size: 18px;
     color: #101010;
 }

 ._options-list {
     padding: 10px 0;

 }
 ._options-list ._option {
     padding: 8px 8px 8px 28px;
     font-size: 16px;
     position: relative;
 }
  ._options-list ._option span:nth-child(1){
      color: #666;
  }
  ._options-list ._option span:nth-child(2) {
       color: #101010; 
  }

  ._options-list ._option:hover span{
      color: #66bb6a;
  }

  ._options-list ._option._right span{
      color: #66bb6a;
  }
  ._options-list ._option._error span{
      color: #ef5350;
  }

  ._options-list ._option i{
      width: 20px;
      height: 20px;
      display: block;
      position: absolute;
      left: 2px;
      top: 8px;
      background-size: 20px;
  }
  ._options-list ._option._right i {
     background-image: url('../assets/right.png');
  }
  ._options-list ._option._error i {
     background-image: url('../assets/error.png');
  }

  ._answer_body {
      width: 100vw;
      height: 100%;
      position: absolute;
      left: 0;
      top: 0;
      z-index: 10;
  }

  ._answer_result {
      position: absolute;
      bottom: 60px;
      padding: 28px;
      font-size: 16px;
  }
  ._answer_result._right {
      color: #66bb6a;
  }
  ._answer_result._error {
      color: #ef5350;
  }
  ._answer_annotation {
      position: absolute;
      bottom: 25px;
      padding: 28px;
      font-size: 12px;
      color: #999;
  }
  ._next_btn {
      position: absolute;
      right: -60px;
      width: 120px;
      height: 120px;
      background: #66bb6a;
      border-radius: 100%;
      opacity: 0.9;
      top: calc((100% - 120px)/2);
  }
  ._next_btn a{
      display: block;
      width: 100%;
      height: 100%;
      color: #fff;
      line-height: 120px;
      text-align: left;
      font-size: 16px;
      padding-left: 8px;
  }

  ._gd_tip {
      height: 30px;
      padding-top: 10px;
      font-size: 20px;
      text-align: center;
      color: #f57c00;
  }
  ._gd_score {
      text-align: center;
      padding-top: 20px;
  }
  ._gd_back {
      text-align: center;
      padding-top: 30px;
      color: #a1887f;
  }

  .fade-enter-active, .fade-leave-active {
      transition: opacity 1s;
  }
  .fade-enter, .fade-leave-to /* .fade-leave-active below version 2.1.8 */ {
     opacity: 0;
  }

</style>

<script>
import TWEEN from '@tweenjs/tween.js'
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
            room: {},
            gameStatus : 0,
            gameCountDown : 3,
            question: {},
            number: 1,
            popAnswerCard: false, //是否弹出答题结果对话框
            openAnswerResult: false,
            openOverGameResult: false,
            resultOptionStyle:{
                A:'',
                B:'',       // _right
                C:'', // _error
                D:''
            },
            resultIsRight: false,
            resultRightChoice: "",
            resultAnnotation: "",
            userRank: [] // 用户排名
        }
    },
    mounted() {
        // 匹配房间类型
        this.type = this.$route.query.type || 1;
        this.init();
        // 初始化 TweenJs 监听
        this.tweenAni();
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
                this.showKpkView = false
            }
            this.userList = data.userList||[]
            this.room.ownerUID = data.ownerUID
            this.room.roomID = data.roomID
            this.room.myUID  = data.myUID
            this.gameStatus = data.status

            if (this.gameStatus == 1) {
                // 倒计时
                let timer = window.setInterval((e) => {
                    
                    this.gameCountDown --
                    if (this.gameCountDown == 1) {
                        // 房主发送发送开始请求
                        if (this.room.myUID == this.room.ownerUID) {
                            this.sendMessage({
                                method : "start",
                            })
                        }
                        window.clearInterval(timer)
                    }

                }, 1000)
            }


        },
        /**
         * 开始游戏响应
         */
        respStart(data) {
            console.log("开始游戏")
            this.showRoomView = false
            this.showKpkView  = true
            // 问题
            this.question = data.question
            // 序号
            this.number = data.number
        },
        /**
         * 返回新的题目
         */
        respNewQuestion(data) {
            // 问题
            this.openAnswerResult = false
            this.resultAnnotation = ''
            this.resultOptionStyle['A'] = ''
            this.resultOptionStyle['B'] = ''
            this.resultOptionStyle['C'] = ''
            this.resultOptionStyle['D'] = ''
            this.resultRightChoice = ''

            this.question = data.question
            // 序号
            this.number = data.number

        },
        /**
         * 战斗状态响应
         */
        respFightDynamic(data) {

            for (let i in data.userList) {
                this._setPetTrack(data.userList[i])
            }

        },
        _setPetTrack(user) {
            let pace = user.pace      // 步进
            let count = user.count    // 答题总数
            let right = user.right    // 正确个数
            let userID = user.userID  // 用户userID

            let startpos = { x: 0, y: 0 }
            let endpos = {x: 0, y: 0}
            let petTarget = document.getElementById("Pet_" + userID)
            startpos.x = parseInt(petTarget.style.left || 0)

            let rulerObj = document.getElementsByClassName("_ruler")[0]
            let allWidth = rulerObj.offsetWidth || 10
            let diffWidth = parseInt(allWidth / 10)

            endpos.x = pace * diffWidth

            new TWEEN.Tween(startpos) // 传入开始位置
            .to(endpos, 500) // 指定时间内完成结束位置
            .easing(TWEEN.Easing.Quadratic.Out) // 缓动方法名
            .onUpdate((pos, elapsed) => {
                petTarget.style.left = pos.x + "px"
                // 上面的值更新时执行的设置
            }).start();// ================================= 不要忘了合适的时候启动动画
        },
        // TweenJs 动画监听
        tweenAni: function () {
            requestAnimationFrame(this.tweenAni);
            TWEEN.update(); // ================================= 关键是这句
        },
        /**
         * 答题结果
         */
        respAnswerResult(data) {
            let choice = data.choice
            let rightChoice = data.rightChoice
            this.resultAnnotation = data.annotation || ""
            this.resultOptionStyle['A'] = ''
            this.resultOptionStyle['B'] = ''
            this.resultOptionStyle['C'] = ''
            this.resultOptionStyle['D'] = ''
            this.resultRightChoice = rightChoice

            if (choice == rightChoice) {
                // 回答正确
                this.resultIsRight = true
                this.resultOptionStyle[rightChoice] = "_right"
            } else {
                this.resultIsRight = false
                this.resultOptionStyle[choice] = "_error"
                this.resultOptionStyle[rightChoice] = "_right"
            }
            this.openAnswerResult = true
        },
        /**
         * 游戏结束响应
         */
        respGameOver(data) {
            this.openOverGameResult = true
            this.userRank = data.userList
            this.gameCountDown = 5
            let timer = window.setInterval((e) => {
                
                this.gameCountDown --
                if (this.gameCountDown == 1) {
                    // 房主发送发送开始请求
                    window.location = "/kpk"
                    window.clearInterval(timer)
                }

            }, 1000)
        },
        /**
         * 错误消息响应
         */
        respError(data) {
            this.$toast.error(data.msg||"错误")
        },

        sendMessage(data) {
            let message = JSON.stringify(data)
            this.socket.send(message)
        },

        selectOption(e) {
            let choice = e.currentTarget.dataset.option
            this.sendMessage({
                method: 'answer',
                cursor: this.number - 1,
                choice: choice,
            })
        },

        nextOption(e) {
            this.sendMessage({
                method: "nextQuestion"
            })
        }
    },
    destroyed () {
        // 销毁监听
        this.socket.close()
    }
}
</script>