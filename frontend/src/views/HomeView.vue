<template>
  <div class="home">
    <!-- 헤더 -->
    <div class="header">
      <h1>HERE : NOTE</h1>
      <div class="header-actions">
        <button class="add-btn" @click="router.push('/add-place')">+ 장소</button>
        <button class="logout-btn" @click="logout">로그아웃</button>
      </div>
    </div>

    <!-- 현재 위치 근처 장소 -->
    <div class="section">
      <h2>📍 지금 여기</h2>
      <div v-if="locating" class="loading">위치 확인 중...</div>
      <div v-else-if="nearbyPlaces.length > 0">
        <div
          v-for="place in nearbyPlaces"
          :key="place.id"
          class="place-card nearby"
          @click="goToPlace(place.id)"
        >
          <div class="place-info">
            <span class="place-name">{{ place.name }}</span>
            <span class="place-desc">{{ place.description }}</span>
          </div>
          <span class="arrow">→</span>
        </div>
      </div>
      <div v-else class="empty">근처에 등록된 장소가 없어요</div>
    </div>

    <!-- 내가 다녀간 장소 -->
    <div class="section">
      <h2>🗺 내가 다녀간 곳</h2>
      <div v-if="myPlaces.length > 0">
        <div
          v-for="place in myPlaces"
          :key="place.id"
          class="place-card"
          @click="router.push({ path: `/places/${place.id}`, query: { mine: 'true' } })"
        >
          <div class="place-info">
            <span class="place-name">{{ place.name }}</span>
            <span class="place-desc">{{ place.description }}</span>
          </div>
          <span class="arrow">→</span>
        </div>
      </div>
      <div v-else class="empty">아직 방명록을 남긴 장소가 없어요</div>
    </div>
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
  // 내가 다녀간 장소 로드
  const userID = getUserID()
  if (userID) {
    try {
      const res = await fetch(`${API}/users/${userID}/guestbooks/places`)
      const data = await res.json()
      myPlaces.value = data.places || []
    } catch {}
  }

  // 현재 위치 근처 장소 조회
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
    () => {
      locating.value = false
    }
  )
})

function goToPlace(id) {
  router.push(`/places/${id}`)
}

function logout() {
  clearUser()
  router.push('/login')
}
</script>

<style scoped>
.home {
  max-width: 600px;
  margin: 0 auto;
  padding: 20px;
}
.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
}
h1 {
  font-size: 1.8rem;
  font-weight: bold;
  color: #1A4FA0;
}
.header-actions {
  display: flex;
  gap: 8px;
  align-items: center;
}
.add-btn {
  background: #1A4FA0;
  color: white;
  border: none;
  border-radius: 8px;
  padding: 8px 16px;
  font-size: 14px;
  cursor: pointer;
}
.logout-btn {
  background: none;
  color: #888;
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 8px 12px;
  font-size: 13px;
  cursor: pointer;
}
.logout-btn:hover {
  background: #f5f5f5;
}
.section {
  margin-bottom: 32px;
}
h2 {
  font-size: 16px;
  font-weight: bold;
  color: #333;
  margin-bottom: 12px;
}
.place-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  border: 1px solid #eee;
  border-radius: 12px;
  padding: 14px 16px;
  margin-bottom: 10px;
  cursor: pointer;
  transition: background 0.15s;
}
.place-card:hover {
  background: #f5f5f5;
}
.place-card.nearby {
  border-color: #1A4FA0;
  background: #f0f4ff;
}
.place-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.place-name {
  font-size: 15px;
  font-weight: 500;
  color: #1A1A2E;
}
.place-desc {
  font-size: 12px;
  color: #888;
}
.arrow {
  color: #1A4FA0;
  font-size: 18px;
}
.empty {
  color: #aaa;
  font-size: 14px;
  padding: 16px 0;
}
.loading {
  color: #888;
  font-size: 14px;
  padding: 16px 0;
}
</style>
