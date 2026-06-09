<template>
  <div style="min-height:100%;display:flex;flex-direction:column;padding:0 28px;background:var(--bg)">
    <!-- hero -->
    <div style="flex:0 0 auto;padding-top:96px;text-align:center">
      <div style="width:76px;height:76px;border-radius:24px;margin:0 auto 22px;background:var(--ink);color:var(--inv-ink);display:flex;align-items:center;justify-content:center;box-shadow:var(--shadow)">
        <HnIcon name="pinFill" :size="40"/>
      </div>
      <h1 style="font-family:var(--font-display);font-weight:var(--disp-weight);font-size:30px;letter-spacing:0.16em;margin:0;color:var(--ink);font-variant:all-small-caps">here note</h1>
      <p style="margin:12px auto 0;max-width:240px;font-size:14.5px;line-height:1.6;color:var(--ink-2);word-break:keep-all">실제 그 장소에 가야만 열리는<br>디지털 방명록</p>
    </div>
    <!-- form -->
    <div style="flex:1;display:flex;flex-direction:column;justify-content:center;max-width:360px;width:100%;margin:0 auto;padding-bottom:40px">
      <div style="display:flex;background:var(--surface-2);border-radius:14px;padding:4px;margin-bottom:16px">
        <button v-for="[m,l] in modes" :key="m" @click="mode=m" :style="{
          flex:1, border:'none', cursor:'pointer', padding:'10px 0', borderRadius:'10px',
          fontSize:'14px', fontWeight:700, fontFamily:'var(--font-sans)',
          background: mode===m ? 'var(--card)' : 'transparent',
          color: mode===m ? 'var(--ink)' : 'var(--ink-3)',
          boxShadow: mode===m ? 'var(--shadow)' : 'none', transition:'all .15s',
        }">{{ l }}</button>
      </div>
      <input v-model="nickname" placeholder="닉네임" :style="inp" @keyup.enter="submit"/>
      <input v-model="password" type="password" :placeholder="mode==='register' ? '비밀번호 (4자 이상)' : '비밀번호'" :style="{...inp, marginTop:'10px'}" @keyup.enter="submit"/>
      <button :disabled="!canSubmit" @click="submit" :style="{
        marginTop:'16px', padding:'15px', borderRadius:'14px', border:'none',
        background:'var(--ink)', color:'var(--inv-ink)', fontSize:'15.5px', fontWeight:700,
        fontFamily:'var(--font-sans)', cursor: canSubmit ? 'pointer' : 'not-allowed',
        opacity: canSubmit ? 1 : 0.32, transition:'opacity .15s',
      }">{{ submitting ? '처리 중...' : mode==='register' ? '시작하기' : '로그인' }}</button>
      <p v-if="error" style="text-align:center;margin-top:10px;font-size:13px;color:var(--ink-2)">{{ error }}</p>
      <p style="text-align:center;margin-top:18px;font-size:12px;color:var(--ink-3);font-family:var(--font-mono)">위치 권한이 필요해요 · GPS로 방문을 인증합니다</p>
    </div>
  </div>
</template>
<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import HnIcon from '@/components/HnIcon.vue'
import { API } from '@/config.js'
import { setToken } from '@/stores/user.js'
const router = useRouter()
const mode = ref('register')
const modes = [['register','처음 왔어요'],['login','돌아왔어요']]
const nickname = ref(''); const password = ref(''); const submitting = ref(false); const error = ref('')
const canSubmit = computed(() => nickname.value.trim() && password.value.length >= 4)
const inp = { width:'100%', boxSizing:'border-box', padding:'15px 16px', borderRadius:'14px', border:'1px solid var(--line)', background:'var(--surface-2)', color:'var(--ink)', fontSize:'15px', fontFamily:'var(--font-sans)', outline:'none' }
async function submit() {
  if (!canSubmit.value || submitting.value) return
  submitting.value = true; error.value = ''
  const url = mode.value === 'register' ? `${API}/users` : `${API}/auth/login`
  const res = await fetch(url, { method:'POST', headers:{'Content-Type':'application/json'}, body:JSON.stringify({ nickname: nickname.value.trim(), password: password.value }) })
  if (res.ok) {
    const data = await res.json()
    if (data.token) setToken(data.token)
    router.push('/')
  } else {
    const data = await res.json().catch(() => ({}))
    error.value = data.error || '오류가 발생했습니다.'
  }
  submitting.value = false
}
</script>
