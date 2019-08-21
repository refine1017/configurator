import Cookies from 'js-cookie'

export function cookieGet(key) {
  return Cookies.get(key)
}

export function cookieSet(key, value) {
  return Cookies.set(key, value)
}

export function cookieRemove(key) {
  return Cookies.remove(key)
}
