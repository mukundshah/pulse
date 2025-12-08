<script setup lang="ts">
interface Props {
  location: string
  size?: 'sm' | 'md' | 'lg'
}

const props = withDefaults(defineProps<Props>(), {
  size: 'md',
})

// Map location codes to flag emojis
const locationFlags: Record<string, string> = {
  'us-east-1': '🇺🇸',
  'us-west-1': '🇺🇸',
  'us-west-2': '🇺🇸',
  'eu-west-1': '🇮🇪',
  'eu-central-1': '🇩🇪',
  'ap-southeast-1': '🇸🇬',
  'ap-southeast-2': '🇦🇺',
  'ap-northeast-1': '🇯🇵',
  'sa-east-1': '🇧🇷',
  'ca-central-1': '🇨🇦',
  'eu-north-1': '🇸🇪',
  'ap-south-1': '🇮🇳',
  'au-syd': '🇦🇺',
  'au-mel': '🇦🇺',
  'id-jkt': '🇮🇩',
  'nz-akl': '🇳🇿',
  'sg-sin': '🇸🇬',
}

const flag = computed(() => {
  const code = props.location.toLowerCase()
  return locationFlags[code] || '🌍'
})

const sizeClasses = {
  sm: 'text-sm',
  md: 'text-base',
  lg: 'text-lg',
}
</script>

<template>
  <span
    class="inline-flex items-center"
    :class="sizeClasses[size]"
    :title="location"
  >
    {{ flag }}
  </span>
</template>

