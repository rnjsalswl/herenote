<template>
  <div class="page">
    <header class="topbar">
      <button class="back-btn" @click="router.push('/')">
        <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round">
          <polyline points="15 18 9 12 15 6"/>
        </svg>
      </button>
      <span class="topbar-title">장소 추가</span>
      <div style="width:36px" />
    </header>

    <div class="content">
      <!-- STEP 1 -->
      <div class="step">
        <div class="step-head">
          <div class="step-num" :class="{ done: !!location }">{{ location ? '✓' : '1' }}</div>
          <span class="step-label">현재 위치 가져오기</span>
        </div>

        <div class="step-body">
          <div v-if="!location" class="loc-empty">
            <p class="loc-hint">버튼을 눌러 현재 위치를 장소 좌표로 등록하세요</p>
            <button class="btn-loc" @click="getLocation" :disabled="locating">
              {{ locating ? '가져오는 중...' : '📍 현재 위치 가져오기' }}
            </button>
            <p v-if="locationError" class="err">{{ locationError }}</p>
          </div>

          <div v-else class="loc-result">
            <div class="loc-row">
              <span class="loc-key">위도</span>
              <span class="loc-val">{{ location.latitude.toFixed(6) }}</span>
            </div>
            <div class="loc-row">
              <span class="loc-key">경도</span>
              <span class="loc-val">{{ location.longitude.toFixed(6) }}</span>
            </div>
            <div class="loc-ok">✓ 위치 준비됨</div>
            <button class="btn-relo" @click="getLocation" :disabled="locating">다시 가져오기</button>
          </div>
        </div>
      </div>

      <div class="line-sep" />

      <!-- STEP 2 -->
      <div class="step" :class="{ disabled: !location }">
        <div class="step-head">
          <div class="step-num">2</div>
          <span class="step-label">장소 정보 입력</span>
        </div>

        <div class="step-body" v-if="location">
          <div class="field-group">
            <label class="field-label">장소 이름 <span class="req">*</span></label>
            <input v-model="form.name" class="field" placeholder="예: 경복궁 정문" />
          </div>

          <div class="field-group">
            <label class="field-label">장소 설명</label>
            <textarea v-model="form.description" class="field ta" placeholder="이 장소에 대해 설명해주세요" />
          </div>

          <div class="field-group">
            <label class="field-label">인증 반경</label>
            <div class="radius-grid">
              <label
                v-for="opt in radiusOptions"
                :key="opt.value"
                class="radius-opt"
                :class="{ selected: form.radius_meters === opt.value }"
              >
                <input type="radio" :value="opt.value" v-model="form.radius_meters" style="display:none" />
                <span class="r-label">{{ opt.label }}</span>
                <span class="r-desc">{{ opt.desc }}</span>
              </label>
            </div>
          </div>

          <button class="btn-submit" @click="submit" :disabled="submitting || !form.name.trim()">
            {{ submitting ? '등록 중...' : '장소 등록하기' }}
          </button>
          <p v-if="submitError" class="err">{{ submitError }}</p>
        </div>

        <div v-else class="step-placeholder">위치를 먼저 가져와주세요</div>
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

const form = ref({ name: '', description: '', radius_meters: 100 })
const submitting = ref(false)
const submitError = ref('')

function getLocation() {
  locating.value = true
  locationError.value = ''
  navigator.geolocation.getCurrentPosition(
    (pos) => {
      location.value = { latitude: pos.coords.latitude, longitude: pos.coords.longitude }
      locating.value = false
    },
    () => {
      locationError.value = '위치 권한을 허용해주세요.'
      locating.value = false
    }
  )
}

async function submit() {
  submitting.value = true
  submitError.value = ''
  const res = await fetch(`${API}/admin/places`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', 'X-Admin-Secret': ADMIN_SECRET },
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
.page {
  max-width: 600px;
  margin: 0 auto;
  min-height: 100vh;
  background: #fff;
}

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
}

.content { padding: 8px 0 48px; }

.line-sep {
  height: 1px;
  background: var(--separator);
  margin: 8px 0;
}

.step { padding: 16px 16px 8px; }
.step.disabled { opacity: 0.4; pointer-events: none; }

.step-head {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 14px;
}

.step-num {
  width: 26px;
  height: 26px;
  border-radius: 50%;
  background: var(--primary);
  color: #fff;
  font-size: 12px;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: background 0.2s;
}

.step-num.done { background: var(--success); }

.step-label { font-size: 15px; font-weight: 600; }

.step-body { padding-left: 36px; }

.loc-empty { display: flex; flex-direction: column; gap: 10px; }
.loc-hint { font-size: 13px; color: var(--text-secondary); }

.btn-loc {
  background: var(--primary);
  color: #fff;
  border: none;
  border-radius: 12px;
  padding: 11px 18px;
  font-size: 14px;
  font-weight: 600;
  align-self: flex-start;
  transition: opacity 0.15s;
}

.btn-loc:disabled { opacity: 0.3; cursor: not-allowed; }

.loc-result { display: flex; flex-direction: column; gap: 8px; }

.loc-row {
  display: flex;
  justify-content: space-between;
  font-size: 13px;
}

.loc-key { color: var(--text-secondary); }
.loc-val { font-weight: 600; font-variant-numeric: tabular-nums; }
.loc-ok { font-size: 13px; font-weight: 600; color: var(--success); }

.btn-relo {
  background: none;
  border: 1px solid var(--border);
  border-radius: 8px;
  padding: 6px 14px;
  font-size: 13px;
  color: var(--text-secondary);
  align-self: flex-start;
}

.field-group { margin-bottom: 18px; }

.field-label {
  display: block;
  font-size: 13px;
  font-weight: 600;
  color: var(--text-secondary);
  margin-bottom: 7px;
}

.req { color: var(--error); }

.field {
  width: 100%;
  background: var(--primary-light);
  border: none;
  border-radius: 12px;
  padding: 12px 14px;
  font-size: 14px;
  color: var(--text-primary);
  outline: none;
  transition: background 0.15s;
}

.field:focus { background: #EBEBEB; }
.field::placeholder { color: var(--text-hint); }
.ta { height: 80px; resize: none; line-height: 1.5; }

.radius-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 8px;
}

.radius-opt {
  border: 1.5px solid var(--separator);
  border-radius: 10px;
  padding: 10px 6px;
  text-align: center;
  cursor: pointer;
  display: flex;
  flex-direction: column;
  gap: 2px;
  transition: border-color 0.15s;
}

.radius-opt.selected {
  border-color: var(--primary);
  background: var(--primary-light);
}

.r-label { font-size: 13px; font-weight: 700; }
.r-desc { font-size: 11px; color: var(--text-hint); }

.btn-submit {
  width: 100%;
  background: var(--primary);
  color: #fff;
  border: none;
  border-radius: 12px;
  padding: 14px;
  font-size: 15px;
  font-weight: 600;
  transition: opacity 0.15s;
}

.btn-submit:disabled { opacity: 0.3; cursor: not-allowed; }

.step-placeholder {
  font-size: 13px;
  color: var(--text-hint);
  padding: 4px 0;
}

.err {
  font-size: 13px;
  color: var(--error);
  margin-top: 8px;
}
</style>
