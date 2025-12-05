<script setup lang="ts">
import { useMediaQuery } from '@vueuse/core'

useHead({ titleTemplate: '%siteName %separator %s' })

const user = {
  name: 'John Doe',
  email: 'john.doe@example.com',
  avatar: 'https://github.com/shadcn.png',
  initials: 'CN',
}

const isMobile = useMediaQuery('(max-width: 768px)')
</script>

<template>
  <SidebarProvider
    :style="{
      '--sidebar-width': 'calc(var(--spacing) * 64)',
      '--header-height': 'calc(var(--spacing) * 15 + 1px)',
    }"
  >
    <Sidebar class="h-auto border-r" collapsible="none">
      <SidebarHeader class="border-b">
        <NuxtLink class="flex items-center gap-2.5 px-2 py-1.5" to="/dashboard">
          <div class="flex h-8 w-8 items-center justify-center rounded-md bg-foreground text-background font-semibold">
            P
          </div>
          <span class="text-base font-semibold">Pulse</span>
        </NuxtLink>
      </SidebarHeader>
      <SidebarContent>
        <SidebarGroup>
          <SidebarMenu>
            <SidebarMenuItem>
              <SidebarMenuButton as-child>
                <NuxtLink to="/dashboard">
                  <Icon name="lucide:layout-dashboard" />
                  <span>Dashboard</span>
                </NuxtLink>
              </SidebarMenuButton>
            </SidebarMenuItem>
            <SidebarMenuItem>
              <SidebarMenuButton as-child>
                <NuxtLink to="/projects">
                  <Icon name="lucide:folder" />
                  <span>Projects</span>
                </NuxtLink>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarGroup>
        <SidebarGroup>
          <SidebarGroupLabel>
            Project 1
            <SidebarGroupAction>
              <Icon name="lucide:settings" /> <span class="sr-only">Settings</span>
            </SidebarGroupAction>
          </SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenuItem>
              <Collapsible default-open class="group/collapsible">
                <CollapsibleTrigger as-child>
                  <SidebarMenuButton as-child>
                    <NuxtLink to="/projects/1">
                      <div class="flex items-center gap-2">
                        <span>Production API</span>
                        <Icon class="ml-auto transition-transform group-data-[state=open]/collapsible:rotate-90" name="lucide:chevron-right" />
                      </div>
                    </NuxtLink>
                  </SidebarMenuButton>
                </CollapsibleTrigger>
                <CollapsibleContent>
                  <SidebarMenuSub>
                    <SidebarMenuItem>
                      <SidebarMenuButton as-child>
                        <NuxtLink to="/projects/1">
                          <Icon name="lucide:activity" />
                          <span>Checks</span>
                        </NuxtLink>
                      </SidebarMenuButton>
                    </SidebarMenuItem>
                    <SidebarMenuItem>
                      <SidebarMenuButton as-child>
                        <NuxtLink to="/alerts">
                          <Icon name="lucide:bell" />
                          <span>Alerts</span>
                        </NuxtLink>
                      </SidebarMenuButton>
                    </SidebarMenuItem>
                    <SidebarMenuItem>
                      <SidebarMenuButton as-child>
                        <NuxtLink to="/status">
                          <Icon name="lucide:globe" />
                          <span>Status</span>
                        </NuxtLink>
                      </SidebarMenuButton>
                    </SidebarMenuItem>
                    <SidebarMenuItem>
                      <SidebarMenuButton as-child>
                        <NuxtLink to="/projects/1/settings">
                          <Icon name="lucide:settings" />
                          <span>Settings</span>
                        </NuxtLink>
                      </SidebarMenuButton>
                    </SidebarMenuItem>
                  </SidebarMenuSub>
                </CollapsibleContent>
              </Collapsible>
            </SidebarMenuItem>
          </SidebarGroupContent>
        </SidebarGroup>
        <SidebarGroup class="mt-auto">
          <!-- <SidebarGroupLabel>Settings</SidebarGroupLabel> -->
          <SidebarMenu>
            <SidebarMenuItem>
              <SidebarMenuButton as-child>
                <NuxtLink to="/settings">
                  <Icon name="lucide:settings" />
                  <span>Settings</span>
                </NuxtLink>
              </SidebarMenuButton>
            </SidebarMenuItem>
            <SidebarMenuItem>
              <SidebarMenuButton as-child>
                <NuxtLink to="/integrations">
                  <Icon name="lucide:plug" />
                  <span>Integrations</span>
                </NuxtLink>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarGroup>
      </SidebarContent>
      <SidebarFooter>
        <SidebarMenu>
          <SidebarMenuItem>
            <DropdownMenu>
              <DropdownMenuTrigger as-child>
                <SidebarMenuButton
                  class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
                  size="lg"
                >
                  <Avatar class="h-8 w-8 rounded-lg grayscale">
                    <AvatarImage :alt="user.name" :src="user.avatar" />
                    <AvatarFallback class="rounded-lg">
                      {{ user.initials }}
                    </AvatarFallback>
                  </Avatar>
                  <div class="grid flex-1 text-left text-sm leading-tight">
                    <span class="truncate font-medium">{{ user.name }}</span>
                    <span class="text-muted-foreground truncate text-xs">
                      {{ user.email }}
                    </span>
                  </div>
                  <Icon class="ml-auto size-4" name="lucide:ellipsis-vertical" />
                </SidebarMenuButton>
              </DropdownMenuTrigger>
              <DropdownMenuContent
                align="end"
                class="w-(--reka-dropdown-menu-trigger-width) min-w-56 rounded-lg"
                :side="isMobile ? 'bottom' : 'right'"
                :side-offset="4"
              >
                <DropdownMenuLabel class="p-0 font-normal">
                  <div class="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
                    <Avatar class="h-8 w-8 rounded-lg">
                      <AvatarImage :alt="user.name" :src="user.avatar" />
                      <AvatarFallback class="rounded-lg">
                        {{ user.initials }}
                      </AvatarFallback>
                    </Avatar>
                    <div class="grid flex-1 text-left text-sm leading-tight">
                      <span class="truncate font-medium">{{ user.name }}</span>
                      <span class="text-muted-foreground truncate text-xs">
                        {{ user.email }}
                      </span>
                    </div>
                  </div>
                </DropdownMenuLabel>
                <DropdownMenuSeparator />
                <DropdownMenuGroup>
                  <DropdownMenuItem>
                    <Icon name="lucide:user-circle" />
                    Account
                  </DropdownMenuItem>
                  <DropdownMenuItem>
                    <Icon name="lucide:credit-card" />
                    Billing
                  </DropdownMenuItem>
                  <DropdownMenuItem>
                    <Icon name="lucide:bell" />
                    Notifications
                  </DropdownMenuItem>
                </DropdownMenuGroup>
                <DropdownMenuSeparator />
                <DropdownMenuItem>
                  <Icon name="lucide:log-out" />
                  Log out
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarFooter>
    </Sidebar>
    <SidebarInset>
      <header class="bg-background/90 sticky top-0 z-10 flex h-(--header-height) shrink-0 items-center gap-2 border-b transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-(--header-height)">
        <div class="flex w-full items-center gap-1 px-4 lg:gap-2 lg:px-6">
          <SidebarTrigger />
        </div>
      </header>
      <div class="flex flex-1 flex-col gap-4 p-4 md:p-6">
        <slot></slot>
      </div>
    </SidebarInset>
  </SidebarProvider>
</template>
