<script setup lang="ts">
import { useInfiniteScroll } from '@vueuse/core'
import { Activity } from 'lucide-vue-next'

const props = withDefaults(defineProps<{
  projectId: string
  checkId: string
  filters?: Record<string, string>
}>(), {
  filters: () => ({}),
})

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

const { $pulseAPI } = useNuxtApp()

const container = useTemplateRef<HTMLElement>('containerRef')

const fetcher = (after?: string, limit: number = 50) => {
  return $pulseAPI('/internal/projects/{projectId}/checks/{checkId}/runs', {
    path: {
      projectId: props.projectId,
      checkId: props.checkId,
    },
    query: {
      after,
      limit,
      ...props.filters,
    },
  })
}

const { data: response, pending, error } = useLazyAsyncData(`check:${props.checkId}:runs`, () => fetcher())

const { isLoading } = useInfiniteScroll(
  container,
  async () => {
    if (!response.value || response.value.next_cursor === null) {
      return
    }

    const { data, next_cursor } = await fetcher(response.value?.next_cursor)
    response.value.data.push(...data)
    response.value.next_cursor = next_cursor
  },
  {
    distance: 100,
    canLoadMore: () => {
      return response.value?.next_cursor !== null
    },
  },
)
</script>

<template>
  <div ref="containerRef" class="divide-y divide-border">
    <Empty v-if="response && response.data.length === 0" class="h-full">
      <EmptyHeader>
        <EmptyMedia class="text-muted-foreground rounded-full bg-muted" variant="icon">
          <Activity />
        </EmptyMedia>
        <EmptyDescription>Check has no run results yet</EmptyDescription>
      </EmptyHeader>
    </Empty>

    <div
      v-for="run in response?.data"
      :key="run.id"
      class="flex items-center gap-3 p-3 hover:bg-accent/50 transition-colors group"
    >
      <Icon
        class="text-lg"
        :class="STATUS_ICON_COLOR_MAP[run.status as keyof typeof STATUS_ICON_COLOR_MAP].color"
        :name="STATUS_ICON_COLOR_MAP[run.status as keyof typeof STATUS_ICON_COLOR_MAP].icon"
      />
      <div class="flex-1 min-w-0">
        <div class="flex items-center gap-2 mb-1">
          <span class="text-sm font-medium text-foreground">
            {{ run.region?.name || 'Unknown' }}
          </span>
        </div>
        <p class="text-xs text-muted-foreground">
          <NuxtTime
            relative
            title
            numeric="auto"
            :datetime="run.created_at"
          />
        </p>
      </div>
      <span class="text-xs font-mono text-muted-foreground">
        {{ run.total_time_ms }}ms
      </span>
    </div>

    <div v-if="pending || isLoading" class="p-3 space-y-2">
      <div
        v-for="i in 5"
        :key="i"
        class="flex items-center gap-3 p-3 rounded-lg"
      >
        <Skeleton class="w-4 h-4 rounded-full shrink-0" />
        <div class="flex-1 min-w-0 space-y-2">
          <div class="flex items-center gap-2">
            <Skeleton class="h-4 w-24" />
          </div>
          <Skeleton class="h-3 w-16" />
        </div>
        <Skeleton class="h-3 w-12" />
      </div>
    </div>

    <div v-if="error" class="p-6 text-center">
      <p class="text-sm text-destructive">
        Failed to load run results
      </p>
    </div>
  </div>
</template>
