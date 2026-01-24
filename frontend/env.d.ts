/// <reference types="vite/client" />

declare module '@splidejs/vue-splide' {
  import { DefineComponent } from 'vue'
  import { Options } from '@splidejs/splide'

  const Splide: DefineComponent<{
    options?: Options
    extensions?: Record<string, any>
    transition?: any
    hasTrack?: boolean
    tag?: string
  }>

  const SplideSlide: DefineComponent<{
    tag?: string
  }>

  const SplideTrack: DefineComponent<{
    tag?: string
  }>

  export { Splide, SplideSlide, SplideTrack, Options }
}
