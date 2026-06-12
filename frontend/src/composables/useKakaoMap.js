// SDK는 한 번만 로드되도록 모듈 레벨에서 상태 관리
let sdkLoaded = false
let sdkLoading = false
const waitQueue = []

export function loadKakaoSDK() {
  return new Promise((resolve, reject) => {
    if (sdkLoaded) { resolve(); return }
    waitQueue.push({ resolve, reject })
    if (sdkLoading) return
    sdkLoading = true

    const script = document.createElement('script')
    script.src = `//dapi.kakao.com/v2/maps/sdk.js?appkey=${import.meta.env.KAKAO_JS_KEY}&autoload=false`
    script.onload = () => {
      window.kakao.maps.load(() => {
        sdkLoaded = true
        waitQueue.forEach(cb => cb.resolve())
        waitQueue.length = 0
      })
    }
    script.onerror = (err) => {
      sdkLoading = false
      waitQueue.forEach(cb => cb.reject(err))
      waitQueue.length = 0
    }
    document.head.appendChild(script)
  })
}
