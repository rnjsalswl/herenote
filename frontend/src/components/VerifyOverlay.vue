<template>
  <div :style="{
    position:'absolute', inset:0, zIndex:80,
    background:'color-mix(in srgb, var(--bg) 94%, transparent)',
    backdropFilter:'blur(8px)', WebkitBackdropFilter:'blur(8px)',
    display:'flex', flexDirection:'column', alignItems:'center', justifyContent:'center',
    padding:'40px', textAlign:'center',
  }">
    <div style="position:relative;width:180px;height:180px;display:flex;align-items:center;justify-content:center">
      <span v-if="phase==='scan'" v-for="i in [0,1,2]" :key="i" :style="{
        position:'absolute', width:'180px', height:'180px', borderRadius:'50%',
        border:'1.5px solid var(--ink)', opacity:0,
        animation:`hnPulse 2.1s ${i*0.55}s infinite ease-out`,
      }"/>
      <div :style="{
        width:'92px', height:'92px', borderRadius:'50%',
        background: phase==='ok' ? 'var(--ink)' : 'var(--surface)',
        color: phase==='ok' ? 'var(--inv-ink)' : 'var(--ink)',
        border:'1px solid var(--line)', boxShadow:'var(--shadow)',
        display:'flex', alignItems:'center', justifyContent:'center',
        transition:'all .4s cubic-bezier(.2,.8,.2,1)',
        transform: phase==='ok' ? 'scale(1.05)' : 'scale(1)',
      }">
        <HnIcon :name="phase==='ok' ? 'check' : 'pinFill'" :size="42" :stroke="3"/>
      </div>
    </div>
    <h2 :style="{
      fontFamily:'var(--font-display)', fontWeight:'var(--disp-weight)', fontSize:'22px',
      letterSpacing:'var(--disp-ls)', color:'var(--ink)', margin:'26px 0 8px',
    }">{{ phase==='ok' ? '인증 완료' : '위치 확인 중...' }}</h2>
    <p style="font-size:14px;color:var(--ink-2);margin:0;line-height:1.5;word-break:keep-all;max-width:240px">
      {{ phase==='ok' ? `${place.name} 반경 ${place.radius_meters}m 안에 있어요.\n방명록이 열렸습니다.` : `${place.name}의 인증 반경 ${place.radius_meters}m 안에 있는지 확인하고 있어요` }}
    </p>
    <div v-if="phase==='ok'" style="margin-top:18px">
      <HnPill solid mono style="padding:6px 12px">
        <HnIcon name="bolt" :size="13" :stroke="2.4"/>스탬프 +1 · 단골
      </HnPill>
    </div>
    <p v-if="errorMsg" style="margin-top:12px;font-size:13px;color:var(--ink-2)">{{ errorMsg }}</p>
    <button v-if="phase==='scan'" @click="$emit('close')" style="margin-top:26px;padding:10px 22px;border-radius:999px;cursor:pointer;border:1px solid var(--line);background:var(--surface);color:var(--ink-2);font-size:13px;font-weight:600">취소</button>
  </div>
</template>
<script setup>
import { ref, onMounted } from 'vue'
import HnIcon from './HnIcon.vue'
import HnPill from './HnPill.vue'
import { API } from '@/config.js'
import { authFetch } from '@/stores/api.js'

const props = defineProps({ place: Object })
const emit = defineEmits(['done', 'close'])
const phase = ref('scan')
const errorMsg = ref('')

onMounted(() => {
  navigator.geolocation.getCurrentPosition(
    async (pos) => {
      try {
        const res = await authFetch(`${API}/places/${props.place.id}/verify`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ latitude: pos.coords.latitude, longitude: pos.coords.longitude }),
        })
        const data = await res.json()
        if (data.verified) {
          phase.value = 'ok'
          setTimeout(() => emit('done', data.token), 1500)
        } else {
          errorMsg.value = data.message || '이 장소에 가야 열려요'
          setTimeout(() => emit('close'), 2000)
        }
      } catch {
        errorMsg.value = '네트워크 오류가 발생했습니다'
        setTimeout(() => emit('close'), 2000)
      }
    },
    () => {
      errorMsg.value = '위치 권한을 허용해주세요'
      setTimeout(() => emit('close'), 2000)
    }
  )
})
</script>
