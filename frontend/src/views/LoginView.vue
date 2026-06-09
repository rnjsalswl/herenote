<template>
  <div class="login-page">
    <div class="brand-area">
      <div class="pin-icon">📍</div>
      <h1 class="brand-name">HERE : NOTE</h1>
      <p class="brand-sub">실제 그 장소에서만 열리는 디지털 방명록</p>
    </div>

    <div class="login-card">
      <!-- 탭 -->
      <div class="tabs">
        <button class="tab" :class="{ active: mode === 'register' }" @click="switchMode('register')">회원가입</button>
        <button class="tab" :class="{ active: mode === 'login' }" @click="switchMode('login')">로그인</button>
      </div>

      <div class="form">
        <div class="field">
          <label class="field-label">닉네임</label>
          <input
            v-model="nickname"
            class="input"
            placeholder="닉네임 입력 (최대 20자)"
            maxlength="20"
            @keyup.enter="submit"
          />
        </div>
        <div class="field">
          <label class="field-label">비밀번호</label>
          <input
            v-model="password"
            class="input"
            type="password"
            :placeholder="mode === 'register' ? '비밀번호 (4자 이상)' : '비밀번호 입력'"
            @keyup.enter="submit"
          />
        </div>

        <button class="btn-primary" @click="submit" :disabled="submitting || !canSubmit">
          {{ submitting ? '처리 중...' : mode === 'register' ? '가입하기 →' : '로그인 →' }}
        </button>
        <p v-if="error" class="error-msg">{{ error }}</p>
      </div>
    </div>

    <p class="footer-note">장소에 직접 방문해야 방명록을 볼 수 있어요</p>
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

const canSubmit = computed(() =>
  nickname.value.trim().length > 0 && password.value.length >= 4
)

function switchMode(m) {
  mode.value = m
  error.value = ''
}

async function submit() {
  if (!canSubmit.value) return
  submitting.value = true
  error.value = ''

  const url = mode.value === 'register' ? `${API}/users` : `${API}/auth/login`
  const res = await fetch(url, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      nickname: nickname.value.trim(),
      password: password.value,
    }),
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
.login-page {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 32px 20px;
  background: linear-gradient(160deg, #1A4FA0 0%, #2563c7 40%, var(--bg) 65%);
}

.brand-area {
  text-align: center;
  margin-bottom: 40px;
}

.pin-icon {
  font-size: 48px;
  margin-bottom: 12px;
  filter: drop-shadow(0 4px 8px rgba(0,0,0,0.15));
}

.brand-name {
  font-size: 2.4rem;
  font-weight: 700;
  color: #fff;
  letter-spacing: 0.04em;
  margin-bottom: 10px;
}

.brand-sub {
  color: rgba(255, 255, 255, 0.8);
  font-size: 14px;
}

.login-card {
  width: 100%;
  max-width: 380px;
  background: var(--surface);
  border-radius: var(--radius-lg);
  overflow: hidden;
  box-shadow: var(--shadow-lg);
}

/* 탭 */
.tabs {
  display: flex;
  border-bottom: 1px solid var(--border);
}

.tab {
  flex: 1;
  background: none;
  border: none;
  padding: 14px;
  font-size: 14px;
  font-weight: 600;
  color: var(--text-hint);
  cursor: pointer;
  transition: color 0.15s, border-bottom 0.15s;
  border-bottom: 2px solid transparent;
}

.tab.active {
  color: var(--primary);
  border-bottom-color: var(--primary);
}

/* 폼 */
.form {
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.field-label {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-secondary);
}

.input {
  width: 100%;
  border: 1.5px solid var(--border);
  border-radius: var(--radius-sm);
  padding: 12px 14px;
  font-size: 15px;
  background: var(--bg);
  color: var(--text-primary);
  outline: none;
  transition: border-color 0.15s, box-shadow 0.15s;
}

.input:focus {
  border-color: var(--primary);
  box-shadow: 0 0 0 3px var(--primary-light);
  background: var(--surface);
}

.input::placeholder {
  color: var(--text-hint);
}

.btn-primary {
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

.btn-primary:hover:not(:disabled) {
  background: var(--primary-dark);
}

.btn-primary:disabled {
  background: var(--text-hint);
  cursor: not-allowed;
}

.error-msg {
  color: var(--error);
  font-size: 13px;
  text-align: center;
  background: var(--error-light);
  border-radius: var(--radius-sm);
  padding: 8px 12px;
}

.footer-note {
  margin-top: 28px;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.6);
  text-align: center;
}
</style>
