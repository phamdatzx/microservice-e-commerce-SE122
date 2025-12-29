export const quantityFormatNumber = (value: number): string => {
  if (value < 1000) return value.toString()

  const units = ['k', 'm', 'b', 't'] // thousand, million, billion, trillion
  let unitIndex = -1
  let formatted = value

  while (formatted >= 1000 && unitIndex < units.length - 1) {
    formatted /= 1000
    unitIndex++
  }

  // Floor it (no decimals), and add "+"
  return Math.floor(formatted) + units[unitIndex]! + '+'
}
