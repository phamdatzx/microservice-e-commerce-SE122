import { ref } from 'vue'

const bus = ref(new Map<string, Function[]>())

export const eventBus = {
  on: (event: string, callback: Function) => {
    if (!bus.value.has(event)) {
      bus.value.set(event, [])
    }
    bus.value.get(event)?.push(callback)
  },
  off: (event: string, callback: Function) => {
    if (bus.value.has(event)) {
      const callbacks = bus.value.get(event)
      if (callbacks) {
        const index = callbacks.indexOf(callback)
        if (index > -1) {
          callbacks.splice(index, 1)
        }
      }
    }
  },
  emit: (event: string, ...args: any[]) => {
    if (bus.value.has(event)) {
      bus.value.get(event)?.forEach((callback) => callback(...args))
    }
  },
}
