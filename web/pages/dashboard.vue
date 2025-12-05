<script setup lang="ts">
useHead({
  title: 'Dashboard',
})
</script>

<template>
  <div class="flex flex-1 flex-col gap-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-semibold tracking-tight">
          Dashboard
        </h1>
        <p class="text-muted-foreground">
          Overview of your monitoring infrastructure
        </p>
      </div>
    </div>

    <!-- Stats Grid -->
    <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">
            Total Monitors
          </CardTitle>
          <Icon class="h-4 w-4 text-muted-foreground" name="lucide:activity" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-semibold">
            12
          </div>
          <p class="text-xs text-muted-foreground">
            +2 from last month
          </p>
        </CardContent>
      </Card>
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">
            Active Alerts
          </CardTitle>
          <Icon class="h-4 w-4 text-muted-foreground" name="lucide:bell" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-semibold">
            3
          </div>
          <p class="text-xs text-muted-foreground">
            2 critical, 1 warning
          </p>
        </CardContent>
      </Card>
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">
            Uptime
          </CardTitle>
          <Icon class="h-4 w-4 text-muted-foreground" name="lucide:trending-up" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-semibold">
            99.9%
          </div>
          <p class="text-xs text-muted-foreground">
            Last 30 days
          </p>
        </CardContent>
      </Card>
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">
            Response Time
          </CardTitle>
          <Icon class="h-4 w-4 text-muted-foreground" name="lucide:zap" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-semibold">
            45ms
          </div>
          <p class="text-xs text-muted-foreground">
            Average response time
          </p>
        </CardContent>
      </Card>
    </div>

    <!-- Main Content Grid -->
    <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
      <!-- Recent Activity -->
      <Card class="col-span-4">
        <CardHeader>
          <CardTitle>Recent Activity</CardTitle>
          <CardDescription>
            Latest monitoring events and status changes
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div class="space-y-4">
            <div class="flex items-center gap-4">
              <div class="flex h-10 w-10 items-center justify-center rounded-full border border-border bg-muted">
                <Icon class="h-5 w-5" name="lucide:check-circle-2" />
              </div>
              <div class="flex-1 space-y-1">
                <p class="text-sm font-medium">
                  API endpoint check passed
                </p>
                <p class="text-xs text-muted-foreground">
                  2 minutes ago
                </p>
              </div>
            </div>
            <div class="flex items-center gap-4">
              <div class="flex h-10 w-10 items-center justify-center rounded-full border border-border bg-muted">
                <Icon class="h-5 w-5 text-destructive" name="lucide:alert-circle" />
              </div>
              <div class="flex-1 space-y-1">
                <p class="text-sm font-medium">
                  Database connection timeout
                </p>
                <p class="text-xs text-muted-foreground">
                  15 minutes ago
                </p>
              </div>
            </div>
            <div class="flex items-center gap-4">
              <div class="flex h-10 w-10 items-center justify-center rounded-full border border-border bg-muted">
                <Icon class="h-5 w-5" name="lucide:check-circle-2" />
              </div>
              <div class="flex-1 space-y-1">
                <p class="text-sm font-medium">
                  Web server health check passed
                </p>
                <p class="text-xs text-muted-foreground">
                  1 hour ago
                </p>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      <!-- Status Overview -->
      <Card class="col-span-3">
        <CardHeader>
          <CardTitle>Status Overview</CardTitle>
          <CardDescription>
            Current system status
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div class="space-y-4">
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2">
                <div class="h-2 w-2 rounded-full bg-green-500"></div>
                <span class="text-sm">Operational</span>
              </div>
              <span class="text-sm font-medium">8</span>
            </div>
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2">
                <div class="h-2 w-2 rounded-full bg-yellow-500"></div>
                <span class="text-sm">Warning</span>
              </div>
              <span class="text-sm font-medium">1</span>
            </div>
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2">
                <div class="h-2 w-2 rounded-full bg-red-500"></div>
                <span class="text-sm">Critical</span>
              </div>
              <span class="text-sm font-medium">2</span>
            </div>
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2">
                <div class="h-2 w-2 rounded-full bg-muted-foreground"></div>
                <span class="text-sm">Maintenance</span>
              </div>
              <span class="text-sm font-medium">1</span>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>

    <!-- Quick Actions -->
    <Card>
      <CardHeader>
        <CardTitle>Quick Actions</CardTitle>
        <CardDescription>
          Common tasks and shortcuts
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div class="flex flex-wrap gap-4">
          <Button as-child variant="outline">
            <NuxtLink to="/monitors/new">
              <Icon class="mr-2 h-4 w-4" name="lucide:plus" />
              Create Monitor
            </NuxtLink>
          </Button>
          <Button as-child variant="outline">
            <NuxtLink to="/alerts/new">
              <Icon class="mr-2 h-4 w-4" name="lucide:bell-plus" />
              New Alert Rule
            </NuxtLink>
          </Button>
          <Button as-child variant="outline">
            <NuxtLink to="/integrations">
              <Icon class="mr-2 h-4 w-4" name="lucide:plug" />
              Add Integration
            </NuxtLink>
          </Button>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
