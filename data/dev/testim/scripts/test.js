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
