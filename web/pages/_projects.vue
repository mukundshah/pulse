<script setup lang="ts">
useHead({
  title: 'Projects',
})

const projects = [
  {
    id: 1,
    name: 'Production API',
    description: 'Main production API monitoring',
    status: 'operational',
    checksCount: 12,
    activeAlerts: 2,
    uptime: '99.9%',
    lastCheck: '2 minutes ago',
  },
  {
    id: 2,
    name: 'Staging Environment',
    description: 'Staging environment monitoring',
    status: 'warning',
    checksCount: 8,
    activeAlerts: 1,
    uptime: '98.5%',
    lastCheck: '5 minutes ago',
  },
  {
    id: 3,
    name: 'Development',
    description: 'Development and testing',
    status: 'operational',
    checksCount: 5,
    activeAlerts: 0,
    uptime: '100%',
    lastCheck: '1 minute ago',
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
    default:
      return 'bg-muted-foreground'
  }
}

const getStatusLabel = (status: string) => {
  switch (status) {
    case 'operational':
      return 'Operational'
    case 'warning':
      return 'Warning'
    case 'critical':
      return 'Critical'
    default:
      return 'Unknown'
  }
}
</script>

<template>
  <div class="flex flex-1 flex-col gap-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-semibold tracking-tight">
          Projects
        </h1>
        <p class="text-muted-foreground">
          Manage your monitoring projects
        </p>
      </div>
      <Button>
        <Icon class="mr-2 h-4 w-4" name="lucide:plus" />
        New Project
      </Button>
    </div>

    <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
      <Card
        v-for="project in projects"
        :key="project.id"
        class="cursor-pointer transition-colors hover:bg-muted/50"
      >
        <NuxtLink :to="`/projects/${project.id}`">
          <CardHeader>
            <div class="flex items-start justify-between">
              <div class="flex-1">
                <CardTitle class="mb-1">{{ project.name }}</CardTitle>
                <CardDescription>
                  {{ project.description }}
                </CardDescription>
              </div>
              <div class="flex items-center gap-2">
                <div
                  class="h-2 w-2 rounded-full"
                  :class="getStatusColor(project.status)"
                ></div>
              </div>
            </div>
          </CardHeader>
          <CardContent>
            <div class="space-y-3">
              <div class="flex items-center justify-between text-sm">
                <span class="text-muted-foreground">Checks</span>
                <span class="font-medium">{{ project.checksCount }}</span>
              </div>
              <div class="flex items-center justify-between text-sm">
                <span class="text-muted-foreground">Active Alerts</span>
                <span
                  class="font-medium"
                  :class="project.activeAlerts > 0 ? 'text-yellow-500' : ''"
                >
                  {{ project.activeAlerts }}
                </span>
              </div>
              <div class="flex items-center justify-between text-sm">
                <span class="text-muted-foreground">Uptime</span>
                <span class="font-medium">{{ project.uptime }}</span>
              </div>
              <div class="flex items-center justify-between text-sm">
                <span class="text-muted-foreground">Last Check</span>
                <span class="text-muted-foreground">{{ project.lastCheck }}</span>
              </div>
            </div>
          </CardContent>
        </NuxtLink>
      </Card>
    </div>
  </div>
</template>
