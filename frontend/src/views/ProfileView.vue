<template>
  <div style="min-height:100vh;background:var(--bg);padding-bottom:96px">
    <HnTopBar title="프로필" big>
      <template #right>
        <button :style="iconBtn" @click="logout">
          <HnIcon name="settings" :size="22" :stroke="1.9"/>
        </button>
      </template>
    </HnTopBar>

    <!-- 프로필 헤더 -->
    <div style="padding:20px 20px 16px;display:flex;align-items:center;gap:16px">
      <HnAvatar :name="nickname" :size="72" :filled="true"/>
      <div style="flex:1;min-width:0">
        <div style="font-family:var(--font-display);font-weight:var(--disp-weight);font-size:20px;letter-spacing:var(--disp-ls);color:var(--ink)">{{ nickname }}</div>
        <div style="font-family:var(--font-mono);font-size:11px;color:var(--ink-3);margin-top:4px">@{{ nickname }}</div>
        <div style="display:flex;gap:8px;margin-top:10px">
          <HnPill mono><HnIcon name="pin" :size="11" :stroke="2.2"/> {{ visitCount }}곳 방문</HnPill>
          <HnPill mono><HnIcon name="bolt" :size="11" :stroke="2.2"/> {{ noteCount }}개 방명록</HnPill>
        </div>
      </div>
    </div>

    <div style="height:8px;background:var(--surface-2);border-top:1px solid var(--line);border-bottom:1px solid var(--line);margin:4px 0 14px"/>

    <!-- 다녀간 장소 -->
    <div style="padding:0 16px 8px">
      <HnSectionLabel kicker="VISITED" title="다녀간 장소"/>
    </div>

    <div v-if="loadingPlaces" style="padding:16px;font-family:var(--font-mono);font-size:12px;color:var(--ink-3)">불러오는 중...</div>
    <div v-else-if="myPlaces.length===0" style="padding:32px 16px;font-size:14px;color:var(--ink-3);text-align:center">아직 방명록을 남긴 장소가 없어요</div>
    <div v-else style="display:flex;flex-direction:column;gap:0">
      <RouterLink v-for="p in myPlaces" :key="p.id"
        :to="{ path: `/places/${p.id}`, query: { mine: 'true' } }"
        :style="{
          display:'flex', alignItems:'center', gap:'12px', padding:'14px 16px',
          borderBottom:'1px solid var(--line)', textDecoration:'none', color:'inherit',
        }">
        <div style="width:44px;height:44px;border-radius:12px;overflow:hidden;flex-shrink:0">
          <HnMapThumb :seed="seedOf(p.id)" :radius="p.radius_meters||100" :height="44" :rounded="false" :show-pin="false"/>
        </div>
        <div style="flex:1;min-width:0">
          <div style="font-weight:700;font-size:15px;color:var(--ink);overflow:hidden;text-overflow:ellipsis;white-space:nowrap">{{ p.name }}</div>
          <div style="font-size:12px;color:var(--ink-3);margin-top:2px;font-family:var(--font-mono)">내 방명록 보기</div>
        </div>
        <HnIcon name="chevR" :size="16" :stroke="2.2" style="color:var(--ink-3);flex-shrink:0"/>
      </RouterLink>
    </div>

    <!-- 로그아웃 -->
    <div style="padding:24px 16px 0">
      <button @click="logout" style="width:100%;padding:13px;border-radius:14px;border:1px solid var(--line);background:var(--surface-2);color:var(--ink-2);font-size:14px;font-weight:600;cursor:pointer">
        로그아웃
      </button>
    </div>

    <HnTabBar active="profile" @tab="goTab" @compose="goCompose" style="position:fixed;bottom:0;left:0;right:0"/>
  </div>
</template>
<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import HnTopBar from '@/components/HnTopBar.vue'
import HnTabBar from '@/components/HnTabBar.vue'
import HnIcon from '@/components/HnIcon.vue'
import HnPill from '@/components/HnPill.vue'
import HnAvatar from '@/components/HnAvatar.vue'
import HnMapThumb from '@/components/HnMapThumb.vue'
import HnSectionLabel from '@/components/HnSectionLabel.vue'
import { API } from '@/config.js'
import { getUserID, clearUser, getNickname } from '@/stores/user.js'
import { authFetch } from '@/stores/api.js'
import { lastVerifiedId } from '@/stores/verified.js'

const router = useRouter()
const userId = getUserID() || ''
const nickname = getNickname() || userId || '나'
const myPlaces = ref([])
const loadingPlaces = ref(true)
const visitCount = ref(0)
const noteCount = ref(0)
const iconBtn = { width:'40px', height:'40px', borderRadius:'50%', border:'none', background:'transparent', color:'var(--ink)', display:'flex', alignItems:'center', justifyContent:'center', cursor:'pointer' }

function seedOf(id) { return id ? id.charCodeAt(0) + (id.charCodeAt(1)||0) : 1 }

onMounted(async () => {
  if (userId) {
    try {
      const res = await authFetch(`${API}/users/${userId}/guestbooks/places`)
      const data = await res.json()
      myPlaces.value = data.places || []
      visitCount.value = myPlaces.value.length
    } catch {}
  }
  loadingPlaces.value = false
})

function logout() {
  clearUser()
  router.push('/login')
}

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
