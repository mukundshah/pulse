<script setup lang="ts">
useHead({ titleTemplate: '%siteName %separator %s' })
</script>

<template>
  <div class="flex min-h-screen flex-col bg-background">
    <MarketingHeader />
    <main class="flex-1">
      <slot></slot>
    </main>
    <MarketingFooter />
  </div>
</template>
