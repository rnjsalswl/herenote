<template>
  <div
    ref="el"
    :style="{
      width: '100%',
      height: height + 'px',
      borderRadius: rounded ? 'var(--radius)' : 0,
      overflow: 'hidden',
      background: 'var(--surface-2)',
    }"
  >
    <div v-if="loadError" style="display:flex;align-items:center;justify-content:center;width:100%;height:100%;font-family:var(--font-mono);font-size:11px;color:var(--ink-3)">
      지도를 불러올 수 없어요
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, watch } from 'vue'
import { loadKakaoSDK } from '@/composables/useKakaoMap.js'

const props = defineProps({
  lat: { type: Number, default: null },
  lng: { type: Number, default: null },
  places: { type: Array, default: () => [] },
  zoom: { type: Number, default: 4 },
  height: { type: Number, default: 260 },
  radius: { type: Number, default: 0 },
  rounded: { type: Boolean, default: false },
})

const emit = defineEmits(['marker-click'])

const el = ref(null)
const loadError = ref(false)
let map = null
let markers = []
let circle = null

async function init() {
  try {
    await loadKakaoSDK()
  } catch {
    loadError.value = true
    return
  }

  const { kakao } = window
  const centerLat = props.lat ?? props.places[0]?.latitude ?? 37.5665
  const centerLng = props.lng ?? props.places[0]?.longitude ?? 126.9780

  map = new kakao.maps.Map(el.value, {
    center: new kakao.maps.LatLng(centerLat, centerLng),
    level: props.zoom,
  })

  drawMarkers()

  if (props.radius > 0 && props.lat != null && props.lng != null) {
    drawCircle()
  }
}

function drawMarkers() {
  if (!map) return
  const { kakao } = window

  markers.forEach(m => m.setMap(null))
  markers = []

  props.places.forEach(p => {
    const marker = new kakao.maps.Marker({
      position: new kakao.maps.LatLng(p.latitude, p.longitude),
      map,
    })
    kakao.maps.event.addListener(marker, 'click', () => emit('marker-click', p))
    markers.push(marker)
  })

  if (props.places.length > 1) {
    const bounds = new kakao.maps.LatLngBounds()
    props.places.forEach(p => bounds.extend(new kakao.maps.LatLng(p.latitude, p.longitude)))
    map.setBounds(bounds)
  } else if (props.places.length === 1 && props.lat == null) {
    map.setCenter(new kakao.maps.LatLng(props.places[0].latitude, props.places[0].longitude))
  }
}

function drawCircle() {
  if (!map) return
  const { kakao } = window
  if (circle) circle.setMap(null)
  circle = new kakao.maps.Circle({
    center: new kakao.maps.LatLng(props.lat, props.lng),
    radius: props.radius,
    strokeWeight: 2,
    strokeColor: '#111111',
    strokeOpacity: 0.65,
    strokeStyle: 'solid',
    fillColor: '#111111',
    fillOpacity: 0.07,
  })
  circle.setMap(map)
}

onMounted(init)

onBeforeUnmount(() => {
  markers.forEach(m => m.setMap(null))
  if (circle) circle.setMap(null)
})

watch(() => props.places, drawMarkers, { deep: true })

watch(() => [props.lat, props.lng, props.radius], ([lat, lng, r]) => {
  if (!map) return
  if (lat != null && lng != null) {
    map.setCenter(new window.kakao.maps.LatLng(lat, lng))
    if (r > 0) drawCircle()
  }
})
</script>
