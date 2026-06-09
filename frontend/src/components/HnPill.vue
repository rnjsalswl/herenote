<template>
  <span :style="{
    display:'inline-flex', alignItems:'center', gap:'5px',
    padding:'4px 10px', borderRadius:'999px', fontSize:'12px', lineHeight:1,
    fontFamily: mono ? 'var(--font-mono)' : 'var(--font-sans)',
    fontWeight:600, letterSpacing: mono ? '0.02em' : 0,
    background: solid ? 'var(--ink)' : 'var(--surface-2)',
    color: solid ? 'var(--inv-ink)' : 'var(--ink-2)',
    border: solid ? 'none' : '1px solid var(--line)',
  }"><slot /></span>
</template>
<script setup>
defineProps({ solid: Boolean, mono: Boolean })
</script>
