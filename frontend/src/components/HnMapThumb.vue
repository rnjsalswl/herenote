<template>
  <div :style="{
    position:'relative', width:'100%', height:height+'px',
    borderRadius: rounded ? 'var(--radius)' : 0, overflow:'hidden',
    background:'var(--surface-2)',
  }">
    <svg viewBox="0 0 100 100" preserveAspectRatio="xMidYMid slice" width="100%" height="100%">
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
    <div v-if="showPin" style="position:absolute;top:52%;left:50%;transform:translate(-50%,-100%);color:var(--ink)">
      <HnIcon name="pinFill" :size="30"/>
    </div>
    <div v-if="label" style="position:absolute;left:10px;bottom:10px;font-family:var(--font-mono);font-size:10px;letter-spacing:0.04em;color:var(--ink-2);background:var(--surface);padding:3px 7px;border-radius:6px;border:1px solid var(--line)">{{ label }}</div>
  </div>
</template>
<script setup>
import HnIcon from './HnIcon.vue'
const props = defineProps({ seed:{default:1}, radius:{default:100}, height:{default:150}, showPin:{default:true}, rounded:{default:true}, label:String })
function r(n) { const x = Math.sin(props.seed*9301+n*49297)*233280; return x-Math.floor(x) }
const roads = Array.from({length:5},(_,i)=>({
  horiz: r(i)>0.5, pos: 12+r(i+10)*76, w: 0.6+r(i+20)*1.4,
}))
</script>
