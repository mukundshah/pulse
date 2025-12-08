<script setup lang="ts">
import ChecksTable from '@/components/checks/ChecksTable.vue'

const route = useRoute()
const projectId = route.params.id?.toString() ?? ''

useHead({
  title: 'Project',
})

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

// Check type definitions
const checkTypes = [
  { id: 'http', name: 'HTTP', icon: 'lucide:globe', implemented: true },
  { id: 'tcp', name: 'TCP', icon: 'lucide:network', implemented: false },
  { id: 'dns', name: 'DNS', icon: 'lucide:server', implemented: false },
  { id: 'browser', name: 'Browser', icon: 'lucide:monitor', implemented: false },
  { id: 'heartbeat', name: 'Heartbeat', icon: 'lucide:heart', implemented: true },
]

// Group check types
const syntheticTesting = computed(() => checkTypes.filter(type => type.id === 'browser'))
const uptimeMonitors = computed(() => checkTypes.filter(type => type.id !== 'browser'))
</script>

<template>
  <div class="flex flex-1 flex-col gap-6">
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
          <DropdownMenuLabel>Synthetic Testing</DropdownMenuLabel>
          <DropdownMenuGroup>
            <DropdownMenuItem
              v-for="type in syntheticTesting"
              :key="type.id"
              as-child
              :disabled="!type.implemented"
            >
              <NuxtLink
                class="flex w-full items-center justify-between"
                :to="`/projects/${projectId}/checks/new/${type.id}`"
              >
                <div class="flex items-center gap-2">
                  <Icon class="h-4 w-4" :name="type.icon" />
                  <span>{{ type.name }}</span>
                </div>
                <Badge v-if="!type.implemented" class="text-xs" variant="secondary">
                  WIP
                </Badge>
              </NuxtLink>
            </DropdownMenuItem>
          </DropdownMenuGroup>
          <DropdownMenuSeparator />
          <DropdownMenuLabel>Uptime Monitors</DropdownMenuLabel>
          <DropdownMenuGroup>
            <DropdownMenuItem
              v-for="type in uptimeMonitors"
              :key="type.id"
              as-child
              :disabled="!type.implemented"
            >
              <NuxtLink
                class="flex w-full items-center justify-between"
                :to="`/projects/${projectId}/checks/new/${type.id}`"
              >
                <div class="flex items-center gap-2">
                  <Icon class="h-4 w-4" :name="type.icon" />
                  <span>{{ type.name }}</span>
                </div>
                <Badge v-if="!type.implemented" class="text-xs" variant="secondary">
                  WIP
                </Badge>
              </NuxtLink>
            </DropdownMenuItem>
          </DropdownMenuGroup>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>

    <!-- Status Overview Cards -->
    <!-- <StatusOverviewCards
      :degraded="statusCounts.degraded"
      :failing="statusCounts.failing"
      :loading="checksResponse.pending.value"
      :passing="statusCounts.passing"
    /> -->

    <!-- Search and Filter Bar -->
    <!-- <SearchBar
      v-model="searchQuery"
      :check-type-filter="checkTypeFilter"
      :status-filter="statusFilter"
      @update:check-type-filter="checkTypeFilter = $event"
      @update:status-filter="statusFilter = $event"
    /> -->

    <!-- Checks Table -->
    <ChecksTable
      :checks="checks || []"
      :loading="checksLoading"
    />
  </div>
</template>
