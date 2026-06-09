<template>
  <article :style="{background:'var(--card)',borderRadius:'var(--radius-lg)',border:'1px solid var(--line)',overflow:'hidden',boxShadow:'var(--shadow)'}">
    <!-- header -->
    <div style="display:flex;align-items:center;gap:10px;padding:12px 14px 10px">
      <HnAvatar :name="note.nickname" :size="40"/>
      <div style="flex:1;min-width:0">
        <div style="display:flex;align-items:center;gap:5px">
          <span style="font-weight:700;font-size:14.5px;color:var(--ink)">{{ note.nickname }}</span>
          <HnBadge :type="note.badge"/>
          <HnPill v-if="note.pinned" style="padding:2px 7px;font-size:10.5px">고정</HnPill>
        </div>
        <button v-if="place" @click="$emit('open')" style="border:none;background:transparent;padding:0;cursor:pointer;display:flex;align-items:center;gap:3px;color:var(--ink-2);font-size:12.5px;font-weight:500;margin-top:1px">
          <HnIcon name="pin" :size="13" :stroke="2.2"/>
          <span style="font-weight:600;color:var(--ink-2)">{{ place }}</span>
        </button>
      </div>
      <span style="font-family:var(--font-mono);font-size:11px;color:var(--ink-3);flex-shrink:0">{{ note.time }}</span>
    </div>
    <!-- photo placeholder -->
    <div v-if="note.photo" style="padding:0 14px 12px">
      <div :style="{height:'196px',borderRadius:'var(--radius)',overflow:'hidden',background:`repeating-linear-gradient(135deg,var(--surface-2),var(--surface-2) 11px,transparent 11px,transparent 22px),var(--surface)`,border:'1px solid var(--line)',display:'flex',alignItems:'flex-end'}">
        <span style="font-family:var(--font-mono);font-size:10.5px;color:var(--ink-3);background:var(--surface);padding:3px 8px;margin:10px;border-radius:6px;border:1px solid var(--line)">사진</span>
      </div>
    </div>
    <!-- content -->
    <p style="padding:0 16px 4px;margin:0;font-size:15px;line-height:1.62;color:var(--ink);white-space:pre-wrap;word-break:keep-all">{{ note.content }}</p>
    <!-- actions -->
    <div style="display:flex;align-items:center;gap:4px;padding:10px 12px 12px">
      <button @click="toggleLike" :style="{display:'inline-flex',alignItems:'center',gap:'6px',border:'none',background:'transparent',cursor:'pointer',padding:'6px 8px',borderRadius:'999px',color:liked?'var(--ink)':'var(--ink-2)'}">
        <span :style="{display:'inline-flex',transform:liked?'scale(1.18)':'scale(1)',transition:'transform .18s cubic-bezier(.2,.9,.3,1.4)'}">
          <HnIcon name="heart" :size="20" :stroke="2" :fill="liked"/>
        </span>
        <span style="font-size:13px;font-weight:700;font-variant-numeric:tabular-nums">{{ count }}</span>
      </button>
      <button style="display:inline-flex;align-items:center;gap:6px;border:none;background:transparent;cursor:pointer;padding:6px 8px;border-radius:999px;color:var(--ink-2)">
        <HnIcon name="comment" :size="19" :stroke="2"/>
        <span style="font-size:13px;font-weight:600">공감</span>
      </button>
      <div style="flex:1"/>
      <button style="width:34px;height:34px;border-radius:50%;border:none;background:transparent;color:var(--ink-2);display:flex;align-items:center;justify-content:center">
        <HnIcon name="bookmark" :size="18" :stroke="2"/>
      </button>
    </div>
  </article>
</template>
<script setup>
import { ref } from 'vue'
import HnIcon from './HnIcon.vue'
import HnAvatar from './HnAvatar.vue'
import HnBadge from './HnBadge.vue'
import HnPill from './HnPill.vue'
const props = defineProps({ note: Object, place: String })
defineEmits(['open'])
const liked = ref(!!props.note.liked)
const count = ref(props.note.likes || 0)
function toggleLike() { liked.value = !liked.value; count.value += liked.value ? 1 : -1 }
</script>
