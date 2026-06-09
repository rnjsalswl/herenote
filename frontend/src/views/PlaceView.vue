<template>
  <div class="page">
    <header class="topbar">
      <button class="back-btn" @click="router.push('/')">
        <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round">
          <polyline points="15 18 9 12 15 6"/>
        </svg>
      </button>
      <span class="topbar-title">{{ loading ? '' : place.name }}</span>
      <div style="width:36px" />
    </header>

    <div v-if="loading" class="status-wrap">
      <p class="status-msg">불러오는 중...</p>
    </div>

    <div v-else>
      <!-- 장소 정보 -->
      <div class="place-info">
        <h1 class="place-name">{{ place.name }}</h1>
        <p v-if="place.description" class="place-desc">{{ place.description }}</p>
        <span class="place-radius">인증 반경 {{ place.radius_meters }}m</span>
      </div>

      <div class="line-sep" />

      <!-- 인증 중 -->
      <div v-if="verifying" class="status-wrap">
        <div class="status-icon">📍</div>
        <p class="status-msg">위치 확인 중...</p>
        <p class="status-sub">잠시만 기다려주세요</p>
      </div>

      <!-- 인증 실패 -->
      <div v-else-if="!verified" class="status-wrap fail">
        <div class="status-icon">🔒</div>
        <p class="status-msg">이 장소에 가야 열립니다</p>
        <p class="status-sub">인증 반경 {{ place.radius_meters }}m 안에 있어야 해요</p>
        <button class="btn-retry" @click="verify">다시 시도</button>
      </div>

      <!-- 인증 성공 -->
      <div v-else>
        <!-- 방명록 작성 (일반 모드만) -->
        <div v-if="!mineMode" class="composer">
          <div class="composer-avatar">✏️</div>
          <div class="composer-main">
            <textarea
              v-model="content"
              class="composer-ta"
              placeholder="방명록을 남겨보세요..."
              maxlength="500"
            />
            <div class="composer-footer">
              <span class="char-count">{{ content.length }}/500</span>
              <button class="btn-post" @click="writeGuestbook" :disabled="writing || !content.trim()">
                {{ writing ? '...' : '게시' }}
              </button>
            </div>
          </div>
        </div>

        <div class="line-sep" />

        <!-- 방명록 목록 -->
        <div class="feed-meta">
          <span class="feed-label">{{ mineMode ? '내 방명록' : '방명록' }}</span>
          <span class="feed-count">{{ displayedGuestbooks.length }}</span>
        </div>

        <div v-if="displayedGuestbooks.length === 0" class="empty">
          아직 방명록이 없어요. 첫 번째로 남겨보세요!
        </div>

        <div
          v-for="g in displayedGuestbooks"
          :key="g.id"
          class="post"
          :class="{ pinned: g.is_pinned }"
        >
          <div class="post-av">{{ g.nickname.charAt(0) }}</div>
          <div class="post-body">
            <div class="post-meta">
              <span class="post-name">{{ g.nickname }}</span>
              <span v-if="g.badge_type === 'BLUE'" class="post-badge">🔵</span>
              <span v-if="g.badge_type === 'GOLD'" class="post-badge">🥇</span>
              <span v-if="g.badge_type === 'SPECIAL'" class="post-badge">⭐</span>
              <span v-if="g.is_pinned" class="pin-tag">고정</span>
              <span class="post-date">{{ formatDate(g.created_at) }}</span>
            </div>
            <p class="post-content">{{ g.content }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { API } from '@/config.js'
import { getUserID } from '@/stores/user.js'

const route = useRoute()
const router = useRouter()

const USER_ID = getUserID()
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

<style scoped>
.page {
  max-width: 600px;
  margin: 0 auto;
  min-height: 100vh;
  background: #fff;
}

/* 탑바 */
.topbar {
  position: sticky;
  top: 0;
  z-index: 10;
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  border-bottom: 1px solid var(--separator);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 8px;
  height: 52px;
}

.back-btn {
  background: none;
  border: none;
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  color: var(--text-primary);
}

.back-btn:active { background: var(--primary-light); }

.topbar-title {
  font-size: 16px;
  font-weight: 600;
  flex: 1;
  text-align: center;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  padding: 0 8px;
}

/* 장소 정보 */
.place-info {
  padding: 20px 16px 16px;
}

.place-name {
  font-size: 22px;
  font-weight: 700;
  letter-spacing: -0.5px;
  margin-bottom: 6px;
}

.place-desc {
  font-size: 14px;
  color: var(--text-secondary);
  margin-bottom: 8px;
  line-height: 1.5;
}

.place-radius {
  font-size: 12px;
  color: var(--text-hint);
  background: var(--primary-light);
  padding: 3px 10px;
  border-radius: 999px;
}

.line-sep {
  height: 1px;
  background: var(--separator);
  margin: 0;
}

/* 상태 */
.status-wrap {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 48px 24px;
  text-align: center;
}

.status-icon { font-size: 40px; margin-bottom: 14px; }
.status-msg { font-size: 16px; font-weight: 600; margin-bottom: 6px; }
.status-sub { font-size: 13px; color: var(--text-secondary); margin-bottom: 20px; }

.btn-retry {
  background: var(--primary);
  color: #fff;
  border: none;
  border-radius: 999px;
  padding: 10px 28px;
  font-size: 14px;
  font-weight: 600;
}

/* 작성 */
.composer {
  display: flex;
  gap: 12px;
  padding: 16px;
  border-bottom: 1px solid var(--separator);
}

.composer-avatar {
  font-size: 18px;
  width: 38px;
  height: 38px;
  background: var(--primary-light);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.composer-main { flex: 1; }

.composer-ta {
  width: 100%;
  border: none;
  outline: none;
  resize: none;
  font-size: 15px;
  color: var(--text-primary);
  background: transparent;
  min-height: 72px;
  line-height: 1.5;
}

.composer-ta::placeholder { color: var(--text-hint); }

.composer-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 8px;
}

.char-count { font-size: 12px; color: var(--text-hint); }

.btn-post {
  background: var(--primary);
  color: #fff;
  border: none;
  border-radius: 999px;
  padding: 7px 18px;
  font-size: 13px;
  font-weight: 600;
  transition: opacity 0.15s;
}

.btn-post:disabled { opacity: 0.3; cursor: not-allowed; }

/* 피드 */
.feed-meta {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 14px 16px 4px;
}

.feed-label { font-size: 13px; font-weight: 600; color: var(--text-secondary); }
.feed-count {
  font-size: 12px;
  color: var(--text-hint);
  background: var(--primary-light);
  padding: 1px 8px;
  border-radius: 999px;
}

.empty {
  padding: 32px 16px;
  font-size: 14px;
  color: var(--text-secondary);
  text-align: center;
}

/* 포스트 */
.post {
  display: flex;
  gap: 12px;
  padding: 14px 16px;
  border-bottom: 1px solid var(--separator);
}

.post.pinned { background: #FAFAFA; }

.post-av {
  width: 38px;
  height: 38px;
  background: var(--primary);
  color: #fff;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 15px;
  font-weight: 700;
  flex-shrink: 0;
}

.post-body { flex: 1; min-width: 0; }

.post-meta {
  display: flex;
  align-items: center;
  gap: 5px;
  margin-bottom: 5px;
  flex-wrap: wrap;
}

.post-name { font-size: 14px; font-weight: 600; }
.post-badge { font-size: 13px; }
.pin-tag {
  font-size: 10px;
  color: #fff;
  background: var(--primary);
  padding: 1px 6px;
  border-radius: 999px;
  font-weight: 600;
}
.post-date { font-size: 12px; color: var(--text-hint); margin-left: auto; }
.post-content { font-size: 14px; color: var(--text-primary); line-height: 1.6; }
</style>
