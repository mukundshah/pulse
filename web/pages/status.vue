<script setup lang="ts">
useHead({
  title: 'Status',
})

const overallStatus = 'operational'

const services = [
  {
    id: 1,
    name: 'API Service',
    status: 'operational',
  },
  {
    id: 2,
    name: 'Database',
    status: 'operational',
  },
  {
    id: 3,
    name: 'Web Server',
    status: 'operational',
  },
  {
    id: 4,
    name: 'Cache Service',
    status: 'warning',
  },
  {
    id: 5,
    name: 'Email Service',
    status: 'operational',
  },
  {
    id: 6,
    name: 'CDN',
    status: 'operational',
  },
]

const updates = [
  {
    id: 1,
    status: 'resolved',
    title: 'Database connection timeout',
    timestamp: 'Jan 15, 2024 at 11:15 AM',
    description: 'Experienced intermittent connection timeouts to the primary database. Issue resolved after failover to secondary node.',
  },
  {
    id: 2,
    status: 'resolved',
    title: 'Increased response times',
    timestamp: 'Jan 10, 2024 at 3:00 PM',
    description: 'API response times increased by 200ms. Identified and resolved performance bottleneck.',
  },
  {
    id: 3,
    status: 'completed',
    title: 'Scheduled maintenance',
    timestamp: 'Jan 5, 2024 at 4:00 AM',
    description: 'Routine database maintenance and optimization completed successfully.',
  },
]

const getStatusColor = (status: string) => {
  switch (status) {
    case 'operational':
      return 'bg-green-500'
    case 'warning':
      return 'bg-yellow-500'
    case 'critical':
      return 'bg-red-500'
    case 'maintenance':
      return 'bg-blue-500'
    default:
      return 'bg-muted-foreground'
  }
}

const getStatusLabel = (status: string) => {
  switch (status) {
    case 'operational':
      return 'Operational'
    case 'warning':
      return 'Degraded Performance'
    case 'critical':
      return 'Major Outage'
    case 'maintenance':
      return 'Maintenance'
    case 'resolved':
      return 'Resolved'
    case 'completed':
      return 'Completed'
    default:
      return 'Unknown'
  }
}
</script>

<template>
  <div class="flex flex-1 flex-col">
    <!-- Main Status Banner -->
    <div class="border-b border-border bg-muted/30 py-12">
      <div class="container mx-auto max-w-4xl px-4">
        <div class="flex flex-col items-center justify-center text-center">
          <div class="mb-4 flex items-center gap-3">
            <div
              class="h-3 w-3 rounded-full"
              :class="getStatusColor(overallStatus)"
            ></div>
            <h1 class="text-4xl font-semibold tracking-tight md:text-5xl">
              {{ getStatusLabel(overallStatus) }}
            </h1>
          </div>
          <p class="text-lg text-muted-foreground">
            All systems are operational
          </p>
        </div>
      </div>
    </div>

    <!-- Services Grid -->
    <div class="container mx-auto max-w-4xl px-4 py-8">
      <div class="mb-6">
        <h2 class="mb-2 text-2xl font-semibold">
          Services
        </h2>
        <p class="text-muted-foreground">
          Status of all services and components
        </p>
      </div>

      <div class="grid gap-4 md:grid-cols-2">
        <div
          v-for="service in services"
          :key="service.id"
          class="flex items-center justify-between rounded-lg border border-border p-4"
        >
          <div class="flex items-center gap-3">
            <div
              class="h-2.5 w-2.5 rounded-full"
              :class="getStatusColor(service.status)"
            ></div>
            <span class="font-medium">{{ service.name }}</span>
          </div>
          <span class="text-sm text-muted-foreground">
            {{ getStatusLabel(service.status) }}
          </span>
        </div>
      </div>
    </div>

    <!-- Updates Timeline -->
    <div class="border-t border-border bg-muted/30">
      <div class="container mx-auto max-w-4xl px-4 py-8">
        <div class="mb-6">
          <h2 class="mb-2 text-2xl font-semibold">
            Recent Updates
          </h2>
          <p class="text-muted-foreground">
            Past incidents and maintenance windows
          </p>
        </div>

        <div class="space-y-6">
          <div
            v-for="update in updates"
            :key="update.id"
            class="flex gap-4"
          >
            <div class="flex flex-col items-center">
              <div
                class="h-2 w-2 rounded-full"
                :class="getStatusColor(update.status === 'resolved' || update.status === 'completed' ? 'operational' : update.status)"
              ></div>
              <div class="h-full w-px bg-border mt-2"></div>
            </div>
            <div class="flex-1 pb-6">
              <div class="mb-1 flex items-center gap-2">
                <span class="font-semibold">{{ update.title }}</span>
                <Badge variant="outline">
                  {{ getStatusLabel(update.status) }}
                </Badge>
              </div>
              <p class="mb-2 text-sm text-muted-foreground">
                {{ update.description }}
              </p>
              <p class="text-xs text-muted-foreground">
                {{ update.timestamp }}
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
