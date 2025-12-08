<script setup lang="ts">
interface Props {
  modelValue?: string
  placeholder?: string
  statusFilter?: string
  checkTypeFilter?: string
  tagsFilter?: string[]
  statusOptions?: Array<{ value: string; label: string }>
  checkTypeOptions?: Array<{ value: string; label: string }>
}

const props = withDefaults(defineProps<Props>(), {
  placeholder: 'Search by name, request url...',
  statusFilter: undefined,
  checkTypeFilter: undefined,
  tagsFilter: () => [],
  statusOptions: () => [
    { value: 'all', label: 'All Status' },
    { value: 'passing', label: 'Passing' },
    { value: 'degraded', label: 'Degraded' },
    { value: 'failing', label: 'Failing' },
  ],
  checkTypeOptions: () => [
    { value: 'all', label: 'All Types' },
    { value: 'http', label: 'HTTP' },
    { value: 'tcp', label: 'TCP' },
    { value: 'dns', label: 'DNS' },
    { value: 'browser', label: 'Browser' },
    { value: 'heartbeat', label: 'Heartbeat' },
  ],
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
  'update:statusFilter': [value: string]
  'update:checkTypeFilter': [value: string]
  'update:tagsFilter': [value: string[]]
}>()

const searchQuery = computed({
  get: () => props.modelValue || '',
  set: (value) => emit('update:modelValue', value),
})

const status = computed({
  get: () => props.statusFilter || 'all',
  set: (value) => emit('update:statusFilter', value),
})

const checkType = computed({
  get: () => props.checkTypeFilter || 'all',
  set: (value) => emit('update:checkTypeFilter', value),
})
</script>

<template>
  <div class="flex flex-col gap-4 md:flex-row md:items-center">
    <!-- Search Input -->
    <div class="relative flex-1">
      <Icon
        name="lucide:search"
        class="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground"
      />
      <Input
        :model-value="searchQuery"
        class="pl-9"
        :placeholder="placeholder"
        @update:model-value="searchQuery = $event"
      />
    </div>

    <!-- Filters -->
    <div class="flex items-center gap-2">
      <Select
        :model-value="status"
        @update:model-value="status = $event"
      >
        <SelectTrigger class="w-[140px]">
          <SelectValue placeholder="Status" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem
            v-for="option in statusOptions"
            :key="option.value"
            :value="option.value"
          >
            {{ option.label }}
          </SelectItem>
        </SelectContent>
      </Select>

      <Select
        :model-value="checkType"
        @update:model-value="checkType = $event"
      >
        <SelectTrigger class="w-[140px]">
          <SelectValue placeholder="Check type" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem
            v-for="option in checkTypeOptions"
            :key="option.value"
            :value="option.value"
          >
            {{ option.label }}
          </SelectItem>
        </SelectContent>
      </Select>

      <Button variant="outline" size="sm">
        <Icon name="lucide:tag" class="mr-2 h-4 w-4" />
        Tags
      </Button>
    </div>
  </div>
</template>


