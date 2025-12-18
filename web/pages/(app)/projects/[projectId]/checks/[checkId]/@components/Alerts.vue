<script setup lang="ts">
import { Bell } from 'lucide-vue-next'

const props = withDefaults(defineProps<{
  projectId: string
  checkId: string
  filters?: Record<string, string>
}>(), {
  filters: () => ({}),
})

const STATUS_MAP = {
  passing: {
    label: 'Recovered',
    icon: 'lucide:circle-check',
    color: 'text-green-500',
  },
  degraded: {
    label: 'Degraded',
    icon: 'lucide:circle-alert',
    color: 'text-amber-500',
  },
  failing: {
    label: 'Failed',
    icon: 'lucide:circle-x',
    color: 'text-red-500',
  },
  unknown: {
    label: 'Unknown',
    icon: 'lucide:circle-minus',
    color: 'text-gray-500',
  },
} as const

const { data: response, pending } = useLazyPulseAPI(`/internal/projects/{projectId}/checks/{checkId}/alerts`, {
  path: {
    projectId: props.projectId,
    checkId: props.checkId,
  },
})
</script>

<template>
  <div class="divide-y divide-border">
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead class="p-0 w-px">
            <span class="sr-only">Link</span>
          </TableHead>
          <TableHead class="w-[30%] min-w-[150px] text-left">
            Status
          </TableHead>
          <TableHead class="w-[25%] min-w-[120px] text-left">
            Region
          </TableHead>
          <TableHead class="w-[45%] min-w-[180px] text-left">
            Timestamp
          </TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <template v-if="pending">
          <TableRow v-for="i in 5" :key="i">
            <TableCell class="p-0" />
            <TableCell>
              <div class="flex items-center gap-2">
                <Skeleton class="h-5 w-5 rounded-full" />
                <Skeleton class="h-4 w-20" />
              </div>
            </TableCell>
            <TableCell>
              <Skeleton class="h-4 w-16" />
            </TableCell>
            <TableCell>
              <Skeleton class="h-4 w-24" />
            </TableCell>
          </TableRow>
        </template>

        <TableRow v-if="response && response.data.length === 0">
          <TableCell colspan="4">
            <Empty class="h-48">
              <EmptyHeader>
                <EmptyMedia class="text-muted-foreground rounded-full bg-muted" variant="icon">
                  <Bell />
                </EmptyMedia>
                <EmptyDescription>No alerts have been triggered yet</EmptyDescription>
              </EmptyHeader>
            </Empty>
          </TableCell>
        </TableRow>

        <TableRow v-for="alert in response?.data" :key="alert.id">
          <TableCell class="p-0">
            <NuxtLink class="absolute h-full w-full inset-0" :to="`/projects/${projectId}/checks/${checkId}/runs/${alert.run_id}`" />
          </TableCell>
          <TableCell>
            <div class="flex items-center gap-2">
              <Icon
                class="text-lg"
                :class="STATUS_MAP[alert.status as keyof typeof STATUS_MAP].color"
                :name="STATUS_MAP[alert.status as keyof typeof STATUS_MAP].icon"
              />
              <span class="text-sm font-medium text-foreground">
                {{ STATUS_MAP[alert.status as keyof typeof STATUS_MAP].label }}
              </span>
            </div>
          </TableCell>
          <TableCell>
            {{ alert.region?.name || 'Unknown' }}
          </TableCell>
          <TableCell>
            <NuxtTime
              relative
              title
              numeric="auto"
              :datetime="alert.created_at"
            />
          </TableCell>
        </TableRow>
      </TableBody>
    </Table>
  </div>
</template>
