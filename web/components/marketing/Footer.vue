<template>
  <footer class="border-t border-border bg-muted/30 py-12">
    <div class="container">
      <div class="flex flex-col items-center justify-between gap-4 md:flex-row">
        <div class="flex items-center gap-2">
          <NuxtLink class="flex items-center gap-2" to="/">
            <div class="flex h-6 w-6 items-center justify-center rounded-md bg-foreground text-background text-xs font-semibold">
              P
            </div>
            <span class="text-sm font-medium">Pulse</span>
          </NuxtLink>
        </div>
        <p class="text-sm text-muted-foreground">
          &copy; {{ new Date().getFullYear() }} Pulse. All rights reserved.
        </p>
      </div>
    </div>
  </footer>
</template>
