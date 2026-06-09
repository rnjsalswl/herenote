<template>
  <div class="wrap">
    <div class="logo">
      <span class="logo-pin">📍</span>
      <h1 class="logo-name">HERE NOTE</h1>
      <p class="logo-sub">실제 그 장소에서만 열리는 디지털 방명록</p>
    </div>

    <div class="form">
      <div class="tabs">
        <button :class="['tab', { active: mode === 'register' }]" @click="switchMode('register')">처음 오셨나요</button>
        <button :class="['tab', { active: mode === 'login' }]" @click="switchMode('login')">돌아오셨나요</button>
      </div>

      <input v-model="nickname" class="field" placeholder="닉네임" maxlength="20" @keyup.enter="submit" />
      <input v-model="password" class="field" type="password" :placeholder="mode === 'register' ? '비밀번호 (4자 이상)' : '비밀번호'" @keyup.enter="submit" />

      <button class="btn-submit" @click="submit" :disabled="submitting || !canSubmit">
        {{ submitting ? '처리 중...' : mode === 'register' ? '가입하기' : '로그인' }}
      </button>
      <p v-if="error" class="err">{{ error }}</p>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { API } from '@/config.js'
import { setUserID } from '@/stores/user.js'

const router = useRouter()
const mode = ref('register')
const nickname = ref('')
const password = ref('')
const submitting = ref(false)
const error = ref('')

const canSubmit = computed(() => nickname.value.trim().length > 0 && password.value.length >= 4)

function switchMode(m) { mode.value = m; error.value = '' }

async function submit() {
  if (!canSubmit.value) return
  submitting.value = true
  error.value = ''
  const url = mode.value === 'register' ? `${API}/users` : `${API}/auth/login`
  const res = await fetch(url, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ nickname: nickname.value.trim(), password: password.value }),
  })
  if (res.ok) {
    const data = await res.json()
    setUserID(mode.value === 'login' ? data.user_id : data.id)
    router.push('/')
  } else {
    const data = await res.json().catch(() => ({}))
    error.value = data.error || '오류가 발생했습니다. 다시 시도해주세요.'
  }
  submitting.value = false
}
</script>

<style scoped>
.wrap {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 28px;
  background: #fff;
}

.logo {
  text-align: center;
  margin-bottom: 52px;
}

.logo-pin { font-size: 48px; display: block; margin-bottom: 14px; }

.logo-name {
  font-size: 30px;
  font-weight: 700;
  letter-spacing: -1px;
  margin-bottom: 8px;
}

.logo-sub {
  font-size: 14px;
  color: var(--text-secondary);
  line-height: 1.5;
}

.form {
  width: 100%;
  max-width: 360px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.tabs {
  display: flex;
  margin-bottom: 6px;
  border-bottom: 1px solid var(--separator);
}

.tab {
  flex: 1;
  background: none;
  border: none;
  padding: 11px 0;
  font-size: 14px;
  font-weight: 500;
  color: var(--text-hint);
  border-bottom: 2px solid transparent;
  margin-bottom: -1px;
  transition: color 0.15s, border-color 0.15s;
}

.tab.active {
  color: var(--text-primary);
  border-bottom-color: var(--text-primary);
}

.field {
  width: 100%;
  background: var(--primary-light);
  border: none;
  border-radius: 12px;
  padding: 14px 16px;
  font-size: 15px;
  color: var(--text-primary);
  outline: none;
  transition: background 0.15s;
}

.field:focus { background: #EBEBEB; }
.field::placeholder { color: var(--text-hint); }

.btn-submit {
  width: 100%;
  background: var(--primary);
  color: #fff;
  border: none;
  border-radius: 12px;
  padding: 14px;
  font-size: 15px;
  font-weight: 600;
  margin-top: 4px;
  transition: opacity 0.15s;
}

.btn-submit:disabled { opacity: 0.3; cursor: not-allowed; }
.btn-submit:not(:disabled):active { opacity: 0.7; }

.err {
  font-size: 13px;
  color: var(--error);
  text-align: center;
  padding: 4px 0;
}
</style>
