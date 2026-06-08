<template>
  <div class="login">
    <h1>HERE : NOTE</h1>
    <p class="sub">실제 그 장소에서만 열리는 디지털 방명록</p>

    <div class="form">
      <p>닉네임을 입력해주세요</p>
      <input
        v-model="nickname"
        placeholder="닉네임"
        maxlength="20"
        @keyup.enter="submit"
      />
      <button @click="submit" :disabled="submitting || !nickname.trim()">
        {{ submitting ? '시작하는 중...' : '시작하기' }}
      </button>
      <p v-if="error" class="error">{{ error }}</p>
    </div>
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
    setUserID(data.id)         // localStorage에 저장
    router.push('/')           // 홈으로 이동
  } else {
    error.value = '오류가 발생했습니다. 다시 시도해주세요.'
  }
  submitting.value = false
}
</script>

<style scoped>
.login {
  max-width: 400px;
  margin: 80px auto;
  padding: 20px;
  text-align: center;
}
h1 {
  font-size: 2.5rem;
  font-weight: bold;
  color: #1A4FA0;
  margin-bottom: 8px;
}
.sub {
  color: #888;
  font-size: 14px;
  margin-bottom: 40px;
}
.form p {
  font-size: 16px;
  margin-bottom: 12px;
}
input {
  width: 100%;
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 12px;
  font-size: 16px;
  box-sizing: border-box;
  margin-bottom: 12px;
  text-align: center;
}
button {
  width: 100%;
  background: #1A4FA0;
  color: white;
  border: none;
  border-radius: 8px;
  padding: 14px;
  font-size: 16px;
  cursor: pointer;
}
button:disabled {
  background: #aaa;
}
.error {
  color: red;
  margin-top: 8px;
}
</style>