<script setup lang="ts">
useHead({
  title: 'Alerts',
})

const filter = ref<'all' | 'active' | 'resolved'>('all')

const alerts = [
  {
    id: 1,
    check: 'API Health Check',
    severity: 'critical',
    status: 'active',
    message: 'Check failed: Connection timeout after 30s',
    triggeredAt: '5 minutes ago',
    acknowledged: false,
  },
  {
    id: 2,
    check: 'Database Connection',
    severity: 'warning',
    status: 'active',
    message: 'Response time exceeded threshold: 500ms',
    triggeredAt: '12 minutes ago',
    acknowledged: true,
  },
  {
    id: 3,
    check: 'Email Service',
    severity: 'critical',
    status: 'resolved',
    message: 'Check failed: SMTP connection refused',
    triggeredAt: '2 hours ago',
    resolvedAt: '1 hour ago',
    acknowledged: true,
  },
  {
    id: 4,
    check: 'Web Server',
    severity: 'warning',
    status: 'active',
    message: 'Uptime dropped below 99% threshold',
    triggeredAt: '1 hour ago',
    acknowledged: false,
  },
  {
    id: 5,
    check: 'Redis Cache',
    severity: 'warning',
    status: 'resolved',
    message: 'Response time spike detected: 200ms',
    triggeredAt: '3 hours ago',
    resolvedAt: '2 hours ago',
    acknowledged: true,
  },
  {
    id: 6,
    check: 'CDN',
    severity: 'critical',
    status: 'resolved',
    message: 'Check failed: HTTP 503 Service Unavailable',
    triggeredAt: '5 hours ago',
    resolvedAt: '4 hours ago',
    acknowledged: true,
  },
]

const filteredAlerts = computed(() => {
  if (filter.value === 'all') {
    return alerts
  }
  return alerts.filter(alert => alert.status === filter.value)
})

const activeCount = computed(() => alerts.filter(a => a.status === 'active').length)
const resolvedCount = computed(() => alerts.filter(a => a.status === 'resolved').length)

const getSeverityColor = (severity: string) => {
  switch (severity) {
    case 'critical':
      return 'bg-red-500'
    case 'warning':
      return 'bg-yellow-500'
    default:
      return 'bg-muted-foreground'
  }
}

const getSeverityLabel = (severity: string) => {
  switch (severity) {
    case 'critical':
      return 'Critical'
    case 'warning':
      return 'Warning'
    default:
      return 'Unknown'
  }
}

const getSeverityBadgeVariant = (severity: string) => {
  switch (severity) {
    case 'critical':
      return 'destructive'
    case 'warning':
      return 'default'
    default:
      return 'outline'
  }
}
</script>

<template>
  <div class="flex flex-1 flex-col gap-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-semibold tracking-tight">
          Alerts
        </h1>
        <p class="text-muted-foreground">
          Monitor and manage alert notifications
        </p>
      </div>
      <div class="flex items-center gap-2">
        <Button size="sm" variant="outline">
          <Icon class="mr-2 h-4 w-4" name="lucide:bell-off" />
          Mute All
        </Button>
        <Button size="sm" variant="outline">
          <Icon class="mr-2 h-4 w-4" name="lucide:settings" />
          Alert Rules
        </Button>
      </div>
    </div>

    <!-- Stats Cards -->
    <div class="grid gap-4 md:grid-cols-3">
      <Card>
        <CardHeader class="pb-3">
          <CardDescription>Total Alerts</CardDescription>
          <CardTitle class="text-2xl">{{ alerts.length }}</CardTitle>
        </CardHeader>
      </Card>
      <Card>
        <CardHeader class="pb-3">
          <CardDescription>Active</CardDescription>
          <CardTitle class="text-2xl text-yellow-500">
            {{ activeCount }}
          </CardTitle>
        </CardHeader>
      </Card>
      <Card>
        <CardHeader class="pb-3">
          <CardDescription>Resolved</CardDescription>
          <CardTitle class="text-2xl text-green-500">
            {{ resolvedCount }}
          </CardTitle>
        </CardHeader>
      </Card>
    </div>

    <!-- Filters -->
    <div class="flex items-center gap-2">
      <Button
        :variant="filter === 'all' ? 'default' : 'outline'"
        size="sm"
        @click="filter = 'all'"
      >
        All ({{ alerts.length }})
      </Button>
      <Button
        :variant="filter === 'active' ? 'default' : 'outline'"
        size="sm"
        @click="filter = 'active'"
      >
        Active ({{ activeCount }})
      </Button>
      <Button
        :variant="filter === 'resolved' ? 'default' : 'outline'"
        size="sm"
        @click="filter = 'resolved'"
      >
        Resolved ({{ resolvedCount }})
      </Button>
    </div>

    <!-- Alerts Table -->
    <Card>
      <CardHeader>
        <div class="flex items-center justify-between">
          <div>
            <CardTitle>Alert History</CardTitle>
            <CardDescription>
              {{ filteredAlerts.length }} alert{{ filteredAlerts.length !== 1 ? 's' : '' }}
            </CardDescription>
          </div>
          <div class="flex items-center gap-2">
            <Button size="sm" variant="outline">
              <Icon class="mr-2 h-4 w-4" name="lucide:download" />
              Export
            </Button>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Status</TableHead>
              <TableHead>Severity</TableHead>
              <TableHead>Check</TableHead>
              <TableHead>Message</TableHead>
              <TableHead>Triggered</TableHead>
              <TableHead>Resolved</TableHead>
              <TableHead class="text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow
              v-for="alert in filteredAlerts"
              :key="alert.id"
              class="cursor-pointer hover:bg-muted/50"
            >
              <TableCell>
                <div class="flex items-center gap-2">
                  <div
                    class="h-2 w-2 rounded-full"
                    :class="alert.status === 'active' ? getSeverityColor(alert.severity) : 'bg-green-500'"
                  ></div>
                  <Badge
                    :variant="alert.status === 'active' ? 'default' : 'outline'"
                  >
                    {{ alert.status === 'active' ? 'Active' : 'Resolved' }}
                  </Badge>
                  <Badge
                    v-if="alert.acknowledged"
                    variant="secondary"
                    class="text-xs"
                  >
                    Acknowledged
                  </Badge>
                </div>
              </TableCell>
              <TableCell>
                <Badge :variant="getSeverityBadgeVariant(alert.severity)">
                  {{ getSeverityLabel(alert.severity) }}
                </Badge>
              </TableCell>
              <TableCell class="font-medium">
                {{ alert.check }}
              </TableCell>
              <TableCell class="text-muted-foreground">
                <span class="text-sm">{{ alert.message }}</span>
              </TableCell>
              <TableCell>
                <span class="text-sm">{{ alert.triggeredAt }}</span>
              </TableCell>
              <TableCell>
                <span
                  v-if="alert.resolvedAt"
                  class="text-sm text-muted-foreground"
                >
                  {{ alert.resolvedAt }}
                </span>
                <span v-else class="text-sm text-muted-foreground">—</span>
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
                    <DropdownMenuItem
                      v-if="alert.status === 'active' && !alert.acknowledged"
                    >
                      <Icon class="mr-2 h-4 w-4" name="lucide:check" />
                      Acknowledge
                    </DropdownMenuItem>
                    <DropdownMenuItem
                      v-if="alert.status === 'active'"
                    >
                      <Icon class="mr-2 h-4 w-4" name="lucide:check-circle-2" />
                      Resolve
                    </DropdownMenuItem>
                    <DropdownMenuItem>
                      <Icon class="mr-2 h-4 w-4" name="lucide:bell-off" />
                      Mute
                    </DropdownMenuItem>
                    <DropdownMenuSeparator />
                    <DropdownMenuItem>
                      <Icon class="mr-2 h-4 w-4" name="lucide:external-link" />
                      View Check
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
</template>
