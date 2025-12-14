<script setup lang="ts">
const route = useRoute()
const { projectId } = route.params as { projectId: string }

await navigateTo(`/projects/${projectId}/checks`)
</script>

<template>
  <div></div>
</template>
