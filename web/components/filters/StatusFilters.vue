<script setup lang="ts">
interface Props {
  selectedStatuses?: string[]
  selectedLocation?: string
  locations?: Array<{ id: string; name: string }>
}

const props = withDefaults(defineProps<Props>(), {
  selectedStatuses: () => [],
  selectedLocation: undefined,
  locations: () => [],
})

const emit = defineEmits<{
  'update:selectedStatuses': [value: string[]]
  'update:selectedLocation': [value: string | undefined]
}>()

const statuses = [
  { value: 'passed', label: 'Passed', color: 'bg-green-500' },
  { value: 'failed', label: 'Failed', color: 'bg-red-500' },
  { value: 'degraded', label: 'Degraded', color: 'bg-yellow-500' },
  { value: 'retries', label: 'Has retries', color: 'bg-blue-500' },
]

const selectedStatuses = computed({
  get: () => props.selectedStatuses || [],
  set: (value) => emit('update:selectedStatuses', value),
})

const selectedLocation = computed({
  get: () => props.selectedLocation,
  set: (value) => emit('update:selectedLocation', value),
})

const toggleStatus = (status: string) => {
  const current = selectedStatuses.value
  if (current.includes(status)) {
    selectedStatuses.value = current.filter(s => s !== status)
  } else {
    selectedStatuses.value = [...current, status]
  }
}
</script>

<template>
  <div class="flex flex-wrap items-center gap-2">
    <Button
      v-for="status in statuses"
      :key="status.value"
      :variant="selectedStatuses.includes(status.value) ? 'default' : 'outline'"
      size="sm"
      @click="toggleStatus(status.value)"
    >
      <div
        class="mr-2 h-2 w-2 rounded-full"
        :class="status.color"
      />
      {{ status.label }}
    </Button>

    <Select
      v-if="locations.length > 0"
      :model-value="selectedLocation"
      @update:model-value="selectedLocation = $event"
    >
      <SelectTrigger class="w-[180px]">
        <SelectValue placeholder="Location" />
      </SelectTrigger>
      <SelectContent>
        <SelectItem value="">
          All Locations
        </SelectItem>
        <SelectItem
          v-for="location in locations"
          :key="location.id"
          :value="location.id"
        >
          {{ location.name }}
        </SelectItem>
      </SelectContent>
    </Select>
  </div>
</template>


