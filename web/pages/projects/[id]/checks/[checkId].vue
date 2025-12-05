<script setup lang="ts">
const route = useRoute()
const projectId = route.params.id as string
const checkId = route.params.checkId as string

useHead({
  title: 'Check Details',
})

const check = {
  id: checkId,
  name: 'API Health Check',
  type: 'HTTP',
  status: 'operational',
  url: 'https://api.example.com/health',
  method: 'GET',
  expectedStatus: 200,
  interval: '60s',
  timeout: '10s',
  lastCheck: '2 minutes ago',
  responseTime: '45ms',
  uptime: '99.9%',
}

const checkRuns = [
  {
    id: 1,
    status: 'success',
    statusCode: 200,
    latencyMs: 45,
    runAt: '2 minutes ago',
    responseTime: '45ms',
  },
  {
    id: 2,
    status: 'success',
    statusCode: 200,
    latencyMs: 42,
    runAt: '3 minutes ago',
    responseTime: '42ms',
  },
  {
    id: 3,
    status: 'success',
    statusCode: 200,
    latencyMs: 48,
    runAt: '4 minutes ago',
    responseTime: '48ms',
  },
  {
    id: 4,
    status: 'fail',
    statusCode: 500,
    latencyMs: 120,
    runAt: '5 minutes ago',
    responseTime: '120ms',
    error: 'Internal server error',
  },
  {
    id: 5,
    status: 'success',
    statusCode: 200,
    latencyMs: 44,
    runAt: '6 minutes ago',
    responseTime: '44ms',
  },
]

const alerts = [
  {
    id: 1,
    severity: 'critical',
    message: 'Check failed: HTTP 500',
    triggeredAt: '5 minutes ago',
    resolvedAt: '4 minutes ago',
    status: 'resolved',
  },
]

const getStatusColor = (status: string) => {
  switch (status) {
    case 'success':
      return 'bg-green-500'
    case 'fail':
      return 'bg-red-500'
    case 'timeout':
      return 'bg-yellow-500'
    default:
      return 'bg-muted-foreground'
  }
}

const getStatusLabel = (status: string) => {
  switch (status) {
    case 'success':
      return 'Success'
    case 'fail':
      return 'Failed'
    case 'timeout':
      return 'Timeout'
    default:
      return 'Unknown'
  }
}
</script>

<template>
  <div class="flex flex-1 flex-col gap-6">
    <div class="flex items-center justify-between">
      <div>
        <div class="mb-2 flex items-center gap-2 text-sm text-muted-foreground">
          <NuxtLink :to="`/projects/${projectId}`" class="hover:underline">
            Project
          </NuxtLink>
          <Icon class="h-4 w-4" name="lucide:chevron-right" />
          <span>Check Details</span>
        </div>
        <h1 class="text-3xl font-semibold tracking-tight">
          {{ check.name }}
        </h1>
        <p class="text-muted-foreground">
          {{ check.url }}
        </p>
      </div>
      <div class="flex items-center gap-2">
        <Button variant="outline">
          <Icon class="mr-2 h-4 w-4" name="lucide:edit" />
          Edit
        </Button>
        <Button>
          <Icon class="mr-2 h-4 w-4" name="lucide:play" />
          Run Now
        </Button>
      </div>
    </div>

    <!-- Check Info Cards -->
    <div class="grid gap-4 md:grid-cols-4">
      <Card>
        <CardHeader class="pb-3">
          <CardDescription>Status</CardDescription>
          <CardTitle class="flex items-center gap-2">
            <div
              class="h-2 w-2 rounded-full"
              :class="getStatusColor(check.status)"
            ></div>
            {{ check.status }}
          </CardTitle>
        </CardHeader>
      </Card>
      <Card>
        <CardHeader class="pb-3">
          <CardDescription>Response Time</CardDescription>
          <CardTitle>{{ check.responseTime }}</CardTitle>
        </CardHeader>
      </Card>
      <Card>
        <CardHeader class="pb-3">
          <CardDescription>Uptime</CardDescription>
          <CardTitle>{{ check.uptime }}</CardTitle>
        </CardHeader>
      </Card>
      <Card>
        <CardHeader class="pb-3">
          <CardDescription>Last Check</CardDescription>
          <CardTitle class="text-sm">{{ check.lastCheck }}</CardTitle>
        </CardHeader>
      </Card>
    </div>

    <!-- Check Configuration -->
    <Card>
      <CardHeader>
        <CardTitle>Configuration</CardTitle>
        <CardDescription>
          Check settings and parameters
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div class="grid gap-4 md:grid-cols-2">
          <div class="space-y-1">
            <label class="text-sm font-medium">Type</label>
            <p class="text-sm text-muted-foreground">{{ check.type }}</p>
          </div>
          <div class="space-y-1">
            <label class="text-sm font-medium">Method</label>
            <p class="text-sm text-muted-foreground">{{ check.method }}</p>
          </div>
          <div class="space-y-1">
            <label class="text-sm font-medium">Expected Status</label>
            <p class="text-sm text-muted-foreground">{{ check.expectedStatus }}</p>
          </div>
          <div class="space-y-1">
            <label class="text-sm font-medium">Interval</label>
            <p class="text-sm text-muted-foreground">{{ check.interval }}</p>
          </div>
          <div class="space-y-1">
            <label class="text-sm font-medium">Timeout</label>
            <p class="text-sm text-muted-foreground">{{ check.timeout }}</p>
          </div>
        </div>
      </CardContent>
    </Card>

    <!-- Check Runs -->
    <Card>
      <CardHeader>
        <div class="flex items-center justify-between">
          <div>
            <CardTitle>Check Runs</CardTitle>
            <CardDescription>
              Recent execution history
            </CardDescription>
          </div>
          <Button size="sm" variant="outline">
            <Icon class="mr-2 h-4 w-4" name="lucide:download" />
            Export
          </Button>
        </div>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Status</TableHead>
              <TableHead>Status Code</TableHead>
              <TableHead>Response Time</TableHead>
              <TableHead>Latency</TableHead>
              <TableHead>Run At</TableHead>
              <TableHead>Error</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow
              v-for="run in checkRuns"
              :key="run.id"
              class="cursor-pointer hover:bg-muted/50"
            >
              <TableCell>
                <div class="flex items-center gap-2">
                  <div
                    class="h-2 w-2 rounded-full"
                    :class="getStatusColor(run.status)"
                  ></div>
                  <span class="text-sm">{{ getStatusLabel(run.status) }}</span>
                </div>
              </TableCell>
              <TableCell>
                <Badge
                  :variant="run.statusCode === 200 ? 'default' : 'destructive'"
                >
                  {{ run.statusCode }}
                </Badge>
              </TableCell>
              <TableCell>
                <span class="text-sm">{{ run.responseTime }}</span>
              </TableCell>
              <TableCell>
                <span class="text-sm">{{ run.latencyMs }}ms</span>
              </TableCell>
              <TableCell class="text-muted-foreground">
                <span class="text-sm">{{ run.runAt }}</span>
              </TableCell>
              <TableCell>
                <span
                  v-if="run.error"
                  class="text-sm text-destructive"
                >
                  {{ run.error }}
                </span>
                <span v-else class="text-sm text-muted-foreground">—</span>
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </CardContent>
    </Card>

    <!-- Alerts -->
    <Card>
      <CardHeader>
        <CardTitle>Alerts</CardTitle>
        <CardDescription>
          Alerts triggered by this check
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div v-if="alerts.length === 0" class="py-8 text-center text-muted-foreground">
          No alerts for this check
        </div>
        <div v-else class="space-y-4">
          <div
            v-for="alert in alerts"
            :key="alert.id"
            class="flex items-start justify-between rounded-lg border border-border p-4"
          >
            <div class="flex-1">
              <div class="mb-2 flex items-center gap-2">
                <Badge variant="destructive">
                  {{ alert.severity }}
                </Badge>
                <Badge variant="outline">
                  {{ alert.status }}
                </Badge>
              </div>
              <p class="text-sm">{{ alert.message }}</p>
              <p class="mt-2 text-xs text-muted-foreground">
                Triggered: {{ alert.triggeredAt }}
                <span v-if="alert.resolvedAt">
                  • Resolved: {{ alert.resolvedAt }}
                </span>
              </p>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
