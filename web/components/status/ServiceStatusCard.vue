<script setup lang="ts">
interface Props {
  project: {
    id: string
    name: string
  }
}

defineProps<Props>()

// Mock status for now
const status = 'operational'
</script>

<template>
  <Card>
    <CardContent class="p-4">
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-2">
          <div
            class="h-2 w-2 rounded-full"
            :class="{
              'bg-green-500': status === 'operational',
              'bg-yellow-500': status === 'degraded',
              'bg-red-500': status === 'major_outage',
            }"
          />
          <span class="font-medium">{{ project.name }}</span>
        </div>
        <Badge variant="outline">
          {{ status }}
        </Badge>
      </div>
    </CardContent>
  </Card>
</template>


