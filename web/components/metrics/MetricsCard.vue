<script setup lang="ts">
interface Props {
  label: string
  value: string | number
  change?: number
  changeLabel?: string
  trend?: 'up' | 'down' | 'neutral'
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  change: undefined,
  changeLabel: undefined,
  trend: 'neutral',
  loading: false,
})

const changeColor = computed(() => {
  if (props.trend === 'up') {
    return props.change && props.change > 0 ? 'text-green-600' : 'text-red-600'
  }
  if (props.trend === 'down') {
    return props.change && props.change < 0 ? 'text-green-600' : 'text-red-600'
  }
  return 'text-muted-foreground'
})

const changeIcon = computed(() => {
  if (props.change === undefined || props.change === 0) {
    return null
  }
  if (props.trend === 'up') {
    return props.change > 0 ? 'lucide:trending-up' : 'lucide:trending-down'
  }
  if (props.trend === 'down') {
    return props.change < 0 ? 'lucide:trending-up' : 'lucide:trending-down'
  }
  return null
})

const formattedChange = computed(() => {
  if (props.change === undefined) {
    return null
  }
  const sign = props.change > 0 ? '+' : ''
  return `${sign}${props.change.toFixed(1)}%`
})
</script>

<template>
  <Card>
    <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
      <CardTitle class="text-sm font-medium">
        {{ label }}
      </CardTitle>
    </CardHeader>
    <CardContent>
      <div class="text-2xl font-bold">
        {{ loading ? '...' : value }}
      </div>
      <div
        v-if="formattedChange"
        class="flex items-center gap-1 text-xs"
        :class="changeColor"
      >
        <Icon
          v-if="changeIcon"
          :name="changeIcon"
          class="h-3 w-3"
        />
        <span>{{ formattedChange }}</span>
        <span v-if="changeLabel" class="text-muted-foreground">
          {{ changeLabel }}
        </span>
      </div>
      <p v-else-if="changeLabel" class="text-xs text-muted-foreground">
        {{ changeLabel }}
      </p>
    </CardContent>
  </Card>
</template>


