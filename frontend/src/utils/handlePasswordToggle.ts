export const handlePasswordToggle = (e: Event) => {
  const passwordToggleIcon = e.target as HTMLElement
  const passwordInput = passwordToggleIcon
    .closest('.input-wrapper')!
    .querySelector('input[type="password"]')

  if (!passwordInput) return
  const type = passwordInput.getAttribute('type') === 'password' ? 'text' : 'password'
  passwordInput.setAttribute('type', type)

  if (type === 'text') {
    passwordToggleIcon.style.opacity = '0.6'
  } else {
    passwordToggleIcon.style.opacity = '1'
  }
}
