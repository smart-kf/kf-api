const token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjYXJkSWQiOiJUTS1FYWZ2ZnA4dmVJIn0.x3AOg5FHKq5xzebl-PPF4gSPfMMoeeD0XPcduTdg_Wc';
const uid = 2;
const room_id = "chat://room1";
const platform = "kf-backend";

new Vue({
    el:"#app",
    data() {
        return {
            token: token,
            status: "failed",
            im: null,
            userInfo: {
                mid: uid,
                room_id: room_id,
                platform: platform,
                accepts: [],
            },
            receiveMessages: [],
            error: "",
            content: '',
        }
    },
    methods: {
        connect() {
            this.im = new window["ImClient"](this.getOptions())
        },
        getOptions() {
            var that = this;
            return {
                wsUrl: "ws://127.0.0.1:3102/sub",
                token: this.token,
                maxRetryTime: 10,
                delay:1500 ,
                platform: "kf-backend",
                userInfo: this.userInfo,
                debug: true,
                onMessage(msg) {
                    that.receiveMessages.push(msg)
                },
                onClose() {
                    that.status = "failed"
                },
                onError(error) {
                    that.error = error
                },
                onAuthSuccess() {
                    that.status = "authSuccess"
                },
            }
        },
        sendMsg(msg) { // 服务器会向此客户端推送一条消息,
            console.log(Object.getOwnPropertyNames(this.im));
            this.im.sendMsg({
                "type": "msg",
                "data": {
                    "msgType": "text",
                    "guestName": "",
                    "guestAvatar": "",
                    "guestId": 2,
                    "msgTime": 0, //时间戳
                    "kfId": 1,
                    "content": this.content,
                    "city": "",
                    "ip": "",
                    "isKf": 1, //   1=客服消息，2=客户消息
                }
            })
            this.content = '';
        }
    }
})