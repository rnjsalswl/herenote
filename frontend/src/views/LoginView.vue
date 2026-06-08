<template>
  <div class="login-page">
    <div class="brand-area">
      <div class="pin-icon">📍</div>
      <h1 class="brand-name">HERE : NOTE</h1>
      <p class="brand-sub">실제 그 장소에서만 열리는 디지털 방명록</p>
    </div>

    <div class="login-card">
      <p class="form-title">닉네임으로 시작하기</p>
      <input
        v-model="nickname"
        class="input"
        placeholder="닉네임 입력 (최대 20자)"
        maxlength="20"
        @keyup.enter="submit"
      />
      <button class="btn-primary" @click="submit" :disabled="submitting || !nickname.trim()">
        {{ submitting ? '시작하는 중...' : '시작하기 →' }}
      </button>
      <p v-if="error" class="error-msg">{{ error }}</p>
    </div>

    <p class="footer-note">장소에 직접 방문해야 방명록을 볼 수 있어요</p>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { API } from '@/config.js'
import { setUserID } from '@/stores/user.js'

const router = useRouter()
const nickname = ref('')
const submitting = ref(false)
const error = ref('')

async function submit() {
  if (!nickname.value.trim()) return
  submitting.value = true
  error.value = ''

  const res = await fetch(`${API}/users`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ nickname: nickname.value.trim() }),
  })

  if (res.ok) {
    const data = await res.json()
    setUserID(data.id)
    router.push('/')
  } else {
    error.value = '오류가 발생했습니다. 다시 시도해주세요.'
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
  font-weight: 400;
}

.login-card {
  width: 100%;
  max-width: 380px;
  background: var(--surface);
  border-radius: var(--radius-lg);
  padding: 28px 24px;
  box-shadow: var(--shadow-lg);
}

.form-title {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-secondary);
  margin-bottom: 14px;
  text-align: center;
}

.input {
  width: 100%;
  border: 1.5px solid var(--border);
  border-radius: var(--radius-sm);
  padding: 13px 16px;
  font-size: 15px;
  background: var(--bg);
  color: var(--text-primary);
  outline: none;
  transition: border-color 0.15s, box-shadow 0.15s;
  text-align: center;
  margin-bottom: 12px;
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
  letter-spacing: 0.02em;
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
  margin-top: 10px;
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
