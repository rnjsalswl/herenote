<template>
  <div style="min-height:100vh;background:var(--bg)">
    <HnTopBar title="장소 추가" :on-back="goBack"/>

    <div style="padding:8px 16px 48px">
      <!-- STEP 1 -->
      <div style="padding:16px 0 8px">
        <div style="display:flex;align-items:center;gap:10px;margin-bottom:14px">
          <div :style="{
            width:'26px', height:'26px', borderRadius:'50%', flexShrink:0,
            background: location ? 'var(--ink)' : 'var(--surface-2)',
            color: location ? 'var(--inv-ink)' : 'var(--ink-2)',
            border:'1px solid var(--line)',
            fontSize:'12px', fontWeight:700, display:'flex', alignItems:'center', justifyContent:'center',
          }">
            <HnIcon v-if="location" name="check" :size="14" :stroke="3"/>
            <span v-else>1</span>
          </div>
          <span style="font-size:15px;font-weight:700;color:var(--ink)">현재 위치 가져오기</span>
        </div>

        <div style="padding-left:36px">
          <div v-if="!location" style="display:flex;flex-direction:column;gap:10px">
            <p style="font-size:13px;color:var(--ink-2)">버튼을 눌러 현재 위치를 장소 좌표로 등록하세요</p>
            <button @click="getLocation" :disabled="locating" :style="{
              alignSelf:'flex-start', padding:'10px 18px', borderRadius:'12px', border:'none',
              background:'var(--ink)', color:'var(--inv-ink)', fontSize:'14px', fontWeight:600, cursor:'pointer',
              opacity: locating ? 0.5 : 1, transition:'opacity .15s',
            }">
              <span style="display:flex;align-items:center;gap:7px">
                <HnIcon name="pin" :size="16" :stroke="2.2"/>
                {{ locating ? '가져오는 중...' : '현재 위치 가져오기' }}
              </span>
            </button>
            <p v-if="locationError" style="font-size:13px;color:var(--ink-2)">{{ locationError }}</p>
          </div>

          <div v-else style="display:flex;flex-direction:column;gap:8px">
            <!-- 위치 미리보기 지도 -->
            <div style="border-radius:12px;overflow:hidden;border:1px solid var(--line)">
              <HnKakaoMap
                :lat="location.latitude"
                :lng="location.longitude"
                :zoom="3"
                :height="160"
                :places="[{ id:'me', latitude: location.latitude, longitude: location.longitude }]"
              />
            </div>
            <div style="display:flex;justify-content:space-between;font-size:13px;padding:10px 12px;background:var(--surface-2);border-radius:10px;border:1px solid var(--line)">
              <span style="color:var(--ink-2)">위도</span>
              <span style="font-weight:600;font-family:var(--font-mono)">{{ location.latitude.toFixed(6) }}</span>
            </div>
            <div style="display:flex;justify-content:space-between;font-size:13px;padding:10px 12px;background:var(--surface-2);border-radius:10px;border:1px solid var(--line)">
              <span style="color:var(--ink-2)">경도</span>
              <span style="font-weight:600;font-family:var(--font-mono)">{{ location.longitude.toFixed(6) }}</span>
            </div>
            <button @click="getLocation" :disabled="locating" style="align-self:flex-start;padding:6px 14px;border-radius:8px;border:1px solid var(--line);background:transparent;color:var(--ink-2);font-size:13px;cursor:pointer">다시 가져오기</button>
          </div>
        </div>
      </div>

      <div style="height:1px;background:var(--line);margin:8px 0"/>

      <!-- STEP 2 -->
      <div :style="{padding:'16px 0 8px', opacity: location ? 1 : 0.4, pointerEvents: location ? 'auto' : 'none'}">
        <div style="display:flex;align-items:center;gap:10px;margin-bottom:14px">
          <div style="width:26px;height:26px;border-radius:50%;flex-shrink:0;background:var(--surface-2);color:var(--ink-2);border:1px solid var(--line);font-size:12px;font-weight:700;display:flex;align-items:center;justify-content:center">2</div>
          <span style="font-size:15px;font-weight:700;color:var(--ink)">장소 정보 입력</span>
        </div>

        <div v-if="location" style="padding-left:36px;display:flex;flex-direction:column;gap:16px">
          <div>
            <label style="display:block;font-size:13px;font-weight:600;color:var(--ink-2);margin-bottom:6px">장소 이름 <span style="color:var(--ink)">*</span></label>
            <input v-model="form.name" :style="fieldStyle" placeholder="예: 경복궁 정문"/>
          </div>

          <div>
            <label style="display:block;font-size:13px;font-weight:600;color:var(--ink-2);margin-bottom:6px">장소 설명</label>
            <textarea v-model="form.description" :style="{...fieldStyle, height:'80px', resize:'none', lineHeight:'1.5'}" placeholder="이 장소에 대해 설명해주세요"/>
          </div>

          <div>
            <label style="display:block;font-size:13px;font-weight:600;color:var(--ink-2);margin-bottom:8px">인증 반경</label>
            <div style="display:grid;grid-template-columns:repeat(4,1fr);gap:8px">
              <label v-for="opt in radiusOptions" :key="opt.value" :style="{
                border: form.radius_meters===opt.value ? '2px solid var(--ink)' : '1.5px solid var(--line)',
                borderRadius:'12px', padding:'10px 6px', textAlign:'center', cursor:'pointer',
                display:'flex', flexDirection:'column', gap:'2px',
                background: form.radius_meters===opt.value ? 'var(--surface-2)' : 'transparent',
                transition:'all .15s',
              }">
                <input type="radio" :value="opt.value" v-model="form.radius_meters" style="display:none"/>
                <span style="font-size:13px;font-weight:700;color:var(--ink)">{{ opt.label }}</span>
                <span style="font-size:11px;color:var(--ink-3)">{{ opt.desc }}</span>
              </label>
            </div>
          </div>

          <button @click="submit" :disabled="submitting || !form.name.trim()" :style="{
            padding:'14px', borderRadius:'14px', border:'none',
            background:'var(--ink)', color:'var(--inv-ink)', fontSize:'15px', fontWeight:700,
            cursor: form.name.trim() ? 'pointer' : 'not-allowed',
            opacity: form.name.trim() ? 1 : 0.3, transition:'opacity .15s',
          }">{{ submitting ? '등록 중...' : '장소 등록하기' }}</button>
          <p v-if="submitError" style="font-size:13px;color:var(--ink-2);text-align:center">{{ submitError }}</p>
        </div>

        <div v-else style="padding-left:36px;font-size:13px;color:var(--ink-3)">위치를 먼저 가져와주세요</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import HnTopBar from '@/components/HnTopBar.vue'
import HnIcon from '@/components/HnIcon.vue'
import HnKakaoMap from '@/components/HnKakaoMap.vue'
import { API } from '@/config.js'

const router = useRouter()
const ADMIN_SECRET = 'herenote-admin'

function goBack() { router.push('/') }

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

const fieldStyle = {
  width:'100%', boxSizing:'border-box', padding:'12px 14px', borderRadius:'12px',
  border:'1px solid var(--line)', background:'var(--surface-2)', color:'var(--ink)',
  fontSize:'14px', fontFamily:'var(--font-sans)', outline:'none',
}

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
