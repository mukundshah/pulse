<script setup lang="ts">
interface Props {
  modelValue?: {
    name?: string
    url?: string
    method?: string
    expected_status?: number
    headers?: Record<string, string>
    assertions?: Array<{ source: string; property: string; comparison: string; target: string }>
    skip_ssl?: boolean
    follow_redirects?: boolean
    should_fail?: boolean
  }
}

const props = defineProps<Props>()
const emit = defineEmits<{
  'update:modelValue': [value: any]
}>()

const formData = computed({
  get: () => props.modelValue || {},
  set: (value) => emit('update:modelValue', value),
})

const methods = ['GET', 'POST', 'PUT', 'DELETE', 'PATCH', 'HEAD', 'OPTIONS']

const assertionSources = [
  { value: 'status_code', label: 'Status code' },
  { value: 'response_time', label: 'Response time' },
  { value: 'body', label: 'Body' },
  { value: 'headers', label: 'Headers' },
]

const comparisons = [
  { value: 'equals', label: 'Equals' },
  { value: 'not_equals', label: 'Not equals' },
  { value: 'contains', label: 'Contains' },
  { value: 'not_contains', label: 'Not contains' },
  { value: 'greater_than', label: 'Greater than' },
  { value: 'less_than', label: 'Less than' },
]

const assertions = computed({
  get: () => formData.value.assertions || [],
  set: (value) => {
    formData.value = { ...formData.value, assertions: value }
  },
})

const addAssertion = () => {
  assertions.value = [
    ...assertions.value,
    { source: 'status_code', property: '', comparison: 'equals', target: '' },
  ]
}

const removeAssertion = (index: number) => {
  assertions.value = assertions.value.filter((_, i) => i !== index)
}
</script>

<template>
  <div class="space-y-6">
    <!-- Section 2: Monitor a URL -->
    <Card>
      <CardHeader>
        <CardTitle>Monitor a URL</CardTitle>
        <CardDescription>
          Configure the HTTP request to monitor. You can use environment variables in your URL.
        </CardDescription>
      </CardHeader>
      <CardContent class="space-y-4">
        <div class="grid gap-4 md:grid-cols-2">
          <div class="space-y-2">
            <Label>IPv4/IPv6</Label>
            <Select v-model="formData.ip_version">
              <SelectTrigger>
                <SelectValue placeholder="IPv4" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="ipv4">IPv4</SelectItem>
                <SelectItem value="ipv6">IPv6</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </div>

        <div class="space-y-2">
          <Label>URL</Label>
          <div class="flex gap-2">
            <Select v-model="formData.method" class="w-32">
              <SelectTrigger>
                <SelectValue placeholder="GET" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem
                  v-for="method in methods"
                  :key="method"
                  :value="method"
                >
                  {{ method }}
                </SelectItem>
              </SelectContent>
            </Select>
            <Input
              v-model="formData.url"
              placeholder="https://api.example.com/health"
              class="flex-1"
            />
          </div>
        </div>

        <div class="flex items-center gap-4">
          <div class="flex items-center gap-2">
            <Checkbox
              v-model="formData.skip_ssl"
              id="skip-ssl"
            />
            <Label for="skip-ssl">Skip SSL</Label>
          </div>
          <div class="flex items-center gap-2">
            <Checkbox
              v-model="formData.follow_redirects"
              id="follow-redirects"
              :checked="formData.follow_redirects !== false"
            />
            <Label for="follow-redirects">Follow redirects</Label>
          </div>
          <div class="flex items-center gap-2">
            <Checkbox
              v-model="formData.should_fail"
              id="should-fail"
            />
            <Label for="should-fail">This request should fail</Label>
          </div>
        </div>

        <Button variant="outline">
          <Icon class="mr-2 h-4 w-4" name="lucide:play" />
          Run request (Ctrl+→)
        </Button>
      </CardContent>
    </Card>

    <!-- Section 3: Assertions -->
    <Card>
      <CardHeader>
        <CardTitle>Assertions</CardTitle>
        <CardDescription>
          Define assertions to validate the response. See documentation for more details.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>SOURCE</TableHead>
              <TableHead>PROPERTY</TableHead>
              <TableHead>COMPARISON</TableHead>
              <TableHead>TARGET</TableHead>
              <TableHead class="w-[100px]">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow
              v-for="(assertion, index) in assertions"
              :key="index"
            >
              <TableCell>
                <Select v-model="assertion.source">
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem
                      v-for="source in assertionSources"
                      :key="source.value"
                      :value="source.value"
                    >
                      {{ source.label }}
                    </SelectItem>
                  </SelectContent>
                </Select>
              </TableCell>
              <TableCell>
                <Input v-model="assertion.property" placeholder="Property" />
              </TableCell>
              <TableCell>
                <Select v-model="assertion.comparison">
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem
                      v-for="comp in comparisons"
                      :key="comp.value"
                      :value="comp.value"
                    >
                      {{ comp.label }}
                    </SelectItem>
                  </SelectContent>
                </Select>
              </TableCell>
              <TableCell>
                <Input v-model="assertion.target" placeholder="Target value" />
              </TableCell>
              <TableCell>
                <Button
                  size="icon"
                  variant="ghost"
                  @click="removeAssertion(index)"
                >
                  <Icon class="h-4 w-4" name="lucide:trash-2" />
                </Button>
              </TableCell>
            </TableRow>
            <TableRow v-if="assertions.length === 0">
              <TableCell colspan="5" class="text-center text-muted-foreground">
                No assertions configured
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
        <Button
          class="mt-4"
          variant="outline"
          @click="addAssertion"
        >
          <Icon class="mr-2 h-4 w-4" name="lucide:plus" />
          Add Assertion
        </Button>
      </CardContent>
    </Card>
  </div>
</template>


