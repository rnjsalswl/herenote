<template>
  <div class="add-place">
    <button class="back" @click="router.push('/')">← 뒤로</button>
    <h1>장소 추가</h1>

    <!-- 위치 가져오기 -->
    <div class="location-box">
      <button @click="getLocation" :disabled="locating">
        {{ locating ? '위치 가져오는 중...' : '📍 현재 위치 가져오기' }}
      </button>
      <div v-if="location" class="location-result">
        <p>위도: {{ location.latitude }}</p>
        <p>경도: {{ location.longitude }}</p>
        <p class="hint">현재 위치가 장소 좌표로 등록돼요</p>
      </div>
      <p v-if="locationError" class="error">{{ locationError }}</p>
    </div>

    <!-- 장소 정보 입력 -->
    <div v-if="location" class="form">
      <div class="field">
        <label>장소 이름</label>
        <input v-model="form.name" placeholder="예: 경복궁 정문" />
      </div>
      <div class="field">
        <label>장소 설명</label>
        <textarea v-model="form.description" placeholder="이 장소에 대해 설명해주세요" />
      </div>
      <div class="field">
        <label>인증 반경 (미터)</label>
        <select v-model="form.radius_meters">
          <option :value="50">50m — 매우 좁음 (실내)</option>
          <option :value="100">100m — 좁음 (광장)</option>
          <option :value="200">200m — 보통 (공원)</option>
          <option :value="500">500m — 넓음 (산 정상)</option>
        </select>
      </div>

      <button class="submit" @click="submit" :disabled="submitting || !form.name">
        {{ submitting ? '등록 중...' : '장소 등록하기' }}
      </button>
      <p v-if="submitError" class="error">{{ submitError }}</p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { API } from '@/config.js'

const router = useRouter()
const ADMIN_SECRET = 'herenote-admin'

const location = ref(null)
const locating = ref(false)
const locationError = ref('')

const form = ref({
  name: '',
  description: '',
  radius_meters: 100,
})

const submitting = ref(false)
const submitError = ref('')

// 현재 위치 가져오기
function getLocation() {
  locating.value = true
  locationError.value = ''

  navigator.geolocation.getCurrentPosition(
    (pos) => {
      location.value = {
        latitude: pos.coords.latitude,
        longitude: pos.coords.longitude,
      }
      locating.value = false
    },
    () => {
      locationError.value = '위치 정보를 가져올 수 없습니다. 위치 권한을 허용해주세요.'
      locating.value = false
    }
  )
}

// 장소 등록
async function submit() {
  submitting.value = true
  submitError.value = ''

  const res = await fetch(`${API}/admin/places`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'X-Admin-Secret': ADMIN_SECRET,
    },
    body: JSON.stringify({
      name: form.value.name,
      description: form.value.description,
      latitude: location.value.latitude,
      longitude: location.value.longitude,
      radius_meters: form.value.radius_meters,
    }),
  })

  if (res.ok) {
    router.push('/')
  } else {
    submitError.value = '등록에 실패했습니다. 다시 시도해주세요.'
  }
  submitting.value = false
}
</script>

<style scoped>
.add-place {
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
.location-box {
  background: #f0f4ff;
  border-radius: 12px;
  padding: 20px;
  margin: 20px 0;
  text-align: center;
}
.location-box button {
  background: #1A4FA0;
  color: white;
  border: none;
  border-radius: 8px;
  padding: 12px 24px;
  font-size: 16px;
  cursor: pointer;
}
.location-result {
  margin-top: 12px;
  font-size: 13px;
  color: #333;
}
.hint {
  color: #0F6E56;
  font-size: 12px;
  margin-top: 4px;
}
.form {
  margin-top: 20px;
}
.field {
  margin-bottom: 16px;
}
label {
  display: block;
  font-size: 14px;
  font-weight: bold;
  margin-bottom: 6px;
  color: #333;
}
input, textarea, select {
  width: 100%;
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 10px 12px;
  font-size: 14px;
  box-sizing: border-box;
}
textarea {
  height: 80px;
  resize: none;
}
.submit {
  width: 100%;
  background: #1A4FA0;
  color: white;
  border: none;
  border-radius: 8px;
  padding: 14px;
  font-size: 16px;
  cursor: pointer;
  margin-top: 8px;
}
.submit:disabled {
  background: #aaa;
}
.error {
  color: red;
  font-size: 13px;
  margin-top: 8px;
}
</style>