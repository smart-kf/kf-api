(function (win) {
    const rawHeaderLen = 16;
    const packetOffset = 0;
    const headerOffset = 4;
    const verOffset = 6;
    const opOffset = 8;
    const seqOffset = 12;
    var defaultOptions = {
        wsUrl: "",      // websocket地址
        maxRetryTime: 3, // 断开重连最大重试次数
        delay: 1500,
        platform: "", // 平台
        token: "", // 用户token
        userInfo: { // 当前用户的连接信息: {"mid":123, "room_id":"live://1000", "platform":"web", "accepts":[1000,1001,1002]}
            mid: 1,
            room_id: "chat://room1",
            platform: "kf-backend",
            accepts: [], // 接收谁的信息.
        },
        debug: false,
        onMessage: function (obj) {
            return;
        },
        onClose: function () {

        },
        onError: function (error) {

        },
        onAuthSuccess: function () {

        },
    };
    var textDecoder = new TextDecoder();
    var textEncoder = new TextEncoder();
    var ws = null;


    function mergeArrayBuffer(ab1, ab2) {
        var u81 = new Uint8Array(ab1),
            u82 = new Uint8Array(ab2),
            res = new Uint8Array(ab1.byteLength + ab2.byteLength);
        res.set(u81, 0);
        res.set(u82, ab1.byteLength);
        return res.buffer;
    }

    function char2ab(str) {
        var buf = new ArrayBuffer(str.length);
        var bufView = new Uint8Array(buf);
        for (var i = 0; i < str.length; i++) {
            bufView[i] = str[i];
        }
        return buf;
    }

    var ImClient = function (options) {
        this.options = Object.assign({},options);
        var maxRetry = 3
        var delay = 1500
        if (options.maxRetryTime) {
            maxRetry = options.maxRetryTime
        }
        if (options.delay) {
            delay = options.delay;
        }
        console.log(this);
        this.createConnect(maxRetry, delay);
    }
    ImClient.prototype.createConnect = function (max, delay) {
        var self = this;
        if (max === 0) {
            return;
        }
        connect();
        var heartbeatInterval;

        function connect() {
            let options = self.options;
            let wsUrl = options.wsUrl + "?token=" + options.token + "&platform=" + options.platform;
            ws = new WebSocket(wsUrl);
            ws.binaryType = 'arraybuffer';
            ws.onopen = function () {
                auth();
            }
            ws.onerror = function (err) {
                if(options.debug) {
                    console.log("onError->", err)
                }
                options.onError && options.onError(err);
            }
            ws.onmessage = function (evt) {
                var data = evt.data;
                var dataView = new DataView(data, 0);
                var packetLen = dataView.getInt32(packetOffset);
                var headerLen = dataView.getInt16(headerOffset);
                var ver = dataView.getInt16(verOffset);
                var op = dataView.getInt32(opOffset);
                var seq = dataView.getInt32(seqOffset);
                console.log(op)
                switch (op) {
                    case 8:
                        heartbeat();
                        heartbeatInterval = setInterval(heartbeat, 30 * 1000);
                        if(options.debug) {
                            console.log("getAuthSuccessMsg, call options.onAuthSuccess()")
                        }
                        options.onAuthSuccess && options.onAuthSuccess();
                        break;
                    case 3:
                        if(options.debug) {
                            console.log("receive server heartBeat")
                        }
                        break;
                    case 9:
                        // batch message
                        for (var offset = rawHeaderLen; offset < data.byteLength; offset += packetLen) {
                            // parse
                            var packetLen = dataView.getInt32(offset);
                            var headerLen = dataView.getInt16(offset + headerOffset);
                            var ver = dataView.getInt16(offset + verOffset);
                            var op = dataView.getInt32(offset + opOffset);
                            var seq = dataView.getInt32(offset + seqOffset);
                            var msgBody = textDecoder.decode(data.slice(offset + headerLen, offset + packetLen));
                            // callback
                            if(options.debug) {
                                console.log("receive-Multiple->Message--->", msgBody)
                            }
                            if(options.onMessage) {
                                options.onMessage(msgBody)
                            }
                        }
                        break;
                    case 100: // 收到新消息.
                        var msgBody = textDecoder.decode(data.slice(headerLen, packetLen));
                        if(options.debug) {
                            console.log("receive-Single->Message100--->", msgBody)
                        }
                        if(options.onMessage) {
                            options.onMessage(msgBody)
                        }
                        break
                    default:
                        var msgBody = textDecoder.decode(data.slice(headerLen, packetLen));
                        if(options.debug) {
                            console.log("receive-Single->Message--->", msgBody)
                        }
                        if(options.onMessage) {
                            options.onMessage(msgBody)
                        }
                        break
                }
            }

            ws.onclose = function (e) {
                if(options.debug) {
                    console.log("onCloseEvent--->", e)
                }
                if (heartbeatInterval) clearInterval(heartbeatInterval);
                setTimeout(reConnect, delay);
            }

            function heartbeat() {
                var headerBuf = new ArrayBuffer(rawHeaderLen);
                var headerView = new DataView(headerBuf, 0);
                headerView.setInt32(packetOffset, rawHeaderLen);
                headerView.setInt16(headerOffset, rawHeaderLen);
                headerView.setInt16(verOffset, 1);
                headerView.setInt32(opOffset, 2);
                headerView.setInt32(seqOffset, 1);
                ws.send(headerBuf);
                if(options.debug) {
                    console.log("send: heartbeat");
                }
            }

            function auth() {
                var token = JSON.stringify(options.userInfo);
                var headerBuf = new ArrayBuffer(rawHeaderLen);
                var headerView = new DataView(headerBuf, 0);
                var bodyBuf = textEncoder.encode(token);
                headerView.setInt32(packetOffset, rawHeaderLen + bodyBuf.byteLength);
                headerView.setInt16(headerOffset, rawHeaderLen);
                headerView.setInt16(verOffset, 1);
                headerView.setInt32(opOffset, 7);
                headerView.setInt32(seqOffset, 1);
                ws.send(mergeArrayBuffer(headerBuf, bodyBuf));
                if(options.debug) {
                    console.log("send: auth info: ", options.userInfo);
                }
            }
        }

        function reConnect() {
           //  self.createConnect(--max, delay * 2);
            if(self.options.debug) {
             //   console.log("time: " + max, ", delay:" + delay *2)
            }
        }
    }

    ImClient.prototype.sendMsg = function(object) {
        var self = this;
        if(!ws) {
            return
        }
        let options = self.options;
        var msg = JSON.stringify(object);
        var headerBuf = new ArrayBuffer(rawHeaderLen);
        var headerView = new DataView(headerBuf, 0);
        var bodyBuf = textEncoder.encode(msg);
        headerView.setInt32(packetOffset, rawHeaderLen + bodyBuf.byteLength);
        headerView.setInt16(headerOffset, rawHeaderLen);
        headerView.setInt16(verOffset, 1);
        headerView.setInt32(opOffset, 100);
        headerView.setInt32(seqOffset, 1);
        ws.send(mergeArrayBuffer(headerBuf, bodyBuf));
        if(options.debug) {
            console.log("sendMsg: ", object);
        }
    }
    window["ImClient"] = ImClient
})(window);

(function(win){
    var someFunction = function() {
        console.log(this.another())
    }
    someFunction.prototype.another = function() {
        return 1;
    }
    someFunction.prototype.another2 = function() {
        return 2;
    }
    window["someFunction"] = someFunction;
})(window)

let x = new window["someFunction"]()
console.log("-->",x.another2())