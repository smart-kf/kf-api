// 顶部提示框
function showToast(message, type = 'success', duration = 3000,) {
  const toast = document.getElementById('top-toast');
  toast.innerText = message; // 设置提示信息
  toast.style.display = 'block'; // 显示提示条
  console.log(toast.className);
  toast.className = 'top-toast ' + type
  setTimeout(() => { // 在指定时间后隐藏提示条
    toast.className = 'top-toast '
    toast.style.display = 'none';
  }, duration);
}
// 加载遮罩层
function showMask(flag) {
  const toast = document.getElementById('loading-overlay');
  toast.style.display = flag ? 'flex' : 'none';
}

// 请求响应状态校验
function statusCheck(res) {
  if (res && res.code !== 200) {
    showToast(res.message || '请求失败')
    return false
  }
  return true
}