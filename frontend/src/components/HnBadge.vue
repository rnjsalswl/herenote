<template>
  <span v-if="type && map[type]" :title="map[type].label" :style="{
    display:'inline-flex', alignItems:'center', justifyContent:'center',
    width: (size+4)+'px', height:(size+4)+'px', borderRadius:'50%',
    background:'var(--ink)', color:'var(--inv-ink)', flexShrink:0,
  }">
    <HnIcon :name="glyph" :size="size-3" :stroke="2.6" />
  </span>
</template>
<script setup>
import HnIcon from './HnIcon.vue'
const props = defineProps({ type: String, size: { default: 14 } })
const map = { GOLD: { label:'단골' }, BLUE: { label:'인증' }, SPECIAL: { label:'첫 발자국' } }
const glyph = props.type === 'GOLD' ? 'bolt' : props.type === 'SPECIAL' ? 'pin' : 'check'
</script>
