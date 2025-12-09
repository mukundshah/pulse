<script setup lang="ts">
const colorMode = useColorMode()

const THEME_ICONS = {
  light: 'lucide:sun',
  dark: 'lucide:moon',
  system: 'lucide:monitor',
} as const

useHead({ titleTemplate: '%siteName %separator %s' })
</script>

<template>
  <div class="magicpattern">
    <header class="py-2 px-2 sm:px-4 flex items-center justify-end">
      <DropdownMenu>
        <DropdownMenuTrigger as-child>
          <Button size="icon" variant="ghost">
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
    </header>
    <slot></slot>
  </div>
</template>

<style>
.magicpattern {
  width: 100%;
  height: 100%;
  background-size: cover;
  background-position: center center;
  background-repeat: repeat;
  background-image: url('assets/svgs/magicpattern.svg');
}

.dark .magicpattern {
  background-image: url('assets/svgs/magicpattern-dark.svg');
}
</style>
