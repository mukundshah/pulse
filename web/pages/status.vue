<script setup lang="ts">
import StatusBanner from '@/components/status/StatusBanner.vue'
import ServiceStatusCard from '@/components/status/ServiceStatusCard.vue'

useHead({
  title: 'Status',
})

// Load all checks to determine overall status
const projectsResponse = useAPI('/projects', {
  protected: true,
})
const projects = computed(() => {
  const data = projectsResponse.data.value
  return Array.isArray(data) ? data : []
})

// Calculate overall status
const overallStatus = computed(() => {
  // This would need to aggregate from all checks
  // For now, return operational
  return 'operational'
})
</script>

<template>
  <div class="flex flex-1 flex-col gap-6">
    <div>
      <h1 class="text-3xl font-semibold tracking-tight">
        Status
      </h1>
      <p class="text-muted-foreground">
        System status and incident history
      </p>
    </div>

    <!-- Status Banner -->
    <StatusBanner :status="overallStatus" />

    <!-- Services Grid -->
    <Card>
      <CardHeader>
        <CardTitle>Services</CardTitle>
        <CardDescription>
          Status of all monitored services
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          <ServiceStatusCard
            v-for="project in projects"
            :key="project.id"
            :project="project"
          />
        </div>
      </CardContent>
    </Card>

    <!-- Recent Updates Timeline -->
    <Card>
      <CardHeader>
        <CardTitle>Recent Updates</CardTitle>
        <CardDescription>
          Past incidents and maintenance
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div class="space-y-4">
          <div class="flex items-start gap-4 border-l-2 border-border pl-4">
            <div class="flex h-8 w-8 items-center justify-center rounded-full bg-green-500 text-white">
              <Icon class="h-4 w-4" name="lucide:check" />
            </div>
            <div class="flex-1">
              <p class="font-medium">All systems operational</p>
              <p class="text-sm text-muted-foreground">
                {{ new Date().toLocaleString() }}
              </p>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
