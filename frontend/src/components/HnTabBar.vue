<template>
  <div :style="{
    flexShrink:0, zIndex:40,
    background:'color-mix(in srgb, var(--bg) 88%, transparent)',
    backdropFilter:'blur(18px) saturate(150%)',
    WebkitBackdropFilter:'blur(18px) saturate(150%)',
    borderTop:'1px solid var(--line)',
    display:'flex', alignItems:'center', gap:'4px', padding:'8px 14px 26px',
    position:'sticky', bottom:0,
  }">
    <button v-for="t in tabs" :key="t.id" @click="$emit('tab', t.id)" :style="{
      flex:1, border:'none', background:'transparent', cursor:'pointer',
      display:'flex', flexDirection:'column', alignItems:'center', gap:'3px',
      color: active===t.id ? 'var(--ink)' : 'var(--ink-3)', paddingTop:'2px',
    }">
      <HnIcon :name="t.icon" :size="25" :stroke="active===t.id ? 2.4 : 1.9" :fill="active===t.id && t.id==='home'"/>
      <span :style="{fontSize:'10.5px', fontWeight: active===t.id ? 700 : 500}">{{ t.label }}</span>
    </button>
    <button @click="$emit('compose')" aria-label="방명록 쓰기" :style="{
      width:'54px', height:'42px', borderRadius:'15px', border:'none', cursor:'pointer',
      background:'var(--ink)', color:'var(--inv-ink)', flexShrink:0,
      display:'flex', alignItems:'center', justifyContent:'center',
      boxShadow:'var(--shadow)',
    }">
      <HnIcon name="plus" :size="24" :stroke="2.6"/>
    </button>
    <button v-for="t in tabs2" :key="t.id" @click="$emit('tab', t.id)" :style="{
      flex:1, border:'none', background:'transparent', cursor:'pointer',
      display:'flex', flexDirection:'column', alignItems:'center', gap:'3px',
      color: active===t.id ? 'var(--ink)' : 'var(--ink-3)', paddingTop:'2px',
    }">
      <HnIcon :name="t.icon" :size="25" :stroke="active===t.id ? 2.4 : 1.9"/>
      <span :style="{fontSize:'10.5px', fontWeight: active===t.id ? 700 : 500}">{{ t.label }}</span>
    </button>
  </div>
</template>
<script setup>
import HnIcon from './HnIcon.vue'
defineProps({ active: String })
defineEmits(['tab', 'compose'])
const tabs = [{ id:'home', icon:'home', label:'홈' }, { id:'explore', icon:'map', label:'지도' }]
const tabs2 = [{ id:'profile', icon:'user', label:'프로필' }]
</script>
