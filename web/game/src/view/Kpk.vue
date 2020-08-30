<template>
    <div>
        <mu-flex class="userarea">

                <!-- 头像 -->
                <div class="user-avtar">
                    <img src="https://ss0.bdstatic.com/70cFvHSh_Q1YnxGkpoWK1HF6hhy/it/u=2448627926,3152967872&fm=26&gp=0.jpg" />
                </div>

                <div class="user-info">
                    <mu-row class="user-name">
                        <mu-col span="12">一世心城</mu-col>
                    </mu-row>
                    <mu-row class="user-level">
                        <span>等级：</span>
                        <span>沧海遗珠</span>
                    </mu-row>    
                </div>    

        </mu-flex> 
        <mu-container class="joinroom">

            <mu-row>
                <mu-col span="6">
                    <div class="grid-cell" @click="joinGame($event, 1)">
                        <div class="p">快速匹配</div>
                        <span>2人PK</span>
                    </div>
                </mu-col>
                <mu-col span="6">
                    <div class="grid-cell _two" @click="joinGame($event, 2)">
                        <div class="p">快速匹配</div>
                        <span>4人PK</span>
                    </div>
                </mu-col>
            </mu-row>

        </mu-container>

        <mu-container class="roomlist">

            <mu-row>
                <mu-col span="4">
                    <div class="room-node _status_1">
                        <ul>
                            <li><img src="https://ss0.bdstatic.com/70cFvHSh_Q1YnxGkpoWK1HF6hhy/it/u=2946010582,1303971542&fm=26&gp=0.jpg" /></li>
                            <li><img src="https://ss0.bdstatic.com/70cFuHSh_Q1YnxGkpoWK1HF6hhy/it/u=1279838220,3182144174&fm=26&gp=0.jpg" /></li>
                            <li><img src="https://ss0.bdstatic.com/70cFuHSh_Q1YnxGkpoWK1HF6hhy/it/u=1175288209,3211190587&fm=26&gp=0.jpg" /></li>
                            <li><img src="https://ss0.bdstatic.com/70cFuHSh_Q1YnxGkpoWK1HF6hhy/it/u=1175288209,3211190587&fm=26&gp=0.jpg" /></li>
                        </ul>    
                        <div>
                        <span>已加入4人</span>
                        <span>未开始</span>
                        </div>
                    </div>    
                </mu-col>
                <mu-col span="4">
                    <div class="room-node _status_2">
                        <ul>
                            <li><img src="https://ss0.bdstatic.com/70cFvHSh_Q1YnxGkpoWK1HF6hhy/it/u=1669872726,1649722464&fm=26&gp=0.jpg" /></li>
                            <li><img src="https://ss0.bdstatic.com/70cFuHSh_Q1YnxGkpoWK1HF6hhy/it/u=1175288209,3211190587&fm=26&gp=0.jpg" /></li>
                        </ul>    
                        <div>
                        <span>已加入4人</span>
                        <span>已开始</span>
                        </div>
                    </div>    
                </mu-col>
                <mu-col span="4">
                    <div class="room-node _status_1">
                        <ul>
                            <li><img src="https://ss0.bdstatic.com/70cFvHSh_Q1YnxGkpoWK1HF6hhy/it/u=1669872726,1649722464&fm=26&gp=0.jpg" /></li>
                            <li><img src="https://ss0.bdstatic.com/70cFuHSh_Q1YnxGkpoWK1HF6hhy/it/u=1175288209,3211190587&fm=26&gp=0.jpg" /></li>
                        </ul>    
                        <div>
                        <span>已加入4人</span>
                        <span>未开始</span>
                        </div>
                    </div>    
                </mu-col>        
            </mu-row>

            <mu-row>
                <mu-col span="4">
                    <div class="room-node _status_2">
                        <ul>
                            <li><img src="https://ss0.bdstatic.com/70cFvHSh_Q1YnxGkpoWK1HF6hhy/it/u=2946010582,1303971542&fm=26&gp=0.jpg" /></li>
                            <li><img src="https://ss0.bdstatic.com/70cFuHSh_Q1YnxGkpoWK1HF6hhy/it/u=1279838220,3182144174&fm=26&gp=0.jpg" /></li>
                            <li><img src="https://ss0.bdstatic.com/70cFuHSh_Q1YnxGkpoWK1HF6hhy/it/u=1175288209,3211190587&fm=26&gp=0.jpg" /></li>
                        </ul>    
                        <div>
                        <span>已加入4人</span>
                        <span>未开始</span>
                        </div>
                    </div>    
                </mu-col>
                <mu-col span="4">
                    <div class="room-node _status_1">
                        <ul>
                            <li><img src="https://ss0.bdstatic.com/70cFvHSh_Q1YnxGkpoWK1HF6hhy/it/u=1669872726,1649722464&fm=26&gp=0.jpg" /></li>
                            <li><img src="https://ss0.bdstatic.com/70cFuHSh_Q1YnxGkpoWK1HF6hhy/it/u=1175288209,3211190587&fm=26&gp=0.jpg" /></li>
                            <li><img src="https://ss0.bdstatic.com/70cFuHSh_Q1YnxGkpoWK1HF6hhy/it/u=1175288209,3211190587&fm=26&gp=0.jpg" /></li>
                        </ul>    
                        <div>
                        <span>已加入4人</span>
                        <span>已开始</span>
                        </div>
                    </div>    
                </mu-col>
                <mu-col span="4">
                    <div class="room-node _status_1">
                        <ul>
                            <li><img src="https://ss0.bdstatic.com/70cFvHSh_Q1YnxGkpoWK1HF6hhy/it/u=1669872726,1649722464&fm=26&gp=0.jpg" /></li>
                            <li><img src="https://ss1.bdstatic.com/70cFvXSh_Q1YnxGkpoWK1HF6hhy/it/u=3011024659,3918111814&fm=26&gp=0.jpg" /></li>
                        </ul>    
                        <div>
                        <span>已加入4人</span>
                        <span>未开始</span>
                        </div>
                    </div>    
                </mu-col>        
            </mu-row>

            <mu-row>
                <mu-col span="4">
                    <div class="room-node _status_1">
                        <ul>
                            <li><img src="https://ss0.bdstatic.com/70cFvHSh_Q1YnxGkpoWK1HF6hhy/it/u=2946010582,1303971542&fm=26&gp=0.jpg" /></li>
                            <li><img src="https://ss1.bdstatic.com/70cFvXSh_Q1YnxGkpoWK1HF6hhy/it/u=3011024659,3918111814&fm=26&gp=0.jpg" /></li>
                            <li><img src="https://ss0.bdstatic.com/70cFuHSh_Q1YnxGkpoWK1HF6hhy/it/u=1175288209,3211190587&fm=26&gp=0.jpg" /></li>
                        </ul>    
                        <div>
                        <span>已加入4人</span>
                        <span>未开始</span>
                        </div>
                    </div>    
                </mu-col>
                <mu-col span="4">
                    <div class="room-node _status_2">
                        <ul>
                            <li><img src="https://ss0.bdstatic.com/70cFvHSh_Q1YnxGkpoWK1HF6hhy/it/u=1669872726,1649722464&fm=26&gp=0.jpg" /></li>
                            <li><img src="https://ss0.bdstatic.com/70cFuHSh_Q1YnxGkpoWK1HF6hhy/it/u=1175288209,3211190587&fm=26&gp=0.jpg" /></li>
                        </ul>    
                        <div>
                        <span>已加入4人</span>
                        <span>已开始</span>
                        </div>
                    </div>    
                </mu-col>
                <mu-col span="4">
                    <div class="room-node _status_2">
                        <ul>
                            <li><img src="https://ss0.bdstatic.com/70cFvHSh_Q1YnxGkpoWK1HF6hhy/it/u=1669872726,1649722464&fm=26&gp=0.jpg" /></li>
                            <li><img src="https://ss0.bdstatic.com/70cFuHSh_Q1YnxGkpoWK1HF6hhy/it/u=1175288209,3211190587&fm=26&gp=0.jpg" /></li>
                        </ul>    
                        <div>
                        <span>已加入4人</span>
                        <span>未开始</span>
                        </div>
                    </div>    
                </mu-col>        
            </mu-row>

        </mu-container> 
        <mu-container class="bangdan">
            <mu-tabs :value.sync="active2" color="#e8e8e8" indicator-color="#43a047" full-width>
                <mu-tab><span>日榜</span></mu-tab>
                <mu-tab><span>周榜</span></mu-tab>
                <mu-tab><span>月榜</span></mu-tab>
                <mu-tab><span>总榜</span></mu-tab>
            </mu-tabs>

            <mu-list v-if="active2 === 0">
                <mu-list-item avatar button :ripple="false">
                    <mu-list-item-action>
                        <span>1</span>
                    </mu-list-item-action>    
                    <mu-list-item-action>
                        <mu-avatar>
                        <img src="https://ss0.bdstatic.com/70cFvHSh_Q1YnxGkpoWK1HF6hhy/it/u=1669872726,1649722464&fm=26&gp=0.jpg">
                        </mu-avatar>
                    </mu-list-item-action>
                    <mu-list-item-title>乐游 <span style="color:#999;font-size:12px;">(伴读书童)</span></mu-list-item-title>
                    <mu-list-item-action>
                        <span>320积分</span>
                    </mu-list-item-action>
                </mu-list-item>
                <mu-list-item avatar button :ripple="false">
                    <mu-list-item-action>
                        <span>2</span>
                    </mu-list-item-action> 
                    <mu-list-item-action>
                        <mu-avatar>
                        <img src="https://ss0.bdstatic.com/70cFvHSh_Q1YnxGkpoWK1HF6hhy/it/u=1669872726,1649722464&fm=26&gp=0.jpg">
                        </mu-avatar>
                    </mu-list-item-action>
                    <mu-list-item-title>二元 <span style="color:#999;font-size:12px;">(天之骄子)</span></mu-list-item-title>
                    <mu-list-item-action>
                        <span>300积分</span>
                    </mu-list-item-action>
                </mu-list-item>
            </mu-list> 
            <mu-list v-if="active2 === 1">
                <mu-list-item avatar button :ripple="false">
                    <mu-list-item-action>
                        <span>1</span>
                    </mu-list-item-action> 
                    <mu-list-item-action>
                        <mu-avatar>
                        <img src="https://ss0.bdstatic.com/70cFvHSh_Q1YnxGkpoWK1HF6hhy/it/u=1669872726,1649722464&fm=26&gp=0.jpg">
                        </mu-avatar>
                    </mu-list-item-action>
                    <mu-list-item-title>Mike Li <span style="color:#999;font-size:12px;">(伴读书童)</span></mu-list-item-title>
                    <mu-list-item-action>
                        <span>320积分</span>
                    </mu-list-item-action>
                </mu-list-item>
            </mu-list>    
            <mu-list v-if="active2 === 2">
                 <mu-list-item avatar button :ripple="false">
                    <mu-list-item-action>
                        <span>1</span>
                    </mu-list-item-action> 
                    <mu-list-item-action>
                        <mu-avatar>
                        <img src="https://ss0.bdstatic.com/70cFvHSh_Q1YnxGkpoWK1HF6hhy/it/u=1669872726,1649722464&fm=26&gp=0.jpg">
                        </mu-avatar>
                    </mu-list-item-action>
                    <mu-list-item-title>Mike Li <span style="color:#999;font-size:12px;">(伴读书童)</span></mu-list-item-title>
                    <mu-list-item-action>
                        <span>320积分</span>
                    </mu-list-item-action>
                </mu-list-item>
            </mu-list>    
            <mu-list v-if="active2 === 3">
                <mu-list-item avatar button :ripple="false">
                    <mu-list-item-action>
                        <span>1</span>
                    </mu-list-item-action> 
                    <mu-list-item-action>
                        <mu-avatar>
                        <img src="https://ss0.bdstatic.com/70cFvHSh_Q1YnxGkpoWK1HF6hhy/it/u=1669872726,1649722464&fm=26&gp=0.jpg">
                        </mu-avatar>
                    </mu-list-item-action>
                    <mu-list-item-title>Mike Li <span style="color:#999;font-size:12px;">(伴读书童)</span></mu-list-item-title>
                    <mu-list-item-action>
                        <span>320积分</span>
                    </mu-list-item-action>
                </mu-list-item>
            </mu-list>    

        </mu-container>

    </div>

</template>

<style scoped>

.userarea {
    height: 80px;
    width: 100%;
    background: linear-gradient(45deg,  #4caf50, #81c784);
}
.user-avtar {
    width: 60px;
    height: 60px;
    border-radius: 50%;
    overflow: hidden;
    margin: 10px 10px 10px 18px;
}
.user-avtar img{
    width: 60px;
    height: 60px;
}
.user-info {
    flex: 1;
}
.user-name {
    padding-top: 12px;
    font-size: 22px;
    color: #fff;
}

.user-level {
    color: #e8f5e9;
}

.joinroom {
    margin-top: 10px;
    text-align: center;
}
.joinroom .grid-cell {
    margin: 5px;
    background: linear-gradient(#ff9800,#ffa726);
    color: #fff;
    height: 100px;
    border-radius: 6px;
    font-size: 16px;
}
.joinroom ._two {
    background: linear-gradient(#f06292,#ec407a);
}
.joinroom .grid-cell span{
    font-size: 12px;
}

.joinroom .p{
    padding-top: 26px;
}
.roomlist {
    margin-top: 10px;
}
.room-node {
    margin: 5px;
    border:1px solid #e8e8e8;
    background: #e8e8e8;
    display: flex;
    flex-direction: column;
    border-radius: 6px;
}
.room-node ul {
    padding: 5px 3px;
    display: flex;
    flex-direction: row;
    margin-left: 6px;

}
.room-node ul li {
    width: 28px;
    height: 28px;
    padding: 1px;
    background: #fff;
    border-radius: 50%;
    overflow: hidden;
    margin-left:-6px;
}
.room-node ul li img{
    width: 100%;
    height: 100%;
    border-radius: 50%;
    display: block;
}

.room-node span{
    display:inline-block;
}
.room-node div {
    display: flex;
    flex-direction: row;
}
.room-node div span:nth-child(1) {
    color: #999;
    font-size: 12px;
    text-align: left;
    padding: 2px 3px;
}
.room-node div span:nth-child(2) {
    font-size: 10px;
    text-align: right;
    color: #43a047;
    flex: 1;
    padding: 2px 3px;
}
.room-node._status_2 {
    border:1px solid #81c784;
    background: #81c784;
}

.room-node._status_2 div span:nth-child(1) {
    color: #fff;
}
.room-node._status_2 div span:nth-child(2) {
    color: #ffe0b2;
}
.bangdan {
    margin-top: 15px;
}
.bangdan .mu-tabs span{
    color: #505050;
}
.bangdan .mu-tabs .mu-tab-active span{
    color: #43a047;
}


</style>

<script>
export default {
    "name" : "kpk",
    data() {
        return {
           active2: 0
        }
    },
    mounted() {

    },
    methods: {
        joinGame(e, type) {
            this.$router.push({path:'/kpkroom',  query: {type}})
        }
    }
}
</script>