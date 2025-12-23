import type { PulseAPIRequestBody } from '#open-fetch'

// @ts-expect-error vfs
import authconfig from '#build/auth.config'

import { getStorage } from './storage'
import { NotImplementedError } from './error'

export const useAuthConfiguration = ({ namespace }: { namespace?: string } = {}) => {
  const { $pulseAPI } = useNuxtApp()

  const getConfiguration = async () => {
    throw new NotImplementedError()
  }

  return {
    getConfiguration,
  }
}

export const useAuth = ({ namespace }: { namespace?: string } = {}) => {
  const { $pulseAPI } = useNuxtApp()
  const storage = getStorage({ namespace })

  const login = async (data: PulseAPIRequestBody<'loginUser'>) => {
    const response = await $pulseAPI('/internal/auth/login', {
      method: 'POST',
      body: data,
    })

    storage.setSessionToken(response.token ?? null)
    storage.setAuthenticationStatus(!!response.token)

    return response
  }

  const requestLoginCode = async (data: unknown) => {
    throw new NotImplementedError()
  }

  const confirmLoginCode = async (data: unknown) => {
    throw new NotImplementedError()
  }

  const logout = async () => {
    const response = await $pulseAPI('/internal/auth/session', {
      method: 'DELETE',
    })

    storage.setSessionToken(null)
    storage.setAuthenticationStatus(false)
    storage.setOnboardingStatus(false)

    return response
  }

  const signup = async (data: PulseAPIRequestBody<'registerUser'>) => {
    const response = await $pulseAPI('/internal/auth/register', {
      method: 'POST',
      body: data,
    })

    return response
  }

  const reauthenticate = async (data: unknown) => {
    throw new NotImplementedError()
  }

  const syncAuthenticationStatus = async () => {
    try {
      await $pulseAPI('/internal/auth/session', { method: 'GET' })
      return true
    } catch (error: any) {
      // If unauthorized, clear auth state
      if (error?.status === 401) {
        storage.setSessionToken(null)
        storage.setAuthenticationStatus(false)
        storage.setOnboardingStatus(false)
        return false
      }
      throw error
    }
  }

  const me = async ({ signal }: { signal?: AbortSignal } = {}) => {
    const response = await $pulseAPI('/internal/account/me', { method: 'GET', signal })
    return response
  }

  return {
    get token() {
      return storage.getSessionToken()
    },
    isAuthenticated: computed(() => storage.isAuthenticated.value),
    isOnboarded: computed(() => storage.isOnboarded.value || !authconfig.onboarding?.enabled),
    login,
    requestLoginCode,
    confirmLoginCode,
    logout,
    signup,
    reauthenticate,
    syncAuthenticationStatus,
    me,
  }
}

export const useSocialAuth = ({ namespace }: { namespace?: string } = {}) => {
  const { $pulseAPI } = useNuxtApp()
  const storage = getStorage({ namespace })
  const url = useRequestURL()

  const getProviderAccounts = async () => {
    throw new NotImplementedError()
  }

  const connect = async (data: unknown) => {
    throw new NotImplementedError()
  }

  const callback = async (provider: string, data: Record<string, any>) => {
    throw new NotImplementedError()
  }

  const disconnect = async (data: unknown) => {
    throw new NotImplementedError()
  }

  return {
    connect,
    callback,
    disconnect,
    getProviderAccounts,
  }
}

export const useMFA = ({ namespace }: { namespace?: string } = {}) => {
  // TODO: Implement MFA

  return { }
}

export const useEmailManagement = ({ namespace }: { namespace?: string } = {}) => {
  const { $pulseAPI } = useNuxtApp()

  const getEmailAddresses = async () => {
    throw new NotImplementedError()
  }

  const addEmail = async (data: unknown) => {
    throw new NotImplementedError()
  }

  const checkEmail = async (data: unknown) => {
    throw new NotImplementedError()
  }

  const deleteEmail = async (data: unknown) => {
    throw new NotImplementedError()
  }

  const markEmailAsPrimary = async (data: unknown) => {
    throw new NotImplementedError()
  }

  const requestEmailVerification = async (data: PulseAPIRequestBody<'resendVerificationEmail'>) => {
    const response = await $pulseAPI('/internal/auth/resend-verification', {
      method: 'POST',
      body: data,
    })
    return response
  }

  const verifyEmail = async (data: PulseAPIRequestBody<'verifyEmail'>) => {
    const response = await $pulseAPI('/internal/auth/verify-email', {
      method: 'POST',
      body: data,
    })
    return response
  }

  return {
    addEmail,
    checkEmail,
    deleteEmail,
    markEmailAsPrimary,
    getEmailAddresses,
    requestEmailVerification,
    verifyEmail,
  }
}

// export const usePhoneManagement = ({ namespace }: { namespace?: string } = {}) => {
//   const storage = getStorage({ namespace })
//   const client = getClient({ storage })

//   const getPhoneNumbers = async (): Promise<PhoneNumber[]> => {
//     const { data: response } = await client.listPhoneNumbers()
//     return response
//   }
// }

export const usePasswordManagement = ({ namespace }: { namespace?: string } = {}) => {
  const { $pulseAPI } = useNuxtApp()

  const requestPasswordReset = async (data: PulseAPIRequestBody<'requestPasswordReset'>) => {
    const response = await $pulseAPI('/internal/auth/password/reset', {
      method: 'POST',
      body: data,
    })
    return response
  }

  const getPasswordResetInfo = async (token: string) => {
    throw new NotImplementedError()
  }

  const resetPassword = async (data: PulseAPIRequestBody<'resetPassword'>) => {
    const response = await $pulseAPI('/internal/auth/password/reset/confirm', {
      method: 'POST',
      body: data,
    })
    return response
  }

  const changePassword = async (data: PulseAPIRequestBody<'changePassword'>) => {
    const response = await $pulseAPI('/internal/account/password', {
      method: 'PATCH',
      body: data,
    })
    return response
  }

  return {
    requestPasswordReset,
    resetPassword,
    getPasswordResetInfo,
    changePassword,
  }
}

export const usePermissions = ({ namespace = 'default' }: { namespace?: string } = {}) => {
  const { $pulseAPI } = useNuxtApp()
  const { token } = useAuth({ namespace })

  const prefix = namespace !== 'default' ? `${namespace}:` : ''

  const _roles = useState<string[]>(`${prefix}roles`, () => [])
  const _permissions = useState<Record<string, Record<string, boolean>>>(`${prefix}permissions`, () => ({}))

  const _loadPermissions = async () => {
    throw new NotImplementedError()
  }

  // const getPermissions = async () => {
  //   const { data: response } = await client.listPermissions()
  //   return response
  // }

  const hasRole = (role: string) => {
    if (!_roles.value?.includes(role)) return false
    return _roles.value.includes(role)
  }

  const hasAnyRole = (roles: string[]) => {
    if (!_roles.value?.length) return false
    return roles.some(role => _roles.value!.includes(role))
  }

  const hasAllRoles = (roles: string[]) => {
    if (!_roles.value.length) return false
    return roles.every(role => _roles.value.includes(role))
  }

  const _hasFullAccess = () => {
    return true
    // return hasAnyRole(authconfig.fullAccessRoles)
  }

  type Permission = `${string}.${string}` | `${string}`

  const hasPermission = (permission: Permission) => {
    // Check for full access roles first
    if (_hasFullAccess()) return true

    if (Object.keys(_permissions.value || {}).length === 0) return false

    if (permission.includes('.')) {
      const [resource, action] = permission.split('.')
      if (!resource || !action) return false
      return _permissions.value?.[resource]?.[action] ?? false
    }

    // Resource-only check - return true if any action is allowed for this resource
    return Object.keys(_permissions.value?.[permission] || {}).some(action => _permissions.value?.[permission]?.[action] === true)
  }

  const hasAnyPermission = (permissions: Permission[]) => {
    // Check for full access roles first
    if (_hasFullAccess()) return true

    if (!_permissions.value) return false
    return permissions.some(permission => hasPermission(permission))
  }

  const hasAllPermissions = (permissions: Permission[]) => {
    // Check for full access roles first
    if (_hasFullAccess()) return true

    if (!_permissions.value) return false
    return permissions.every(permission => hasPermission(permission))
  }

  return {
    hasRole,
    hasAnyRole,
    hasAllRoles,
    hasPermission,
    hasAnyPermission,
    hasAllPermissions,
  }
}
