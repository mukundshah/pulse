<script setup lang="ts">
import CheckFormHttp from '@/components/checks/CheckFormHttp.vue'

const route = useRoute()
const projectId = route.params.id as string
const checkType = route.params.type as string

useHead({
  title: `New ${checkType.toUpperCase()} Check`,
})

const formData = ref({
  name: '',
  type: checkType,
  url: '',
  method: 'GET',
  expected_status: 200,
  is_enabled: true,
  is_muted: false,
  timeout_ms: 10000,
  interval_seconds: 60,
  degraded_threshold_ms: 3000,
  failed_threshold_ms: 5000,
  locations: [],
})

const isSubmitting = ref(false)

const createCheck = async () => {
  isSubmitting.value = true
  try {
    const nuxtApp = useNuxtApp()
    const newCheck = await nuxtApp.$api(`/projects/${projectId}/checks`, {
      method: 'POST',
      body: formData.value,
      protected: true,
    })

    // Navigate to the new check
    if (newCheck && typeof newCheck === 'object' && 'id' in newCheck) {
      await navigateTo(`/projects/${projectId}/checks/${(newCheck as any).id}`)
    }
  } catch (error) {
    console.error('Failed to create check:', error)
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <div class="flex flex-1 flex-col gap-6">
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-4">
        <Button as-child variant="ghost" size="icon">
          <NuxtLink :to="`/projects/${projectId}`">
            <Icon class="h-4 w-4" name="lucide:arrow-left" />
          </NuxtLink>
        </Button>
        <div>
          <h1 class="text-3xl font-semibold tracking-tight">
            New {{ checkType.toUpperCase() }} Monitor
          </h1>
        </div>
      </div>
      <Button :disabled="isSubmitting" @click="createCheck">
        <Icon v-if="!isSubmitting" class="mr-2 h-4 w-4" name="lucide:save" />
        <Spinner v-else class="mr-2 h-4 w-4" />
        {{ isSubmitting ? 'Creating...' : 'Create Monitor' }} (Ctrl+S)
      </Button>
    </div>

    <!-- Section 1: Monitor Name & Tags -->
    <Card>
      <CardHeader>
        <CardTitle>Monitor Name & Tags</CardTitle>
      </CardHeader>
      <CardContent class="space-y-4">
        <div class="space-y-2">
          <Label>Monitor Name</Label>
          <Input
            v-model="formData.name"
            placeholder="HTTP Monitor #2"
            class="text-lg"
          />
        </div>
        <div class="space-y-2">
          <Label>Tags</Label>
          <Input
            placeholder="Type a tag, hit enter"
            class="w-full"
          />
        </div>
        <div class="flex items-center gap-4">
          <div class="flex items-center gap-2">
            <Checkbox
              v-model="formData.is_enabled"
              id="activated"
              :checked="formData.is_enabled !== false"
            />
            <Label for="activated">Activated</Label>
          </div>
          <div class="flex items-center gap-2">
            <Checkbox
              v-model="formData.is_muted"
              id="muted"
            />
            <Label for="muted">Muted</Label>
          </div>
        </div>
      </CardContent>
    </Card>

    <!-- Type-specific form -->
    <CheckFormHttp
      v-if="checkType === 'http'"
      v-model="formData"
    />

    <!-- Section 4: Response Time Limits -->
    <Card>
      <CardHeader>
        <CardTitle>Response Time Limits</CardTitle>
      </CardHeader>
      <CardContent class="space-y-4">
        <div class="grid gap-4 md:grid-cols-2">
          <div class="space-y-2">
            <Label>DEGRADED AFTER</Label>
            <div class="flex gap-2">
              <Input
                v-model.number="formData.degraded_threshold_ms"
                type="number"
                placeholder="3000"
              />
              <Select>
                <SelectTrigger class="w-32">
                  <SelectValue>Milliseconds</SelectValue>
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="ms">Milliseconds</SelectItem>
                  <SelectItem value="s">Seconds</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
          <div class="space-y-2">
            <Label>FAILED AFTER</Label>
            <div class="flex gap-2">
              <Input
                v-model.number="formData.failed_threshold_ms"
                type="number"
                placeholder="5000"
              />
              <Select>
                <SelectTrigger class="w-32">
                  <SelectValue>Milliseconds</SelectValue>
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="ms">Milliseconds</SelectItem>
                  <SelectItem value="s">Seconds</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
        </div>
        <p class="text-sm text-muted-foreground">
          Note: Checks have a 30-second timeout cap regardless of these settings.
        </p>
      </CardContent>
    </Card>

    <!-- Section 5: Scheduling Strategy -->
    <Card>
      <CardHeader>
        <CardTitle>Scheduling Strategy</CardTitle>
      </CardHeader>
      <CardContent>
        <RadioGroup :default-value="'round-robin'">
          <div class="flex items-center gap-2">
            <RadioGroupItem value="parallel" id="parallel" disabled />
            <Label for="parallel">Parallel runs (ENTERPRISE)</Label>
          </div>
          <div class="flex items-center gap-2">
            <RadioGroupItem value="round-robin" id="round-robin" checked />
            <Label for="round-robin">Round-robin</Label>
          </div>
        </RadioGroup>
        <p class="mt-2 text-sm text-muted-foreground">
          Round-robin distributes runs across selected locations sequentially.
        </p>
      </CardContent>
    </Card>

    <!-- Section 6: Locations -->
    <Card>
      <CardHeader>
        <CardTitle>Locations</CardTitle>
        <CardDescription>
          Select 2-3 locations for optimal monitoring coverage
        </CardDescription>
      </CardHeader>
      <CardContent>
        <p class="mb-4 text-sm text-muted-foreground">
          2 locations selected
        </p>
        <div class="grid gap-4 md:grid-cols-3">
          <!-- Location grid would go here -->
          <div class="rounded-lg border p-4 text-center">
            <p class="text-sm text-muted-foreground">Location selector component</p>
          </div>
        </div>
      </CardContent>
    </Card>
  </div>
</template>


