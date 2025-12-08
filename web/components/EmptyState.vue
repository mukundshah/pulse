<script setup lang="ts">
interface Props {
  title: string
  description?: string
  icon?: string
  actionLabel?: string
  action?: () => void
}

withDefaults(defineProps<Props>(), {
  description: undefined,
  icon: 'lucide:inbox',
  actionLabel: undefined,
  action: undefined,
})
</script>

<template>
  <div class="flex flex-col items-center justify-center py-12 text-center">
    <div class="mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-muted">
      <Icon :name="icon" class="h-8 w-8 text-muted-foreground" />
    </div>
    <h3 class="mb-2 text-lg font-semibold">
      {{ title }}
    </h3>
    <p
      v-if="description"
      class="mb-6 max-w-sm text-sm text-muted-foreground"
    >
      {{ description }}
    </p>
    <Button
      v-if="action && actionLabel"
      @click="action"
    >
      {{ actionLabel }}
    </Button>
  </div>
</template>

