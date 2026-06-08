<template>
  <div class="add-page">
    <!-- 헤더 -->
    <header class="page-header">
      <div class="header-inner">
        <button class="btn-back" @click="router.push('/')">
          <span class="back-arrow">←</span> 뒤로
        </button>
        <span class="header-title">장소 추가</span>
        <div style="width: 52px;"></div>
      </div>
    </header>

    <div class="content">
      <!-- STEP 1: 위치 -->
      <div class="step">
        <div class="step-label">
          <span class="step-number" :class="{ done: !!location }">
            {{ location ? '✓' : '1' }}
          </span>
          <span class="step-title">현재 위치 가져오기</span>
        </div>

        <div class="location-card" :class="{ 'has-location': !!location }">
          <div v-if="!location">
            <p class="location-hint">장소 좌표로 현재 위치가 등록돼요</p>
            <button class="btn-location" @click="getLocation" :disabled="locating">
              {{ locating ? '위치 가져오는 중...' : '📍 현재 위치 가져오기' }}
            </button>
            <p v-if="locationError" class="error-msg">{{ locationError }}</p>
          </div>
          <div v-else class="location-result">
            <div class="location-row">
              <span class="location-label">위도</span>
              <span class="location-value">{{ location.latitude.toFixed(6) }}</span>
            </div>
            <div class="location-row">
              <span class="location-label">경도</span>
              <span class="location-value">{{ location.longitude.toFixed(6) }}</span>
            </div>
            <p class="location-ok">✓ 위치 정보가 준비됐어요</p>
            <button class="btn-relocate" @click="getLocation" :disabled="locating">
              다시 가져오기
            </button>
          </div>
        </div>
      </div>

      <!-- STEP 2: 장소 정보 -->
      <div class="step" :class="{ disabled: !location }">
        <div class="step-label">
          <span class="step-number">2</span>
          <span class="step-title">장소 정보 입력</span>
        </div>

        <div v-if="location" class="form-card">
          <div class="field">
            <label class="field-label">장소 이름 <span class="required">*</span></label>
            <input
              v-model="form.name"
              class="input"
              placeholder="예: 경복궁 정문"
            />
          </div>

          <div class="field">
            <label class="field-label">장소 설명</label>
            <textarea
              v-model="form.description"
              class="textarea"
              placeholder="이 장소에 대해 설명해주세요"
            />
          </div>

          <div class="field">
            <label class="field-label">인증 반경</label>
            <div class="radius-options">
              <label
                v-for="opt in radiusOptions"
                :key="opt.value"
                class="radius-option"
                :class="{ selected: form.radius_meters === opt.value }"
              >
                <input
                  type="radio"
                  :value="opt.value"
                  v-model="form.radius_meters"
                  style="display: none;"
                />
                <span class="radius-name">{{ opt.label }}</span>
                <span class="radius-desc">{{ opt.desc }}</span>
              </label>
            </div>
          </div>

          <button
            class="btn-submit"
            @click="submit"
            :disabled="submitting || !form.name.trim()"
          >
            {{ submitting ? '등록 중...' : '장소 등록하기' }}
          </button>
          <p v-if="submitError" class="error-msg">{{ submitError }}</p>
        </div>

        <div v-else class="step-placeholder">
          위치를 먼저 가져오면 장소 정보를 입력할 수 있어요
        </div>
      </div>
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

const radiusOptions = [
  { value: 50,  label: '50m',  desc: '실내' },
  { value: 100, label: '100m', desc: '광장' },
  { value: 200, label: '200m', desc: '공원' },
  { value: 500, label: '500m', desc: '산 정상' },
]

const form = ref({
  name: '',
  description: '',
  radius_meters: 100,
})

const submitting = ref(false)
const submitError = ref('')

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
.add-page {
  min-height: 100vh;
  background: var(--bg);
}

/* 헤더 */
.page-header {
  background: var(--surface);
  border-bottom: 1px solid var(--border);
  position: sticky;
  top: 0;
  z-index: 10;
  box-shadow: var(--shadow-sm);
}

.header-inner {
  max-width: 640px;
  margin: 0 auto;
  padding: 12px 20px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.btn-back {
  background: none;
  border: none;
  font-size: 14px;
  color: var(--primary);
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 0;
}

.back-arrow {
  font-size: 18px;
}

.header-title {
  font-size: 16px;
  font-weight: 700;
  color: var(--text-primary);
}

/* 본문 */
.content {
  max-width: 640px;
  margin: 0 auto;
  padding: 24px 20px 48px;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

/* 스텝 */
.step {
  transition: opacity 0.2s;
}

.step.disabled {
  opacity: 0.5;
  pointer-events: none;
}

.step-label {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 12px;
}

.step-number {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: var(--primary);
  color: #fff;
  font-size: 13px;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: background 0.2s;
}

.step-number.done {
  background: var(--success);
}

.step-title {
  font-size: 16px;
  font-weight: 700;
  color: var(--text-primary);
}

/* 위치 카드 */
.location-card {
  background: var(--surface);
  border-radius: var(--radius);
  padding: 20px;
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--border);
  text-align: center;
}

.location-card.has-location {
  border-color: var(--success);
  background: var(--success-light);
  text-align: left;
}

.location-hint {
  font-size: 13px;
  color: var(--text-hint);
  margin-bottom: 14px;
}

.btn-location {
  background: var(--primary);
  color: #fff;
  border: none;
  border-radius: var(--radius-sm);
  padding: 13px 24px;
  font-size: 15px;
  font-weight: 600;
  width: 100%;
}

.btn-location:hover:not(:disabled) {
  background: var(--primary-dark);
}

.btn-location:disabled {
  background: var(--text-hint);
  cursor: not-allowed;
}

.location-result {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.location-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 13px;
}

.location-label {
  color: var(--text-hint);
}

.location-value {
  font-weight: 600;
  color: var(--text-primary);
  font-variant-numeric: tabular-nums;
}

.location-ok {
  font-size: 13px;
  font-weight: 600;
  color: var(--success);
  margin-top: 4px;
}

.btn-relocate {
  background: none;
  border: 1px solid var(--success);
  color: var(--success);
  border-radius: var(--radius-sm);
  padding: 7px 14px;
  font-size: 13px;
  align-self: flex-start;
  margin-top: 4px;
}

/* 폼 카드 */
.form-card {
  background: var(--surface);
  border-radius: var(--radius);
  padding: 20px;
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--border);
}

.field {
  margin-bottom: 18px;
}

.field-label {
  display: block;
  font-size: 13px;
  font-weight: 600;
  color: var(--text-secondary);
  margin-bottom: 7px;
}

.required {
  color: var(--error);
}

.input, .textarea {
  width: 100%;
  border: 1.5px solid var(--border);
  border-radius: var(--radius-sm);
  padding: 11px 14px;
  font-size: 14px;
  background: var(--bg);
  color: var(--text-primary);
  outline: none;
  transition: border-color 0.15s, box-shadow 0.15s;
}

.input:focus, .textarea:focus {
  border-color: var(--primary);
  box-shadow: 0 0 0 3px var(--primary-light);
  background: var(--surface);
}

.input::placeholder, .textarea::placeholder {
  color: var(--text-hint);
}

.textarea {
  height: 90px;
  resize: none;
  line-height: 1.6;
}

/* 반경 옵션 */
.radius-options {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 8px;
}

.radius-option {
  border: 1.5px solid var(--border);
  border-radius: var(--radius-sm);
  padding: 10px 6px;
  text-align: center;
  cursor: pointer;
  transition: border-color 0.15s, background 0.15s;
  display: flex;
  flex-direction: column;
  gap: 3px;
}

.radius-option.selected {
  border-color: var(--primary);
  background: var(--primary-light);
}

.radius-name {
  font-size: 13px;
  font-weight: 700;
  color: var(--text-primary);
}

.radius-option.selected .radius-name {
  color: var(--primary);
}

.radius-desc {
  font-size: 11px;
  color: var(--text-hint);
}

/* 등록 버튼 */
.btn-submit {
  width: 100%;
  background: var(--primary);
  color: #fff;
  border: none;
  border-radius: var(--radius-sm);
  padding: 14px;
  font-size: 15px;
  font-weight: 600;
  margin-top: 4px;
}

.btn-submit:hover:not(:disabled) {
  background: var(--primary-dark);
}

.btn-submit:disabled {
  background: var(--text-hint);
  cursor: not-allowed;
}

/* 에러 */
.error-msg {
  color: var(--error);
  font-size: 13px;
  margin-top: 10px;
  background: var(--error-light);
  border-radius: var(--radius-sm);
  padding: 8px 12px;
  text-align: center;
}

.step-placeholder {
  background: var(--surface);
  border-radius: var(--radius);
  padding: 20px;
  border: 1.5px dashed var(--border);
  font-size: 13px;
  color: var(--text-hint);
  text-align: center;
}
</style>
