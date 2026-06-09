import { reactive } from 'vue'
const state = reactive({ map: {}, lastId: null })
export function setVerified(placeId, token) { state.map[placeId] = token; state.lastId = placeId }
export function getToken(placeId) { return state.map[placeId] || '' }
export function isVerified(placeId) { return !!state.map[placeId] }
export function lastVerifiedId() { return state.lastId }
