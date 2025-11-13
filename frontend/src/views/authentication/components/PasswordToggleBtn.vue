<template>
  <button ref="passwordToggle" type="button" class="password-toggle">
    <svg width="20" height="20" viewBox="0 0 20 20" fill="none" class="eye-icon" ref="eyeIcon">
      <path
        d="M1 10s4-8 9-8 9 8 9 8-4 8-9 8-9-8-9-8z"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
      />
      <circle
        cx="10"
        cy="10"
        r="3"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
      />
    </svg>
  </button>
</template>

<style scoped>
.password-toggle {
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  cursor: pointer;
  color: #9ca3af;
  padding: 4px;
}

.password-toggle:hover {
  color: #666;
}
</style>
