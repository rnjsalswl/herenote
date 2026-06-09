<template>
  <div class="page">
    <header class="topbar">
      <span class="topbar-title">HERE NOTE</span>
      <div class="topbar-actions">
        <button class="icon-btn" @click="router.push('/add-place')">
          <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round">
            <line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/>
          </svg>
        </button>
        <button class="icon-btn" @click="logout">
          <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/>
            <polyline points="16 17 21 12 16 7"/><line x1="21" y1="12" x2="9" y2="12"/>
          </svg>
        </button>
      </div>
    </header>

    <!-- 지금 여기 -->
    <section>
      <p class="section-label">지금 여기</p>
      <div v-if="locating" class="state-msg">위치 확인 중...</div>
      <template v-else-if="nearbyPlaces.length > 0">
        <RouterLink
          v-for="place in nearbyPlaces"
          :key="place.id"
          :to="`/places/${place.id}`"
          class="row"
        >
          <div class="row-icon">📍</div>
          <div class="row-body">
            <span class="row-name">{{ place.name }}</span>
            <span class="row-desc">{{ place.description }}</span>
            <span class="row-meta">반경 {{ place.radius_meters }}m</span>
          </div>
          <svg class="row-chev" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
            <polyline points="9 18 15 12 9 6"/>
          </svg>
        </RouterLink>
      </template>
      <div v-else class="state-msg">근처에 등록된 장소가 없어요</div>
    </section>

    <div class="block-sep" />

    <!-- 내가 다녀간 곳 -->
    <section>
      <p class="section-label">내가 다녀간 곳</p>
      <template v-if="myPlaces.length > 0">
        <RouterLink
          v-for="place in myPlaces"
          :key="place.id"
          :to="{ path: `/places/${place.id}`, query: { mine: 'true' } }"
          class="row"
        >
          <div class="row-icon">🗺️</div>
          <div class="row-body">
            <span class="row-name">{{ place.name }}</span>
            <span class="row-desc">{{ place.description }}</span>
          </div>
          <svg class="row-chev" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
            <polyline points="9 18 15 12 9 6"/>
          </svg>
        </RouterLink>
      </template>
      <div v-else class="state-msg">아직 방명록을 남긴 장소가 없어요</div>
    </section>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { API } from '@/config.js'
import { getUserID, clearUser } from '@/stores/user.js'

const router = useRouter()
const nearbyPlaces = ref([])
const myPlaces = ref([])
const locating = ref(true)

onMounted(async () => {
  const userID = getUserID()
  if (userID) {
    try {
      const res = await fetch(`${API}/users/${userID}/guestbooks/places`)
      const data = await res.json()
      myPlaces.value = data.places || []
    } catch {}
  }
  navigator.geolocation.getCurrentPosition(
    async (pos) => {
      try {
        const { latitude, longitude } = pos.coords
        const res = await fetch(`${API}/places/nearby?lat=${latitude}&lng=${longitude}`)
        const data = await res.json()
        nearbyPlaces.value = data.places || []
      } catch {}
      locating.value = false
    },
    () => { locating.value = false }
  )
})

function logout() {
  clearUser()
  router.push('/login')
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
  padding: 10px 16px;
  height: 52px;
}

.topbar-title {
  font-size: 17px;
  font-weight: 700;
  letter-spacing: -0.3px;
}

.topbar-actions { display: flex; gap: 2px; }

.icon-btn {
  background: none;
  border: none;
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  color: var(--text-primary);
  transition: background 0.1s;
}

.icon-btn:active { background: var(--primary-light); }

/* 섹션 */
.section-label {
  padding: 18px 16px 6px;
  font-size: 12px;
  font-weight: 600;
  color: var(--text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.6px;
}

.block-sep {
  height: 8px;
  background: var(--primary-light);
  border-top: 1px solid var(--separator);
  border-bottom: 1px solid var(--separator);
}

/* 장소 row */
.row {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border-bottom: 1px solid var(--separator);
  text-decoration: none;
  color: inherit;
  -webkit-tap-highlight-color: transparent;
}

.row:active { background: var(--primary-light); }

.row-icon {
  width: 42px;
  height: 42px;
  background: var(--primary-light);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  flex-shrink: 0;
}

.row-body {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.row-name {
  font-size: 15px;
  font-weight: 600;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.row-desc {
  font-size: 13px;
  color: var(--text-secondary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.row-meta {
  font-size: 12px;
  color: var(--text-hint);
  margin-top: 2px;
}

.row-chev { color: var(--text-hint); flex-shrink: 0; }

.state-msg {
  padding: 16px 16px;
  font-size: 14px;
  color: var(--text-secondary);
}
</style>
