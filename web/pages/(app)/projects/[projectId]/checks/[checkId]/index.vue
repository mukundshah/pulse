<script setup lang="ts">
import StatusBadge from '../@components/StatusBadge.vue'
import PerformanceChart from './@components/PerformanceChart.vue'
import RunResults from './@components/RunResults.vue'
import UptimeChart from './@components/UptimeChart.vue'

const route = useRoute()
const { checkId } = route.params as { checkId: string }

const { data: check } = await usePulseAPI('/v1/checks/{checkId}', {
  path: {
    checkId,
  },
})

useHead({
  title: `Check ${check.value?.name}`,
})
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
            <StatusBadge status="passing" />
          </div>

          <p class="text-sm text-muted-foreground font-mono">
            {{ check?.method }} {{ check?.host }} {{ check?.path }}
          </p>

          <div>
            <Badge class="text-xs" variant="outline">
              Last updated <NuxtTime :datetime="check?.updated_at ?? new Date()" />
            </Badge>
          </div>
        </div>

        <div class="flex items-center gap-2.5">
          <Button size="sm" variant="outline">
            Edit
          </Button>

          <Button size="sm">
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
            <UptimeChart />
          </div>
        </CardContent>
      </Card>

      <!-- Alerts Table -->
      <!-- <Card class="bg-card border-border p-6 mb-6">
        <h2 class="text-base font-semibold text-foreground mb-4">
          Alerts
        </h2>
        <AlertsList />
      </Card> -->

      <!-- Locations -->
      <!-- <Card class="bg-card border-border p-6 mb-6">
        <h2 class="text-base font-semibold text-foreground mb-4">
          Locations
        </h2>
        <LocationsTable />
      </Card> -->

      <!-- Performance Chart -->
      <Card>
        <CardHeader>
          <CardTitle>
            Performance
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div>
            <PerformanceChart />
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
        <RunResults class="flex-1 overflow-y-auto" :check-id="checkId" />
      </div>
    </aside>
  </div>
</template>
