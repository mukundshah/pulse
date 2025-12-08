<script setup lang="ts">
interface Props {
  runs: Array<{
    id: string
    status: string
    created_at: string
    region?: { name: string; code: string }
    network_timings?: { total_time_ms?: number }
  }>
}

const props = defineProps<Props>()

const getStatusColor = (status: string) => {
  if (status === 'success') {
    return 'text-green-600'
  }
  if (status === 'fail' || status === 'error') {
    return 'text-red-600'
  }
  if (status === 'timeout') {
    return 'text-yellow-600'
  }
  return 'text-gray-600'
}

const formatDuration = (run) => {
  if (run.network_timings?.total_time_ms) {
    return `${run.network_timings.total_time_ms}ms`
  }
  return '--'
}
</script>

<template>
  <Card class="sticky top-4">
    <CardHeader>
      <CardTitle>Run Results</CardTitle>
      <CardDescription>
        Recent check executions
      </CardDescription>
    </CardHeader>
    <CardContent>
      <div class="space-y-4">
        <div
          v-for="run in runs.slice(0, 20)"
          :key="run.id"
          class="flex items-start justify-between border-b pb-4 last:border-0"
        >
          <div class="flex-1 space-y-1">
            <div class="flex items-center gap-2">
              <div
                class="h-2 w-2 rounded-full"
                :class="{
                  'bg-green-500': run.status === 'success',
                  'bg-red-500': run.status === 'fail' || run.status === 'error',
                  'bg-yellow-500': run.status === 'timeout',
                  'bg-gray-500': !['success', 'fail', 'error', 'timeout'].includes(run.status),
                }"
              />
              <span class="text-sm font-medium" :class="getStatusColor(run.status)">
                {{ run.status.toUpperCase() }}
              </span>
            </div>
            <div class="text-xs text-muted-foreground">
              {{ run.region?.name || run.region?.code || 'Unknown' }}
            </div>
            <div class="text-xs text-muted-foreground">
              {{ formatDuration(run) }}
            </div>
          </div>
          <div class="text-xs text-muted-foreground">
            {{ new Date(run.created_at).toLocaleTimeString() }}
          </div>
        </div>
        <div v-if="runs.length === 0" class="text-center text-sm text-muted-foreground">
          No runs yet
        </div>
      </div>
    </CardContent>
  </Card>
</template>


