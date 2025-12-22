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
    <div class="container grid h-16 grid-cols-[1fr_auto_1fr] items-center gap-4">
      <div class="flex items-center gap-2">
        <NuxtLink class="flex items-center gap-2" to="/">
          <AppIcon class="text-foreground size-10" />
          <span class="text-lg font-semibold tracking-wider">Pulse</span>
        </NuxtLink>
      </div>
      <nav class="hidden items-center gap-6 md:flex">
        <a class="text-sm text-muted-foreground transition-colors hover:text-foreground" href="/#story">Story</a>
        <a class="text-sm text-muted-foreground transition-colors hover:text-foreground" href="/#features">Features</a>
      </nav>
      <div class="flex items-center justify-end gap-4">
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
          <Button size="sm" variant="ghost" @click="handleLogout">
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
