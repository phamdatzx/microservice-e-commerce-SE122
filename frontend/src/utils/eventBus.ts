import { ref } from 'vue'

const bus = ref(new Map())

export const eventBus = {
  on: (event: string, callback: any) => {
    bus.value.set(event, callback)
  },
  emit: (event: string, ...args: any[]) => {
    if (bus.value.has(event)) {
      bus.value.get(event)(...args)
    }
  },
}
