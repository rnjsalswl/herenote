<template>
  <div style="display:flex;align-items:flex-end;justify-content:space-between">
    <div>
      <div v-if="kicker" style="font-family:var(--font-mono);font-size:10.5px;letter-spacing:0.18em;color:var(--ink-3);margin-bottom:3px">{{ kicker }}</div>
      <div style="font-family:var(--font-display);font-weight:var(--disp-weight);font-size:20px;letter-spacing:var(--disp-ls);color:var(--ink)">{{ title }}</div>
    </div>
    <span v-if="right" style="font-size:12px;color:var(--ink-3);font-family:var(--font-mono)">{{ right }}</span>
  </div>
</template>
<script setup>
defineProps({ kicker: String, title: String, right: String })
</script>
