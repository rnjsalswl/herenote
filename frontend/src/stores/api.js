import { getToken } from './user.js'

// 인증 헤더가 포함된 fetch 래퍼
export function authFetch(url, options = {}) {
  const token = getToken()
  return fetch(url, {
    ...options,
    headers: {
      ...options.headers,
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
    },
  })
}
