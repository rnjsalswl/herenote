<template>
  <div :style="{
    width:size+'px', height:size+'px', borderRadius:'50%', flexShrink:0,
    display:'flex', alignItems:'center', justifyContent:'center',
    background: filled ? 'var(--ink)' : 'var(--surface-2)',
    color: filled ? 'var(--inv-ink)' : 'var(--ink)',
    fontFamily:'var(--font-display)', fontWeight:700,
    fontSize: Math.round(size*0.42)+'px',
    border: filled ? 'none' : '1px solid var(--line)',
  }">{{ (name||'?').charAt(0) }}</div>
</template>
<script setup>
defineProps({ name: { default:'?' }, size: { default:38 }, filled: { default:true } })
</script>
