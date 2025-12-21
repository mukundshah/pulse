<script setup lang="ts">
import type { PulseAPIResponse } from '#open-fetch'

import { h } from 'vue'
import { toast } from 'vue-sonner'

import StatusBadge from '../@components/StatusBadge.vue'
import AlertsTable from './@components/Alerts.vue'
import HTTPPerformanceChart from './@components/HTTPPerformanceChart.vue'
import RunResults from './@components/RunResults.vue'
import TCPPerformanceChart from './@components/TCPPerformanceChart.vue'
import UptimeChart from './@components/UptimeChart.vue'

const route = useRoute()
const { projectId, checkId } = route.params as { projectId: string, checkId: string }

const { data: check } = await usePulseAPI('/internal/projects/{projectId}/checks/{checkId}', {
  path: {
    projectId,
    checkId,
  },
})

useHead({
  title: `Check ${check.value?.name}`,
})

useLayoutContext({
  breadcrumbOverrides: computed(() => [
    undefined, // Root
    undefined, // Projects
    {
      label: check.value?.project?.name || 'Project',
      to: `/projects/${projectId}/checks`,
    }, // Project
    {
      label: check.value?.name || 'Check',
      to: `/projects/${projectId}/checks/${checkId}`,
      active: true,
    }, // Check
    false, // false to hide the current breadcrumb item
  ]),
})

const { $pulseAPI } = useNuxtApp()

const handleTriggerCheck = () => {
  const promise = $pulseAPI('/internal/projects/{projectId}/checks/{checkId}/runs/trigger', {
    method: 'POST',
    path: {
      projectId,
      checkId,
    },
  })

  toast.promise(promise, {
    loading: 'Running check',
    success: (data: PulseAPIResponse<'triggerCheckRun'>) => {
      const labels = {
        passing: ['passed', 'success'],
        failing: ['failed', 'error'],
        degraded: ['degraded', 'warning'],
        unknown: ['unknown', 'info'],
      } as const
      return {
        type: labels[data.status as keyof typeof labels][1],
        message: data.status === 'unknown' ? 'Unknown status' : `Check ${labels[data.status as keyof typeof labels][0]}`,
        description: data.total_time_ms ? () => h('span', { class: 'text-muted-foreground text-xs font-mono' }, formatDuration(data.total_time_ms! * 1000)) : undefined,
      }
    },
    error: (error: Error) => error.message,
  })
}
</script>

<template>
  <div class="flex">
    <main class="flex-1 p-4 md:p-6 overflow-y-auto h-[calc(100svh-61px)] flex flex-col gap-y-6">
      <div class="flex items-start justify-between">
        <div class="flex flex-col gap-y-2">
          <div class="flex items-center gap-3">
            <h1 class="text-2xl font-semibold text-balance text-foreground">
              {{ check?.name }}
            </h1>
            <StatusBadge :status="check?.last_status ?? 'unknown'" />
          </div>

          <div class="flex items-center gap-2">
            <Badge class="text-xs font-mono" variant="secondary">
              {{ check?.type.toUpperCase() }}
            </Badge>

            <div v-if="check?.type === 'http'" class="text-sm text-muted-foreground font-mono flex items-center gap-2">
              <Badge class="text-xs" variant="secondary">
                {{ check?.method }}
              </Badge>

              <Badge class="text-xs" variant="secondary">
                {{ constructURL({ host: check?.host!, port: check?.port!, path: check?.path!, queryParams: check?.query_params as Record<string, string> | undefined, secure: check?.secure }) }}
              </Badge>
            </div>
            <div v-else-if="check?.type === 'tcp'" class="text-sm text-muted-foreground font-mono flex items-center gap-2">
              <Badge class="text-xs" variant="secondary">
                {{ check?.host }}{{ ':' }}{{ check?.port }}
              </Badge>
            </div>
            <div v-else-if="check?.type === 'dns'" class="text-sm text-muted-foreground font-mono flex items-center gap-2">
              <Badge class="text-xs" variant="secondary">
                {{ check?.dns_record_type?.toUpperCase() }}
              </Badge>
              <Badge class="text-xs" variant="secondary">
                {{ check?.host }}
              </Badge>
            </div>
          </div>
        </div>

        <div class="flex items-center gap-2.5">
          <Button as-child size="sm" variant="outline">
            <NuxtLink :to="`/projects/${projectId}/checks/${checkId}/edit`">
              Edit
            </NuxtLink>
          </Button>

          <Button loading-auto size="sm" @click="handleTriggerCheck">
            Schedule now
          </Button>
        </div>
      </div>

      <!-- Availability Stats -->
      <!-- <div class="grid grid-cols-2 gap-4 mb-6">
        <Card class="bg-card border-border p-4">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-lg bg-success/10 flex items-center justify-center">
              <Icon class="w-5 h-5 text-success" name="lucide:activity" />
            </div>
            <div>
              <p class="text-sm text-muted-foreground mb-1">
                Availability
              </p>
              <p class="text-2xl font-semibold text-foreground">
                100%
              </p>
            </div>
          </div>
        </Card>
        <Card class="bg-card border-border p-4">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-lg bg-destructive/10 flex items-center justify-center">
              <Icon class="w-5 h-5 text-destructive" name="lucide:activity" />
            </div>
            <div>
              <p class="text-sm text-muted-foreground mb-1">
                Failure Alerts
              </p>
              <p class="text-2xl font-semibold text-foreground">
                0
              </p>
            </div>
          </div>
        </Card>
      </div> -->

      <!-- Uptime Chart -->
      <Card>
        <CardHeader>
          <CardTitle>
            Availability
          </CardTitle>
          <CardDescription>
            Uptime over the last 7 days
          </CardDescription>
          <CardAction>
            <div class="flex items-center justify-end gap-2">
              <Button class="text-xs h-7" size="sm" variant="ghost">
                Today
              </Button>
              <Button class="text-xs h-7" size="sm" variant="ghost">
                1hr
              </Button>
              <Button class="text-xs h-7" size="sm" variant="ghost">
                3hr
              </Button>
              <Button class="text-xs h-7" size="sm" variant="ghost">
                24hr
              </Button>
              <Button class="text-xs h-7" size="sm" variant="secondary">
                7d
              </Button>
            </div>
          </CardAction>
        </CardHeader>
        <CardContent>
          <div>
            <UptimeChart :check-id="checkId" :project-id="projectId" />
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>
            Alerts
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div>
            <AlertsTable :check-id="checkId" :project-id="projectId" />
          </div>
        </CardContent>
      </Card>

      <!-- Locations -->
      <!-- <Card class="bg-card border-border p-6 mb-6">
        <h2 class="text-base font-semibold text-foreground mb-4">
          Locations
        </h2>
        <LocationsTable />
      </Card> -->

      <!-- Performance Chart -->
      <Card v-if="check?.type === 'http' || check?.type === 'tcp'">
        <CardHeader>
          <CardTitle>
            Performance
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div>
            <HTTPPerformanceChart v-if="check?.type === 'http'" :check-id="checkId" :project-id="projectId" />
            <TCPPerformanceChart v-else-if="check?.type === 'tcp'" :check-id="checkId" :project-id="projectId" />
          </div>
        </CardContent>
      </Card>
    </main>

    <aside class="w-96 shrink-0 border-l border-border h-[calc(100svh-61px)]">
      <div class="flex flex-col h-full">
        <div class="p-4 border-b border-border">
          <h2 class="text-sm font-semibold text-foreground mb-1">
            Run results
          </h2>
          <p class="text-xs text-muted-foreground">
            Last 7 days
          </p>
        </div>
        <RunResults class="flex-1 overflow-y-auto" :check-id="checkId" :project-id="projectId" />
      </div>
    </aside>
  </div>
</template>
