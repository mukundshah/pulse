<script setup lang="ts">
import AlertCard from '@/components/alerts/AlertCard.vue'
import StatusBadge from '~/components/StatusBadge.vue'

useHead({
  title: 'Alerts',
})

const filter = ref<'all' | 'active' | 'resolved'>('all')

// Load alerts (would need to aggregate from all checks or have a dedicated endpoint)
const alerts = ref([])
const loading = ref(false)

// Mock stats for now
const stats = computed(() => {
  const active = alerts.value.filter(a => a.status === 'active').length
  const resolved = alerts.value.filter(a => a.status === 'resolved').length
  return {
    total: alerts.value.length,
    active,
    resolved,
  }
})

const filteredAlerts = computed(() => {
  if (filter.value === 'all') {
    return alerts.value
  }
  return alerts.value.filter(a => a.status === filter.value)
})
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
        <Button
          :variant="filter === 'all' ? 'default' : 'outline'"
          size="sm"
          @click="filter = 'all'"
        >
          All
        </Button>
        <Button
          :variant="filter === 'active' ? 'default' : 'outline'"
          size="sm"
          @click="filter = 'active'"
        >
          Active
        </Button>
        <Button
          :variant="filter === 'resolved' ? 'default' : 'outline'"
          size="sm"
          @click="filter = 'resolved'"
        >
          Resolved
        </Button>
      </div>
    </div>

    <!-- Stats Cards -->
    <div class="grid gap-4 md:grid-cols-3">
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">
            Total Alerts
          </CardTitle>
          <Icon class="h-4 w-4 text-muted-foreground" name="lucide:bell" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">
            {{ stats.total }}
          </div>
        </CardContent>
      </Card>
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">
            Active
          </CardTitle>
          <Icon class="h-4 w-4 text-muted-foreground" name="lucide:alert-circle" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold text-red-600">
            {{ stats.active }}
          </div>
        </CardContent>
      </Card>
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">
            Resolved
          </CardTitle>
          <Icon class="h-4 w-4 text-muted-foreground" name="lucide:check-circle-2" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold text-green-600">
            {{ stats.resolved }}
          </div>
        </CardContent>
      </Card>
    </div>

    <!-- Alerts Table -->
    <Card>
      <CardHeader>
        <CardTitle>Alerts</CardTitle>
        <CardDescription>
          All alert notifications
        </CardDescription>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Status</TableHead>
              <TableHead>Severity</TableHead>
              <TableHead>Check Name</TableHead>
              <TableHead>Message</TableHead>
              <TableHead>Triggered At</TableHead>
              <TableHead>Resolved At</TableHead>
              <TableHead class="text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-if="filteredAlerts.length === 0">
              <TableCell colspan="7" class="text-center">
                <EmptyState
                  title="No alerts"
                  description="All systems operational"
                />
              </TableCell>
            </TableRow>
            <TableRow
              v-for="alert in filteredAlerts"
              :key="alert.id"
            >
              <TableCell>
                <StatusBadge
                  :status="alert.status === 'active' ? 'failing' : 'passing'"
                  :label="alert.status"
                />
              </TableCell>
              <TableCell>
                <Badge
                  :variant="alert.severity === 'critical' ? 'destructive' : 'default'"
                >
                  {{ alert.severity }}
                </Badge>
              </TableCell>
              <TableCell class="font-medium">
                {{ alert.check_name }}
              </TableCell>
              <TableCell>
                {{ alert.message }}
              </TableCell>
              <TableCell class="text-muted-foreground">
                {{ alert.triggered_at ? new Date(alert.triggered_at).toLocaleString() : '--' }}
              </TableCell>
              <TableCell class="text-muted-foreground">
                {{ alert.resolved_at ? new Date(alert.resolved_at).toLocaleString() : '--' }}
              </TableCell>
              <TableCell class="text-right">
                <DropdownMenu>
                  <DropdownMenuTrigger as-child>
                    <Button size="icon" variant="ghost">
                      <Icon class="h-4 w-4" name="lucide:more-horizontal" />
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end">
                    <DropdownMenuItem>
                      <Icon class="mr-2 h-4 w-4" name="lucide:eye" />
                      View Details
                    </DropdownMenuItem>
                    <DropdownMenuItem v-if="alert.status === 'active'">
                      <Icon class="mr-2 h-4 w-4" name="lucide:check" />
                      Resolve
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
