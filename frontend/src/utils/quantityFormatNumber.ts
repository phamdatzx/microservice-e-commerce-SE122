export const quantityFormatNumber = (value: number): string => {
  if (value < 1000) return value.toString()

  const units = ['k', 'm', 'b', 't']
  let unitIndex = -1
  let divisor = 1

  while (value / (divisor * 1000) >= 1 && unitIndex < units.length - 1) {
    divisor *= 1000
    unitIndex++
  }

  const rawValue = value / divisor
  // Use floor with 1 decimal place precision
  const formattedValue = Math.floor(rawValue * 10) / 10
  const unit = units[unitIndex]!

  // Check if there is any remainder after formatting to 1 decimal place
  const hasRemainder = value > formattedValue * divisor
  const displayValue = formattedValue.toString().replace(/\.0$/, '')

  return displayValue + unit + (hasRemainder ? '+' : '')
}
