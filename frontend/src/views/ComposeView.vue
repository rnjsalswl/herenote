<template>
  <div style="min-height:100vh;background:var(--bg)">
    <HnTopBar title="방명록 쓰기" :on-back="() => router.back()"/>

    <div v-if="loading" style="padding:40px 16px;font-family:var(--font-mono);font-size:12px;color:var(--ink-3);text-align:center">불러오는 중...</div>

    <div v-else style="padding:16px">
      <!-- 장소 표시 -->
      <div style="display:flex;align-items:center;gap:10px;padding:12px 14px;border-radius:var(--radius);background:var(--surface-2);border:1px solid var(--line);margin-bottom:16px">
        <HnIcon name="pinFill" :size="20" style="color:var(--ink);flex-shrink:0"/>
        <div style="flex:1;min-width:0">
          <div style="font-weight:700;font-size:15px;color:var(--ink);overflow:hidden;text-overflow:ellipsis;white-space:nowrap">{{ place.name }}</div>
          <div style="font-family:var(--font-mono);font-size:11px;color:var(--ink-3);margin-top:1px">인증 반경 {{ place.radius_meters }}m</div>
        </div>
        <HnPill solid mono><HnIcon name="check" :size="11" :stroke="3"/>인증됨</HnPill>
      </div>

      <!-- 작성 영역 -->
      <div style="background:var(--card);border:1px solid var(--line);border-radius:var(--radius-lg);padding:16px;box-shadow:var(--shadow)">
        <div style="display:flex;align-items:center;gap:10px;margin-bottom:14px">
          <HnAvatar :name="nickname" :size="38"/>
          <div>
            <div style="font-weight:700;font-size:14px;color:var(--ink)">{{ nickname }}</div>
            <div style="font-family:var(--font-mono);font-size:10.5px;color:var(--ink-3)">방명록 남기기</div>
          </div>
        </div>
        <textarea
          v-model="content"
          style="width:100%;border:none;outline:none;resize:none;font-size:15.5px;color:var(--ink);background:transparent;min-height:140px;line-height:1.65;font-family:var(--font-sans)"
          placeholder="이 장소에서의 순간을 기록해보세요..."
          maxlength="500"
          autofocus
        />
        <div style="display:flex;align-items:center;justify-content:space-between;padding-top:10px;border-top:1px solid var(--line);margin-top:8px">
          <span style="font-family:var(--font-mono);font-size:11px;color:var(--ink-3)">{{ content.length }}/500</span>
          <button @click="submit" :disabled="submitting || !content.trim()" :style="{
            padding:'9px 22px', borderRadius:'999px', border:'none',
            background:'var(--ink)', color:'var(--inv-ink)', fontSize:'14px', fontWeight:700,
            cursor: content.trim() ? 'pointer' : 'not-allowed',
            opacity: content.trim() ? 1 : 0.3, transition:'opacity .15s',
          }">{{ submitting ? '게시 중...' : '게시하기' }}</button>
        </div>
      </div>

      <p v-if="error" style="text-align:center;margin-top:12px;font-size:13px;color:var(--ink-2)">{{ error }}</p>
    </div>
  </div>
</template>
<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import HnTopBar from '@/components/HnTopBar.vue'
import HnIcon from '@/components/HnIcon.vue'
import HnPill from '@/components/HnPill.vue'
import HnAvatar from '@/components/HnAvatar.vue'
import { API } from '@/config.js'
import { getUserID } from '@/stores/user.js'
import { getToken } from '@/stores/verified.js'

const route = useRoute()
const router = useRouter()
const placeId = route.params.id
const userId = getUserID()
const nickname = userId || '나'
const token = getToken(placeId)

const place = ref({})
const loading = ref(true)
const content = ref('')
const submitting = ref(false)
const error = ref('')

onMounted(async () => {
  try {
    const res = await fetch(`${API}/places/${placeId}`)
    place.value = await res.json()
  } catch {
    error.value = '장소 정보를 불러올 수 없어요'
  }
  loading.value = false
})

async function submit() {
  if (!content.value.trim() || submitting.value) return
  submitting.value = true
  error.value = ''
  try {
    const res = await fetch(`${API}/places/${placeId}/guestbooks`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', 'X-Location-Token': token },
      body: JSON.stringify({ content: content.value }),
    })
    if (res.ok) {
      router.push(`/places/${placeId}`)
    } else {
      const data = await res.json().catch(() => ({}))
      error.value = data.error || '게시에 실패했습니다.'
    }
  } catch {
    error.value = '네트워크 오류가 발생했습니다.'
  }
  submitting.value = false
}
</script>
