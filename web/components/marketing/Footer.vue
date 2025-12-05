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
        <nav class="flex items-center gap-6">
          <NuxtLink class="text-sm text-muted-foreground transition-colors hover:text-foreground" to="/terms">
            Terms
          </NuxtLink>
          <NuxtLink class="text-sm text-muted-foreground transition-colors hover:text-foreground" to="/privacy">
            Privacy
          </NuxtLink>
          <NuxtLink class="text-sm text-muted-foreground transition-colors hover:text-foreground" to="/contact">
            Contact
          </NuxtLink>
        </nav>
        <p class="text-sm text-muted-foreground">
          &copy; {{ new Date().getFullYear() }} Pulse. All rights reserved.
        </p>
      </div>
    </div>
  </footer>
</template>
