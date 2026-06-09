<template>
  <RouterView />
</template>

<style>
*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }

:root {
  --bg: #FFFFFF; --surface: #FFFFFF; --surface-2: #F4F4F5; --card: #FFFFFF;
  --ink: #0A0A0B; --ink-2: #6F6F76; --ink-3: #A8A8AE; --line: #ECECEE;
  --inv-ink: #FFFFFF; --shadow: 0 1px 2px rgba(0,0,0,0.04), 0 8px 24px rgba(0,0,0,0.05);
  --radius: 16px; --radius-lg: 24px;
  --font-sans: 'Pretendard', -apple-system, system-ui, sans-serif;
  --font-display: 'Pretendard', -apple-system, system-ui, sans-serif;
  --font-mono: 'Spline Sans Mono', ui-monospace, monospace;
  --disp-weight: 800; --disp-ls: -0.03em;
}

html, body { height: 100%; }

body {
  font-family: var(--font-sans);
  background: var(--bg); color: var(--ink);
  -webkit-font-smoothing: antialiased;
}

#app { min-height: 100vh; }

button, input, textarea, select { font-family: inherit; }
button { cursor: pointer; -webkit-tap-highlight-color: transparent; }

@keyframes hnPulse {
  0% { transform: scale(.5); opacity: .55; }
  100% { transform: scale(1); opacity: 0; }
}
@keyframes hnPop {
  0% { transform: scale(.4); opacity: 0; }
  60% { transform: scale(1.12); }
  100% { transform: scale(1); opacity: 1; }
}
@keyframes hnFade { from { opacity: 0; } to { opacity: 1; } }
@keyframes hnUp { from { transform: translateY(100%); } to { transform: translateY(0); } }

@media (prefers-reduced-motion: reduce) {
  *, *::before, *::after { animation-duration: 0.01ms !important; }
}
</style>
