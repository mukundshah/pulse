<script setup lang="ts">
import SparklineChart from '@/components/metrics/SparklineChart.vue'
import StatusBadge from '~/components/StatusBadge.vue'

interface Check {
  id: string
  name: string
  type: string
  last_status: string
  last_run_at?: string
  url?: string
  project_id: string
}

interface Props {
  checks: Check[]
  loading?: boolean
  onCheckClick?: (check: Check) => void
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  onCheckClick: undefined,
})

// Generate mock sparkline data for now (will be replaced with real data)
const getSparklineData = (check: Check) => {
  // Mock data - in real implementation, this would come from check runs
  const data = []
  for (let i = 0; i < 20; i++) {
    data.push({
      timestamp: Date.now() - (20 - i) * 3600000, // Last 20 hours
      value: Math.random() * 100 + 50, // Mock response time
    })
  }
  return data
}

const getStatusFromLastStatus = (status: string) => {
  if (status === 'success') {
    return 'passing'
  }
  if (status === 'fail' || status === 'error') {
    return 'failing'
  }
  if (status === 'timeout') {
    return 'degraded'
  }
  return 'unknown'
}

const formatUptime = (check: Check) => {
  // Mock uptime calculation - in real implementation, calculate from check runs
  return '99.9%'
}

const formatResponseTime = (check: Check) => {
  // Mock response time - in real implementation, calculate from check runs
  return '45ms'
}

const formatLastCheck = (check: Check) => {
  if (!check.last_run_at) {
    return 'Never'
  }
  const date = new Date(check.last_run_at)
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffMins = Math.floor(diffMs / 60000)
  const diffHours = Math.floor(diffMs / 3600000)
  const diffDays = Math.floor(diffMs / 86400000)

  if (diffMins < 1) {
    return 'Just now'
  }
  if (diffMins < 60) {
    return `${diffMins}m ago`
  }
  if (diffHours < 24) {
    return `${diffHours}h ago`
  }
  return `${diffDays}d ago`
}
</script>

<template>
  <Card>
    <CardHeader>
      <CardTitle>Checks</CardTitle>
      <CardDescription>
        Monitor and manage your checks
      </CardDescription>
    </CardHeader>
    <CardContent>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>NAME</TableHead>
            <TableHead>TYPE</TableHead>
            <TableHead>LAST RESULTS</TableHead>
            <TableHead>24H</TableHead>
            <TableHead>7D</TableHead>
            <TableHead>AVG</TableHead>
            <TableHead>P95</TableHead>
            <TableHead>ΔT</TableHead>
            <TableHead class="text-right">Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <template v-if="loading">
            <TableRow v-for="i in 5" :key="i">
              <TableCell colspan="9">
                <div class="flex items-center gap-2">
                  <Spinner class="h-4 w-4" />
                  <span class="text-sm text-muted-foreground">Loading...</span>
                </div>
              </TableCell>
            </TableRow>
          </template>
          <template v-else-if="checks.length === 0">
            <TableRow>
              <TableCell colspan="9">
                <EmptyState
                  title="No checks configured"
                  description="Create a check to start monitoring"
                />
              </TableCell>
            </TableRow>
          </template>
          <TableRow
            v-else
            v-for="check in checks"
            :key="check.id"
            class="cursor-pointer hover:bg-muted/50"
            @click="onCheckClick && onCheckClick(check)"
          >
            <TableCell>
              <div class="flex items-center gap-2">
                <StatusBadge :status="getStatusFromLastStatus(check.last_status)" />
                <div>
                  <div class="font-medium">
                    {{ check.name }}
                  </div>
                  <div class="text-xs text-muted-foreground">
                    {{ formatLastCheck(check) }}
                  </div>
                </div>
              </div>
            </TableCell>
            <TableCell>
              <Badge variant="outline">
                {{ check.type.toUpperCase() }}
              </Badge>
            </TableCell>
            <TableCell>
              <div class="w-24">
                <SparklineChart :data="getSparklineData(check)" />
              </div>
            </TableCell>
            <TableCell>
              <span class="text-sm">{{ formatUptime(check) }}</span>
            </TableCell>
            <TableCell>
              <span class="text-sm">{{ formatUptime(check) }}</span>
            </TableCell>
            <TableCell>
              <span class="text-sm">{{ formatResponseTime(check) }}</span>
            </TableCell>
            <TableCell>
              <span class="text-sm">{{ formatResponseTime(check) }}</span>
            </TableCell>
            <TableCell>
              <span class="text-sm text-muted-foreground">--</span>
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
</template>

