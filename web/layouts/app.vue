<script setup lang="ts">
import { useMediaQuery } from '@vueuse/core'

useHead({ titleTemplate: '%siteName %separator %s' })

const { me, logout } = useAuth()

const colorMode = useColorMode()
const THEME_ICONS = {
  light: 'lucide:sun',
  dark: 'lucide:moon',
  system: 'lucide:monitor',
} as const

const isMobile = useMediaQuery('(max-width: 768px)')

const { breadcrumbOverrides, actions } = useLayoutContext()
const breadcrumbs = useBreadcrumbItems({
  hideRoot: true,
  // @ts-expect-error - TODO: fix the typing of the breadcrumb overrides
  overrides: breadcrumbOverrides,
})

const { data: user, pending: userPending } = useAsyncData('user', () => me())

const showProjectInput = ref(false)

const { data: projects, pending: projectsLoading, refresh: refreshProjects } = useLazyPulseAPI('/internal/projects')

const handleLogout = async () => {
  await logout()
  await navigateTo('/auth/login')
}
</script>

<template>
  <Body class="overflow-hidden">
    <SidebarProvider
      :style="{
        '--sidebar-width': 'calc(var(--spacing) * 80)',
        '--header-height': 'calc(var(--spacing) * 15 + 1px)',
      }"
    >
      <Sidebar class="h-auto border-r" collapsible="offcanvas">
        <SidebarHeader class="border-b">
          <NuxtLink class="flex items-center gap-2.5 px-2 py-1.5" to="/dashboard">
            <div class="flex h-8 w-8 items-center justify-center rounded-md bg-foreground text-background font-semibold">
              P
            </div>
            <span class="text-base font-semibold">Pulse</span>
          </NuxtLink>
        </SidebarHeader>
        <SidebarContent class="gap-y-0">
          <SidebarGroup class="mt-2 -mb-2">
            <SidebarMenu>
              <SidebarMenuItem>
                <SidebarMenuButton as-child :is-active="$route.path === '/dashboard'">
                  <NuxtLink to="/dashboard">
                    <Icon name="lucide:layout-dashboard" />
                    <span>Dashboard</span>
                  </NuxtLink>
                </SidebarMenuButton>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarGroup>
          <SidebarGroup>
            <SidebarGroupLabel>
              Projects
              <SidebarGroupAction @click="showProjectInput = true">
                <Icon name="lucide:plus" /> <span class="sr-only">Add project</span>
              </SidebarGroupAction>
            </SidebarGroupLabel>
            <SidebarGroupContent>
              <ProjectInlineForm
                v-if="showProjectInput"
                @cancel="showProjectInput = false"
                @success="async ($event) => {
                  showProjectInput = false
                  await refreshProjects()
                  await navigateTo({ name: 'projects-projectId', params: { projectId: $event.id } })
                }"
              />
              <SidebarMenu>
                <template v-if="projectsLoading">
                  <SidebarMenuItem v-for="item in [1, 2, 3, 4, 5]" :key="item">
                    <SidebarMenuSkeleton />
                  </SidebarMenuItem>
                </template>
                <template v-else-if="projects && projects.length > 0">
                  <SidebarMenuItem v-for="project in projects" :key="project.id">
                    <SidebarMenuButton as-child :is-active="$route.fullPath.startsWith(`/projects/${project.id}`)">
                      <NuxtLink :to="`/projects/${project.id}`">
                        {{ project.name }}
                      </NuxtLink>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                </template>
                <template v-else>
                  <SidebarMenuItem>
                    <div class="px-2 py-2 text-sm text-muted-foreground">
                      No projects yet
                    </div>
                  </SidebarMenuItem>
                </template>
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        </SidebarContent>
        <SidebarFooter>
          <SidebarMenu>
            <SidebarMenuItem>
              <template v-if="userPending">
                <SidebarMenuButton
                  class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
                  size="lg"
                >
                  <Skeleton class="h-8 w-8 rounded-lg" />
                  <div class="grid flex-1 text-left text-sm leading-tight">
                    <Skeleton class="h-4 max-w-(--skeleton-width) flex-1" :style="{ '--skeleton-width': '40%' }" />
                    <Skeleton class="h-3 max-w-(--skeleton-width) flex-1 mt-1" :style="{ '--skeleton-width': '70%' }" />
                  </div>
                  <Skeleton class="ml-auto size-6 rounded" />
                </SidebarMenuButton>
              </template>
              <DropdownMenu v-else>
                <DropdownMenuTrigger as-child>
                  <SidebarMenuButton
                    class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
                    size="lg"
                  >
                    <Avatar class="h-8 w-8 rounded-lg grayscale">
                      <AvatarImage :alt="user?.name" :src="user?.avatar_url ?? ''" />
                      <AvatarFallback class="rounded-lg">
                        {{ user?.initials }}
                      </AvatarFallback>
                    </Avatar>
                    <div class="grid flex-1 text-left text-sm leading-tight">
                      <span class="truncate font-medium">{{ user?.name }}</span>
                      <span class="text-muted-foreground truncate text-xs">
                        {{ user?.email }}
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
                        <AvatarImage :alt="user?.name" :src="user?.avatar_url ?? ''" />
                        <AvatarFallback class="rounded-lg">
                          {{ user?.initials }}
                        </AvatarFallback>
                      </Avatar>
                      <div class="grid flex-1 text-left text-sm leading-tight">
                        <span class="truncate font-medium">{{ user?.name }}</span>
                        <span class="text-muted-foreground truncate text-xs">
                          {{ user?.email }}
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
                  </DropdownMenuGroup>
                  <DropdownMenuSub>
                    <DropdownMenuSubTrigger>
                      <div class="flex items-center gap-2">
                        <Icon class="h-4 w-4" :name="THEME_ICONS[colorMode.preference as keyof typeof THEME_ICONS]" />
                        <span>Theme</span>
                      </div>
                    </DropdownMenuSubTrigger>
                    <DropdownMenuPortal>
                      <DropdownMenuSubContent>
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
                      </DropdownMenuSubContent>
                    </DropdownMenuPortal>
                  </DropdownMenuSub>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem @click="handleLogout">
                    <Icon name="lucide:log-out" />
                    Log out
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarFooter>
      </Sidebar>
      <SidebarInset class="overflow-auto h-svh">
        <header class="bg-background/90 sticky top-0 z-10 flex h-(--header-height) shrink-0 items-center gap-2 border-b transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-(--header-height) backdrop-blur-md">
          <div class="flex items-center gap-2 justify-between w-full px-4 lg:px-6">
            <div class="flex items-center gap-2">
              <SidebarTrigger />
              <!-- Breadcrumbs -->
              <Breadcrumb>
                <BreadcrumbList>
                  <template v-for="(breadcrumb, idx) in breadcrumbs" :key="breadcrumb.id">
                    <BreadcrumbItem>
                      <BreadcrumbLink v-if="!breadcrumb.active" as-child>
                        <NuxtLink :to="breadcrumb.to">
                          {{ breadcrumb.label }}
                        </NuxtLink>
                      </BreadcrumbLink>
                      <BreadcrumbPage v-else as="h1">
                        {{ breadcrumb.label }}
                      </BreadcrumbPage>
                    </BreadcrumbItem>
                    <BreadcrumbSeparator v-if="idx < breadcrumbs.length - 1" />
                  </template>
                </BreadcrumbList>
              </Breadcrumb>
            </div>
            <div class="flex items-center gap-2">
              <template v-for="action in actions" :key="action.label">
                <template v-if="action.children">
                  <DropdownMenu>
                    <DropdownMenuTrigger as-child>
                      <Button v-bind="action.props">
                        <Icon v-if="action.icon" :name="action.icon" />
                        {{ action.label }}
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end" class="min-w-56">
                      <template v-for="(child, idx) in action.children" :key="child.label">
                        <template v-if="'children' in child">
                          <DropdownMenuGroup>
                            <DropdownMenuLabel>
                              {{ child.label }}
                            </DropdownMenuLabel>
                            <DropdownMenuItem v-for="c in child.children" :key="c.label" :as-child="!!c.to">
                              <template v-if="c.to">
                                <NuxtLink :to="c.to">
                                  <Icon v-if="c.icon" :name="c.icon" />
                                  {{ c.label }}
                                </NuxtLink>
                              </template>
                              <template v-else>
                                <Icon v-if="c.icon" :name="c.icon" />
                                {{ c.label }}
                              </template>
                            </DropdownMenuItem>
                          </DropdownMenuGroup>
                          <DropdownMenuSeparator v-if="idx < action.children.length - 1" />
                        </template>
                        <template v-else>
                          <DropdownMenuItem :key="child.label" :as-child="!!child.to">
                            <template v-if="child.to">
                              <NuxtLink :to="child.to">
                                <Icon v-if="child.icon" :name="child.icon" />
                                {{ child.label }}
                              </NuxtLink>
                            </template>
                            <template v-else>
                              <Icon v-if="child.icon" :name="child.icon" />
                              {{ child.label }}
                            </template>
                          </DropdownMenuItem>
                        </template>
                      </template>
                    </DropdownMenuContent>
                  </DropdownMenu>
                </template>
                <template v-else>
                  <Button
                    :as-child="!!action.to"
                    v-bind="action.props"
                    @click="action.onClick"
                  >
                    <template v-if="action.to">
                      <NuxtLink :to="action.to">
                        <Icon v-if="action.icon" :name="action.icon" />
                        {{ action.label }}
                      </NuxtLink>
                    </template>
                    <template v-else>
                      <Icon v-if="action.icon" :name="action.icon" />
                      {{ action.label }}
                    </template>
                  </Button>
                </template>
              </template>
            </div>
          </div>
        </header>
        <div class="flex flex-1 flex-col gap-4">
          <slot></slot>
        </div>
      </SidebarInset>
    </SidebarProvider>
  </Body>
</template>
