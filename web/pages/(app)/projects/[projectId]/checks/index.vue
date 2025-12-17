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
      { id: 'tcp', name: 'TCP', icon: 'lucide:network', implemented: false },
      { id: 'dns', name: 'DNS', icon: 'lucide:server', implemented: false },
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

useHead({
  title: `Checks for ${project.value?.name}`,
})
</script>

<template>
  <div class="flex flex-1 flex-col gap-6 p-4 md:p-6">
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

    <Table>
      <TableHeader>
        <TableRow>
          <TableHead class="p-0 w-px">
            <span class="sr-only">Link</span>
          </TableHead>
          <TableHead class="pl-9 w-[calc(45%-64px)] min-w-[200px] text-left">
            Name
          </TableHead>
          <TableHead class="w-[8%] min-w-[80px] text-center">
            Type
          </TableHead>
          <TableHead class="w-[12%] min-w-[240px] text-center">
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
            <span class="text-sm">{{ check.avg_response_time_24h_ms }} ms</span>
          </TableCell>
          <TableCell class="text-center">
            <span class="text-sm">{{ check.p95_response_time_24h_ms }} ms</span>
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
                <DropdownMenuItem>
                  <Icon class="mr-2 h-4 w-4" name="lucide:edit" />
                  Edit
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
