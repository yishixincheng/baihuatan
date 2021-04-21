const TokenKey = 'Admin-Token'
const RefreshTokenKey = 'Admin-Refresh-Token'

export function getToken() {
  return localStorage.getItem(TokenKey)
}

export function getRefreshToken() {
  return localStorage.getItem(RefreshTokenKey)
}

export function setToken(token) {
  return localStorage.setItem(TokenKey, token)
}

export function setRefreshToken(refreshToken) {
  return localStorage.setItem(RefreshTokenKey, refreshToken)
}

export function removeToken() {
  localStorage.removeItem(TokenKey)
  localStorage.removeItem(RefreshTokenKey)
}
