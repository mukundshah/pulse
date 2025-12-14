<script setup lang="ts">
import { useInfiniteScroll } from '@vueuse/core'
import { Activity, CheckCircle2 } from 'lucide-vue-next'

const props = withDefaults(defineProps<{
  checkId: string
  filters?: Record<string, string>
}>(), {
  filters: () => ({}),
})

const { $pulseAPI } = useNuxtApp()

const container = useTemplateRef<HTMLElement>('containerRef')

const fetcher = (after?: string, limit: number = 50) => {
  return $pulseAPI('/v1/checks/{checkId}/runs', {
    path: {
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
  <div ref="containerRef">
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
      class="flex items-center gap-3 p-3 rounded-lg hover:bg-accent/50 transition-colors group"
    >
      <CheckCircle2 class="w-4 h-4 text-success shrink-0" />
      <div class="flex-1 min-w-0">
        <div class="flex items-center gap-2 mb-1">
          <span class="text-sm font-medium text-foreground">
            {{ run.region?.name || 'Unknown' }}
          </span>
          <Badge v-if="run.region?.code" class="text-xs px-1.5 py-0 h-5" variant="outline">
            {{ run.region.code }}
          </Badge>
        </div>
        <p class="text-xs text-muted-foreground">
          <NuxtTime :datetime="run.created_at" />
        </p>
      </div>
      <span class="text-xs font-mono text-muted-foreground">
        {{ run.metrics?.duration }}
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
            <Skeleton class="h-5 w-12 rounded-full" />
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
