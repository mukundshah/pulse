<script setup lang="ts">
import TimeRangeFilters from '@/components/filters/TimeRangeFilters.vue'
import StatusFilters from '@/components/filters/StatusFilters.vue'
import MetricsCard from '@/components/metrics/MetricsCard.vue'
import PerformanceGraph from '@/components/metrics/PerformanceGraph.vue'
import RunResultsSidebar from '@/components/checks/RunResultsSidebar.vue'
import StatusBadge from '~/components/StatusBadge.vue'
import LocationFlag from '~/components/LocationFlag.vue'

const route = useRoute()
const projectId = route.params.id as string
const checkId = route.params.checkId as string

useHead({
  title: 'Check Details',
})

// Load check details
const checkResponse = useAPI(`/checks/${checkId}`, {
  protected: true,
})
const check = computed(() => checkResponse.data.value)

// Load check runs
const runsResponse = useAPI(`/checks/${checkId}/runs?limit=100`, {
  protected: true,
})
const runs = computed(() => {
  const data = runsResponse.data.value
  return Array.isArray(data) ? data : []
})

// Time range filter
const timeRange = ref('24h')

// Calculate metrics from runs
const metrics = computed(() => {
  if (!runs.value.length) {
    return {
      availability: 100,
      retries: 0,
      p50: 0,
      p95: 0,
      failureAlerts: 0,
      spanErrors: 0,
    }
  }

  const successful = runs.value.filter(r => r.status === 'success').length
  const availability = (successful / runs.value.length) * 100

  // Calculate response times (mock for now - would need actual timing data)
  const responseTimes = runs.value.map(() => Math.random() * 1000 + 500)
  responseTimes.sort((a, b) => a - b)
  const p50 = responseTimes[Math.floor(responseTimes.length * 0.5)]
  const p95 = responseTimes[Math.floor(responseTimes.length * 0.95)]

  return {
    availability,
    retries: 0,
    p50: Math.round(p50),
    p95: Math.round(p95),
    failureAlerts: runs.value.filter(r => r.status === 'fail' || r.status === 'error').length,
    spanErrors: 0,
  }
})

const getStatusFromLastStatus = (status: string) => {
  if (status === 'success') {
    return 'passing'
  }
  if (status === 'fail' || status === 'error') {
    return 'failing'
  }
  if (status === 'timeout') {
    return 'degraded'
  }
  return 'unknown'
}
</script>

<template>
  <div class="flex flex-1 gap-6">
    <!-- Main Content -->
    <div class="flex-1 space-y-6">
      <!-- Header -->
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-4">
          <StatusBadge
            :status="check ? getStatusFromLastStatus(check.last_status) : 'unknown'"
          />
          <div>
            <h1 class="text-3xl font-semibold tracking-tight">
              {{ check?.name || 'Check' }}
            </h1>
            <p class="text-muted-foreground">
              Last updated {{ check?.last_run_at ? new Date(check.last_run_at).toLocaleString() : 'Never' }}
            </p>
          </div>
        </div>
        <div class="flex items-center gap-2">
          <LocationFlag
            v-for="region in check?.regions || []"
            :key="region.id"
            :location="region.code || region.name"
          />
          <Button as-child variant="outline">
            <NuxtLink :to="`/projects/${projectId}/checks/${checkId}/edit`">
              <Icon class="mr-2 h-4 w-4" name="lucide:edit" />
              Edit
            </NuxtLink>
          </Button>
        </div>
      </div>

      <!-- Time Range and Status Filters -->
      <div class="flex items-center justify-between">
        <TimeRangeFilters v-model="timeRange" />
        <StatusFilters />
      </div>

      <!-- Key Metrics Cards -->
      <div class="grid gap-4 md:grid-cols-3 lg:grid-cols-6">
        <MetricsCard
          label="Availability"
          :value="`${metrics.availability.toFixed(1)}%`"
          :change="0"
        />
        <MetricsCard
          label="Retries"
          :value="`${metrics.retries}%`"
          :change="0"
        />
        <MetricsCard
          label="P50"
          :value="`${metrics.p50}ms`"
          :change="0"
        />
        <MetricsCard
          label="P95"
          :value="`${metrics.p95}ms`"
          :change="0"
        />
        <MetricsCard
          label="Failure Alerts"
          :value="metrics.failureAlerts"
        />
        <MetricsCard
          label="Span Errors"
          :value="metrics.spanErrors"
        />
      </div>

      <!-- Historical Graph -->
      <Card>
        <CardHeader>
          <CardTitle>Check Results Over Time</CardTitle>
        </CardHeader>
        <CardContent>
          <PerformanceGraph
            :data="runs.map((run, i) => ({
              timestamp: run.created_at,
              success: run.status === 'success' ? 1 : 0,
              fail: run.status === 'fail' || run.status === 'error' ? 1 : 0,
            }))"
            :metrics="[
              { key: 'success', label: 'Success', color: 'hsl(142, 76%, 36%)' },
              { key: 'fail', label: 'Failed', color: 'hsl(0, 84%, 60%)' },
            ]"
            type="bar"
          />
        </CardContent>
      </Card>

      <!-- Alerts Table -->
      <Card>
        <CardHeader>
          <CardTitle>Recent Alerts</CardTitle>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>STATUS</TableHead>
                <TableHead>LOCATION</TableHead>
                <TableHead>NOTIFICATIONS</TableHead>
                <TableHead>TIMESTAMP</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow v-if="runs.length === 0">
                <TableCell colspan="4" class="text-center text-muted-foreground">
                  No alerts yet
                </TableCell>
              </TableRow>
              <TableRow
                v-for="run in runs.filter(r => r.status === 'fail' || r.status === 'error').slice(0, 10)"
                :key="run.id"
              >
                <TableCell>
                  <StatusBadge :status="getStatusFromLastStatus(run.status)" />
                </TableCell>
                <TableCell>
                  <LocationFlag
                    v-if="run.region"
                    :location="run.region.code || run.region.name"
                  />
                </TableCell>
                <TableCell>--</TableCell>
                <TableCell>
                  {{ new Date(run.created_at).toLocaleString() }}
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>

    <!-- Right Sidebar -->
    <div class="w-80">
      <RunResultsSidebar :runs="runs" />
    </div>
  </div>
</template>
