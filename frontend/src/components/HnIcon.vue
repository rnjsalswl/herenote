<template>
  <svg :width="size" :height="size" viewBox="0 0 24 24" aria-hidden="true"
    style="display:inline-block;vertical-align:middle;flex-shrink:0">
    <!-- home -->
    <template v-if="name==='home'">
      <path d="M3 10.5 12 3l9 7.5" :stroke-width="stroke" v-bind="p"/>
      <path d="M5 9.5V20a1 1 0 0 0 1 1h12a1 1 0 0 0 1-1V9.5" :stroke-width="stroke" v-bind="p"/>
    </template>
    <!-- map -->
    <template v-else-if="name==='map'">
      <path d="M9 4 3 6.5v13L9 17l6 2.5 6-2.5v-13L15 6.5 9 4Z" :stroke-width="stroke" v-bind="p"/>
      <path d="M9 4v13M15 6.5v13" :stroke-width="stroke" v-bind="p"/>
    </template>
    <!-- user -->
    <template v-else-if="name==='user'">
      <circle cx="12" cy="8" r="4" :stroke-width="stroke" v-bind="p"/>
      <path d="M4 21c0-4 3.5-6 8-6s8 2 8 6" :stroke-width="stroke" v-bind="p"/>
    </template>
    <!-- plus -->
    <template v-else-if="name==='plus'">
      <path d="M12 5v14M5 12h14" :stroke-width="stroke" v-bind="p"/>
    </template>
    <!-- heart -->
    <template v-else-if="name==='heart'">
      <path d="M12 20s-7-4.5-9.5-9C1 8 2.5 4.5 6 4.5c2 0 3.2 1.2 4 2.3.8-1.1 2-2.3 4-2.3 3.5 0 5 3.5 3.5 6.5C19 15.5 12 20 12 20Z"
        :fill="fill ? 'currentColor' : 'none'" stroke="currentColor" :stroke-width="stroke" stroke-linejoin="round"/>
    </template>
    <!-- comment -->
    <template v-else-if="name==='comment'">
      <path d="M21 11.5a8.5 8.5 0 0 1-12.5 7.5L3 21l1.9-4.5A8.5 8.5 0 1 1 21 11.5Z" :stroke-width="stroke" v-bind="p"/>
    </template>
    <!-- pin -->
    <template v-else-if="name==='pin'">
      <path d="M12 21s7-6 7-11a7 7 0 0 0-14 0c0 5 7 11 7 11Z" :stroke-width="stroke" v-bind="p"/>
      <circle cx="12" cy="10" r="2.5" :stroke-width="stroke" v-bind="p"/>
    </template>
    <!-- pinFill -->
    <template v-else-if="name==='pinFill'">
      <path d="M12 21s7-6 7-11a7 7 0 0 0-14 0c0 5 7 11 7 11Z" fill="currentColor" stroke="none"/>
      <circle cx="12" cy="10" r="2.4" fill="var(--surface)" stroke="none"/>
    </template>
    <!-- lock -->
    <template v-else-if="name==='lock'">
      <rect x="5" y="11" width="14" height="9" rx="2" :stroke-width="stroke" v-bind="p"/>
      <path d="M8 11V8a4 4 0 0 1 8 0v3" :stroke-width="stroke" v-bind="p"/>
    </template>
    <!-- check -->
    <template v-else-if="name==='check'">
      <path d="M5 12.5 10 17 19 7" :stroke-width="stroke" v-bind="p"/>
    </template>
    <!-- chevL -->
    <template v-else-if="name==='chevL'">
      <path d="M15 5 8 12l7 7" :stroke-width="stroke" v-bind="p"/>
    </template>
    <!-- chevR -->
    <template v-else-if="name==='chevR'">
      <path d="M9 5l7 7-7 7" :stroke-width="stroke" v-bind="p"/>
    </template>
    <!-- more -->
    <template v-else-if="name==='more'">
      <circle cx="5" cy="12" r="1.6" fill="currentColor" stroke="none"/>
      <circle cx="12" cy="12" r="1.6" fill="currentColor" stroke="none"/>
      <circle cx="19" cy="12" r="1.6" fill="currentColor" stroke="none"/>
    </template>
    <!-- camera -->
    <template v-else-if="name==='camera'">
      <path d="M3 8a2 2 0 0 1 2-2h2l1.5-2h7L17 6h2a2 2 0 0 1 2 2v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8Z" :stroke-width="stroke" v-bind="p"/>
      <circle cx="12" cy="12.5" r="3.5" :stroke-width="stroke" v-bind="p"/>
    </template>
    <!-- image -->
    <template v-else-if="name==='image'">
      <rect x="3" y="4" width="18" height="16" rx="2" :stroke-width="stroke" v-bind="p"/>
      <circle cx="8.5" cy="9.5" r="1.8" :stroke-width="stroke" v-bind="p"/>
      <path d="m4 18 5-5 4 3 3-3 4 4" :stroke-width="stroke" v-bind="p"/>
    </template>
    <!-- settings -->
    <template v-else-if="name==='settings'">
      <circle cx="12" cy="12" r="3" :stroke-width="stroke" v-bind="p"/>
      <path d="M19 12a7 7 0 0 0-.1-1.3l2-1.5-2-3.4-2.3 1a7 7 0 0 0-2.2-1.3L14 2.5h-4l-.4 2.7a7 7 0 0 0-2.2 1.3l-2.3-1-2 3.4 2 1.5A7 7 0 0 0 7 12c0 .4 0 .9.1 1.3l-2 1.5 2 3.4 2.3-1a7 7 0 0 0 2.2 1.3l.4 2.7h4l.4-2.7a7 7 0 0 0 2.2-1.3l2.3 1 2-3.4-2-1.5c.1-.4.1-.9.1-1.3Z" :stroke-width="stroke" v-bind="p"/>
    </template>
    <!-- search -->
    <template v-else-if="name==='search'">
      <circle cx="11" cy="11" r="7" :stroke-width="stroke" v-bind="p"/>
      <path d="m20 20-3.5-3.5" :stroke-width="stroke" v-bind="p"/>
    </template>
    <!-- clock -->
    <template v-else-if="name==='clock'">
      <circle cx="12" cy="12" r="8.5" :stroke-width="stroke" v-bind="p"/>
      <path d="M12 7v5.5l3.5 2" :stroke-width="stroke" v-bind="p"/>
    </template>
    <!-- bolt -->
    <template v-else-if="name==='bolt'">
      <path d="M13 2 4 14h6l-1 8 9-12h-6l1-8Z" :stroke-width="stroke" v-bind="p"/>
    </template>
    <!-- bookmark -->
    <template v-else-if="name==='bookmark'">
      <path d="M6 4h12v17l-6-4-6 4V4Z" :stroke-width="stroke" v-bind="p"/>
    </template>
    <!-- walk -->
    <template v-else-if="name==='walk'">
      <circle cx="13" cy="4.5" r="1.8" :stroke-width="stroke" v-bind="p"/>
      <path d="M11 22l1.5-6-2.5-2 1-5 3 2 2 1M9 22l2-6" :stroke-width="stroke" v-bind="p"/>
    </template>
    <!-- x -->
    <template v-else-if="name==='x'">
      <path d="M6 6l12 12M18 6 6 18" :stroke-width="stroke" v-bind="p"/>
    </template>
    <!-- grid -->
    <template v-else-if="name==='grid'">
      <rect x="4" y="4" width="7" height="7" rx="1.5" :stroke-width="stroke" v-bind="p"/>
      <rect x="13" y="4" width="7" height="7" rx="1.5" :stroke-width="stroke" v-bind="p"/>
      <rect x="4" y="13" width="7" height="7" rx="1.5" :stroke-width="stroke" v-bind="p"/>
      <rect x="13" y="13" width="7" height="7" rx="1.5" :stroke-width="stroke" v-bind="p"/>
    </template>
    <!-- list -->
    <template v-else-if="name==='list'">
      <path d="M8 6h12M8 12h12M8 18h12" :stroke-width="stroke" v-bind="p"/>
      <circle cx="4" cy="6" r="1" fill="currentColor" stroke="none"/>
      <circle cx="4" cy="12" r="1" fill="currentColor" stroke="none"/>
      <circle cx="4" cy="18" r="1" fill="currentColor" stroke="none"/>
    </template>
    <!-- send -->
    <template v-else-if="name==='send'">
      <path d="M4 12 20 4l-7 16-2.5-6.5L4 12Z" :stroke-width="stroke" v-bind="p"/>
    </template>
  </svg>
</template>
<script setup>
defineProps({ name: String, size: { default: 24 }, stroke: { default: 2 }, fill: Boolean })
const p = { fill: 'none', stroke: 'currentColor', 'stroke-linecap': 'round', 'stroke-linejoin': 'round' }
</script>
