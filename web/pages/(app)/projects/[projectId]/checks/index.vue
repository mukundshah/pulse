<script setup lang="ts">
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

const { data: project } = await usePulseAPI('/v1/projects/{projectId}', {
  path: {
    projectId,
  },
})

const { data: checks, pending: checksLoading } = useLazyPulseAPI('/v1/projects/{projectId}/checks', {
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
          <TableHead class="p-0">
            <span class="sr-only">Link</span>
          </TableHead>
          <TableHead class="pl-14">
            Name
          </TableHead>
          <TableHead>Type</TableHead>
          <TableHead>Last Results</TableHead>
          <TableHead>24h</TableHead>
          <TableHead>7d</TableHead>
          <TableHead>Avg</TableHead>
          <TableHead>P95</TableHead>
          <TableHead>ΔT</TableHead>
          <TableHead class="text-right">
            Actions
          </TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <template v-if="checksLoading">
          <TableRow v-for="i in 5" :key="i">
            <TableCell colspan="10">
              <div class="flex items-center gap-2">
                <Spinner class="h-4 w-4" />
                <span class="text-sm text-muted-foreground">Loading...</span>
              </div>
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
              <StatusBadge :status="check.status" />
              <div>
                <div class="font-medium">
                  {{ check.name }}
                </div>
                <div class="text-xs text-muted-foreground">
                  <span class="sr-only">Last run at:</span>
                  <NuxtTime v-if="check.last_run_at" :datetime="check.last_run_at" />
                  <span v-else>Never ran</span>
                </div>
              </div>
            </div>
          </TableCell>
          <TableCell>
            <Badge variant="outline">
              {{ check.type.toUpperCase() }}
            </Badge>
          </TableCell>
          <TableCell>
            sparkline
          </TableCell>
          <TableCell>
            <!-- <span class="text-sm">{{ formatUptime(check) }}</span> -->
          </TableCell>
          <TableCell>
            <!-- <span class="text-sm">{{ formatUptime(check) }}</span> -->
          </TableCell>
          <TableCell>
            <!-- <span class="text-sm">{{ formatResponseTime(check) }}</span> -->
          </TableCell>
          <TableCell>
            <!-- <span class="text-sm">{{ formatResponseTime(check) }}</span> -->
          </TableCell>
          <TableCell>
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
