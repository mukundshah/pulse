<script setup lang="ts">
const route = useRoute()
const projectId = route.params.id as string

useHead({
  title: 'Project',
})

const project = {
  id: projectId,
  name: 'Production API',
  description: 'Main production API monitoring',
  status: 'operational',
}

const tabs = [
  {
    id: 'checks',
    label: 'Checks',
    icon: 'lucide:activity',
  },
  {
    id: 'settings',
    label: 'Settings',
    icon: 'lucide:settings',
  },
]

const activeTab = ref('checks')

const checks = [
  {
    id: 1,
    name: 'API Health Check',
    type: 'HTTP',
    status: 'operational',
    url: 'https://api.example.com/health',
    lastCheck: '2 minutes ago',
    responseTime: '45ms',
    uptime: '99.9%',
  },
  {
    id: 2,
    name: 'Database Connection',
    type: 'TCP',
    status: 'warning',
    url: 'db.example.com:5432',
    lastCheck: '5 minutes ago',
    responseTime: '120ms',
    uptime: '98.5%',
  },
  {
    id: 3,
    name: 'Web Server',
    type: 'HTTP',
    status: 'operational',
    url: 'https://www.example.com',
    lastCheck: '1 minute ago',
    responseTime: '32ms',
    uptime: '100%',
  },
]

const getStatusColor = (status: string) => {
  switch (status) {
    case 'operational':
      return 'bg-green-500'
    case 'warning':
      return 'bg-yellow-500'
    case 'critical':
      return 'bg-red-500'
    default:
      return 'bg-muted-foreground'
  }
}

const getStatusLabel = (status: string) => {
  switch (status) {
    case 'operational':
      return 'Operational'
    case 'warning':
      return 'Warning'
    case 'critical':
      return 'Critical'
    default:
      return 'Unknown'
  }
}
</script>

<template>
  <div class="flex flex-1 flex-col gap-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-semibold tracking-tight">
          {{ project.name }}
        </h1>
        <p class="text-muted-foreground">
          {{ project.description }}
        </p>
      </div>
      <Button>
        <Icon class="mr-2 h-4 w-4" name="lucide:plus" />
        New Check
      </Button>
    </div>

    <!-- Tabs -->
    <div class="border-b border-border">
      <div class="flex gap-2">
        <button
          v-for="tab in tabs"
          :key="tab.id"
          class="flex items-center gap-2 border-b-2 px-4 py-2 text-sm font-medium transition-colors"
          :class="activeTab === tab.id ? 'border-primary text-foreground' : 'border-transparent text-muted-foreground hover:text-foreground'"
          @click="activeTab = tab.id"
        >
          <Icon :name="tab.icon" class="h-4 w-4" />
          {{ tab.label }}
        </button>
      </div>
    </div>

    <!-- Checks Tab -->
    <div v-if="activeTab === 'checks'">
      <Card>
        <CardHeader>
          <div class="flex items-center justify-between">
            <div>
              <CardTitle>Checks</CardTitle>
              <CardDescription>
                {{ checks.length }} active checks
              </CardDescription>
            </div>
            <div class="flex items-center gap-2">
              <Button size="sm" variant="outline">
                <Icon class="mr-2 h-4 w-4" name="lucide:filter" />
                Filter
              </Button>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Status</TableHead>
                <TableHead>Name</TableHead>
                <TableHead>Type</TableHead>
                <TableHead>URL</TableHead>
                <TableHead>Response Time</TableHead>
                <TableHead>Uptime</TableHead>
                <TableHead>Last Check</TableHead>
                <TableHead class="text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow
                v-for="check in checks"
                :key="check.id"
                class="cursor-pointer hover:bg-muted/50"
              >
                <TableCell>
                  <div class="flex items-center gap-2">
                    <div
                      class="h-2 w-2 rounded-full"
                      :class="getStatusColor(check.status)"
                    ></div>
                    <span class="text-sm">{{ getStatusLabel(check.status) }}</span>
                  </div>
                </TableCell>
                <TableCell class="font-medium">
                  <NuxtLink
                    :to="`/projects/${projectId}/checks/${check.id}`"
                    class="hover:underline"
                  >
                    {{ check.name }}
                  </NuxtLink>
                </TableCell>
                <TableCell>
                  <Badge variant="outline">
                    {{ check.type }}
                  </Badge>
                </TableCell>
                <TableCell class="text-muted-foreground">
                  {{ check.url }}
                </TableCell>
                <TableCell>
                  <span class="text-sm">{{ check.responseTime }}</span>
                </TableCell>
                <TableCell>
                  <span class="text-sm">{{ check.uptime }}</span>
                </TableCell>
                <TableCell class="text-muted-foreground">
                  <span class="text-sm">{{ check.lastCheck }}</span>
                </TableCell>
                <TableCell class="text-right">
                  <DropdownMenu>
                    <DropdownMenuTrigger as-child>
                      <Button size="icon" variant="ghost">
                        <Icon class="h-4 w-4" name="lucide:more-horizontal" />
                        <span class="sr-only">Open menu</span>
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end">
                      <DropdownMenuItem>
                        <Icon class="mr-2 h-4 w-4" name="lucide:edit" />
                        Edit
                      </DropdownMenuItem>
                      <DropdownMenuItem>
                        <Icon class="mr-2 h-4 w-4" name="lucide:copy" />
                        Duplicate
                      </DropdownMenuItem>
                      <DropdownMenuSeparator />
                      <DropdownMenuItem class="text-destructive">
                        <Icon class="mr-2 h-4 w-4" name="lucide:trash-2" />
                        Delete
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>

    <!-- Settings Tab -->
    <div v-if="activeTab === 'settings'">
      <NuxtPage />
    </div>
  </div>
</template>
