<script setup lang="ts">
import { ASSERTION_PROPERTIES } from '@/constants/http'
import { constructURL } from '@/utils/url'

interface AssertionResult {
  source?: string
  property?: string
  comparison?: string
  target?: string | number
  actual?: string | number
  passed?: boolean
  name?: string
}

const route = useRoute()
const { projectId, checkId, runId } = route.params as { projectId: string, checkId: string, runId: string }

const STATUS_ICON_COLOR_MAP = {
  passing: {
    icon: 'lucide:circle-check',
    color: 'text-green-500',
  },
  degraded: {
    icon: 'lucide:circle-alert',
    color: 'text-amber-500',
  },
  failing: {
    icon: 'lucide:circle-x',
    color: 'text-red-500',
  },
  unknown: {
    icon: 'lucide:circle-minus',
    color: 'text-gray-500',
  },
} as const

const { data: run } = await usePulseAPI('/internal/projects/{projectId}/checks/{checkId}/runs/{runId}', {
  path: {
    projectId,
    checkId,
    runId,
  },
})

useHead({
  title: `Run #${run.value?.id.slice(0, 8)}`,
})

const formatTime = (ms: number) => {
  if (ms < 1) {
    return `${Math.round(ms * 1000)}µs`
  }
  if (ms < 1000) {
    return `${Math.round(ms)}ms`
  }
  return `${(ms / 1000).toFixed(2)}s`
}

const assertions = computed(() => {
  if (run.value?.assertion_results && Array.isArray(run.value.assertion_results) && run.value.assertion_results.length > 0) {
    return run.value.assertion_results as AssertionResult[]
  }

  return []
})

const networkTimings = computed(() => {
  if (!run.value?.network_timings || typeof run.value.network_timings !== 'object') {
    return null
  }
  const timings = run.value.network_timings as Record<string, unknown>
  return {
    dns: typeof timings.dns_lookup_ms === 'number' ? timings.dns_lookup_ms : null,
    tcp: typeof timings.tcp_connection_ms === 'number' ? timings.tcp_connection_ms : null,
    tls: typeof timings.tls_handshake_ms === 'number' ? timings.tls_handshake_ms : null,
    firstByte: typeof timings.time_to_first_byte_ms === 'number' ? timings.time_to_first_byte_ms : null,
    download: typeof timings.response_time_ms === 'number' && typeof timings.time_to_first_byte_ms === 'number'
      ? timings.response_time_ms - timings.time_to_first_byte_ms
      : null,
  }
})

const maxTiming = computed(() => {
  if (!networkTimings.value) return 0
  const timings = networkTimings.value
  return Math.max(
    timings.dns || 0,
    timings.tcp || 0,
    timings.tls || 0,
    timings.firstByte || 0,
    timings.download || 0,
  )
})

const checkURL = computed(() => {
  if (!run.value?.check) return ''
  return constructURL({
    host: run.value.check.host,
    port: run.value.check.port,
    path: run.value.check.path,
    queryParams: run.value.check.query_params as Record<string, string> | undefined,
    secure: run.value.check.secure,
  })
})

const statusText = computed(() => {
  const status = run.value?.status
  if (status === 'passing') return 'Passed'
  if (status === 'failing') return 'Failed'
  if (status === 'degraded') return 'Degraded'
  return 'Unknown'
})
</script>

<template>
  <div class="flex">
    <main class="flex-1 p-4 md:p-6 overflow-y-auto h-[calc(100svh-61px)] flex flex-col gap-y-6">
      <!-- Header -->
      <div class="space-y-2.5">
        <div class="flex items-center gap-3">
          <Icon
            class="text-lg"
            :class="STATUS_ICON_COLOR_MAP[run?.status as keyof typeof STATUS_ICON_COLOR_MAP].color"
            :name="STATUS_ICON_COLOR_MAP[run?.status as keyof typeof STATUS_ICON_COLOR_MAP].icon"
          />

          <h1 class="text-lg font-semibold text-foreground">
            {{ run?.region?.name || 'Unknown region' }}
            <span v-if="run?.region?.flag" class="ml-2">{{ run.region.flag }}</span>
          </h1>
        </div>

        <div class="text-sm text-muted-foreground flex items-center gap-2">
          <div>
            {{ statusText }} on
            <NuxtTime
              :datetime="run?.created_at ?? new Date()"
              v-bind="{
                locale: 'en-US',
                dateStyle: 'long',
                timeStyle: 'short',
              }"
            />
          </div>

          <div v-if="run?.check?.ip_version">
            <span class="mr-2">•</span>
            <!-- FIXME: use proper label mapping -->
            <span>{{ run.check.ip_version.toUpperCase().replace('V', 'v') }}</span>
          </div>
        </div>
        <div class="flex items-center justify-between">
          <div class="text-sm font-mono text-foreground">
            <Badge class="mr-2" variant="secondary">
              {{ run?.check?.method || 'GET' }}
            </Badge>
            {{ checkURL }}
          </div>
          <div class="flex items-center gap-4">
            <div class="flex items-center gap-2">
              <Badge
                class="text-base font-mono px-3"
                :variant="run?.response_status && run.response_status >= 200 && run.response_status < 300 ? 'default' : 'destructive'"
              >
                {{ run?.response_status || '—' }}
              </Badge>
            </div>
            <div class="flex items-center gap-2 text-sm font-mono">
              <Icon class="w-4 h-4 text-muted-foreground" name="lucide:clock" />
              <span>{{ run?.total_time_ms ? formatTime(run.total_time_ms) : '—' }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Assertions -->
      <Card v-if="assertions.length > 0">
        <CardHeader>
          <div class="flex items-center justify-between">
            <CardTitle>Assertions</CardTitle>
          </div>
        </CardHeader>
        <CardContent>
          <div class="overflow-x-auto">
            <table class="w-full text-sm">
              <thead>
                <tr class="border-b border-border">
                  <th class="text-left py-2 px-3 text-muted-foreground font-medium">
                    SOURCE
                  </th>
                  <th class="text-left py-2 px-3 text-muted-foreground font-medium">
                    PROPERTY
                  </th>
                  <th class="text-left py-2 px-3 text-muted-foreground font-medium">
                    COMPARISON
                  </th>
                  <th class="text-left py-2 px-3 text-muted-foreground font-medium">
                    TARGET
                  </th>
                  <th class="text-left py-2 px-3 text-muted-foreground font-medium">
                    ACTUAL
                  </th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="(assertion, index) in assertions"
                  :key="index"
                  class="border-b border-border/50"
                >
                  <td class="py-2 px-3">
                    <div class="flex items-center gap-2">
                      <Icon
                        v-if="assertion.passed !== false"
                        class="w-4 h-4 text-green-500 shrink-0"
                        name="lucide:check"
                      />
                      <span v-if="assertion.source">{{ ASSERTION_PROPERTIES[assertion.source as keyof typeof ASSERTION_PROPERTIES].label || 'Response' }}</span>
                    </div>
                  </td>
                  <td class="py-2 px-3">
                    {{ assertion.property || assertion.name || '—' }}
                  </td>
                  <td class="py-2 px-3">
                    {{ assertion.comparison || '—' }}
                  </td>
                  <td class="py-2 px-3">
                    <span v-if="assertion.target !== undefined" class="font-mono">{{ assertion.target }}</span>
                    <span v-else>—</span>
                  </td>
                  <td class="py-2 px-3">
                    <span v-if="assertion.actual !== undefined" class="font-mono">{{ assertion.actual }}</span>
                    <span v-else>—</span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </CardContent>
      </Card>

      <!-- Timing -->
      <Card v-if="networkTimings && maxTiming > 0">
        <CardHeader>
          <CardTitle>Timing</CardTitle>
        </CardHeader>
        <CardContent>
          <div class="space-y-4">
            <!-- Connection Start -->
            <div>
              <div class="text-xs font-medium text-muted-foreground mb-2">
                CONNECTION START
              </div>
              <div class="space-y-2">
                <div v-if="networkTimings.dns !== null">
                  <div class="flex items-center justify-between mb-1">
                    <span class="text-sm">DNS</span>
                    <span class="text-sm font-mono">{{ formatTime(networkTimings.dns) }}</span>
                  </div>
                  <div class="h-2 bg-muted rounded-full overflow-hidden">
                    <div
                      class="h-full bg-purple-500"
                      :style="{ width: `${(networkTimings.dns / maxTiming) * 100}%` }"
                    ></div>
                  </div>
                </div>
                <div v-if="networkTimings.tcp !== null">
                  <div class="flex items-center justify-between mb-1">
                    <span class="text-sm">TCP</span>
                    <span class="text-sm font-mono">{{ formatTime(networkTimings.tcp) }}</span>
                  </div>
                  <div class="h-2 bg-muted rounded-full overflow-hidden">
                    <div
                      class="h-full bg-orange-500"
                      :style="{ width: `${(networkTimings.tcp / maxTiming) * 100}%` }"
                    ></div>
                  </div>
                </div>
                <div v-if="networkTimings.tls !== null">
                  <div class="flex items-center justify-between mb-1">
                    <span class="text-sm">TLS</span>
                    <span class="text-sm font-mono">{{ formatTime(networkTimings.tls) }}</span>
                  </div>
                  <div class="h-2 bg-muted rounded-full overflow-hidden">
                    <div
                      class="h-full bg-yellow-500"
                      :style="{ width: `${(networkTimings.tls / maxTiming) * 100}%` }"
                    ></div>
                  </div>
                </div>
              </div>
            </div>

            <!-- Request / Response -->
            <div>
              <div class="text-xs font-medium text-muted-foreground mb-2">
                REQUEST / RESPONSE
              </div>
              <div class="space-y-2">
                <div v-if="networkTimings.firstByte !== null">
                  <div class="flex items-center justify-between mb-1">
                    <span class="text-sm">First Byte</span>
                    <span class="text-sm font-mono">{{ formatTime(networkTimings.firstByte) }}</span>
                  </div>
                  <div class="h-2 bg-muted rounded-full overflow-hidden">
                    <div
                      class="h-full bg-blue-500"
                      :style="{ width: `${(networkTimings.firstByte / maxTiming) * 100}%` }"
                    ></div>
                  </div>
                </div>
                <div v-if="networkTimings.download !== null && networkTimings.download > 0">
                  <div class="flex items-center justify-between mb-1">
                    <span class="text-sm">Download</span>
                    <span class="text-sm font-mono">{{ formatTime(networkTimings.download) }}</span>
                  </div>
                  <div class="h-2 bg-muted rounded-full overflow-hidden">
                    <div
                      class="h-full bg-green-500"
                      :style="{ width: `${(networkTimings.download / maxTiming) * 100}%` }"
                    ></div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </main>
  </div>
</template>
