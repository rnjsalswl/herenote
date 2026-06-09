<template>
  <div style="min-height:100vh;background:var(--bg);padding-bottom:96px">
    <HnTopBar title="HERE NOTE" big mono>
      <template #right>
        <button :style="iconBtn"><HnIcon name="search" :size="23" :stroke="2"/></button>
      </template>
    </HnTopBar>

    <!-- 지금 여기 -->
    <div style="padding:14px 16px 4px">
      <HnSectionLabel kicker="NOW" title="지금 여기" :right="locating ? '' : nearbyPlaces.length + '곳 열림'"/>
    </div>
    <div v-if="locating" style="padding:16px 16px;font-family:var(--font-mono);font-size:12px;color:var(--ink-3)">위치 확인 중...</div>
    <div v-else-if="nearbyPlaces.length===0" style="padding:16px 16px;font-size:14px;color:var(--ink-3)">근처에 열린 장소가 없어요</div>
    <div v-else style="display:flex;gap:12px;overflow-x:auto;padding:4px 16px 14px;scrollbar-width:none">
      <RouterLink v-for="p in nearbyPlaces" :key="p.id" :to="`/places/${p.id}`" :style="{
        flexShrink:0, width:'84%', textDecoration:'none', color:'inherit',
        border:'1px solid var(--line)', background:'var(--card)',
        borderRadius:'var(--radius-lg)', overflow:'hidden', boxShadow:'var(--shadow)', display:'block',
      }">
        <div style="position:relative">
          <HnMapThumb :seed="seedOf(p.id)" :radius="p.radius_meters" :height="132" :rounded="false"/>
          <div style="position:absolute;top:10px;left:10px">
            <HnPill solid mono style="gap:4px"><HnIcon name="check" :size="12" :stroke="3"/>여기 도착 · {{p.radius_meters}}m</HnPill>
          </div>
        </div>
        <div style="padding:12px 14px 14px">
          <div style="font-weight:700;font-size:16px;color:var(--ink)">{{ p.name }}</div>
          <div style="font-size:12.5px;color:var(--ink-3);margin-top:2px;font-family:var(--font-mono)">{{ p.description }}</div>
          <div style="margin-top:12px;display:flex;align-items:center;justify-content:space-between;padding:10px 12px;border-radius:12px;background:var(--ink);color:var(--inv-ink)">
            <span style="font-size:13.5px;font-weight:700">방명록 열기</span>
            <HnIcon name="chevR" :size="18" :stroke="2.6"/>
          </div>
        </div>
      </RouterLink>
    </div>

    <div style="height:8px;background:var(--surface-2);border-top:1px solid var(--line);border-bottom:1px solid var(--line);margin:6px 0 14px"/>

    <!-- 피드 -->
    <div style="padding:0 16px 4px">
      <HnSectionLabel kicker="FEED" title="다녀간 곳의 새 글"/>
    </div>
    <div style="display:flex;flex-direction:column;gap:14px;padding:10px 16px 0">
      <HnNoteCard v-for="n in feed" :key="n.id" :note="n" :place="n.place" @open="router.push(`/places/${n.placeId}`)"/>
    </div>

    <HnTabBar active="home" @tab="goTab" @compose="goCompose" style="position:fixed;bottom:0;left:0;right:0"/>
  </div>
</template>
<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import HnTopBar from '@/components/HnTopBar.vue'
import HnTabBar from '@/components/HnTabBar.vue'
import HnIcon from '@/components/HnIcon.vue'
import HnPill from '@/components/HnPill.vue'
import HnMapThumb from '@/components/HnMapThumb.vue'
import HnNoteCard from '@/components/HnNoteCard.vue'
import HnSectionLabel from '@/components/HnSectionLabel.vue'
import { API } from '@/config.js'
import { getUserID } from '@/stores/user.js'
import { authFetch } from '@/stores/api.js'
import { lastVerifiedId } from '@/stores/verified.js'

const router = useRouter()
const nearbyPlaces = ref([])
const myPlaces = ref([])
const locating = ref(true)
const iconBtn = { width:'40px', height:'40px', borderRadius:'50%', border:'none', background:'transparent', color:'var(--ink)', display:'flex', alignItems:'center', justifyContent:'center', cursor:'pointer' }

const feed = ref([
  { id:'f1', place:'경복궁 향원정', placeId:'p1', nickname:'도윤', badge:'GOLD', time:'2시간 전', likes:48, liked:false, photo:'pond', content:'비 갠 뒤 향원정. 물 냄새가 좋아서 한참을 앉아 있었다.' },
  { id:'f2', place:'연남동 경의선숲길', placeId:'p2', nickname:'준호', badge:'BLUE', time:'1시간 전', likes:33, liked:false, photo:null, content:'돗자리 깔고 맥주 한 캔. 기차가 다니던 길이라니 신기하다.' },
  { id:'f3', place:'경복궁 향원정', placeId:'p1', nickname:'서연', badge:'BLUE', time:'5시간 전', likes:21, liked:true, photo:null, content:'여기서 프로포즈 받았어요. 매년 와서 한 줄씩 남기는 중 🙂' },
])

function seedOf(id) { return id ? id.charCodeAt(0) + (id.charCodeAt(1)||0) : 1 }

onMounted(async () => {
  const uid = getUserID()
  if (uid) {
    try {
      const res = await authFetch(`${API}/users/${uid}/guestbooks/places`)
      const data = await res.json()
      myPlaces.value = data.places || []
    } catch {}
  }
  navigator.geolocation.getCurrentPosition(
    async (pos) => {
      try {
        const res = await fetch(`${API}/places/nearby?lat=${pos.coords.latitude}&lng=${pos.coords.longitude}`)
        const data = await res.json()
        nearbyPlaces.value = data.places || []
      } catch {}
      locating.value = false
    },
    () => { locating.value = false }
  )
})

function goTab(id) {
  if (id === 'home') router.push('/')
  else if (id === 'explore') router.push('/explore')
  else if (id === 'profile') router.push('/profile')
}
function goCompose() {
  const id = lastVerifiedId()
  if (id) router.push(`/compose/${id}`)
  else router.push('/explore')
}
</script>
