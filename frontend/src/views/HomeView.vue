<template>
<button class="add-btn" @click="router.push('/add-place')">
  + 장소 추가
</button>
  <div class="home">
    <h1>HERE : NOTE</h1>
    <p>장소를 선택해서 방명록을 확인하세요</p>
    <div class="place-list">
      <div v-if="loading">불러오는 중...</div>
      <div v-else-if="places.length === 0">장소가 없습니다</div>
      <div
        v-for="place in places"
        :key="place.id"
        class="place-card"
        @click="goToPlace(place.id)"
      >
        <h2>{{ place.name }}</h2>
        <p>{{ place.description }}</p>
        <span class="radius">인증 반경 {{ place.radius_meters }}m</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { API } from '@/config.js'

const router = useRouter()
const places = ref([])
const loading = ref(true)

onMounted(async () => {
  const res = await fetch(`${API}/places`)
  const data = await res.json()
  places.value = data.places
  loading.value = false
})

function goToPlace(id) {
  router.push(`/places/${id}`)
}
</script>

<style scoped>
.home {
  max-width: 600px;
  margin: 0 auto;
  padding: 20px;
}

h1 {
  font-size: 2rem;
  font-weight: bold;
  margin-bottom: 8px;
}

.place-card {
  border: 1px solid #ddd;
  border-radius: 12px;
  padding: 16px;
  margin: 12px 0;
  cursor: pointer;
  transition: background 0.2s;
}

.place-card:hover {
  background: #f5f5f5;
}

.radius {
  font-size: 12px;
  color: #888;
}

.add-btn {
  background: #1A4FA0;
  color: white;
  border: none;
  border-radius: 8px;
  padding: 10px 20px;
  font-size: 14px;
  cursor: pointer;
  margin-bottom: 16px;
}
</style>