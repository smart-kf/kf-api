const token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjYXJkSWQiOiJUTS1FYWZ2ZnA4dmVJIn0.x3AOg5FHKq5xzebl-PPF4gSPfMMoeeD0XPcduTdg_Wc';
const uid = 1;
const room_id = "chat://room1";
const platform = "kf-backend";

new Vue({
    el: "#app",
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
            let opt = this.getOptions();
            opt.userInfo.mid = parseInt(opt.userInfo.mid);
            this.im = new window["ImClient"](opt)
        },
        getOptions() {
            var that = this;
            return {
                wsUrl: "ws://127.0.0.1:3102/sub",
                token: this.token,
                maxRetryTime: 10,
                delay: 1500,
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
            if (this.status == "failed") {
                alert("请先创建websocket链接")
                return
            }
            this.im.sendMsg({
                "type": "msg",
                "data": {
                    "msgType": "text",
                    "guestId": 2,
                    "kfId": parseInt(this.userInfo.mid),
                    "content": this.content,
                    "isKf": 1, //   1=客服消息，2=客户消息
                }
            })
            this.content = '';
        },
        sendToMyself(msg) { // 服务器会向此客户端推送一条消息,
            if (this.status == "failed") {
                alert("请先创建websocket链接")
                return
            }

            var url = "http://localhost:8081/api/dev/push" ; // 本地的服务端域名, 线上是 https://api.smartkf.top
            let data = {
                "content": "hello world",
                "guestId": 2,
                "isKf": 1,
                "kfId": parseInt(this.userInfo.mid),
                "msgType": "msg"
            }

            fetch(url, {
                method: 'POST', // 或 'PUT' 取决于你的需求
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(data) // 将对象转换为 JSON 字符串
            })
                .then(response => {
                    if (!response.ok) {
                        throw new Error('Network response was not ok');
                    }
                    return response.json(); // 解析 JSON 响应
                })
                .then(data => {
                    console.log('Success:', data); // 处理成功的响应
                })
                .catch(error => {
                    console.error('Error:', error); // 处理错误
                });
            this.content = '';
        }
    }
})