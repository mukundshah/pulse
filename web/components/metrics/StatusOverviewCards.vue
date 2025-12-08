<script setup lang="ts">
interface Props {
  passing?: number
  degraded?: number
  failing?: number
  loading?: boolean
}

withDefaults(defineProps<Props>(), {
  passing: 0,
  degraded: 0,
  failing: 0,
  loading: false,
})
</script>

<template>
  <div class="grid gap-4 md:grid-cols-3">
    <!-- PASSING Card -->
    <Card>
      <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle class="text-sm font-medium">
          PASSING
        </CardTitle>
        <div class="h-2 w-2 rounded-full bg-green-500" />
      </CardHeader>
      <CardContent>
        <div class="text-2xl font-bold text-green-600">
          {{ loading ? '...' : passing }}
        </div>
        <p class="text-xs text-muted-foreground">
          Checks passing
        </p>
      </CardContent>
    </Card>

    <!-- DEGRADED Card -->
    <Card>
      <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle class="text-sm font-medium">
          DEGRADED
        </CardTitle>
        <div class="h-2 w-2 rounded-full bg-yellow-500" />
      </CardHeader>
      <CardContent>
        <div class="text-2xl font-bold text-yellow-600">
          {{ loading ? '...' : degraded }}
        </div>
        <p class="text-xs text-muted-foreground">
          Checks degraded
        </p>
      </CardContent>
    </Card>

    <!-- FAILING Card -->
    <Card>
      <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle class="text-sm font-medium">
          FAILING
        </CardTitle>
        <div class="h-2 w-2 rounded-full bg-red-500" />
      </CardHeader>
      <CardContent>
        <div class="text-2xl font-bold text-red-600">
          {{ loading ? '...' : failing }}
        </div>
        <p class="text-xs text-muted-foreground">
          Checks failing
        </p>
      </CardContent>
    </Card>
  </div>
</template>


