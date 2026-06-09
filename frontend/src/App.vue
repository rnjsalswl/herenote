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
  --bg: #FFFFFF;
  --surface: #FFFFFF;
  --separator: #EFEFEF;
  --border: #DBDBDB;

  --text-primary: #000000;
  --text-secondary: #737373;
  --text-hint: #AAAAAA;

  --primary: #000000;
  --primary-hover: #1A1A1A;
  --primary-light: #F5F5F5;

  --success: #29A86C;
  --success-light: #EDF7F2;
  --error: #FF3040;
  --error-light: #FFF0F1;
}

html, body { height: 100%; }

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', sans-serif;
  background: var(--bg);
  color: var(--text-primary);
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  font-size: 15px;
  line-height: 1.5;
}

#app { min-height: 100vh; }

button, input, textarea, select { font-family: inherit; }
button { cursor: pointer; }
</style>
