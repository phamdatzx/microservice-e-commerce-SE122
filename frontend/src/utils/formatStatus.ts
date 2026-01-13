export const formatStatus = (status: string): string => {
  const map: Record<string, string> = {
    TO_CONFIRM: 'To Confirm',
    TO_PAY: 'To Pay',
    TO_PICKUP: 'To Pickup',
    SHIPPING: 'Shipping',
    COMPLETED: 'Completed',
    CANCELLED: 'Cancelled',
    RETURNED: 'Returned',
  }
  return map[status] || status
}
