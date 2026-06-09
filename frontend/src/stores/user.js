const TOKEN_KEY = 'herenote_token'

function decodeJWT(token) {
  try {
    const b64 = token.split('.')[1].replace(/-/g, '+').replace(/_/g, '/')
    return JSON.parse(atob(b64))
  } catch { return null }
}

export function setToken(token) {
  localStorage.setItem(TOKEN_KEY, token)
}

export function getToken() {
  return localStorage.getItem(TOKEN_KEY) || ''
}

export function getUserID() {
  const payload = decodeJWT(getToken())
  return payload?.user_id || null
}

export function getNickname() {
  const payload = decodeJWT(getToken())
  return payload?.nickname || null
}

export function isLoggedIn() {
  const token = getToken()
  if (!token) return false
  const payload = decodeJWT(token)
  if (!payload) return false
  if (payload.exp && payload.exp * 1000 < Date.now()) {
    clearUser()
    return false
  }
  return true
}

export function clearUser() {
  localStorage.removeItem(TOKEN_KEY)
}

// 하위 호환 — LoginView가 이전에 setUserID를 사용했으므로 유지
export function setUserID(id) {
  localStorage.setItem('herenote_user_id', id)
}
