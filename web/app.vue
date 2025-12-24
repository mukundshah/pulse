<script setup lang="ts">
import { Toaster } from '@/components/ui/sonner'
import { FAVICONS } from '@/constants/favicon'
import 'vue-sonner/style.css'

useHead({
  link: FAVICONS.default,
})
</script>

<template>
  <NuxtRouteAnnouncer />
  <NuxtLoadingIndicator color="#000" />
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>
  <Toaster />
</template>
