<script setup lang="ts">
interface Props {
  modelValue?: string
}

const props = defineProps<Props>()
const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const timeRanges = [
  { value: 'custom', label: 'Custom' },
  { value: 'today', label: 'Today' },
  { value: '1h', label: '1hr' },
  { value: '3h', label: '3hr' },
  { value: '24h', label: '24hr' },
  { value: '7d', label: '7d' },
]

const selectedRange = computed({
  get: () => props.modelValue || '24h',
  set: (value) => emit('update:modelValue', value),
})
</script>

<template>
  <div class="flex items-center gap-2">
    <Button
      v-for="range in timeRanges"
      :key="range.value"
      :variant="selectedRange === range.value ? 'default' : 'outline'"
      size="sm"
      @click="selectedRange = range.value"
    >
      {{ range.label }}
    </Button>
  </div>
</template>


