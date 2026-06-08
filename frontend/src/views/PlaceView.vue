<template>
  <div class="place">
    <button class="back" @click="router.push('/')">← 뒤로</button>

    <div v-if="loading">불러오는 중...</div>
    <div v-else>
      <h1>{{ place.name }}</h1>
      <p class="desc">{{ place.description }}</p>
      <p class="radius">인증 반경 {{ place.radius_meters }}m</p>

      <!-- 위치 인증 전 -->
      <div v-if="!verified" class="verify-box">
        <p>📍 현장에 도착하셨나요?</p>
        <button @click="verify" :disabled="verifying">
          {{ verifying ? '인증 중...' : '위치 인증하기' }}
        </button>
        <p v-if="verifyError" class="error">{{ verifyError }}</p>
      </div>

      <!-- 위치 인증 후 -->
      <div v-else>
        <p class="success">✓ 위치 인증 완료</p>

        <!-- 방명록 작성 -->
        <div class="write-box">
          <textarea
            v-model="content"
            placeholder="방명록을 작성해주세요 (최대 500자)"
            maxlength="500"
          />
          <button @click="writeGuestbook" :disabled="writing">
            {{ writing ? '작성 중...' : '방명록 남기기' }}
          </button>
        </div>

        <!-- 방명록 목록 -->
        <div class="guestbook-list">
          <h2>방명록 ({{ guestbooks.length }})</h2>
          <div v-if="guestbooks.length === 0">아직 방명록이 없습니다</div>
          <div
            v-for="g in guestbooks"
            :key="g.id"
            class="guestbook-card"
            :class="{ pinned: g.is_pinned }"
          >
            <div class="card-header">
              <span class="nickname">{{ g.nickname }}</span>
              <span v-if="g.badge_type !== 'NONE'" class="badge">
                {{ g.badge_type === 'BLUE' ? '🔵' : g.badge_type === 'GOLD' ? '🥇' : '⭐' }}
              </span>
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

const USER_ID = getUserID()

const place = ref({})
const loading = ref(true)
const verified = ref(false)
const verifying = ref(false)
const verifyError = ref('')
const token = ref('')
const guestbooks = ref([])
const content = ref('')
const writing = ref(false)

onMounted(async () => {
  const res = await fetch(`${API}/places/${route.params.id}`)
  place.value = await res.json()
  loading.value = false
})

// GPS 위치 인증
async function verify() {
  verifying.value = true
  verifyError.value = ''

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

      alert(JSON.stringify(data))  // 임시 추가

      if (data.verified) {
        verified.value = true
        token.value = data.token
        loadGuestbooks()
      } else {
        verifyError.value = data.message
      }
      verifying.value = false
    },
    (err) => {
      verifyError.value = '위치 정보를 가져올 수 없습니다. 위치 권한을 허용해주세요.'
      verifying.value = false
    }
  )
}

// 방명록 목록
async function loadGuestbooks() {
  const res = await fetch(`${API}/places/${route.params.id}/guestbooks`, {
    headers: { 'X-Location-Token': token.value },
  })
  const data = await res.json()
  guestbooks.value = data.guestbooks
}

// 방명록 작성
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
    year: 'numeric', month: 'long', day: 'numeric',
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
.desc { color: #555; }
.radius { font-size: 12px; color: #888; }
.verify-box {
  background: #f0f4ff;
  border-radius: 12px;
  padding: 20px;
  text-align: center;
  margin: 20px 0;
}
.verify-box button {
  background: #1A4FA0;
  color: white;
  border: none;
  border-radius: 8px;
  padding: 12px 24px;
  font-size: 16px;
  cursor: pointer;
  margin-top: 12px;
}
.success { color: #0F6E56; font-weight: bold; margin: 12px 0; }
.error { color: red; margin-top: 8px; }
.write-box {
  margin: 20px 0;
}
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
.guestbook-card {
  border: 1px solid #eee;
  border-radius: 12px;
  padding: 14px;
  margin: 10px 0;
}
.pinned {
  border-color: #1A4FA0;
  background: #f0f4ff;
}
.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}
.nickname { font-weight: bold; }
.date { font-size: 12px; color: #888; margin-left: auto; }
.badge { font-size: 16px; }
.card-content { color: #333; }
</style>