<script setup lang="ts">
const colorMode = useColorMode()

const THEME_ICONS = {
  light: 'lucide:sun',
  dark: 'lucide:moon',
  system: 'lucide:monitor',
} as const

const { isAuthenticated, logout } = useAuth()

const handleLogout = async () => {
  await logout()
}
</script>

<template>
  <header class="sticky top-0 z-50 w-full border-b border-border/40 bg-background/95 backdrop-blur supports-backdrop-filter:bg-background/60">
    <div class="bg-foreground text-background">
      <div class="container py-2">
        <p class="text-center text-xs tracking-wide font-medium flex items-center justify-center gap-2">
          <Icon class="h-4 w-4" name="lucide:alert-triangle" /> The contents shown are dummy data and for illustration purposes only.
        </p>
      </div>
    </div>
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
            <Button class="-mr-2" size="icon" variant="ghost">
              <Icon class="h-4 w-4" :name="THEME_ICONS[colorMode.preference as keyof typeof THEME_ICONS]" />
              <span class="sr-only">Toggle theme</span>
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuItem @click="colorMode.preference = 'light'">
              <Icon class="h-4 w-4" :name="THEME_ICONS.light" />
              Light
            </DropdownMenuItem>
            <DropdownMenuItem @click="colorMode.preference = 'dark'">
              <Icon class="h-4 w-4" :name="THEME_ICONS.dark" />
              Dark
            </DropdownMenuItem>
            <DropdownMenuItem @click="colorMode.preference = 'system'">
              <Icon class="h-4 w-4" :name="THEME_ICONS.system" />
              System
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
        <template v-if="!isAuthenticated">
        <Button as-child size="sm" variant="ghost">
          <NuxtLink to="/auth/login">
            Sign in
          </NuxtLink>
        </Button>
        <Button as-child size="sm">
          <NuxtLink to="/auth/register">
            Get started
          </NuxtLink>
        </Button>
        </template>
        <template v-else>
          <Button size="sm" @click="handleLogout" variant="ghost">
            Logout
          </Button>
          <Button as-child size="sm" variant="outline">
            <NuxtLink to="/dashboard">
              Dashboard
            </NuxtLink>
          </Button>
        </template>
      </div>
    </div>
  </header>
</template>
