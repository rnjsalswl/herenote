<template>
  <div style="min-height:100vh;background:var(--bg);position:relative">
    <HnTopBar :title="loading ? '' : place.name" :on-back="() => router.push('/')"/>

    <div v-if="loading" style="display:flex;align-items:center;justify-content:center;padding:64px 24px;font-family:var(--font-mono);font-size:12px;color:var(--ink-3)">불러오는 중...</div>

    <div v-else>
      <!-- 장소 헤더 -->
      <div style="padding:0 16px 16px">
        <HnMapThumb :seed="seedOf(place.id)" :radius="place.radius_meters" :height="160" :rounded="true"/>
        <div style="margin-top:14px">
          <h1 style="font-family:var(--font-display);font-weight:var(--disp-weight);font-size:22px;letter-spacing:var(--disp-ls);color:var(--ink);margin:0 0 6px">{{ place.name }}</h1>
          <p v-if="place.description" style="font-size:14px;color:var(--ink-2);margin:0 0 10px;line-height:1.5">{{ place.description }}</p>
          <HnPill mono><HnIcon name="pin" :size="11" :stroke="2.2"/> 인증 반경 {{ place.radius_meters }}m</HnPill>
        </div>
      </div>

      <div style="height:1px;background:var(--line);margin:0 0 0"/>

      <!-- 인증 중 -->
      <div v-if="verifying" style="display:flex;flex-direction:column;align-items:center;padding:48px 24px;text-align:center">
        <div style="font-family:var(--font-mono);font-size:12px;color:var(--ink-3)">위치 확인 중...</div>
      </div>

      <!-- 인증 실패 -->
      <div v-else-if="!verified" style="display:flex;flex-direction:column;align-items:center;padding:48px 24px;text-align:center;gap:12px">
        <div style="width:72px;height:72px;border-radius:50%;background:var(--surface-2);border:1px solid var(--line);display:flex;align-items:center;justify-content:center;color:var(--ink-2)">
          <HnIcon name="lock" :size="32" :stroke="1.8"/>
        </div>
        <h2 style="font-family:var(--font-display);font-weight:var(--disp-weight);font-size:19px;letter-spacing:var(--disp-ls);color:var(--ink);margin:0">이 장소에 가야 열립니다</h2>
        <p style="font-size:13.5px;color:var(--ink-2);margin:0;max-width:220px;line-height:1.5">인증 반경 {{ place.radius_meters }}m 안에 있어야 해요</p>
        <button @click="verify" style="margin-top:6px;padding:11px 24px;border-radius:999px;border:1px solid var(--line);background:var(--surface);color:var(--ink);font-size:14px;font-weight:600;cursor:pointer">다시 시도</button>
      </div>

      <!-- 인증 성공 -->
      <div v-else>
        <!-- 방명록 작성 (일반 모드만) -->
        <div v-if="!mineMode" style="padding:14px 16px;border-bottom:1px solid var(--line);display:flex;gap:12px">
          <HnAvatar :name="nickname" :size="38"/>
          <div style="flex:1;min-width:0">
            <textarea
              v-model="content"
              style="width:100%;border:none;outline:none;resize:none;font-size:15px;color:var(--ink);background:transparent;min-height:72px;line-height:1.62;font-family:var(--font-sans)"
              placeholder="방명록을 남겨보세요..."
              maxlength="500"
            />
            <div style="display:flex;align-items:center;justify-content:space-between;margin-top:6px">
              <span style="font-family:var(--font-mono);font-size:11px;color:var(--ink-3)">{{ content.length }}/500</span>
              <button @click="writeGuestbook" :disabled="writing || !content.trim()" :style="{
                padding:'7px 18px', borderRadius:'999px', border:'none', cursor: content.trim() ? 'pointer' : 'not-allowed',
                background:'var(--ink)', color:'var(--inv-ink)', fontSize:'13px', fontWeight:700,
                opacity: content.trim() ? 1 : 0.3, transition:'opacity .15s',
              }">{{ writing ? '...' : '게시' }}</button>
            </div>
          </div>
        </div>

        <!-- 방명록 목록 헤더 -->
        <div style="padding:14px 16px 6px;display:flex;align-items:center;gap:8px">
          <span style="font-family:var(--font-display);font-weight:var(--disp-weight);font-size:16px;letter-spacing:var(--disp-ls);color:var(--ink)">{{ mineMode ? '내 방명록' : '방명록' }}</span>
          <HnPill mono>{{ displayedGuestbooks.length }}</HnPill>
        </div>

        <div v-if="displayedGuestbooks.length === 0" style="padding:32px 16px;font-size:14px;color:var(--ink-3);text-align:center">
          아직 방명록이 없어요. 첫 번째로 남겨보세요!
        </div>

        <div style="display:flex;flex-direction:column;gap:0">
          <div
            v-for="g in displayedGuestbooks"
            :key="g.id"
            :style="{
              display:'flex', gap:'12px', padding:'14px 16px',
              borderBottom:'1px solid var(--line)',
              background: g.is_pinned ? 'var(--surface-2)' : 'var(--card)',
            }"
          >
            <HnAvatar :name="g.nickname" :size="38"/>
            <div style="flex:1;min-width:0">
              <div style="display:flex;align-items:center;gap:5px;margin-bottom:5px;flex-wrap:wrap">
                <span style="font-size:14px;font-weight:700;color:var(--ink)">{{ g.nickname }}</span>
                <HnBadge :type="g.badge_type"/>
                <HnPill v-if="g.is_pinned" style="padding:2px 7px;font-size:10.5px">고정</HnPill>
                <span style="font-family:var(--font-mono);font-size:11px;color:var(--ink-3);margin-left:auto">{{ formatDate(g.created_at) }}</span>
              </div>
              <p style="font-size:14px;color:var(--ink);line-height:1.6;margin:0;white-space:pre-wrap;word-break:keep-all">{{ g.content }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import HnTopBar from '@/components/HnTopBar.vue'
import HnIcon from '@/components/HnIcon.vue'
import HnPill from '@/components/HnPill.vue'
import HnAvatar from '@/components/HnAvatar.vue'
import HnBadge from '@/components/HnBadge.vue'
import HnMapThumb from '@/components/HnMapThumb.vue'
import { API } from '@/config.js'
import { getUserID } from '@/stores/user.js'
import { authFetch } from '@/stores/api.js'

const route = useRoute()
const router = useRouter()

const USER_ID = getUserID()
const nickname = USER_ID || '나'
const mineMode = computed(() => route.query.mine === 'true')
const displayedGuestbooks = computed(() =>
  mineMode.value ? guestbooks.value.filter(g => g.user_id === USER_ID) : guestbooks.value
)

const place = ref({})
const loading = ref(true)
const verified = ref(false)
const verifying = ref(true)
const token = ref('')
const guestbooks = ref([])
const content = ref('')
const writing = ref(false)

function seedOf(id) { return id ? id.charCodeAt(0) + (id.charCodeAt(1)||0) : 1 }

onMounted(async () => {
  try {
    const res = await fetch(`${API}/places/${route.params.id}`)
    place.value = await res.json()
  } finally {
    loading.value = false
  }
  verify()
})

async function verify() {
  verifying.value = true
  navigator.geolocation.getCurrentPosition(
    async (pos) => {
      try {
        const res = await authFetch(`${API}/places/${route.params.id}/verify`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            latitude: pos.coords.latitude,
            longitude: pos.coords.longitude,
          }),
        })
        const data = await res.json()
        if (data.verified) {
          verified.value = true
          token.value = data.token
          await loadGuestbooks()
        }
      } finally {
        verifying.value = false
      }
    },
    () => { verifying.value = false }
  )
}

async function loadGuestbooks() {
  try {
    const res = await fetch(`${API}/places/${route.params.id}/guestbooks`, {
      headers: { 'X-Location-Token': token.value },
    })
    const data = await res.json()
    guestbooks.value = data.guestbooks || []
  } catch {
    guestbooks.value = []
  }
}

async function writeGuestbook() {
  if (!content.value.trim()) return
  writing.value = true
  await fetch(`${API}/places/${route.params.id}/guestbooks`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', 'X-Location-Token': token.value },
    body: JSON.stringify({ content: content.value }),
  })
  content.value = ''
  await loadGuestbooks()
  writing.value = false
}

function formatDate(dateStr) {
  return new Date(dateStr).toLocaleDateString('ko-KR', {
    month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit',
  })
}
</script>
