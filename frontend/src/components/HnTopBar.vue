<template>
  <div :style="{
    position:'sticky',top:0,zIndex:30,
    paddingTop:'54px',
    background: transparent ? 'transparent' : 'color-mix(in srgb, var(--bg) 82%, transparent)',
    backdropFilter: transparent ? 'none' : 'blur(18px) saturate(150%)',
    WebkitBackdropFilter: transparent ? 'none' : 'blur(18px) saturate(150%)',
    borderBottom: transparent ? '1px solid transparent' : '1px solid var(--line)',
  }">
    <div style="height:50px;display:flex;align-items:center;gap:6px;padding:0 10px">
      <button v-if="onBack" @click="onBack" :style="iconBtn">
        <HnIcon name="chevL" :size="24" :stroke="2.4"/>
      </button>
      <div :style="{
        flex:1, minWidth:0, display:'flex', alignItems:'center',
        justifyContent: onBack ? 'center' : 'flex-start',
        paddingLeft: onBack ? 0 : '8px',
      }">
        <span v-if="big" :style="{
          fontFamily:'var(--font-display)', fontWeight:'var(--disp-weight)',
          fontSize:'21px', letterSpacing: mono ? '0.14em' : 'var(--disp-ls)',
          fontVariant: mono ? 'all-small-caps' : 'normal',
          color:'var(--ink)', whiteSpace:'nowrap', overflow:'hidden', textOverflow:'ellipsis',
        }">{{ title }}</span>
        <span v-else :style="{
          fontFamily:'var(--font-sans)', fontWeight:700, fontSize:'16px',
          color:'var(--ink)', whiteSpace:'nowrap', overflow:'hidden', textOverflow:'ellipsis',
        }">{{ title }}</span>
      </div>
      <div style="display:flex;align-items:center;gap:2px">
        <slot name="right"/>
      </div>
    </div>
  </div>
</template>
<script setup>
import HnIcon from './HnIcon.vue'
defineProps({ title:String, big:Boolean, mono:Boolean, transparent:Boolean, onBack:Function })
const iconBtn = { width:'40px', height:'40px', borderRadius:'50%', border:'none', background:'transparent', color:'var(--ink)', display:'flex', alignItems:'center', justifyContent:'center', cursor:'pointer', flexShrink:0 }
</script>
