<script lang="ts">
import type { PulseAPIRequestBody, PulseAPIResponse } from '#open-fetch'
import { toTypedSchema } from '@vee-validate/zod'
import { VisuallyHidden } from 'reka-ui'
import { useForm } from 'vee-validate'
import { z } from 'zod'

import { ASSERTION_COMPARISONS, ASSERTION_PROPERTIES, ASSERTION_SOURCES, HTTP_METHODS, IP_VERSIONS, STATUS_CODES } from '~/constants/http'

const TYPE_TITLE_MAP = {
  http: 'HTTP Monitor',
  tcp: 'TCP Monitor',
  dns: 'DNS Monitor',
  browser: 'Browser Check',
  heartbeat: 'Heartbeat',
} as const

const CHECK_TYPE_VALUES = Object.keys(TYPE_TITLE_MAP) as (keyof typeof TYPE_TITLE_MAP)[]

const INTERVAL_OPTIONS = [
  { value: '5s', label: '5 seconds', locked: true },
  { value: '10s', label: '10 seconds', locked: true },
  { value: '15s', label: '15 seconds', locked: true },
  { value: '30s', label: '30 seconds', locked: true },
  { value: '1m', label: '1 minute', locked: true },
  { value: '2m', label: '2 minutes', locked: true },
  { value: '5m', label: '5 minutes', locked: false },
  { value: '10m', label: '10 minutes', locked: false },
  { value: '15m', label: '15 minutes', locked: false },
  { value: '30m', label: '30 minutes', locked: false },
  { value: '1h', label: '1 hour', locked: false },
  { value: '2h', label: '2 hours', locked: false },
  { value: '3h', label: '3 hours', locked: false },
  { value: '6h', label: '6 hours', locked: false },
  { value: '12h', label: '12 hours', locked: false },
  { value: '24h', label: '24 hours', locked: false },
] as const

const routeSchema = z.object({
  projectId: z.uuidv7(),
  type: z.enum(CHECK_TYPE_VALUES),
})

// Base schema with common fields
const baseSchema = z.object({
  type: z.enum(CHECK_TYPE_VALUES),
  name: z.string().min(1),
  tag_ids: z.uuidv7().array(),
  region_ids: z.uuidv7().array(),
  is_enabled: z.boolean(),
  is_muted: z.boolean(),
  should_fail: z.boolean(),
  interval: z.enum(INTERVAL_OPTIONS.map(option => option.value)),
  degraded_threshold: z.number().int().min(1),
  degraded_threshold_unit: z.enum(['ms', 's']),
  failed_threshold: z.number().int().min(1),
  failed_threshold_unit: z.enum(['ms', 's']),
})

// HTTP-specific fields
const httpSchema = baseSchema.extend({
  type: z.literal('http'),
  method: z.enum(HTTP_METHODS),
  url: z.url({ protocol: /^https?$/ }),
  ip_version: z.enum(['ipv4', 'ipv6']),
  ssl_verification: z.boolean(),
  follow_redirects: z.boolean(),
  headers: z.object({
    key: z.string(),
    value: z.string(),
  }).array().optional(),
  body: z.any().optional(),
  pre_script: z.string().optional(),
  post_script: z.string().optional(),
  assertions: z.object({
    source: z.enum(ASSERTION_SOURCES),
    property: z.string().optional(),
    comparison: z.enum(ASSERTION_COMPARISONS),
    target: z.string(),
  }).array().optional(),
})

// Browser-specific fields
const browserSchema = baseSchema.extend({
  type: z.literal('browser'),
  url: z.url({ protocol: /^https?$/ }),
  pre_script: z.string().optional(),
  post_script: z.string().optional(),
  playwright_script: z.string(),
})

// TCP-specific fields
const tcpSchemaBase = baseSchema.extend({
  type: z.literal('tcp'),
  port: z.number().int().min(1).max(65535),
})

const tcpSchema = z.discriminatedUnion('ip_version', [
  tcpSchemaBase.extend({
    ip_version: z.literal('ipv4'),
    host: z.union([
      z.ipv4(),
      z.hostname(),
    ]),
  }),
  tcpSchemaBase.extend({
    ip_version: z.literal('ipv6'),
    host: z.union([
      z.ipv6(),
      z.hostname(),
    ]),
  }),
])

// DNS-specific fields
const dnsSchema = baseSchema.extend({
  type: z.literal('dns'),
  domain: z.hostname(),
})

// Heartbeat-specific fields
const heartbeatSchema = baseSchema.extend({
  type: z.literal('heartbeat'),
})

// Retry configuration schema with conditional validation
const retrySchema = z.discriminatedUnion('retries', [
  // No retries
  z.object({
    retries: z.literal('none'),
  }),

  z.object({
    retries: z.literal('fixed'),
    retries_count: z.number().int().min(1),
    retries_delay: z.number().int().min(1),
    retries_delay_unit: z.enum(['ms', 's']),
  }),

  // Linear retries
  z.object({
    retries: z.literal('linear'),
    retries_count: z.number().int().min(1),
    retries_delay: z.number().int().min(1),
    retries_delay_unit: z.enum(['ms', 's']),
    retries_factor: z.number().optional(),
    retries_max_delay: z.number().int().min(1).optional(),
    retries_max_delay_unit: z.enum(['ms', 's']).optional(),
  }),

  // Exponential retries
  z.object({
    retries: z.literal('exponential'),
    retries_count: z.number().int().min(1),
    retries_delay: z.number().int().min(1),
    retries_delay_unit: z.enum(['ms', 's']),
    retries_factor: z.number().min(1),
    retries_jitter: z.enum(['none', 'full', 'equal', 'decorrelated']).optional(),
    retries_jitter_factor: z.number().optional(),
    retries_max_delay: z.number().int().min(1).optional(),
    retries_max_delay_unit: z.enum(['ms', 's']).optional(),
    retries_timeout: z.number().int().min(1).optional(),
    retries_timeout_unit: z.enum(['ms', 's']).optional(),
  }),
])

const schema = z.intersection(
  z.discriminatedUnion('type', [
    httpSchema,
    browserSchema,
    tcpSchema,
    dnsSchema,
    heartbeatSchema,
  ]),
  retrySchema,
)

const getDefaultValues = (type: keyof typeof TYPE_TITLE_MAP) => {
  const baseValues = {
    type,
    tag_ids: [],
    region_ids: [],
    is_enabled: true,
    is_muted: false,
    should_fail: false,
    interval: '10m' as const,
    degraded_threshold: 3000,
    degraded_threshold_unit: 'ms' as const,
    failed_threshold: 5000,
    failed_threshold_unit: 'ms' as const,
    retries: 'none' as const,
  }
  if (type === 'http') {
    return {
      ...baseValues,
      type: 'http' as const,
      ip_version: 'ipv4' as const,
      method: 'GET' as const,
      ssl_verification: true,
      follow_redirects: false,
      headers: [{ key: '', value: '' }],
      assertions: [{ source: 'status_code' as const, property: '', comparison: 'equals' as const, target: '' }],
    }
  }
  if (type === 'tcp') {
    return {
      ...baseValues,
      type: 'tcp' as const,
      ip_version: 'ipv4' as const,
      host: '',
      port: 80,
    }
  }
  if (type === 'dns') {
    return {
      ...baseValues,
      type: 'dns' as const,
      domain: '',
    }
  }
  if (type === 'browser') {
    return {
      ...baseValues,
      type: 'browser' as const,
    }
  }
  if (type === 'heartbeat') {
    return {
      ...baseValues,
      type: 'heartbeat' as const,
    }
  }
}
</script>

<script setup lang="ts">
const route = useRoute()

const { $pulseAPI } = useNuxtApp()

const { success, data: params } = routeSchema.safeParse(route.params)

if (!success || !params) {
  throw createError({
    statusCode: 400,
    statusMessage: 'Invalid route parameters',
  })
}

const { projectId, type } = params

useHead({
  title: `New ${TYPE_TITLE_MAP[type as keyof typeof TYPE_TITLE_MAP]}`,
})

const { handleSubmit, isSubmitting, values } = useForm({
  validationSchema: toTypedSchema(schema),
  initialValues: getDefaultValues(type as keyof typeof TYPE_TITLE_MAP),
})

const onSubmit = handleSubmit(async (data) => {
  const payload: PulseAPIRequestBody<'createProjectCheck'> = {
    host: '', // will be set later
    name: data.name,
    type: data.type,
    tag_ids: data.tag_ids,
    region_ids: data.region_ids,
    is_enabled: data.is_enabled,
    is_muted: data.is_muted,
    should_fail: data.should_fail,
    interval: data.interval,
    degraded_threshold: data.degraded_threshold,
    degraded_threshold_unit: data.degraded_threshold_unit,
    failed_threshold: data.failed_threshold,
    failed_threshold_unit: data.failed_threshold_unit,
  }

  if (data.type === 'http') {
    const url = new URL(data.url)
    const https = url.protocol === 'https:'

    payload.method = data.method
    payload.host = url.hostname
    payload.port = url.port ? Number.parseInt(url.port) : (https ? 443 : 80)
    payload.path = url.pathname
    payload.query_params = Object.fromEntries(url.searchParams.entries())
    payload.secure = https

    payload.ip_version = data.ip_version
    payload.ssl_verification = data.ssl_verification
    payload.follow_redirects = data.follow_redirects
    payload.headers = Object.fromEntries(data.headers?.map(header => [header.key, header.value]) ?? [])
    payload.body = data.body
    payload.assertions = data.assertions as PulseAPIRequestBody<'createProjectCheck'>['assertions']
  } else if (data.type === 'tcp') {
    payload.host = data.host
    payload.port = data.port
    payload.ip_version = data.ip_version
  } else if (data.type === 'dns') {
    payload.host = data.domain
  } else if (data.type === 'browser') {
    const url = new URL(data.url)
    const https = url.protocol === 'https:'
    payload.host = url.hostname
    payload.port = url.port ? Number.parseInt(url.port) : (https ? 443 : 80)
    payload.path = url.pathname
    payload.query_params = Object.fromEntries(url.searchParams.entries())
    payload.secure = https

    payload.pre_script = data.pre_script
    payload.post_script = data.post_script
    payload.playwright_script = data.playwright_script
  }

  await $pulseAPI('/v1/projects/{projectId}/checks', {
    method: 'POST',
    path: { projectId },
    body: payload,
  })
})
</script>

<template>
  <div>
    <form @submit="onSubmit">
      <div class="space-y-6">
        <div class="flex flex-row items-center gap-2 justify-between">
          <h1 class="text-2xl font-bold">
            New {{ TYPE_TITLE_MAP[type as keyof typeof TYPE_TITLE_MAP] }}
          </h1>
          <Button type="submit" variant="default" :loading="isSubmitting">
            Create Monitor
          </Button>
        </div>

        <VisuallyHidden>
          <FormField v-slot="{ componentField }" name="type">
            <FormItem>
              <FormLabel>Type</FormLabel>
              <FormControl>
                <Input v-bind="componentField" />
              </FormControl>
            </FormItem>
          </FormField>
        </VisuallyHidden>

        <Card>
          <CardContent>
            <div class="space-y-6">
              <FormField v-slot="{ componentField }" name="name">
                <FormItem>
                  <FormLabel>Name</FormLabel>
                  <FormControl>
                    <Input v-bind="componentField" />
                  </FormControl>
                </FormItem>
              </FormField>
              <FormField v-slot="{ componentField }" name="tag_ids">
                <FormItem>
                  <FormLabel>Tags</FormLabel>
                  <FormControl>
                    <Input v-bind="componentField" />
                  </FormControl>
                </FormItem>
              </FormField>
              <FormField v-slot="{ componentField }" name="region_ids">
                <FormItem>
                  <FormLabel>Regions</FormLabel>
                  <FormControl>
                    <Input v-bind="componentField" />
                  </FormControl>
                </FormItem>
              </FormField>
              <div class="flex flex-row gap-4">
                <FormField v-slot="{ componentField }" name="is_enabled">
                  <FormItem class="flex flex-row items-center gap-2">
                    <FormControl>
                      <Checkbox v-bind="componentField" />
                    </FormControl>
                    <FormLabel>Enabled</FormLabel>
                  </FormItem>
                </FormField>
                <FormField v-slot="{ componentField }" name="is_muted">
                  <FormItem class="flex flex-row items-center gap-2">
                    <FormControl>
                      <Checkbox v-bind="componentField" />
                    </FormControl>
                    <FormLabel>Muted</FormLabel>
                  </FormItem>
                </FormField>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card v-if="values.type === 'http'">
          <CardHeader>
            URL Configuration
          </CardHeader>
          <CardContent>
            <div class="flex flex-col gap-4">
              <div class="flex flex-col gap-4">
                <div class="flex flex-row gap-2">
                  <FormField v-slot="{ componentField }" name="ip_version">
                    <FormItem>
                      <FormLabel class="sr-only">
                        IP Version
                      </FormLabel>
                      <FormControl>
                        <Select v-bind="componentField">
                          <SelectTrigger>
                            <SelectValue placeholder="Select IP Version" />
                          </SelectTrigger>
                          <SelectContent>
                            <SelectItem value="ipv4">
                              IPv4
                            </SelectItem>
                            <SelectItem value="ipv6">
                              IPv6
                            </SelectItem>
                          </SelectContent>
                        </Select>
                      </FormControl>
                    </FormItem>
                  </FormField>
                  <FormField v-slot="{ componentField }" name="method">
                    <FormItem>
                      <FormLabel class="sr-only">
                        HTTP Method
                      </FormLabel>
                      <FormControl>
                        <Select v-bind="componentField">
                          <SelectTrigger>
                            <SelectValue placeholder="Select HTTP Method" />
                          </SelectTrigger>
                          <SelectContent>
                            <SelectItem v-for="method in HTTP_METHODS" :key="method" :value="method">
                              {{ method }}
                            </SelectItem>
                          </SelectContent>
                        </Select>
                      </FormControl>
                    </FormItem>
                  </FormField>
                  <FormField v-slot="{ componentField }" name="url">
                    <FormItem class="w-full">
                      <FormLabel class="sr-only">
                        URL
                      </FormLabel>
                      <FormControl>
                        <Input v-bind="componentField" placeholder="https://example.com" />
                      </FormControl>
                    </FormItem>
                  </FormField>
                </div>
              </div>
              <div class="flex flex-row gap-4 justify-end">
                <div class="flex flex-row gap-4">
                  <FormField v-slot="{ componentField }" name="ssl_verification">
                    <FormItem class="flex flex-row items-center gap-2">
                      <FormControl>
                        <Checkbox v-bind="componentField" />
                      </FormControl>
                      <FormLabel>
                        Skip SSL
                      </FormLabel>
                    </FormItem>
                  </FormField>

                  <FormField v-slot="{ componentField }" name="follow_redirects">
                    <FormItem class="flex flex-row items-center gap-2">
                      <FormControl>
                        <Checkbox v-bind="componentField" />
                      </FormControl>
                      <FormLabel>
                        Follow Redirects
                      </FormLabel>
                    </FormItem>
                  </FormField>

                  <FormField v-slot="{ componentField }" name="should_fail">
                    <FormItem class="flex flex-row items-center gap-2">
                      <FormControl>
                        <Checkbox v-bind="componentField" />
                      </FormControl>
                      <FormLabel>
                        Should fail
                      </FormLabel>
                    </FormItem>
                  </FormField>
                </div>
              </div>

              <Separator class="my-4" />

              <div class="flex flex-col gap-4">
                <FormFieldArray v-slot="{ fields, push, remove }" name="headers">
                  <div class="flex flex-row items-center gap-2 justify-between">
                    <Label>Headers</Label>
                    <Button
                      size="sm"
                      type="button"
                      variant="outline"
                      @click="push({ key: '', value: '' })"
                    >
                      <Icon name="lucide:plus" />
                      Add Header
                    </Button>
                  </div>
                  <div v-for="(field, idx) in fields" :key="field.key" class="grid grid-cols-[1fr_1fr_auto] gap-4">
                    <FormField v-slot="{ componentField }" :name="`headers[${idx}].key`">
                      <FormItem>
                        <FormLabel class="sr-only">
                          Key
                        </FormLabel>
                        <FormControl>
                          <Input v-bind="componentField" placeholder="Key" />
                        </FormControl>
                      </FormItem>
                    </FormField>
                    <FormField v-slot="{ componentField }" :name="`headers[${idx}].value`">
                      <FormItem>
                        <FormLabel class="sr-only">
                          Value
                        </FormLabel>
                        <FormControl>
                          <Input v-bind="componentField" placeholder="Value" />
                        </FormControl>
                      </FormItem>
                    </FormField>
                    <Button type="button" variant="outline" @click="remove(idx)">
                      Remove
                    </Button>
                  </div>
                </FormFieldArray>
              </div>

              <template v-if="values.type === 'http' && values.method && ['POST', 'PUT', 'PATCH', 'DELETE'].includes(values.method)">
                <Separator class="my-4" />

                <div class="flex flex-col gap-4">
                  <FormField v-slot="{ componentField }" name="body">
                    <FormItem>
                      <FormLabel>Body</FormLabel>
                      <FormControl>
                        <Input v-bind="componentField" placeholder="Body" />
                      </FormControl>
                    </FormItem>
                  </FormField>
                </div>
              </template>
            </div>
          </CardContent>
        </Card>

        <Card v-if="values.type === 'tcp'">
          <CardHeader>
            TCP Configuration
          </CardHeader>
          <CardContent>
            <div class="flex flex-col gap-4">
              <div class="flex flex-row gap-2">
                <FormField v-slot="{ componentField }" name="ip_version">
                  <FormItem>
                    <FormLabel>
                      IP Version
                    </FormLabel>
                    <FormControl>
                      <Select v-bind="componentField">
                        <SelectTrigger>
                          <SelectValue placeholder="Select IP Version" />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="ipv4">
                            IPv4
                          </SelectItem>
                          <SelectItem value="ipv6">
                            IPv6
                          </SelectItem>
                        </SelectContent>
                      </Select>
                    </FormControl>
                  </FormItem>
                </FormField>
                <FormField v-slot="{ componentField }" name="host">
                  <FormItem class="flex-1">
                    <FormLabel>
                      Host
                    </FormLabel>
                    <FormControl>
                      <Input
                        v-bind="componentField"
                        :placeholder="values.ip_version === 'ipv6' ? '2001:0db8::1 or example.com' : '192.168.1.1 or example.com'"
                      />
                    </FormControl>
                  </FormItem>
                </FormField>
                <FormField v-slot="{ componentField }" name="port">
                  <FormItem>
                    <FormLabel>
                      Port
                    </FormLabel>
                    <FormControl>
                      <Input
                        v-bind="componentField"
                        placeholder="80"
                        type="number"
                        :max="65535"
                        :min="1"
                      />
                    </FormControl>
                  </FormItem>
                </FormField>
              </div>
              <div class="flex flex-row gap-4 justify-end">
                <FormField v-slot="{ componentField }" name="should_fail">
                  <FormItem class="flex flex-row items-center gap-2">
                    <FormControl>
                      <Checkbox v-bind="componentField" />
                    </FormControl>
                    <FormLabel>
                      Should fail
                    </FormLabel>
                  </FormItem>
                </FormField>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card v-if="values.type === 'dns'">
          <CardHeader>
            DNS Configuration
          </CardHeader>
          <CardContent>
            <div class="flex flex-col gap-4">
              <FormField v-slot="{ componentField }" name="domain">
                <FormItem>
                  <FormLabel>
                    Domain
                  </FormLabel>
                  <FormControl>
                    <Input v-bind="componentField" placeholder="example.com" />
                  </FormControl>
                </FormItem>
              </FormField>
              <div class="flex flex-row gap-4 justify-end">
                <FormField v-slot="{ componentField }" name="should_fail">
                  <FormItem class="flex flex-row items-center gap-2">
                    <FormControl>
                      <Checkbox v-bind="componentField" />
                    </FormControl>
                    <FormLabel>
                      Should fail
                    </FormLabel>
                  </FormItem>
                </FormField>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card v-if="values.type === 'http'">
          <FormFieldArray v-slot="{ fields, push, remove }" name="assertions">
            <CardHeader>
              <div class="flex flex-row items-center gap-2 justify-between">
                <h2 class="text-lg font-medium">
                  Assertions
                </h2>
                <Button
                  size="sm"
                  type="button"
                  variant="outline"
                  @click="push({ source: 'status_code', property: '', comparison: 'equals', target: '' })"
                >
                  <Icon name="lucide:plus" />
                  Add Assertion
                </Button>
              </div>
            </CardHeader>
            <CardContent>
              <div v-for="(field, idx) in fields" :key="field.key" class="grid grid-cols-[180px_1fr_220px_1fr_auto] gap-4">
                <FormField v-slot="{ componentField }" :name="`assertions[${idx}].source`">
                  <FormItem>
                    <FormLabel class="sr-only">
                      Source
                    </FormLabel>
                    <FormControl>
                      <Select v-bind="componentField">
                        <SelectTrigger class="w-full" size="sm">
                          <SelectValue placeholder="Select Source" />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem v-for="source in ASSERTION_SOURCES" :key="source" :value="source">
                            {{ ASSERTION_PROPERTIES[source as keyof typeof ASSERTION_PROPERTIES].label }}
                          </SelectItem>
                        </SelectContent>
                      </Select>
                    </FormControl>
                  </FormItem>
                </FormField>
                <FormField v-slot="{ componentField }" :name="`assertions[${idx}].property`">
                  <FormItem>
                    <FormLabel class="sr-only">
                      Property
                    </FormLabel>
                    <FormControl>
                      <Input
                        v-bind="componentField"
                        class="w-full"
                        placeholder="Property"
                        size="sm"
                      />
                    </FormControl>
                  </FormItem>
                </FormField>
                <FormField v-slot="{ componentField }" :name="`assertions[${idx}].comparison`">
                  <FormItem>
                    <FormLabel class="sr-only">
                      Comparison
                    </FormLabel>
                    <FormControl>
                      <Select v-bind="componentField">
                        <SelectTrigger class="w-full" size="sm">
                          <SelectValue placeholder="Select Comparison" />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem v-for="comparison in (ASSERTION_PROPERTIES[values.assertions?.[idx]?.source as keyof typeof ASSERTION_PROPERTIES]!.operators as string[])" :key="comparison" :value="comparison">
                            {{ comparison.replace(/_/g, ' ') }}
                          </SelectItem>
                        </SelectContent>
                      </Select>
                    </FormControl>
                  </FormItem>
                </FormField>
                <FormField v-slot="{ componentField }" :name="`assertions[${idx}].target`">
                  <FormItem>
                    <FormLabel class="sr-only">
                      Target
                    </FormLabel>
                    <FormControl>
                      <Input
                        v-bind="componentField"
                        class="w-full"
                        placeholder="Target"
                        size="sm"
                      />
                    </FormControl>
                  </FormItem>
                </FormField>
                <Button
                  size="sm"
                  type="button"
                  variant="outline"
                  @click="remove(idx)"
                >
                  Remove
                </Button>
              </div>
            </CardContent>
          </FormFieldArray>
        </Card>

        <Card v-if="values.type === 'browser'">
          <CardHeader>
            Scripts
          </CardHeader>
          <CardContent>
            <div class="space-y-6">
              <FormField v-slot="{ componentField }" name="pre_script">
                <FormItem>
                  <FormLabel>Pre Script</FormLabel>
                  <FormControl>
                    <Textarea v-bind="componentField" placeholder="Pre-execution script" rows="4" />
                  </FormControl>
                </FormItem>
              </FormField>

              <FormField v-slot="{ componentField }" name="playwright_script">
                <FormItem>
                  <FormLabel>Playwright Script</FormLabel>
                  <FormControl>
                    <Textarea v-bind="componentField" placeholder="Playwright script for browser checks" rows="4" />
                  </FormControl>
                </FormItem>
              </FormField>

              <FormField v-slot="{ componentField }" name="post_script">
                <FormItem>
                  <FormLabel>Post Script</FormLabel>
                  <FormControl>
                    <Textarea v-bind="componentField" placeholder="Post-execution script" rows="4" />
                  </FormControl>
                </FormItem>
              </FormField>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            Response Time Limits
          </CardHeader>
          <CardContent>
            <div class="space-y-6">
              <div class="flex flex-row items-end gap-4">
                <FormField v-slot="{ componentField }" name="degraded_threshold">
                  <FormItem>
                    <FormLabel>
                      Degraded after
                    </FormLabel>
                    <FormControl>
                      <Input
                        v-bind="componentField"
                        class="w-24"
                        type="number"
                      />
                    </FormControl>
                  </FormItem>
                </FormField>
                <FormField v-slot="{ componentField }" name="degraded_threshold_unit">
                  <FormItem>
                    <FormLabel class="sr-only">
                      Unit
                    </FormLabel>
                    <FormControl>
                      <ToggleGroup
                        v-bind="componentField"
                        type="single"
                        variant="outline"
                        :spacing="0"
                      >
                        <ToggleGroupItem value="ms">
                          Milliseconds
                        </ToggleGroupItem>
                        <ToggleGroupItem value="s">
                          Seconds
                        </ToggleGroupItem>
                      </ToggleGroup>
                    </FormControl>
                  </FormItem>
                </FormField>
              </div>
              <div class="flex flex-row items-end gap-4">
                <FormField v-slot="{ componentField }" name="failed_threshold">
                  <FormItem>
                    <FormLabel>
                      Failed after
                    </FormLabel>
                    <FormControl>
                      <Input
                        v-bind="componentField"
                        class="w-24"
                        type="number"
                      />
                    </FormControl>
                  </FormItem>
                </FormField>
                <FormField v-slot="{ componentField: unitField }" name="failed_threshold_unit">
                  <FormItem>
                    <FormLabel class="sr-only">
                      Unit
                    </FormLabel>
                    <FormControl>
                      <ToggleGroup
                        v-bind="unitField"
                        type="single"
                        variant="outline"
                      >
                        <ToggleGroupItem value="ms">
                          Milliseconds
                        </ToggleGroupItem>
                        <ToggleGroupItem value="s">
                          Seconds
                        </ToggleGroupItem>
                      </ToggleGroup>
                    </FormControl>
                  </FormItem>
                </FormField>
              </div>

              <p class="text-sm text-muted-foreground">
                Checks are hard capped at a timeout of 30 seconds, this means the <span class="font-medium">Fail after</span> threshold has a maximum value of 30 seconds.
              </p>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            Frequency & Retries
          </CardHeader>
          <CardContent>
            <div class="space-y-6">
              <FormField v-slot="{ componentField }" name="interval">
                <FormItem>
                  <FormLabel>Interval</FormLabel>
                  <FormControl>
                    <Select v-bind="componentField">
                      <SelectTrigger>
                        <SelectValue placeholder="Select Interval" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem
                          v-for="option in INTERVAL_OPTIONS"
                          :key="option.value"
                          :disabled="option.locked"
                          :value="option.value"
                        >
                          <span>
                            {{ option.label }}
                          </span>
                          <Badge v-if="option.locked" variant="outline">
                            <Icon name="lucide:lock" />
                          </Badge>
                        </SelectItem>
                      </SelectContent>
                    </Select>
                  </FormControl>
                </FormItem>
              </FormField>
              <FormField v-slot="{ componentField }" name="retries">
                <FormItem>
                  <FormLabel>Retries</FormLabel>
                  <FormControl>
                    <RadioGroup v-bind="componentField" class="flex flex-row gap-6">
                      <div class="flex flex-row items-center gap-2">
                        <RadioGroupItem value="none" />
                        <Label>None</Label>
                      </div>
                      <div class="flex flex-row items-center gap-2">
                        <RadioGroupItem value="fixed" />
                        <Label>Fixed</Label>
                      </div>
                      <div class="flex flex-row items-center gap-2">
                        <RadioGroupItem value="linear" />
                        <Label>Linear</Label>
                      </div>
                      <div class="flex flex-row items-center gap-2">
                        <RadioGroupItem value="exponential" />
                        <Label>Exponential</Label>
                      </div>
                    </RadioGroup>
                  </FormControl>
                </FormItem>
              </FormField>

              <template v-if="values.retries !== 'none'">
                <FormField v-slot="{ componentField }" name="retries_count">
                  <FormItem>
                    <FormLabel>Number of retries</FormLabel>
                    <FormControl>
                      <Input v-bind="componentField" type="number" />
                    </FormControl>
                  </FormItem>
                </FormField>

                <div class="flex flex-row items-end gap-4">
                  <FormField v-slot="{ componentField }" name="retries_delay">
                    <FormItem>
                      <FormLabel>Delay between retries</FormLabel>
                      <FormControl>
                        <Input v-bind="componentField" type="number" />
                      </FormControl>
                    </FormItem>
                  </FormField>
                  <FormField v-slot="{ componentField }" name="retries_delay_unit">
                    <FormItem>
                      <FormLabel class="sr-only">
                        Unit
                      </FormLabel>
                      <FormControl>
                        <ToggleGroup
                          v-bind="componentField"
                          type="single"
                          variant="outline"
                          :spacing="0"
                        >
                          <ToggleGroupItem value="ms">
                            Milliseconds
                          </ToggleGroupItem>
                          <ToggleGroupItem value="s">
                            Seconds
                          </ToggleGroupItem>
                        </ToggleGroup>
                      </FormControl>
                    </FormItem>
                  </FormField>
                </div>

                <template v-if="values.retries === 'linear' || values.retries === 'exponential'">
                  <FormField v-slot="{ componentField }" name="retries_factor">
                    <FormItem>
                      <FormLabel>Factor</FormLabel>
                      <FormControl>
                        <Input v-bind="componentField" step="0.1" type="number" />
                      </FormControl>
                    </FormItem>
                  </FormField>
                </template>

                <template v-if="values.retries === 'exponential'">
                  <FormField v-slot="{ componentField }" name="retries_jitter">
                    <FormItem>
                      <FormLabel>Jitter</FormLabel>
                      <FormControl>
                        <RadioGroup v-bind="componentField" class="flex flex-row gap-6">
                          <div class="flex flex-row items-center gap-2">
                            <RadioGroupItem value="none" />
                            <Label>None</Label>
                          </div>
                          <div class="flex flex-row items-center gap-2">
                            <RadioGroupItem value="full" />
                            <Label>Full</Label>
                          </div>
                          <div class="flex flex-row items-center gap-2">
                            <RadioGroupItem value="equal" />
                            <Label>Equal</Label>
                          </div>
                          <div class="flex flex-row items-center gap-2">
                            <RadioGroupItem value="decorrelated" />
                            <Label>Decorrelated</Label>
                          </div>
                        </RadioGroup>
                      </FormControl>
                    </FormItem>
                  </FormField>

                  <template v-if="values.retries_jitter && values.retries_jitter !== 'none'">
                    <FormField v-slot="{ componentField }" name="retries_jitter_factor">
                      <FormItem>
                        <FormLabel>Jitter Factor</FormLabel>
                        <FormControl>
                          <Input v-bind="componentField" step="0.1" type="number" />
                        </FormControl>
                      </FormItem>
                    </FormField>
                  </template>

                  <div class="flex flex-row items-end gap-4">
                    <FormField v-slot="{ componentField }" name="retries_max_delay">
                      <FormItem>
                        <FormLabel>Max delay</FormLabel>
                        <FormControl>
                          <Input v-bind="componentField" type="number" />
                        </FormControl>
                      </FormItem>
                    </FormField>
                    <FormField v-slot="{ componentField }" name="retries_max_delay_unit">
                      <FormItem>
                        <FormLabel class="sr-only">
                          Unit
                        </FormLabel>
                        <FormControl>
                          <ToggleGroup
                            v-bind="componentField"
                            type="single"
                            variant="outline"
                            :spacing="0"
                          >
                            <ToggleGroupItem value="ms">
                              Milliseconds
                            </ToggleGroupItem>
                            <ToggleGroupItem value="s">
                              Seconds
                            </ToggleGroupItem>
                          </ToggleGroup>
                        </FormControl>
                      </FormItem>
                    </FormField>
                  </div>

                  <div class="flex flex-row items-end gap-4">
                    <FormField v-slot="{ componentField }" name="retries_timeout">
                      <FormItem>
                        <FormLabel>Timeout</FormLabel>
                        <FormControl>
                          <Input v-bind="componentField" type="number" />
                        </FormControl>
                      </FormItem>
                    </FormField>
                    <FormField v-slot="{ componentField }" name="retries_timeout_unit">
                      <FormItem>
                        <FormLabel class="sr-only">
                          Unit
                        </FormLabel>
                        <FormControl>
                          <ToggleGroup
                            v-bind="componentField"
                            type="single"
                            variant="outline"
                            :spacing="0"
                          >
                            <ToggleGroupItem value="ms">
                              Milliseconds
                            </ToggleGroupItem>
                            <ToggleGroupItem value="s">
                              Seconds
                            </ToggleGroupItem>
                          </ToggleGroup>
                        </FormControl>
                      </FormItem>
                    </FormField>
                  </div>
                </template>
              </template>
            </div>
          </CardContent>
        </Card>
      </div>
    </form>
  </div>
</template>
