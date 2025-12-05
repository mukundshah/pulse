<script setup lang="ts">
const colorMode = useColorMode()
</script>

<template>
  <header class="sticky top-0 z-50 w-full border-b border-border/40 bg-background/95 backdrop-blur supports-backdrop-filter:bg-background/60">
    <div class="container flex h-16 items-center justify-between">
      <div class="flex items-center gap-2">
        <NuxtLink class="flex items-center gap-2" to="/">
          <div class="flex h-8 w-8 items-center justify-center rounded-md bg-foreground text-background font-semibold">
            P
          </div>
          <span class="text-lg font-semibold">Pulse</span>
        </NuxtLink>
      </div>
      <nav class="hidden items-center gap-6 md:flex">
        <a class="text-sm text-muted-foreground transition-colors hover:text-foreground" href="/#features">Features</a>
        <a class="text-sm text-muted-foreground transition-colors hover:text-foreground" href="/#how-it-works">How it works</a>
        <a class="text-sm text-muted-foreground transition-colors hover:text-foreground" href="/#integrations">Integrations</a>
      </nav>
      <div class="flex items-center gap-4">
        <DropdownMenu>
          <DropdownMenuTrigger as-child>
          <Button size="icon" variant="ghost" class="-mr-2">
            <Icon class="h-[1.2rem] w-[1.2rem] rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0" name="lucide:moon" />
            <Icon class="absolute h-[1.2rem] w-[1.2rem] rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100" name="lucide:sun" />
            <span class="sr-only">Toggle theme</span>
          </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
          <DropdownMenuItem @click="colorMode.preference = 'light'">
            <Icon name="lucide:sun" class="h-4 w-4" />
            Light
          </DropdownMenuItem>
          <DropdownMenuItem @click="colorMode.preference = 'dark'">
            <Icon name="lucide:moon" class="h-4 w-4" />
            Dark
          </DropdownMenuItem>
          <DropdownMenuItem @click="colorMode.preference = 'system'">
            <Icon name="lucide:monitor" class="h-4 w-4" />
            System
          </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
        <Button as-child size="sm" variant="ghost">
          <NuxtLink to="/auth/login">
            Sign in
          </NuxtLink>
        </Button>
        <Button as-child size="sm">
          <NuxtLink to="/auth/login">
            Get started
          </NuxtLink>
        </Button>
      </div>
    </div>
  </header>
</template>
