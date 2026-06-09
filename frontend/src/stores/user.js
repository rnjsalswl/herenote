export function getUserID() {
  return localStorage.getItem('herenote_user_id')
}

export function setUserID(id) {
  localStorage.setItem('herenote_user_id', id)
}

export function clearUser() {
  localStorage.removeItem('herenote_user_id')
}

export function isLoggedIn() {
  return !!getUserID()
}