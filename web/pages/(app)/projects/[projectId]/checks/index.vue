<script setup lang="ts">
import SparklineChart from './@components/SparklineChart.vue'
import StatusBadge from './@components/StatusBadge.vue'

const route = useRoute()
const { projectId } = route.params as { projectId: string }

// Check type definitions
const checkTypes = [
  {
    id: 'synthetic',
    name: 'Synthetic',
    checks: [
      { id: 'browser', name: 'Browser', icon: 'lucide:monitor', implemented: false },
    ],
  },
  {
    id: 'uptime',
    name: 'Uptime',
    checks: [
      { id: 'http', name: 'HTTP', icon: 'lucide:globe', implemented: true },
      { id: 'tcp', name: 'TCP', icon: 'lucide:network', implemented: true },
      { id: 'dns', name: 'DNS', icon: 'lucide:server', implemented: true },
    ],
  },
]

const { data: project } = await usePulseAPI('/internal/projects/{projectId}', {
  path: {
    projectId,
  },
})

const { data: checks, pending: checksLoading } = useLazyPulseAPI('/internal/projects/{projectId}/checks', {
  path: {
    projectId,
  },
})

const { data: counts, pending: statusCountsLoading } = useLazyPulseAPI('/internal/projects/{projectId}/checks/status/counts', {
  path: {
    projectId,
  },
})

useHead({
  title: `Checks for ${project.value?.name}`,
})

useLayoutContext({
  breadcrumbOverrides: computed(() => [
    undefined, // Root
    undefined, // Projects
    {
      label: project?.value?.name || 'Project',
      to: `/projects/${projectId}/checks`,
    }, // Project
    false, // Checks
  ]),
})
</script>

<template>
  <div class="@container/main flex flex-1 flex-col gap-6 p-4 md:p-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-semibold tracking-tight">
          {{ project?.name }}
        </h1>
      </div>
      <DropdownMenu>
        <DropdownMenuTrigger>
          <Button size="sm" variant="outline">
            <Icon class="mr-2 h-4 w-4" name="lucide:plus" />
            New Check
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end" class="w-56">
          <template v-for="(type, idx) in checkTypes" :key="type.id">
            <DropdownMenuGroup>
              <DropdownMenuLabel>{{ type.name }}</DropdownMenuLabel>
              <DropdownMenuItem
                v-for="check in type.checks"
                :key="check.id"
                as-child
                :disabled="!check.implemented"
              >
                <NuxtLink
                  class="flex w-full items-center justify-between"
                  :to="`/projects/${projectId}/checks/new/${check.id}`"
                >
                  <div class="flex items-center gap-2">
                    <Icon class="h-4 w-4" :name="check.icon" />
                    <span>{{ check.name }}</span>
                  </div>
                  <Badge v-if="!check.implemented" class="text-xs" variant="secondary">
                    WIP
                  </Badge>
                </NuxtLink>
              </DropdownMenuItem>
            </DropdownMenuGroup>
            <DropdownMenuSeparator v-if="idx < checkTypes.length - 1" />
          </template>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>

    <div class="grid grid-cols-1 gap-4 @xl/main:grid-cols-3 py-4">
      <Card class="@container/card from-green-200/30 to-green-200/5 dark:from-green-900/30 dark:to-green-900/5 bg-linear-to-t shadow-none border-green-500/20 dark:border-green-900/30 py-4 *:data-[slot=card-header]:px-4">
        <CardHeader>
          <CardDescription class="text-green-600 dark:text-green-400">
            Passing
          </CardDescription>
          <Skeleton v-if="statusCountsLoading" class="h-9 w-24" />
          <CardTitle v-else class="font-semibold tabular-nums text-3xl text-green-600 dark:text-green-400">
            {{ counts?.passing }}
          </CardTitle>
          <CardAction>
            <Icon class="size-4 text-green-600 dark:text-green-400" name="lucide:check-circle" />
          </CardAction>
        </CardHeader>
      </Card>
      <Card class="@container/card from-yellow-200/30 to-yellow-200/5 dark:from-yellow-900/30 dark:to-yellow-900/5 bg-linear-to-t shadow-none border-yellow-500/30 dark:border-yellow-900/30 py-4 *:data-[slot=card-header]:px-4">
        <CardHeader>
          <CardDescription class="text-yellow-600 dark:text-yellow-400">
            Degraded
          </CardDescription>
          <Skeleton v-if="statusCountsLoading" class="h-9 w-24" />
          <CardTitle v-else class="font-semibold tabular-nums text-3xl text-yellow-600 dark:text-yellow-400">
            {{ counts?.degraded }}
          </CardTitle>
          <CardAction>
            <Icon class="size-4 text-yellow-600 dark:text-yellow-400" name="lucide:alert-circle" />
          </CardAction>
        </CardHeader>
      </Card>
      <Card class="@container/card from-red-200/30 to-red-200/5 dark:from-red-900/30 dark:to-red-900/5 bg-linear-to-t shadow-none border-red-500/20 dark:border-red-900/30 py-4 *:data-[slot=card-header]:px-4">
        <CardHeader>
          <CardDescription class="text-red-600 dark:text-red-400">
            Failing
          </CardDescription>
          <Skeleton v-if="statusCountsLoading" class="h-9 w-24" />
          <CardTitle v-else class="font-semibold tabular-nums text-3xl text-red-600 dark:text-red-400">
            {{ counts?.failing }}
          </CardTitle>
          <CardAction>
            <Icon class="size-4 text-red-600 dark:text-red-400" name="lucide:x-circle" />
          </CardAction>
        </CardHeader>
      </Card>
    </div>

    <Table>
      <TableHeader>
        <TableRow>
          <TableHead class="p-0 w-px">
            <span class="sr-only">Link</span>
          </TableHead>
          <TableHead class="pl-9 w-[calc(47%-64px)] min-w-[200px] text-left">
            Name
          </TableHead>
          <TableHead class="w-[8%] min-w-[80px] text-center">
            Type
          </TableHead>
          <TableHead class="w-[10%] min-w-[160px] text-center">
            Last Results
          </TableHead>
          <TableHead class="w-[7%] min-w-[70px] text-center">
            24h
          </TableHead>
          <TableHead class="w-[7%] min-w-[70px] text-center">
            7d
          </TableHead>
          <TableHead class="w-[7%] min-w-[70px] text-center">
            Avg
          </TableHead>
          <TableHead class="w-[7%] min-w-[70px] text-center">
            P95
          </TableHead>
          <TableHead class="w-[7%] min-w-[70px] text-center">
            ΔT
          </TableHead>
          <TableHead class="w-16">
            <span class="sr-only">Actions</span>
          </TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <template v-if="checksLoading">
          <TableRow v-for="i in 5" :key="i">
            <TableCell class="p-0" />
            <TableCell>
              <div class="flex items-center gap-4">
                <Skeleton class="h-2.5 w-2.5 rounded-full" />
                <div class="space-y-2">
                  <Skeleton class="h-4 w-32" />
                  <Skeleton class="h-3 w-24" />
                </div>
              </div>
            </TableCell>
            <TableCell>
              <Skeleton class="h-5 w-12 rounded-full" />
            </TableCell>
            <TableCell>
              <Skeleton class="h-16 w-full rounded-md" />
            </TableCell>
            <TableCell>
              <Skeleton class="h-4 w-12" />
            </TableCell>
            <TableCell>
              <Skeleton class="h-4 w-12" />
            </TableCell>
            <TableCell>
              <Skeleton class="h-4 w-12" />
            </TableCell>
            <TableCell>
              <Skeleton class="h-4 w-12" />
            </TableCell>
            <TableCell>
              <Skeleton class="h-4 w-12" />
            </TableCell>
            <TableCell class="text-right">
              <Skeleton class="h-8 w-8 rounded-md" />
            </TableCell>
          </TableRow>
        </template>
        <template v-else-if="checks && checks.length === 0">
          <TableRow>
            <TableCell colspan="10">
              <Empty>
                <EmptyHeader>
                  <EmptyMedia variant="icon">
                    <Icon name="lucide:clipboard-list" />
                  </EmptyMedia>
                  <EmptyTitle>No checks configured</EmptyTitle>
                  <EmptyDescription>Create a check to start monitoring</EmptyDescription>
                </EmptyHeader>
              </Empty>
            </TableCell>
          </TableRow>
        </template>
        <TableRow
          v-for="check in checks || []"
          v-else
          :key="check.id"
          class="relative cursor-pointer hover:bg-muted/50"
        >
          <TableCell class="p-0">
            <NuxtLink class="absolute h-full w-full inset-0" :to="`/projects/${projectId}/checks/${check.id}`" />
          </TableCell>
          <TableCell>
            <div class="flex items-center justify-start gap-4">
              <StatusBadge :status="check.last_status" />
              <div>
                <div class="font-medium">
                  {{ check.name }}
                </div>
                <div class="text-xs text-muted-foreground">
                  <span class="sr-only">Last run at:</span>
                  <NuxtTime
                    v-if="check.last_run"
                    relative
                    title
                    :datetime="check.last_run"
                  />
                  <span v-else>Never ran</span>
                </div>
              </div>
            </div>
          </TableCell>
          <TableCell class="text-center">
            <Badge variant="outline">
              {{ check.type.toUpperCase() }}
            </Badge>
          </TableCell>
          <TableCell>
            <SparklineChart :runs="check.last_24_runs" />
          </TableCell>
          <TableCell class="text-center">
            <FormattedNumber class="text-sm" :options="{ style: 'unit', unit: 'percent' }" :value="check.uptime_24h" />
          </TableCell>
          <TableCell class="text-center">
            <FormattedNumber class="text-sm" :options="{ style: 'unit', unit: 'percent' }" :value="check.uptime_7d" />
          </TableCell>
          <TableCell class="text-center">
            <FormattedNumber class="text-sm" :options="{ style: 'unit', unit: 'millisecond' }" :value="check.avg_response_time_24h_ms" />
          </TableCell>
          <TableCell class="text-center">
            <FormattedNumber class="text-sm" :options="{ style: 'unit', unit: 'millisecond' }" :value="check.p95_response_time_24h_ms" />
          </TableCell>
          <TableCell class="text-center">
            <span class="text-sm">{{ check.interval }}</span>
          </TableCell>
          <TableCell class="text-right">
            <DropdownMenu>
              <DropdownMenuTrigger as-child>
                <Button size="icon" variant="ghost">
                  <Icon name="lucide:more-vertical" />
                  <span class="sr-only">Open menu</span>
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuItem as-child>
                  <NuxtLink :to="`/projects/${projectId}/checks/${check.id}/edit`">
                    <Icon class="mr-2 h-4 w-4" name="lucide:edit" />
                    Edit
                  </NuxtLink>
                </DropdownMenuItem>
                <DropdownMenuItem>
                  <Icon class="mr-2 h-4 w-4" name="lucide:copy" />
                  Duplicate
                </DropdownMenuItem>
                <DropdownMenuSeparator />
                <DropdownMenuItem class="text-destructive">
                  <Icon class="mr-2 h-4 w-4" name="lucide:trash-2" />
                  Delete
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </TableCell>
        </TableRow>
      </TableBody>
    </Table>
  </div>
</template>
