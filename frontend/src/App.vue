<template>
  <RouterView />
</template>

<style>
*, *::before, *::after {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

:root {
  --primary: #1A4FA0;
  --primary-dark: #153D7A;
  --primary-light: #EEF3FF;
  --success: #0B6E4F;
  --success-light: #E8F7F2;
  --error: #DC2626;
  --error-light: #FEF2F2;

  --bg: #F4F7FC;
  --surface: #FFFFFF;
  --border: #E2E8F0;

  --text-primary: #111827;
  --text-secondary: #4B5563;
  --text-hint: #9CA3AF;

  --shadow-sm: 0 1px 4px rgba(0, 0, 0, 0.07);
  --shadow-md: 0 4px 20px rgba(0, 0, 0, 0.10);
  --shadow-lg: 0 8px 32px rgba(0, 0, 0, 0.12);

  --radius-sm: 8px;
  --radius: 14px;
  --radius-lg: 20px;
}

html, body {
  height: 100%;
}

body {
  font-family: 'Noto Sans KR', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
  background: var(--bg);
  color: var(--text-primary);
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  font-size: 15px;
  line-height: 1.6;
}

#app {
  min-height: 100vh;
}

button {
  font-family: inherit;
  cursor: pointer;
  transition: opacity 0.15s, transform 0.1s;
}

button:active:not(:disabled) {
  transform: scale(0.97);
}

input, textarea, select {
  font-family: inherit;
}
</style>
