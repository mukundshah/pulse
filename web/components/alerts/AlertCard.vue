<script setup lang="ts">
interface Props {
  alert: {
    id: string
    status: string
    severity: string
    message: string
    triggered_at: string
    resolved_at?: string
  }
}

defineProps<Props>()
</script>

<template>
  <Card>
    <CardContent class="p-4">
      <div class="flex items-start justify-between">
        <div class="flex-1">
          <div class="mb-2 flex items-center gap-2">
            <Badge
              :variant="alert.severity === 'critical' ? 'destructive' : 'default'"
            >
              {{ alert.severity }}
            </Badge>
            <Badge variant="outline">
              {{ alert.status }}
            </Badge>
          </div>
          <p class="text-sm font-medium">
            {{ alert.message }}
          </p>
          <p class="mt-2 text-xs text-muted-foreground">
            Triggered: {{ new Date(alert.triggered_at).toLocaleString() }}
            <span v-if="alert.resolved_at">
              • Resolved: {{ new Date(alert.resolved_at).toLocaleString() }}
            </span>
          </p>
        </div>
      </div>
    </CardContent>
  </Card>
</template>


