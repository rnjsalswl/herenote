<template>
  <div :style="{
    position:'relative', width:'100%', height:height+'px',
    borderRadius: rounded ? 'var(--radius)' : 0, overflow:'hidden',
    background:'var(--surface-2)',
  }">
    <!-- 실제 좌표가 있으면 카카오 정적 지도 -->
    <img
      v-if="lat && lng && !imgError"
      :src="staticUrl"
      :alt="label || '지도'"
      style="width:100%;height:100%;object-fit:cover;display:block"
      @error="imgError = true"
    />

    <!-- 좌표 없거나 이미지 실패 시 SVG 폴백 -->
    <svg v-else viewBox="0 0 100 100" preserveAspectRatio="xMidYMid slice" width="100%" height="100%">
      <g v-for="g in [20,40,60,80]" :key="g" stroke="var(--line)" stroke-width="0.5">
        <line :x1="g" y1="0" :x2="g" y2="100"/>
        <line x1="0" :y1="g" x2="100" :y2="g"/>
      </g>
      <template v-for="(rd,i) in roads" :key="i">
        <line v-if="rd.horiz" x1="0" :y1="rd.pos" x2="100" :y2="rd.pos" stroke="var(--ink-3)" stroke-opacity="0.5" :stroke-width="rd.w"/>
        <line v-else :x1="rd.pos" y1="0" :x2="rd.pos" y2="100" stroke="var(--ink-3)" stroke-opacity="0.5" :stroke-width="rd.w"/>
      </template>
      <rect :x="10+r(3)*50" :y="10+r(4)*50" width="26" height="22" rx="2" fill="var(--ink-3)" opacity="0.12"/>
      <circle v-if="showPin" cx="50" cy="52" :r="Math.min(34,10+radius/16)" fill="var(--ink)" opacity="0.06" stroke="var(--ink)" stroke-opacity="0.25" stroke-width="0.6" stroke-dasharray="2 2"/>
    </svg>

    <div v-if="showPin && !(lat && lng && !imgError)" style="position:absolute;top:52%;left:50%;transform:translate(-50%,-100%);color:var(--ink)">
      <HnIcon name="pinFill" :size="30"/>
    </div>
    <div v-if="label" style="position:absolute;left:10px;bottom:10px;font-family:var(--font-mono);font-size:10px;letter-spacing:0.04em;color:var(--ink-2);background:var(--surface);padding:3px 7px;border-radius:6px;border:1px solid var(--line)">{{ label }}</div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import HnIcon from './HnIcon.vue'

const props = defineProps({
  lat: { type: Number, default: null },
  lng: { type: Number, default: null },
  seed: { default: 1 },
  radius: { default: 100 },
  height: { default: 150 },
  showPin: { default: true },
  rounded: { default: true },
  label: String,
})

const imgError = ref(false)
const REST_KEY = import.meta.env.VITE_KAKAO_REST_KEY

// 반경에 따른 줌 레벨 계산 (카카오 level: 1=소, 14=대)
const zoomLevel = computed(() => {
  const r = props.radius
  if (r <= 50) return 2
  if (r <= 100) return 3
  if (r <= 200) return 4
  if (r <= 500) return 5
  return 6
})

const staticUrl = computed(() => {
  if (!props.lat || !props.lng) return ''
  const w = 400
  const h = Math.max(200, props.height * 2)
  return `https://dapi.kakao.com/v2/maps/staticmap?appkey=${REST_KEY}&center=${props.lng},${props.lat}&level=${zoomLevel.value}&size=${w}x${h}&markers=color:red%7C${props.lng},${props.lat}`
})

function r(n) { const x = Math.sin(props.seed*9301+n*49297)*233280; return x-Math.floor(x) }
const roads = Array.from({length:5},(_,i)=>({
  horiz: r(i)>0.5, pos: 12+r(i+10)*76, w: 0.6+r(i+20)*1.4,
}))
</script>
