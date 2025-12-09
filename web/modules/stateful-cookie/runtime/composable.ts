import type { CookieOptions } from 'nuxt/app'
import type { Ref, WatchCallback } from 'vue'
import type { StatefulCookieOptions } from './types'

import { useCookie, useState } from 'nuxt/app'
import { shallowReadonly, watch } from 'vue'

export function useStatefulCookie<T = string | null | undefined>(name: string, options: StatefulCookieOptions<T> & { readonly: true }): Readonly<Ref<T>>
export function useStatefulCookie<T = string | null | undefined>(name: string, options?: StatefulCookieOptions<T>): Ref<T>
export function useStatefulCookie<T = string | null | undefined>(name: string, options: StatefulCookieOptions<T> = {}) {
  const _opts = { ...options, readonly: false } as CookieOptions<T> & { readonly: false }

  const cookie = useCookie<T>(name, _opts)

  // Initialize state with current cookie value
  // On server, useCookie reads from request headers immediately
  // On client, useCookie reads from document.cookie
  const state = useState<T>(`cookies:${name}`, () => cookie.value)

  // Ensure state is synced with cookie value on initialization
  // This is critical for server-side where useState might initialize before cookie is fully read
  // or where the cookie value needs to be reflected in the state
  // Note: On server, if cookie isn't in request headers, cookie.value will be the default value
  if (cookie.value !== undefined && cookie.value !== state.value) {
    state.value = cookie.value
  }

  // Sync state -> cookie (when state changes, update cookie)
  // This is the primary direction: state changes should persist to cookie
  const stateToCookieCallback: WatchCallback<T, T> = (newValue) => {
    if (cookie.value !== newValue) {
      cookie.value = newValue
    }
  }
  watch(state, stateToCookieCallback, { deep: true })

  // Sync cookie -> state (when cookie changes externally, update state)
  // This handles cases where cookie is modified outside of this composable
  // or when cookie is updated by the server/another process
  const cookieToStateCallback: WatchCallback<T, T> = (newValue) => {
    if (state.value !== newValue) {
      state.value = newValue
    }
  }
  watch(cookie, cookieToStateCallback, { deep: true })

  return options.readonly ? shallowReadonly(state) : state
}
