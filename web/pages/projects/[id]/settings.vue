<script setup lang="ts">
const route = useRoute()
const projectId = route.params.id as string

useHead({
  title: 'Project Settings',
})

const settingsTabs = [
  {
    id: 'status-page',
    label: 'Status Page',
    icon: 'lucide:globe',
  },
  {
    id: 'notifications',
    label: 'Notifications',
    icon: 'lucide:bell',
  },
  {
    id: 'members',
    label: 'Members',
    icon: 'lucide:users',
  },
]

const activeSettingsTab = ref('status-page')

const statusPageUrl = computed(() => {
  return `https://status.example.com/projects/${projectId}`
})

const embedCode = computed(() => {
  return `<iframe src="${statusPageUrl.value}" width="100%" height="600" frameborder="0"></iframe>`
})

const webhooks = [
  {
    id: 1,
    name: 'Slack Notifications',
    url: 'https://hooks.slack.com/services/xxx',
    events: ['check_failed', 'check_recovered'],
    status: 'active',
    lastAttempt: '2 minutes ago',
    attempts: 150,
    failures: 2,
  },
  {
    id: 2,
    name: 'Discord Webhook',
    url: 'https://discord.com/api/webhooks/xxx',
    events: ['check_failed'],
    status: 'active',
    lastAttempt: '5 minutes ago',
    attempts: 89,
    failures: 0,
  },
  {
    id: 3,
    name: 'Custom API',
    url: 'https://api.example.com/webhooks/pulse',
    events: ['check_failed', 'check_recovered', 'alert_created'],
    status: 'inactive',
    lastAttempt: '1 hour ago',
    attempts: 45,
    failures: 5,
  },
]

const webhookAttempts = [
  {
    id: 1,
    webhookId: 1,
    webhookName: 'Slack Notifications',
    status: 'success',
    statusCode: 200,
    responseTime: '45ms',
    timestamp: '2 minutes ago',
    payload: { check_id: 1, status: 'failed', message: 'Check failed' },
  },
  {
    id: 2,
    webhookId: 1,
    webhookName: 'Slack Notifications',
    status: 'failed',
    statusCode: 500,
    responseTime: '120ms',
    timestamp: '15 minutes ago',
    error: 'Internal server error',
    payload: { check_id: 2, status: 'failed', message: 'Check failed' },
  },
  {
    id: 3,
    webhookId: 2,
    webhookName: 'Discord Webhook',
    status: 'success',
    statusCode: 204,
    responseTime: '32ms',
    timestamp: '1 hour ago',
    payload: { check_id: 3, status: 'recovered', message: 'Check recovered' },
  },
]

const members = [
  {
    id: 1,
    name: 'John Doe',
    email: 'john.doe@example.com',
    role: 'owner',
    avatar: 'https://github.com/shadcn.png',
    joinedAt: 'Jan 1, 2024',
  },
  {
    id: 2,
    name: 'Jane Smith',
    email: 'jane.smith@example.com',
    role: 'admin',
    avatar: 'https://github.com/shadcn.png',
    joinedAt: 'Jan 15, 2024',
  },
  {
    id: 3,
    name: 'Bob Johnson',
    email: 'bob.johnson@example.com',
    role: 'member',
    avatar: 'https://github.com/shadcn.png',
    joinedAt: 'Feb 1, 2024',
  },
]

const getRoleBadgeVariant = (role: string) => {
  switch (role) {
    case 'owner':
      return 'default'
    case 'admin':
      return 'secondary'
    default:
      return 'outline'
  }
}
</script>

<template>
  <div class="flex flex-1 flex-col gap-6">
    <!-- Settings Sub-tabs -->
    <div class="border-b border-border">
      <div class="flex gap-2">
        <button
          v-for="tab in settingsTabs"
          :key="tab.id"
          class="flex items-center gap-2 border-b-2 px-4 py-2 text-sm font-medium transition-colors"
          :class="activeSettingsTab === tab.id ? 'border-primary text-foreground' : 'border-transparent text-muted-foreground hover:text-foreground'"
          @click="activeSettingsTab = tab.id"
        >
          <Icon :name="tab.icon" class="h-4 w-4" />
          {{ tab.label }}
        </button>
      </div>
    </div>

    <!-- Status Page Tab -->
    <div v-if="activeSettingsTab === 'status-page'">
      <Card>
        <CardHeader>
          <CardTitle>Status Page</CardTitle>
          <CardDescription>
            Configure your public status page and embed code
          </CardDescription>
        </CardHeader>
        <CardContent class="space-y-6">
          <div class="space-y-2">
            <label class="text-sm font-medium">Status Page URL</label>
            <div class="flex items-center gap-2">
              <Input
                :model-value="statusPageUrl"
                readonly
                class="font-mono"
              />
              <Button size="icon" variant="outline">
                <Icon class="h-4 w-4" name="lucide:copy" />
              </Button>
            </div>
          </div>
          <div class="space-y-2">
            <label class="text-sm font-medium">Embed Code</label>
            <div class="flex items-center gap-2">
              <Input
                :model-value="embedCode"
                readonly
                class="font-mono text-xs"
              />
              <Button size="icon" variant="outline">
                <Icon class="h-4 w-4" name="lucide:copy" />
              </Button>
            </div>
            <p class="text-xs text-muted-foreground">
              Copy and paste this code into your website to embed the status page
            </p>
          </div>
          <div class="flex items-center gap-2">
            <Button>
              <Icon class="mr-2 h-4 w-4" name="lucide:external-link" />
              View Status Page
            </Button>
            <Button variant="outline">
              <Icon class="mr-2 h-4 w-4" name="lucide:settings" />
              Configure
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>

    <!-- Notifications Tab -->
    <div v-if="activeSettingsTab === 'notifications'">
      <div class="space-y-6">
        <div class="flex items-center justify-between">
          <div>
            <h2 class="text-2xl font-semibold">Webhooks</h2>
            <p class="text-muted-foreground">
              Configure webhook notifications for alerts
            </p>
          </div>
          <Button>
            <Icon class="mr-2 h-4 w-4" name="lucide:plus" />
            Add Webhook
          </Button>
        </div>

        <Card>
          <CardHeader>
            <CardTitle>Webhooks</CardTitle>
            <CardDescription>
              {{ webhooks.length }} webhook{{ webhooks.length !== 1 ? 's' : '' }} configured
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Name</TableHead>
                  <TableHead>URL</TableHead>
                  <TableHead>Events</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Last Attempt</TableHead>
                  <TableHead>Success Rate</TableHead>
                  <TableHead class="text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                <TableRow
                  v-for="webhook in webhooks"
                  :key="webhook.id"
                  class="cursor-pointer hover:bg-muted/50"
                >
                  <TableCell class="font-medium">
                    {{ webhook.name }}
                  </TableCell>
                  <TableCell class="text-muted-foreground font-mono text-xs">
                    {{ webhook.url }}
                  </TableCell>
                  <TableCell>
                    <div class="flex flex-wrap gap-1">
                      <Badge
                        v-for="event in webhook.events"
                        :key="event"
                        variant="outline"
                        class="text-xs"
                      >
                        {{ event }}
                      </Badge>
                    </div>
                  </TableCell>
                  <TableCell>
                    <Badge
                      :variant="webhook.status === 'active' ? 'default' : 'secondary'"
                    >
                      {{ webhook.status }}
                    </Badge>
                  </TableCell>
                  <TableCell class="text-muted-foreground">
                    <span class="text-sm">{{ webhook.lastAttempt }}</span>
                  </TableCell>
                  <TableCell>
                    <span
                      class="text-sm"
                      :class="(webhook.attempts - webhook.failures) / webhook.attempts > 0.95 ? 'text-green-500' : 'text-yellow-500'"
                    >
                      {{ Math.round(((webhook.attempts - webhook.failures) / webhook.attempts) * 100) }}%
                    </span>
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
                          <Icon class="mr-2 h-4 w-4" name="lucide:eye" />
                          View Logs
                        </DropdownMenuItem>
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

        <Card>
          <CardHeader>
            <div class="flex items-center justify-between">
              <div>
                <CardTitle>Webhook Attempt Logs</CardTitle>
                <CardDescription>
                  Recent webhook delivery attempts
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
                  <TableHead>Webhook</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Status Code</TableHead>
                  <TableHead>Response Time</TableHead>
                  <TableHead>Timestamp</TableHead>
                  <TableHead class="text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                <TableRow
                  v-for="attempt in webhookAttempts"
                  :key="attempt.id"
                  class="cursor-pointer hover:bg-muted/50"
                >
                  <TableCell class="font-medium">
                    {{ attempt.webhookName }}
                  </TableCell>
                  <TableCell>
                    <Badge
                      :variant="attempt.status === 'success' ? 'default' : 'destructive'"
                    >
                      {{ attempt.status }}
                    </Badge>
                  </TableCell>
                  <TableCell>
                    <span class="text-sm">{{ attempt.statusCode }}</span>
                  </TableCell>
                  <TableCell>
                    <span class="text-sm">{{ attempt.responseTime }}</span>
                  </TableCell>
                  <TableCell class="text-muted-foreground">
                    <span class="text-sm">{{ attempt.timestamp }}</span>
                  </TableCell>
                  <TableCell class="text-right">
                    <Button size="sm" variant="ghost">
                      <Icon class="h-4 w-4" name="lucide:eye" />
                    </Button>
                  </TableCell>
                </TableRow>
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      </div>
    </div>

    <!-- Members Tab -->
    <div v-if="activeSettingsTab === 'members'">
      <div class="flex items-center justify-between">
        <div>
          <h2 class="text-2xl font-semibold">Project Members</h2>
          <p class="text-muted-foreground">
            Manage who has access to this project
          </p>
        </div>
        <Button>
          <Icon class="mr-2 h-4 w-4" name="lucide:user-plus" />
          Invite Member
        </Button>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Members</CardTitle>
          <CardDescription>
            {{ members.length }} member{{ members.length !== 1 ? 's' : '' }}
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Member</TableHead>
                <TableHead>Role</TableHead>
                <TableHead>Joined</TableHead>
                <TableHead class="text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow
                v-for="member in members"
                :key="member.id"
                class="cursor-pointer hover:bg-muted/50"
              >
                <TableCell>
                  <div class="flex items-center gap-3">
                    <Avatar class="h-8 w-8">
                      <AvatarImage :alt="member.name" :src="member.avatar" />
                      <AvatarFallback>
                        {{ member.name.split(' ').map(n => n[0]).join('') }}
                      </AvatarFallback>
                    </Avatar>
                    <div>
                      <div class="font-medium">{{ member.name }}</div>
                      <div class="text-sm text-muted-foreground">
                        {{ member.email }}
                      </div>
                    </div>
                  </div>
                </TableCell>
                <TableCell>
                  <Badge :variant="getRoleBadgeVariant(member.role)">
                    {{ member.role }}
                  </Badge>
                </TableCell>
                <TableCell class="text-muted-foreground">
                  <span class="text-sm">{{ member.joinedAt }}</span>
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
                        Change Role
                      </DropdownMenuItem>
                      <DropdownMenuSeparator />
                      <DropdownMenuItem class="text-destructive">
                        <Icon class="mr-2 h-4 w-4" name="lucide:user-x" />
                        Remove
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
  </div>
</template>
