<script setup lang="ts">
interface Props {
  status: 'operational' | 'degraded' | 'major_outage'
}

const props = defineProps<Props>()

const statusConfig = {
  operational: {
    label: 'All Systems Operational',
    color: 'bg-green-500',
    textColor: 'text-green-600',
    icon: 'lucide:check-circle-2',
  },
  degraded: {
    label: 'Degraded Performance',
    color: 'bg-yellow-500',
    textColor: 'text-yellow-600',
    icon: 'lucide:alert-triangle',
  },
  major_outage: {
    label: 'Major Outage',
    color: 'bg-red-500',
    textColor: 'text-red-600',
    icon: 'lucide:alert-circle',
  },
}

const config = statusConfig[props.status]
</script>

<template>
  <Card :class="config.color">
    <CardContent class="flex items-center gap-4 p-6">
      <Icon :name="config.icon" class="h-8 w-8 text-white" />
      <div>
        <h2 class="text-2xl font-bold text-white">
          {{ config.label }}
        </h2>
        <p class="text-white/80">
          All services are running normally
        </p>
      </div>
    </CardContent>
  </Card>
</template>


