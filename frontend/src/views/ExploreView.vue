<template>
  <div style="min-height:100vh;background:var(--bg);padding-bottom:96px">
    <HnTopBar title="지도 탐색" big>
      <template #right>
        <button :style="iconBtn"><HnIcon name="search" :size="23" :stroke="2"/></button>
      </template>
    </HnTopBar>

    <!-- 카카오맵 -->
    <div style="margin:12px 16px;border-radius:var(--radius-lg);overflow:hidden;border:1px solid var(--line);box-shadow:var(--shadow)">
      <HnKakaoMap
        :places="places"
        :zoom="6"
        :height="260"
        @marker-click="p => router.push(`/places/${p.id}`)"
      />
    </div>

    <!-- 장소 목록 -->
    <div style="padding:14px 16px 4px">
      <HnSectionLabel kicker="PLACES" title="모든 장소"/>
    </div>

    <div v-if="loading" style="padding:16px;font-family:var(--font-mono);font-size:12px;color:var(--ink-3)">불러오는 중...</div>

    <div v-else-if="places.length===0" style="padding:32px 16px;font-size:14px;color:var(--ink-3);text-align:center">등록된 장소가 없어요</div>

    <div v-else style="display:flex;flex-direction:column;gap:0;margin:8px 0 0">
      <RouterLink v-for="p in places" :key="p.id" :to="`/places/${p.id}`" :style="{
        display:'flex', alignItems:'center', gap:'12px', padding:'14px 16px',
        borderBottom:'1px solid var(--line)', textDecoration:'none', color:'inherit',
      }">
        <div style="width:46px;height:46px;border-radius:14px;overflow:hidden;flex-shrink:0">
          <HnMapThumb
            :lat="p.latitude"
            :lng="p.longitude"
            :seed="seedOf(p.id)"
            :radius="p.radius_meters"
            :height="46"
            :rounded="false"
            :show-pin="false"
          />
        </div>
        <div style="flex:1;min-width:0">
          <div style="font-weight:700;font-size:15px;color:var(--ink);overflow:hidden;text-overflow:ellipsis;white-space:nowrap">{{ p.name }}</div>
          <div style="font-size:12.5px;color:var(--ink-3);margin-top:2px;font-family:var(--font-mono);overflow:hidden;text-overflow:ellipsis;white-space:nowrap">{{ p.description || '설명 없음' }}</div>
        </div>
        <div style="display:flex;align-items:center;gap:6px;flex-shrink:0">
          <HnPill mono>{{ p.radius_meters }}m</HnPill>
          <HnIcon name="chevR" :size="16" :stroke="2.2" style="color:var(--ink-3)"/>
        </div>
      </RouterLink>
    </div>

    <div style="padding:16px;display:flex;justify-content:center">
      <RouterLink to="/add-place" :style="{
        display:'inline-flex', alignItems:'center', gap:'8px', padding:'11px 20px',
        borderRadius:'999px', border:'1px solid var(--line)', textDecoration:'none',
        color:'var(--ink)', fontSize:'14px', fontWeight:600,
      }">
        <HnIcon name="plus" :size="18" :stroke="2.2"/>
        장소 추가하기
      </RouterLink>
    </div>

    <HnTabBar active="explore" @tab="goTab" @compose="goCompose" style="position:fixed;bottom:0;left:0;right:0"/>
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
import HnKakaoMap from '@/components/HnKakaoMap.vue'
import HnSectionLabel from '@/components/HnSectionLabel.vue'
import { API } from '@/config.js'
import { lastVerifiedId } from '@/stores/verified.js'

const router = useRouter()
const places = ref([])
const loading = ref(true)
const iconBtn = { width:'40px', height:'40px', borderRadius:'50%', border:'none', background:'transparent', color:'var(--ink)', display:'flex', alignItems:'center', justifyContent:'center', cursor:'pointer' }

function seedOf(id) { return id ? id.charCodeAt(0) + (id.charCodeAt(1)||0) : 1 }

onMounted(async () => {
  try {
    const res = await fetch(`${API}/places`)
    const data = await res.json()
    places.value = data.places || data || []
  } catch {}
  loading.value = false
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
