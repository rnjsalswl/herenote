<template>
  <div class="place">
    <button class="back" @click="router.push('/')">← 뒤로</button>

  <div v-if="loading" class="loading">불러오는 중...</div>
  <div v-else>
    <h1>{{ place.name }}</h1>
    <p class="desc">{{ place.description }}</p>

    <!-- 인증 중 -->
    <div v-if="verifying" class="status-box">
      <p>📍 위치 확인 중...</p>
    </div>

    <!-- 인증 실패 -->
    <div v-else-if="!verified" class="status-box fail">
      <p>🔒 이 장소에 가야 방명록이 열립니다</p>
      <p class="radius">인증 반경 {{ place.radius_meters }}m</p>
      <button @click="verify">다시 시도</button>
  </div>

  <!-- 인증 성공 -->
  <div v-else>
    <p class="success">✓ 위치 인증 완료</p>

    <!-- 방명록 작성 -->
    <div class="write-box">
      <textarea
              v-model="content"
              placeholder="방명록을 작성해주세요 (최대 500자)"
              maxlength="500"
      />
      <button @click="writeGuestbook" :disabled="writing || !content.trim()">
      {{ writing ? '작성 중...' : '방명록 남기기' }}
    </button>
  </div>

  <!-- 방명록 목록 -->
  <div class="guestbook-list">
    <h2>방명록 ({{ guestbooks.length }})</h2>
    <div v-if="guestbooks.length === 0" class="empty">
      첫 번째 방명록을 남겨보세요
    </div>
    <div
            v-for="g in guestbooks"
            :key="g.id"
            class="guestbook-card"
            :class="{ pinned: g.is_pinned }"
    >
      <div class="card-header">
        <span class="nickname">{{ g.nickname }}</span>
        <span v-if="g.badge_type === 'BLUE'">🔵</span>
        <span v-if="g.badge_type === 'GOLD'">🥇</span>
        <span v-if="g.badge_type === 'SPECIAL'">⭐</span>
        <span class="date">{{ formatDate(g.created_at) }}</span>
      </div>
      <p class="card-content">{{ g.content }}</p>
    </div>
  </div>
</div>
        </div>
        </div>
        </template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { API } from '@/config.js'
import { getUserID } from '@/stores/user.js'

const route = useRoute()
const router = useRouter()

const place = ref({})
const loading = ref(true)
const verified = ref(false)
const verifying = ref(true)
const token = ref('')
const guestbooks = ref([])
const content = ref('')
const writing = ref(false)

const USER_ID = getUserID()

onMounted(async () => {
// 장소 정보 로드
const res = await fetch(`${API}/places/${route.params.id}`)
place.value = await res.json()
loading.value = false

// GPS 자동 인증 시도
verify()
})

async function verify() {
verifying.value = true

navigator.geolocation.getCurrentPosition(
async (pos) => {
const res = await fetch(`${API}/places/${route.params.id}/verify`, {
method: 'POST',
headers: { 'Content-Type': 'application/json' },
body: JSON.stringify({
latitude: pos.coords.latitude,
longitude: pos.coords.longitude,
user_id: USER_ID,
}),
})
const data = await res.json()

if (data.verified) {
verified.value = true
token.value = data.token
await loadGuestbooks()
}
verifying.value = false
},
() => {
verifying.value = false
}
)
}

async function loadGuestbooks() {
try {
const res = await fetch(`${API}/places/${route.params.id}/guestbooks`, {
headers: { 'X-Location-Token': token.value },
})
const data = await res.json()
guestbooks.value = data.guestbooks || []
} catch (e) {
guestbooks.value = []
}
}

async function writeGuestbook() {
if (!content.value.trim()) return
writing.value = true

await fetch(`${API}/places/${route.params.id}/guestbooks`, {
method: 'POST',
headers: {
'Content-Type': 'application/json',
'X-Location-Token': token.value,
},
body: JSON.stringify({ content: content.value }),
})

content.value = ''
await loadGuestbooks()
writing.value = false
}

function formatDate(dateStr) {
return new Date(dateStr).toLocaleDateString('ko-KR', {
month: 'long', day: 'numeric',
hour: '2-digit', minute: '2-digit',
})
}
</script>

<style scoped>
.place {
max-width: 600px;
margin: 0 auto;
padding: 20px;
}
.back {
background: none;
border: none;
cursor: pointer;
font-size: 16px;
margin-bottom: 16px;
padding: 0;
}
.loading { color: #888; padding: 20px 0; }
h1 { font-size: 1.5rem; font-weight: bold; margin-bottom: 6px; }
.desc { color: #555; font-size: 14px; margin-bottom: 20px; }
.status-box {
background: #f0f4ff;
border-radius: 12px;
padding: 24px;
text-align: center;
margin: 20px 0;
}
.status-box.fail {
background: #fff5f5;
}
.status-box button {
background: #1A4FA0;
color: white;
border: none;
border-radius: 8px;
padding: 10px 24px;
margin-top: 12px;
cursor: pointer;
}
.radius { font-size: 12px; color: #888; margin-top: 4px; }
.success { color: #0F6E56; font-weight: bold; margin-bottom: 16px; }
.write-box { margin-bottom: 24px; }
.write-box textarea {
width: 100%;
height: 100px;
border: 1px solid #ddd;
border-radius: 8px;
padding: 12px;
font-size: 14px;
resize: none;
box-sizing: border-box;
}
.write-box button {
background: #1A4FA0;
color: white;
border: none;
border-radius: 8px;
padding: 10px 20px;
cursor: pointer;
margin-top: 8px;
}
.write-box button:disabled { background: #aaa; }
h2 { font-size: 16px; font-weight: bold; margin-bottom: 12px; }
.empty { color: #aaa; font-size: 14px; }
.guestbook-card {
border: 1px solid #eee;
border-radius: 12px;
padding: 14px;
margin-bottom: 10px;
}
.pinned {
border-color: #1A4FA0;
background: #f0f4ff;
}
.card-header {
display: flex;
align-items: center;
gap: 6px;
margin-bottom: 8px;
}
.nickname { font-weight: bold; font-size: 14px; }
.date { font-size: 11px; color: #888; margin-left: auto; }
.card-content { font-size: 14px; color: #333; }
</style>