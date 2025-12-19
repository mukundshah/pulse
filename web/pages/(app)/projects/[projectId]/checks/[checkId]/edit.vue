<script lang="ts">
import type { PulseAPIRequestBody } from '#open-fetch'
import { toTypedSchema } from '@vee-validate/zod'
import { ChevronDown, ChevronUp } from 'lucide-vue-next'
import { VisuallyHidden } from 'reka-ui'
import { useForm, Field as VeeField, FieldArray as VeeFieldArray } from 'vee-validate'
import { z } from 'zod'

import { constructURL } from '@/utils/url'
import { ASSERTION_COMPARISONS, ASSERTION_PROPERTIES, ASSERTION_SOURCES, DNS_RECORD_TYPES, DNS_RESOLVER_PROTOCOLS, HTTP_METHODS, IP_VERSIONS, STATUS_CODES } from '~/constants/http'

const TYPE_TITLE_MAP = {
  http: 'HTTP Monitor',
  tcp: 'TCP Monitor',
  dns: 'DNS Monitor',
  browser: 'Browser Check',
  heartbeat: 'Heartbeat',
} as const

const CHECK_TYPE_VALUES = Object.keys(TYPE_TITLE_MAP) as (keyof typeof TYPE_TITLE_MAP)[]

const INTERVAL_MAPPING = {
  0: '5s',
  1: '10s',
  2: '15s',
  3: '30s',
  4: '1m',
  5: '2m',
  6: '5m',
  7: '10m',
  8: '15m',
  9: '30m',
  10: '1h',
  11: '2h',
  12: '3h',
  13: '6h',
  14: '12h',
  15: '24h',
} as const

const routeSchema = z.object({
  projectId: z.uuidv7(),
  checkId: z.uuidv7(),
})

// Base schema with common fields
const baseSchema = z.object({
  type: z.enum(CHECK_TYPE_VALUES),
  name: z.string().min(1),
  tag_ids: z.uuidv7().array(),
  region_ids: z.uuidv7().array().min(1),
  is_enabled: z.boolean(),
  is_muted: z.boolean(),
  should_fail: z.boolean(),
  interval: z.coerce.number().int().min(0).max(15),
  degraded_threshold: z.coerce.number().int().min(1),
  degraded_threshold_unit: z.enum(['ms', 's']),
  failed_threshold: z.coerce.number().int().min(1),
  failed_threshold_unit: z.enum(['ms', 's']),
})

// HTTP-specific fields
const httpSchema = baseSchema.extend({
  type: z.literal('http'),
  method: z.enum(HTTP_METHODS),
  url: z.url({ protocol: /^https?$/ }),
  ip_version: z.enum(IP_VERSIONS),
  skip_ssl_verification: z.boolean(),
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
  port: z.coerce.number().int().min(1).max(65535),
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
  dns_record_type: z.enum(DNS_RECORD_TYPES),
  dns_resolver: z.string().optional(),
  dns_resolver_port: z.coerce.number().int().min(1).max(65535).optional(),
  dns_resolver_protocol: z.enum(DNS_RESOLVER_PROTOCOLS),
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
    retries_count: z.coerce.number().int().min(1),
    retries_delay: z.coerce.number().int().min(1),
    retries_delay_unit: z.enum(['ms', 's']),
  }),

  // Linear retries
  z.object({
    retries: z.literal('linear'),
    retries_count: z.coerce.number().int().min(1),
    retries_delay: z.coerce.number().int().min(1),
    retries_delay_unit: z.enum(['ms', 's']),
    retries_factor: z.coerce.number().optional(),
    retries_max_delay: z.coerce.number().int().min(1).optional(),
    retries_max_delay_unit: z.enum(['ms', 's']).optional(),
  }),

  // Exponential retries
  z.object({
    retries: z.literal('exponential'),
    retries_count: z.coerce.number().int().min(1),
    retries_delay: z.coerce.number().int().min(1),
    retries_delay_unit: z.enum(['ms', 's']),
    retries_factor: z.coerce.number().min(1),
    retries_jitter: z.enum(['none', 'full', 'equal', 'decorrelated']).optional(),
    retries_jitter_factor: z.coerce.number().optional(),
    retries_max_delay: z.number().int().min(1).optional(),
    retries_max_delay_unit: z.enum(['ms', 's']).optional(),
    retries_timeout: z.coerce.number().int().min(1).optional(),
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

// Helper function to convert interval string to index
const intervalStringToIndex = (interval: string): number => {
  const index = Object.entries(INTERVAL_MAPPING).find(([, value]) => value === interval)?.[0]
  return index ? Number.parseInt(index) : 6 // default to 5m
}

// Transform check data from API format to form format
const transformCheckToFormValues = (check: any) => {
  const baseValues = {
    type: check.type,
    name: check.name,
    tag_ids: check.tags?.map((tag: any) => tag.id) ?? [],
    region_ids: check.regions?.map((region: any) => region.id) ?? [],
    is_enabled: check.is_enabled ?? true,
    is_muted: check.is_muted ?? false,
    should_fail: check.should_fail ?? false,
    interval: intervalStringToIndex(check.interval),
    degraded_threshold: check.degraded_threshold ?? 3000,
    degraded_threshold_unit: (check.degraded_threshold_unit ?? 'ms') as 'ms' | 's',
    failed_threshold: check.failed_threshold ?? 5000,
    failed_threshold_unit: (check.failed_threshold_unit ?? 'ms') as 'ms' | 's',
    retries: (check.retries ?? 'none') as 'none' | 'fixed' | 'linear' | 'exponential',
  }

  // Add retry-specific fields
  if (check.retries && check.retries !== 'none') {
    Object.assign(baseValues, {
      retries_count: check.retries_count,
      retries_delay: check.retries_delay,
      retries_delay_unit: check.retries_delay_unit,
    })

    if (check.retries === 'linear' || check.retries === 'exponential') {
      Object.assign(baseValues, {
        retries_factor: check.retries_factor,
        retries_max_delay: check.retries_max_delay,
        retries_max_delay_unit: check.retries_max_delay_unit,
      })
    }

    if (check.retries === 'exponential') {
      Object.assign(baseValues, {
        retries_jitter: check.retries_jitter ?? 'none',
        retries_jitter_factor: check.retries_jitter_factor,
        retries_timeout: check.retries_timeout,
        retries_timeout_unit: check.retries_timeout_unit,
      })
    }
  }

  if (check.type === 'http') {
    const url = constructURL({
      host: check.host,
      port: check.port,
      path: check.path,
      queryParams: check.query_params as Record<string, string> | undefined,
      secure: check.secure,
    })

    const headers = check.headers
      ? Object.entries(check.headers as Record<string, string>).map(([key, value]) => ({ key, value }))
      : [{ key: '', value: '' }]

    const assertions = check.assertions && Array.isArray(check.assertions) && check.assertions.length > 0
      ? check.assertions
      : [{ source: 'status_code' as const, property: '', comparison: 'equals' as const, target: '' }]

    return {
      ...baseValues,
      type: 'http' as const,
      ip_version: (check.ip_version ?? 'ipv4') as 'ipv4' | 'ipv6',
      method: (check.method ?? 'GET') as string,
      url,
      skip_ssl_verification: check.skip_ssl_verification ?? false,
      follow_redirects: check.follow_redirects ?? false,
      headers,
      body: check.body,
      pre_script: check.pre_script ?? '',
      post_script: check.post_script ?? '',
      assertions,
    }
  }

  if (check.type === 'tcp') {
    return {
      ...baseValues,
      type: 'tcp' as const,
      ip_version: (check.ip_version ?? 'ipv4') as 'ipv4' | 'ipv6',
      host: check.host,
      port: check.port ?? 80,
    }
  }

  if (check.type === 'dns') {
    return {
      ...baseValues,
      type: 'dns' as const,
      domain: check.host,
      dns_record_type: (check.dns_record_type ?? 'A') as string,
      dns_resolver: check.dns_resolver ?? '',
      dns_resolver_port: check.dns_resolver_port,
      dns_resolver_protocol: (check.dns_resolver_protocol ?? 'udp') as 'udp' | 'tcp',
    }
  }

  if (check.type === 'browser') {
    const url = constructURL({
      host: check.host,
      port: check.port,
      path: check.path,
      queryParams: check.query_params as Record<string, string> | undefined,
      secure: check.secure,
    })

    return {
      ...baseValues,
      type: 'browser' as const,
      url,
      pre_script: check.pre_script ?? '',
      post_script: check.post_script ?? '',
      playwright_script: check.playwright_script ?? '',
    }
  }

  if (check.type === 'heartbeat') {
    return {
      ...baseValues,
      type: 'heartbeat' as const,
    }
  }

  return baseValues
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

const { projectId, checkId } = params

const { data: check, refresh: refreshCheck } = await usePulseAPI('/internal/projects/{projectId}/checks/{checkId}', {
  path: {
    projectId,
    checkId,
  },
})

if (!check.value) {
  throw createError({
    statusCode: 404,
    statusMessage: 'Check not found',
  })
}

useHead({
  title: `Edit ${TYPE_TITLE_MAP[check.value.type as keyof typeof TYPE_TITLE_MAP]}`,
})

const { data: regions, pending: isLoadingRegions, error: regionsFetchError } = await useLazyPulseAPI('/internal/regions')

const { handleSubmit, isSubmitting, values } = useForm({
  validationSchema: toTypedSchema(schema),
  initialValues: transformCheckToFormValues(check.value),
})

const onSubmit = handleSubmit(async (data) => {
  const payload: PulseAPIRequestBody<'updateCheck'> = {
    host: '', // will be set later
    name: data.name,
    type: data.type,
    tag_ids: data.tag_ids,
    region_ids: data.region_ids,
    is_enabled: data.is_enabled,
    is_muted: data.is_muted,
    should_fail: data.should_fail,
    interval: INTERVAL_MAPPING[data.interval as keyof typeof INTERVAL_MAPPING],
    degraded_threshold: data.degraded_threshold,
    degraded_threshold_unit: data.degraded_threshold_unit,
    failed_threshold: data.failed_threshold,
    failed_threshold_unit: data.failed_threshold_unit,
  }

  if (data.type === 'http') {
    const { host, port, path, queryParams, secure } = decomposeURL(data.url)

    payload.method = data.method
    payload.host = host
    payload.port = port
    payload.path = path
    payload.query_params = queryParams
    payload.secure = secure

    payload.ip_version = data.ip_version
    payload.skip_ssl_verification = data.skip_ssl_verification
    payload.follow_redirects = data.follow_redirects
    payload.headers = Object.fromEntries(data.headers?.filter(header => header.key && header.value).map(header => [header.key, header.value]) ?? [])
    payload.body = data.body
    payload.assertions = data.assertions?.filter(assertion => assertion.source && assertion.comparison && assertion.target).map(assertion => ({
      source: assertion.source,
      property: assertion.property,
      comparison: assertion.comparison,
      target: assertion.target,
    })) as PulseAPIRequestBody<'updateCheck'>['assertions']
  } else if (data.type === 'tcp') {
    payload.host = data.host
    payload.port = data.port
    payload.ip_version = data.ip_version
  } else if (data.type === 'dns') {
    payload.host = data.domain
    payload.dns_record_type = data.dns_record_type
    payload.dns_resolver = data.dns_resolver
    payload.dns_resolver_port = data.dns_resolver_port
    payload.dns_resolver_protocol = data.dns_resolver_protocol
  } else if (data.type === 'browser') {
    const { host, port, path, queryParams, secure } = decomposeURL(data.url)
    payload.host = host
    payload.port = port
    payload.path = path
    payload.query_params = queryParams
    payload.secure = secure

    payload.pre_script = data.pre_script
    payload.post_script = data.post_script
    payload.playwright_script = data.playwright_script
  }

  await $pulseAPI('/internal/projects/{projectId}/checks/{checkId}', {
    method: 'PUT',
    path: { projectId, checkId },
    body: payload,
  })

  await refreshCheck()
  await navigateTo(`/projects/${projectId}/checks/${checkId}`)
})
</script>

<template>
  <div class="flex flex-col gap-6 p-4 md:p-6">
    <div class="flex flex-row items-center gap-2 justify-between">
      <h1 class="text-2xl font-bold">
        Edit {{ TYPE_TITLE_MAP[values.type as keyof typeof TYPE_TITLE_MAP] }}
      </h1>
      <Button
        form="edit-check-form"
        type="submit"
        variant="default"
        :loading="isSubmitting"
      >
        Update Monitor
      </Button>
    </div>

    <form id="edit-check-form" class="flex flex-col gap-6" @submit="onSubmit">
      <VisuallyHidden>
        <VeeField v-slot="{ componentField, errors }" name="type">
          <Field :data-invalid="!!errors.length">
            <FieldLabel for="type">
              Type
            </FieldLabel>
            <Input id="type" v-bind="componentField" :aria-invalid="!!errors.length" />
            <FieldError v-if="errors.length" :errors="errors" />
          </Field>
        </VeeField>
      </VisuallyHidden>

      <Card>
        <CardContent>
          <FieldGroup>
            <VeeField v-slot="{ componentField, errors }" name="name">
              <Field :data-invalid="!!errors.length">
                <FieldLabel for="name">
                  Name
                </FieldLabel>
                <Input id="name" v-bind="componentField" :aria-invalid="!!errors.length" />
                <FieldError v-if="errors.length" :errors="errors" />
              </Field>
            </VeeField>

            <!-- TODO: Add tags back in -->
            <!-- <VeeField v-slot="{ field, errors }" name="tag_ids">
              <Field :data-invalid="!!errors.length">
                <FieldLabel for="tag_ids">
                  Tags
                </FieldLabel>
                <Input id="tag_ids" v-bind="field" :aria-invalid="!!errors.length" />
                <FieldError v-if="errors.length" :errors="errors" />
              </Field>
            </VeeField> -->

            <FieldGroup class="flex flex-row gap-4">
              <VeeField v-slot="{ field, errors }" name="is_enabled">
                <Field class="w-fit" orientation="horizontal">
                  <Checkbox
                    id="is_enabled"
                    :aria-invalid="!!errors.length"
                    :model-value="field.value"
                    :name="field.name"
                    @update:model-value="field.onChange"
                  />
                  <FieldContent>
                    <FieldLabel for="is_enabled">
                      Enabled
                    </FieldLabel>
                  </FieldContent>
                </Field>
              </VeeField>
              <VeeField v-slot="{ field, errors }" name="is_muted">
                <Field class="w-fit" orientation="horizontal">
                  <Checkbox
                    id="is_muted"
                    :aria-invalid="!!errors.length"
                    :model-value="field.value"
                    :name="field.name"
                    @update:model-value="field.onChange"
                  />
                  <FieldContent>
                    <FieldLabel for="is_muted">
                      Muted
                    </FieldLabel>
                  </FieldContent>
                </Field>
              </VeeField>
            </FieldGroup>
          </FieldGroup>
        </CardContent>
      </Card>

      <Card v-if="values.type === 'http'">
        <CardHeader>
          URL Configuration
        </CardHeader>
        <CardContent>
          <div class="flex flex-col gap-4">
            <FieldGroup>
              <div class="flex flex-row gap-2">
                <VeeField v-slot="{ field, errors }" name="ip_version">
                  <Field class="w-fit" :data-invalid="!!errors.length">
                    <FieldLabel class="sr-only" for="ip_version">
                      IP Version
                    </FieldLabel>
                    <Select
                      :model-value="field.value"
                      :name="field.name"
                      @update:model-value="field.onChange"
                    >
                      <SelectTrigger
                        id="ip_version"
                        :aria-invalid="!!errors.length"
                      >
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
                    <FieldError v-if="errors.length" :errors="errors" />
                  </Field>
                </VeeField>
                <VeeField v-slot="{ field, errors }" name="method">
                  <Field class="w-fit" :data-invalid="!!errors.length">
                    <FieldLabel class="sr-only" for="method">
                      HTTP Method
                    </FieldLabel>
                    <Select
                      :model-value="field.value"
                      :name="field.name"
                      @update:model-value="field.onChange"
                    >
                      <SelectTrigger
                        id="method"
                        :aria-invalid="!!errors.length"
                      >
                        <SelectValue placeholder="Select HTTP Method" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem v-for="method in HTTP_METHODS" :key="method" :value="method">
                          {{ method }}
                        </SelectItem>
                      </SelectContent>
                    </Select>
                    <FieldError v-if="errors.length" :errors="errors" />
                  </Field>
                </VeeField>
                <VeeField v-slot="{ componentField, errors }" name="url">
                  <Field class="flex-1" :data-invalid="!!errors.length">
                    <FieldLabel class="sr-only" for="url">
                      URL
                    </FieldLabel>
                    <Input
                      id="url"
                      v-bind="componentField"
                      placeholder="https://example.com"
                      :aria-invalid="!!errors.length"
                    />
                    <FieldError v-if="errors.length" :errors="errors" />
                  </Field>
                </VeeField>
              </div>
            </FieldGroup>
            <FieldGroup>
              <div class="flex flex-row justify-end gap-4">
                <VeeField v-slot="{ field, errors }" name="skip_ssl_verification">
                  <Field class="w-fit" orientation="horizontal" :data-invalid="!!errors.length">
                    <Checkbox
                      id="skip_ssl_verification"
                      :aria-invalid="!!errors.length"
                      :model-value="field.value"
                      :name="field.name"
                      @update:model-value="field.onChange"
                    />
                    <FieldContent>
                      <FieldLabel for="skip_ssl_verification">
                        Skip SSL verification
                      </FieldLabel>
                    </FieldContent>
                    <FieldError v-if="errors.length" :errors="errors" />
                  </Field>
                </VeeField>

                <VeeField v-slot="{ field, errors }" name="follow_redirects">
                  <Field class="w-fit" orientation="horizontal" :data-invalid="!!errors.length">
                    <Checkbox
                      id="follow_redirects"
                      :aria-invalid="!!errors.length"
                      :model-value="field.value"
                      :name="field.name"
                      @update:model-value="field.onChange"
                    />
                    <FieldContent>
                      <FieldLabel for="follow_redirects">
                        Follow redirects
                      </FieldLabel>
                    </FieldContent>
                    <FieldError v-if="errors.length" :errors="errors" />
                  </Field>
                </VeeField>

                <VeeField v-slot="{ field, errors }" name="should_fail">
                  <Field class="w-fit" orientation="horizontal" :data-invalid="!!errors.length">
                    <Checkbox
                      id="should_fail"
                      :aria-invalid="!!errors.length"
                      :model-value="field.value"
                      :name="field.name"
                      @update:model-value="field.onChange"
                    />
                    <FieldContent>
                      <FieldLabel for="should_fail">
                        Should fail
                      </FieldLabel>
                    </FieldContent>
                    <FieldError v-if="errors.length" :errors="errors" />
                  </Field>
                </VeeField>
              </div>
            </FieldGroup>

            <Separator class="my-4" />

            <div class="flex flex-col gap-4">
              <VeeFieldArray v-slot="{ fields, push, remove }" name="headers">
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
                <FieldGroup
                  v-for="(field, idx) in fields"
                  :key="field.key"
                  class="grid gap-4"
                  :class="fields.length > 1 ? 'grid-cols-[1fr_2fr_auto]' : 'grid-cols-[1fr_2fr]'"
                >
                  <VeeField
                    v-slot="{ componentField: keyField, errors: keyErrors }"
                    :name="`headers[${idx}].key`"
                  >
                    <Field :data-invalid="!!keyErrors.length">
                      <FieldLabel class="sr-only" :for="`headers[${idx}].key`">
                        Key
                      </FieldLabel>
                      <Input
                        id="`headers[${idx}].key`"
                        v-bind="keyField"
                        placeholder="Key"
                        size="sm"
                        :aria-invalid="!!keyErrors.length"
                      />
                      <FieldError v-if="keyErrors.length" :errors="keyErrors" />
                    </Field>
                  </VeeField>
                  <VeeField
                    v-slot="{ componentField: valueField, errors: valueErrors }"
                    :name="`headers[${idx}].value`"
                  >
                    <Field :data-invalid="!!valueErrors.length">
                      <FieldLabel class="sr-only" :for="`headers[${idx}].value`">
                        Value
                      </FieldLabel>
                      <Input
                        id="`headers[${idx}].value`"
                        v-bind="valueField"
                        placeholder="Value"
                        size="sm"
                        :aria-invalid="!!valueErrors.length"
                      />
                      <FieldError v-if="valueErrors.length" :errors="valueErrors" />
                    </Field>
                  </VeeField>
                  <Button
                    v-if="fields.length > 1"
                    size="sm"
                    type="button"
                    variant="outline"
                    @click="remove(idx)"
                  >
                    Remove
                  </Button>
                </FieldGroup>
              </VeeFieldArray>
            </div>

            <template v-if="values.type === 'http' && values.method && ['POST', 'PUT', 'PATCH', 'DELETE'].includes(values.method)">
              <Separator class="my-4" />

              <div class="flex flex-col gap-4">
                <VeeField v-slot="{ componentField, errors }" name="body">
                  <Field :data-invalid="!!errors.length">
                    <FieldLabel for="body">
                      Body
                    </FieldLabel>
                    <Textarea
                      id="body"
                      v-bind="componentField"
                      class="min-h-[120px]"
                      :aria-invalid="!!errors.length"
                    />
                    <FieldError v-if="errors.length" :errors="errors" />
                  </Field>
                </VeeField>
              </div>
            </template>
          </div>
        </CardContent>
      </Card>

      <Card v-if="values.type === 'http'">
        <CardContent>
          <VeeFieldArray v-slot="{ fields, push, remove }" name="assertions">
            <div class="flex flex-col gap-4">
              <div class="flex flex-row items-center gap-2 justify-between mb-4">
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

              <FieldGroup
                v-for="(field, idx) in fields"
                :key="field.key"
                class="grid gap-4"
                :class="fields.length > 1 ? 'grid-cols-[180px_1fr_220px_1fr_auto]' : 'grid-cols-[180px_1fr_220px_1fr]'"
              >
                <VeeField
                  v-slot="{ field: sourceField, errors: sourceErrors }"
                  :name="`assertions[${idx}].source`"
                >
                  <Field :data-invalid="!!sourceErrors.length">
                    <FieldLabel class="sr-only" :for="`assertions[${idx}].source`">
                      Source
                    </FieldLabel>
                    <Select
                      :model-value="sourceField.value"
                      :name="sourceField.name"
                      @update:model-value="sourceField.onChange"
                    >
                      <SelectTrigger
                        :id="`assertions[${idx}].source`"
                        class="w-full"
                        size="sm"
                        :aria-invalid="!!sourceErrors.length"
                      >
                        <SelectValue placeholder="Select Source" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem v-for="source in ASSERTION_SOURCES" :key="source" :value="source">
                          {{ ASSERTION_PROPERTIES[source as keyof typeof ASSERTION_PROPERTIES].label }}
                        </SelectItem>
                      </SelectContent>
                    </Select>
                    <FieldError v-if="sourceErrors.length" :errors="sourceErrors" />
                  </Field>
                </VeeField>
                <VeeField
                  v-slot="{ componentField: propertyField, errors: propertyErrors }"
                  :name="`assertions[${idx}].property`"
                >
                  <Field :data-invalid="!!propertyErrors.length">
                    <FieldLabel class="sr-only" :for="`assertions[${idx}].property`">
                      Property
                    </FieldLabel>
                    <Input
                      :id="`assertions[${idx}].property`"
                      v-bind="propertyField"
                      placeholder="Property"
                      size="sm"
                      :aria-invalid="!!propertyErrors.length"
                    />
                    <FieldError v-if="propertyErrors.length" :errors="propertyErrors" />
                  </Field>
                </VeeField>
                <VeeField
                  v-slot="{ field: comparisonField, errors: comparisonErrors }"
                  :name="`assertions[${idx}].comparison`"
                >
                  <Field :data-invalid="!!comparisonErrors.length">
                    <FieldLabel class="sr-only" :for="`assertions[${idx}].comparison`">
                      Comparison
                    </FieldLabel>
                    <Select
                      :model-value="comparisonField.value"
                      :name="comparisonField.name"
                      @update:model-value="comparisonField.onChange"
                    >
                      <SelectTrigger
                        :id="`assertions[${idx}].comparison`"
                        class="w-full"
                        size="sm"
                        :aria-invalid="!!comparisonErrors.length"
                      >
                        <SelectValue placeholder="Select Comparison" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem v-for="comparison in (ASSERTION_PROPERTIES[values.assertions?.[idx]?.source as keyof typeof ASSERTION_PROPERTIES]!.operators as string[])" :key="comparison" :value="comparison">
                          {{ comparison.replace(/_/g, ' ') }}
                        </SelectItem>
                      </SelectContent>
                    </Select>
                    <FieldError v-if="comparisonErrors.length" :errors="comparisonErrors" />
                  </Field>
                </VeeField>
                <VeeField
                  v-slot="{ componentField: targetField, errors: targetErrors }"
                  :name="`assertions[${idx}].target`"
                >
                  <Field :data-invalid="!!targetErrors.length">
                    <FieldLabel class="sr-only" :for="`assertions[${idx}].target`">
                      Target
                    </FieldLabel>
                    <Input
                      :id="`assertions[${idx}].target`"
                      v-bind="targetField"
                      placeholder="Target"
                      size="sm"
                      :aria-invalid="!!targetErrors.length"
                    />
                    <FieldError v-if="targetErrors.length" :errors="targetErrors" />
                  </Field>
                </VeeField>
                <Button
                  v-if="fields.length > 1"
                  size="sm"
                  type="button"
                  variant="outline"
                  @click="remove(idx)"
                >
                  Remove
                </Button>
              </FieldGroup>
            </div>
          </VeeFieldArray>
        </CardContent>
      </Card>

      <Card v-if="values.type === 'tcp'">
        <CardHeader>
          TCP Configuration
        </CardHeader>
        <CardContent>
          <div class="flex flex-col gap-4">
            <div class="flex flex-row gap-2">
              <VeeField v-slot="{ field, errors }" name="ip_version">
                <Field class="w-fit" :data-invalid="!!errors.length">
                  <FieldLabel for="ip_version">
                    IP Version
                  </FieldLabel>
                  <Select
                    :model-value="field.value"
                    :name="field.name"
                    @update:model-value="field.onChange"
                  >
                    <SelectTrigger id="ip_version" :aria-invalid="!!errors.length">
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
                  <FieldError v-if="errors.length" :errors="errors" />
                </Field>
              </VeeField>
              <VeeField v-slot="{ componentField, errors }" name="host">
                <Field class="flex-1" :data-invalid="!!errors.length">
                  <FieldLabel for="host">
                    Host
                  </FieldLabel>
                  <Input
                    id="host"
                    v-bind="componentField"
                    :aria-invalid="!!errors.length"
                    :placeholder="values.ip_version === 'ipv6' ? '2001:0db8::1 or example.com' : '192.168.1.1 or example.com'"
                  />
                  <FieldError v-if="errors.length" :errors="errors" />
                </Field>
              </VeeField>
              <VeeField v-slot="{ componentField, errors }" name="port">
                <Field class="w-36" :data-invalid="!!errors.length">
                  <FieldLabel for="port">
                    Port
                  </FieldLabel>
                  <Input
                    id="port"
                    v-bind="componentField"
                    placeholder="1-65535"
                    type="number"
                    :aria-invalid="!!errors.length"
                    :max="65535"
                    :min="1"
                  />
                  <FieldError v-if="errors.length" :errors="errors" />
                </Field>
              </VeeField>
            </div>
            <div class="flex flex-row gap-4 justify-start">
              <VeeField v-slot="{ field, errors }" name="should_fail">
                <Field class="w-fit" orientation="horizontal" :data-invalid="!!errors.length">
                  <Checkbox
                    id="should_fail"
                    :aria-invalid="!!errors.length"
                    :model-value="field.value"
                    :name="field.name"
                    @update:model-value="field.onChange"
                  />
                  <FieldContent>
                    <FieldLabel for="should_fail">
                      Should fail
                    </FieldLabel>
                  </FieldContent>
                  <FieldError v-if="errors.length" :errors="errors" />
                </Field>
              </VeeField>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card v-if="values.type === 'dns'">
        <CardHeader>
          DNS Configuration
        </CardHeader>
        <CardContent>
          <div class="flex flex-col gap-6">
            <div class="flex flex-row gap-4">
              <VeeField v-slot="{ field, errors }" name="dns_record_type">
                <Field class="w-fit" :data-invalid="!!errors.length">
                  <FieldLabel for="dns_record_type">
                    Type
                  </FieldLabel>
                  <Select
                    :model-value="field.value"
                    :name="field.name"
                    @update:model-value="field.onChange"
                  >
                    <SelectTrigger id="ip_version" :aria-invalid="!!errors.length">
                      <SelectValue placeholder="Select Type" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem v-for="type in DNS_RECORD_TYPES" :key="type" :value="type">
                        {{ type.toUpperCase() }}
                      </SelectItem>
                    </SelectContent>
                  </Select>
                  <FieldError v-if="errors.length" :errors="errors" />
                </Field>
              </VeeField>
              <VeeField v-slot="{ componentField, errors }" name="domain">
                <Field class="flex-1" :data-invalid="!!errors.length">
                  <FieldLabel for="domain">
                    Domain
                  </FieldLabel>
                  <Input
                    id="domain"
                    v-bind="componentField"
                    placeholder="example.com"
                    :aria-invalid="!!errors.length"
                  />
                  <FieldError v-if="errors.length" :errors="errors" />
                </Field>
              </VeeField>
            </div>

            <div class="grid grid-cols-4 gap-4">
              <VeeField v-slot="{ componentField, errors }" name="dns_resolver">
                <Field class="col-span-2" :data-invalid="!!errors.length">
                  <FieldLabel for="dns_resolver">
                    DNS Resolver
                  </FieldLabel>
                  <Input
                    id="dns_resolver"
                    v-bind="componentField"
                    placeholder="1.1.1.1"
                    :aria-invalid="!!errors.length"
                  />
                  <FieldError v-if="errors.length" :errors="errors" />
                </Field>
              </VeeField>
              <VeeField v-slot="{ componentField, errors }" name="dns_resolver_port">
                <Field :data-invalid="!!errors.length">
                  <FieldLabel for="dns_resolver_port">
                    Port
                  </FieldLabel>
                  <Input
                    id="dns_resolver_port"
                    v-bind="componentField"
                    placeholder="53"
                    :aria-invalid="!!errors.length"
                  />
                  <FieldError v-if="errors.length" :errors="errors" />
                </Field>
              </VeeField>
              <VeeField v-slot="{ field, errors }" name="dns_resolver_protocol">
                <Field :data-invalid="!!errors.length">
                  <FieldLabel for="dns_resolver_protocol">
                    Protocol
                  </FieldLabel>
                  <Select
                    :model-value="field.value"
                    :name="field.name"
                    @update:model-value="field.onChange"
                  >
                    <SelectTrigger id="dns_resolver_protocol" :aria-invalid="!!errors.length">
                      <SelectValue placeholder="Select Protocol" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem v-for="protocol in DNS_RESOLVER_PROTOCOLS" :key="protocol" :value="protocol">
                        {{ protocol.toUpperCase() }}
                      </SelectItem>
                    </SelectContent>
                  </Select>
                  <FieldError v-if="errors.length" :errors="errors" />
                </Field>
              </VeeField>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card v-if="values.type === 'browser'">
        <CardHeader>
          Scripts
        </CardHeader>
        <CardContent>
          <div class="space-y-6">
            <VeeField v-slot="{ componentField, errors }" name="pre_script">
              <Field :data-invalid="!!errors.length">
                <FieldLabel for="pre_script">
                  Pre Script
                </FieldLabel>
                <Textarea
                  id="pre_script"
                  v-bind="componentField"
                  placeholder="Pre-execution script"
                  rows="4"
                  :aria-invalid="!!errors.length"
                />
                <FieldError v-if="errors.length" :errors="errors" />
              </Field>
            </VeeField>

            <VeeField v-slot="{ componentField, errors }" name="playwright_script">
              <Field :data-invalid="!!errors.length">
                <FieldLabel for="playwright_script">
                  Playwright Script
                </FieldLabel>
                <Textarea
                  id="playwright_script"
                  v-bind="componentField"
                  placeholder="Playwright script for browser checks"
                  rows="4"
                  :aria-invalid="!!errors.length"
                />
                <FieldError v-if="errors.length" :errors="errors" />
              </Field>
            </VeeField>

            <VeeField v-slot="{ componentField, errors }" name="post_script">
              <Field :data-invalid="!!errors.length">
                <FieldLabel for="post_script">
                  Post Script
                </FieldLabel>
                <Textarea
                  id="post_script"
                  v-bind="componentField"
                  placeholder="Post-execution script"
                  rows="4"
                  :aria-invalid="!!errors.length"
                />
                <FieldError v-if="errors.length" :errors="errors" />
              </Field>
            </VeeField>
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
              <VeeField v-slot="{ field, errors }" name="degraded_threshold">
                <Field class="[&>div]:w-36 w-fit" :data-invalid="!!errors.length">
                  <FieldLabel for="degraded_threshold">
                    Degraded after
                  </FieldLabel>
                  <NumberField
                    id="degraded_threshold"
                    :format-options="{ useGrouping: false }"
                    :min="0"
                    :model-value="field.value"
                    :name="field.name"
                    @update:model-value="field.onChange"
                  >
                    <NumberFieldContent class="*:data-[slot=input]:has-data-[slot=increment]:pr-0 *:data-[slot=input]:has-data-[slot=decrement]:pl-0">
                      <NumberFieldInput class="pe-7" :aria-invalid="!!errors.length" />
                      <NumberFieldIncrement class="border-b left-[unset] right-px -translate-y-full border-l border-input py-0 px-1.25 h-[calc(50%-1px)]">
                        <ChevronUp class="h-4 w-4" />
                      </NumberFieldIncrement>
                      <NumberFieldDecrement class="left-[unset] right-px translate-y-0 border-l border-input py-0 px-1.25 h-[calc(50%-1px)]">
                        <ChevronDown class="h-4 w-4" />
                      </NumberFieldDecrement>
                    </NumberFieldContent>
                  </NumberField>
                  <FieldError v-if="errors.length" :errors="errors" />
                </Field>
              </VeeField>
              <VeeField v-slot="{ field, errors }" name="degraded_threshold_unit">
                <Field :data-invalid="!!errors.length">
                  <FieldLabel class="sr-only" for="degraded_threshold_unit">
                    Unit
                  </FieldLabel>
                  <ToggleGroup
                    id="degraded_threshold_unit"
                    type="single"
                    variant="outline"
                    :aria-invalid="!!errors.length"
                    :model-value="field.value"
                    :name="field.name"
                    :spacing="0"
                    @update:model-value="field.onChange"
                  >
                    <ToggleGroupItem value="ms">
                      Milliseconds
                    </ToggleGroupItem>
                    <ToggleGroupItem value="s">
                      Seconds
                    </ToggleGroupItem>
                  </ToggleGroup>
                  <FieldError v-if="errors.length" :errors="errors" />
                </Field>
              </VeeField>
            </div>
            <div class="flex flex-row items-end gap-4">
              <VeeField v-slot="{ field, errors }" name="failed_threshold">
                <Field class="[&>div]:w-36 w-fit" :data-invalid="!!errors.length">
                  <FieldLabel for="failed_threshold">
                    Failed after
                  </FieldLabel>
                  <NumberField
                    id="failed_threshold"
                    :format-options="{ useGrouping: false }"
                    :min="0"
                    :model-value="field.value"
                    :name="field.name"
                    @update:model-value="field.onChange"
                  >
                    <NumberFieldContent class="*:data-[slot=input]:has-data-[slot=increment]:pr-0 *:data-[slot=input]:has-data-[slot=decrement]:pl-0">
                      <NumberFieldInput class="pe-7" :aria-invalid="!!errors.length" />
                      <NumberFieldIncrement class="border-b left-[unset] right-px -translate-y-full border-l border-input py-0 px-1.25 h-[calc(50%-1px)]">
                        <ChevronUp class="h-4 w-4" />
                      </NumberFieldIncrement>
                      <NumberFieldDecrement class="left-[unset] right-px translate-y-0 border-l border-input py-0 px-1.25 h-[calc(50%-1px)]">
                        <ChevronDown class="h-4 w-4" />
                      </NumberFieldDecrement>
                    </NumberFieldContent>
                  </NumberField>
                  <FieldError v-if="errors.length" :errors="errors" />
                </Field>
              </VeeField>
              <VeeField v-slot="{ field: unitField, errors }" name="failed_threshold_unit">
                <Field :data-invalid="!!errors.length">
                  <FieldLabel class="sr-only" for="failed_threshold_unit">
                    Unit
                  </FieldLabel>
                  <ToggleGroup
                    id="failed_threshold_unit"
                    type="single"
                    variant="outline"
                    :aria-invalid="!!errors.length"
                    :model-value="unitField.value"
                    :name="unitField.name"
                    @update:model-value="unitField.onChange"
                  >
                    <ToggleGroupItem value="ms">
                      Milliseconds
                    </ToggleGroupItem>
                    <ToggleGroupItem value="s">
                      Seconds
                    </ToggleGroupItem>
                  </ToggleGroup>
                  <FieldError v-if="errors.length" :errors="errors" />
                </Field>
              </VeeField>
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
            <VeeField v-slot="{ field, errors }" name="interval">
              <Field :data-invalid="!!errors.length">
                <FieldLabel for="interval">
                  Interval
                </FieldLabel>
                <Slider
                  class="w-full"
                  :max="15"
                  :model-value="[field.value]"
                  :step="1"
                  :ticks="INTERVAL_MAPPING"
                  @update:model-value="(value?: number[]) => field.onChange(value?.[0] ?? 0)"
                />
                <FieldError v-if="errors.length" :errors="errors" />
              </Field>
            </VeeField>
            <VeeField v-slot="{ field, errors }" name="retries">
              <FieldSet>
                <FieldLegend>Retries</FieldLegend>
                <RadioGroup
                  class="flex flex-row gap-6"
                  :model-value="field.value"
                  @update:model-value="field.onChange"
                >
                  <FieldLabel for="retries-none">
                    <Field orientation="horizontal" :data-invalid="!!errors.length">
                      <FieldContent>
                        <FieldTitle>None</FieldTitle>
                      </FieldContent>
                      <RadioGroupItem
                        id="retries-none"
                        value="none"
                        :aria-invalid="!!errors.length"
                      />
                    </Field>
                  </FieldLabel>
                  <FieldLabel for="retries-fixed">
                    <Field orientation="horizontal" :data-invalid="!!errors.length">
                      <FieldContent>
                        <FieldTitle>Fixed</FieldTitle>
                      </FieldContent>
                      <RadioGroupItem
                        id="retries-fixed"
                        value="fixed"
                        :aria-invalid="!!errors.length"
                      />
                    </Field>
                  </FieldLabel>
                  <FieldLabel for="retries-linear">
                    <Field orientation="horizontal" :data-invalid="!!errors.length">
                      <FieldContent>
                        <FieldTitle>Linear</FieldTitle>
                      </FieldContent>
                      <RadioGroupItem
                        id="retries-linear"
                        value="linear"
                        :aria-invalid="!!errors.length"
                      />
                    </Field>
                  </FieldLabel>
                  <FieldLabel for="retries-exponential">
                    <Field orientation="horizontal" :data-invalid="!!errors.length">
                      <FieldContent>
                        <FieldTitle>Exponential</FieldTitle>
                      </FieldContent>
                      <RadioGroupItem
                        id="retries-exponential"
                        value="exponential"
                        :aria-invalid="!!errors.length"
                      />
                    </Field>
                  </FieldLabel>
                </RadioGroup>
                <FieldError v-if="errors.length" :errors="errors" />
              </FieldSet>
            </VeeField>

            <template v-if="values.retries !== 'none'">
              <VeeField v-slot="{ field, errors }" name="retries_count">
                <Field class="[&>div]:w-36 w-fit" :data-invalid="!!errors.length">
                  <FieldLabel for="retries_count">
                    Number of retries
                  </FieldLabel>
                  <NumberField
                    id="retries_count"
                    :format-options="{ useGrouping: false }"
                    :min="0"
                    :model-value="field.value"
                    :name="field.name"
                    @update:model-value="field.onChange"
                  >
                    <NumberFieldContent class="*:data-[slot=input]:has-data-[slot=increment]:pr-0 *:data-[slot=input]:has-data-[slot=decrement]:pl-0">
                      <NumberFieldInput class="pe-7" :aria-invalid="!!errors.length" />
                      <NumberFieldIncrement class="border-b left-[unset] right-px -translate-y-full border-l border-input py-0 px-1.25 h-[calc(50%-1px)]">
                        <ChevronUp class="h-4 w-4" />
                      </NumberFieldIncrement>
                      <NumberFieldDecrement class="left-[unset] right-px translate-y-0 border-l border-input py-0 px-1.25 h-[calc(50%-1px)]">
                        <ChevronDown class="h-4 w-4" />
                      </NumberFieldDecrement>
                    </NumberFieldContent>
                  </NumberField>
                  <FieldError v-if="errors.length" :errors="errors" />
                </Field>
              </VeeField>

              <div class="flex flex-row items-end gap-4">
                <VeeField v-slot="{ field, errors }" name="retries_delay">
                  <Field class="[&>div]:w-36 w-fit" :data-invalid="!!errors.length">
                    <FieldLabel for="retries_delay">
                      Delay between retries
                    </FieldLabel>
                    <NumberField
                      id="retries_delay"
                      :format-options="{ useGrouping: false }"
                      :min="0"
                      :model-value="field.value"
                      :name="field.name"
                      @update:model-value="field.onChange"
                    >
                      <NumberFieldContent class="*:data-[slot=input]:has-data-[slot=increment]:pr-0 *:data-[slot=input]:has-data-[slot=decrement]:pl-0">
                        <NumberFieldInput class="pe-7" :aria-invalid="!!errors.length" />
                        <NumberFieldIncrement class="border-b left-[unset] right-px -translate-y-full border-l border-input py-0 px-1.25 h-[calc(50%-1px)]">
                          <ChevronUp class="h-4 w-4" />
                        </NumberFieldIncrement>
                        <NumberFieldDecrement class="left-[unset] right-px translate-y-0 border-l border-input py-0 px-1.25 h-[calc(50%-1px)]">
                          <ChevronDown class="h-4 w-4" />
                        </NumberFieldDecrement>
                      </NumberFieldContent>
                    </NumberField>
                    <FieldError v-if="errors.length" :errors="errors" />
                  </Field>
                </VeeField>
                <VeeField v-slot="{ field, errors }" name="retries_delay_unit">
                  <Field :data-invalid="!!errors.length">
                    <FieldLabel class="sr-only" for="retries_delay_unit">
                      Unit
                    </FieldLabel>
                    <ToggleGroup
                      id="retries_delay_unit"
                      type="single"
                      variant="outline"
                      :aria-invalid="!!errors.length"
                      :model-value="field.value"
                      :name="field.name"
                      :spacing="0"
                      @update:model-value="field.onChange"
                    >
                      <ToggleGroupItem value="ms">
                        Milliseconds
                      </ToggleGroupItem>
                      <ToggleGroupItem value="s">
                        Seconds
                      </ToggleGroupItem>
                    </ToggleGroup>
                    <FieldError v-if="errors.length" :errors="errors" />
                  </Field>
                </VeeField>
              </div>

              <template v-if="values.retries === 'linear' || values.retries === 'exponential'">
                <VeeField v-slot="{ field, errors }" name="retries_factor">
                  <Field class="[&>div]:w-36 w-fit" :data-invalid="!!errors.length">
                    <FieldLabel for="retries_factor">
                      Factor
                    </FieldLabel>
                    <NumberField
                      id="retries_factor"
                      :format-options="{ useGrouping: false }"
                      :min="0"
                      :model-value="field.value"
                      :name="field.name"
                      @update:model-value="field.onChange"
                    >
                      <NumberFieldContent class="*:data-[slot=input]:has-data-[slot=increment]:pr-0 *:data-[slot=input]:has-data-[slot=decrement]:pl-0">
                        <NumberFieldInput class="pe-7" :aria-invalid="!!errors.length" />
                        <NumberFieldIncrement class="border-b left-[unset] right-px -translate-y-full border-l border-input py-0 px-1.25 h-[calc(50%-1px)]">
                          <ChevronUp class="h-4 w-4" />
                        </NumberFieldIncrement>
                        <NumberFieldDecrement class="left-[unset] right-px translate-y-0 border-l border-input py-0 px-1.25 h-[calc(50%-1px)]">
                          <ChevronDown class="h-4 w-4" />
                        </NumberFieldDecrement>
                      </NumberFieldContent>
                    </NumberField>
                    <FieldError v-if="errors.length" :errors="errors" />
                  </Field>
                </VeeField>
              </template>

              <template v-if="values.retries === 'exponential'">
                <VeeField v-slot="{ field, errors }" name="retries_jitter">
                  <FieldSet>
                    <FieldLegend>Jitter</FieldLegend>
                    <RadioGroup
                      class="flex flex-row gap-6"
                      :model-value="field.value"
                      @update:model-value="field.onChange"
                    >
                      <FieldLabel for="retries_jitter-none">
                        <Field orientation="horizontal" :data-invalid="!!errors.length">
                          <FieldContent>
                            <FieldTitle>None</FieldTitle>
                          </FieldContent>
                          <RadioGroupItem
                            id="retries_jitter-none"
                            value="none"
                            :aria-invalid="!!errors.length"
                          />
                        </Field>
                      </FieldLabel>
                      <FieldLabel for="retries_jitter-full">
                        <Field orientation="horizontal" :data-invalid="!!errors.length">
                          <FieldContent>
                            <FieldTitle>Full</FieldTitle>
                          </FieldContent>
                          <RadioGroupItem
                            id="retries_jitter-full"
                            value="full"
                            :aria-invalid="!!errors.length"
                          />
                        </Field>
                      </FieldLabel>
                      <FieldLabel for="retries_jitter-equal">
                        <Field orientation="horizontal" :data-invalid="!!errors.length">
                          <FieldContent>
                            <FieldTitle>Equal</FieldTitle>
                          </FieldContent>
                          <RadioGroupItem
                            id="retries_jitter-equal"
                            value="equal"
                            :aria-invalid="!!errors.length"
                          />
                        </Field>
                      </FieldLabel>
                      <FieldLabel for="retries_jitter-decorrelated">
                        <Field orientation="horizontal" :data-invalid="!!errors.length">
                          <FieldContent>
                            <FieldTitle>Decorrelated</FieldTitle>
                          </FieldContent>
                          <RadioGroupItem
                            id="retries_jitter-decorrelated"
                            value="decorrelated"
                            :aria-invalid="!!errors.length"
                          />
                        </Field>
                      </FieldLabel>
                    </RadioGroup>
                    <FieldError v-if="errors.length" :errors="errors" />
                  </FieldSet>
                </VeeField>

                <template v-if="values.retries_jitter && values.retries_jitter !== 'none'">
                  <VeeField v-slot="{ field, errors }" name="retries_jitter_factor">
                    <Field class="[&>div]:w-36 w-fit" :data-invalid="!!errors.length">
                      <FieldLabel for="retries_jitter_factor">
                        Jitter Factor
                      </FieldLabel>
                      <NumberField
                        id="retries_jitter_factor"
                        :format-options="{ useGrouping: false }"
                        :min="0"
                        :model-value="field.value"
                        :name="field.name"
                        @update:model-value="field.onChange"
                      >
                        <NumberFieldContent class="*:data-[slot=input]:has-data-[slot=increment]:pr-0 *:data-[slot=input]:has-data-[slot=decrement]:pl-0">
                          <NumberFieldInput class="pe-7" :aria-invalid="!!errors.length" />
                          <NumberFieldIncrement class="border-b left-[unset] right-px -translate-y-full border-l border-input py-0 px-1.25 h-[calc(50%-1px)]">
                            <ChevronUp class="h-4 w-4" />
                          </NumberFieldIncrement>
                          <NumberFieldDecrement class="left-[unset] right-px translate-y-0 border-l border-input py-0 px-1.25 h-[calc(50%-1px)]">
                            <ChevronDown class="h-4 w-4" />
                          </NumberFieldDecrement>
                        </NumberFieldContent>
                      </NumberField>
                      <FieldError v-if="errors.length" :errors="errors" />
                    </Field>
                  </VeeField>
                </template>

                <div class="flex flex-row items-end gap-4">
                  <VeeField v-slot="{ field, errors }" name="retries_max_delay">
                    <Field class="[&>div]:w-36 w-fit" :data-invalid="!!errors.length">
                      <FieldLabel for="retries_max_delay">
                        Max delay
                      </FieldLabel>
                      <NumberField
                        id="retries_max_delay"
                        :format-options="{ useGrouping: false }"
                        :min="0"
                        :model-value="field.value"
                        :name="field.name"
                        @update:model-value="field.onChange"
                      >
                        <NumberFieldContent class="*:data-[slot=input]:has-data-[slot=increment]:pr-0 *:data-[slot=input]:has-data-[slot=decrement]:pl-0">
                          <NumberFieldInput class="pe-7" :aria-invalid="!!errors.length" />
                          <NumberFieldIncrement class="border-b left-[unset] right-px -translate-y-full border-l border-input py-0 px-1.25 h-[calc(50%-1px)]">
                            <ChevronUp class="h-4 w-4" />
                          </NumberFieldIncrement>
                          <NumberFieldDecrement class="left-[unset] right-px translate-y-0 border-l border-input py-0 px-1.25 h-[calc(50%-1px)]">
                            <ChevronDown class="h-4 w-4" />
                          </NumberFieldDecrement>
                        </NumberFieldContent>
                      </NumberField>
                      <FieldError v-if="errors.length" :errors="errors" />
                    </Field>
                  </VeeField>
                  <VeeField v-slot="{ field, errors }" name="retries_max_delay_unit">
                    <Field :data-invalid="!!errors.length">
                      <FieldLabel class="sr-only" for="retries_max_delay_unit">
                        Unit
                      </FieldLabel>
                      <ToggleGroup
                        id="retries_max_delay_unit"
                        type="single"
                        variant="outline"
                        :aria-invalid="!!errors.length"
                        :model-value="field.value"
                        :name="field.name"
                        :spacing="0"
                        @update:model-value="field.onChange"
                      >
                        <ToggleGroupItem value="ms">
                          Milliseconds
                        </ToggleGroupItem>
                        <ToggleGroupItem value="s">
                          Seconds
                        </ToggleGroupItem>
                      </ToggleGroup>
                      <FieldError v-if="errors.length" :errors="errors" />
                    </Field>
                  </VeeField>
                </div>

                <div class="flex flex-row items-end gap-4">
                  <VeeField v-slot="{ field, errors }" name="retries_timeout">
                    <Field class="[&>div]:w-36 w-fit" :data-invalid="!!errors.length">
                      <FieldLabel for="retries_timeout">
                        Timeout
                      </FieldLabel>
                      <NumberField
                        id="retries_timeout"
                        :format-options="{ useGrouping: false }"
                        :min="0"
                        :model-value="field.value"
                        :name="field.name"
                        @update:model-value="field.onChange"
                      >
                        <NumberFieldContent class="*:data-[slot=input]:has-data-[slot=increment]:pr-0 *:data-[slot=input]:has-data-[slot=decrement]:pl-0">
                          <NumberFieldInput class="pe-7" :aria-invalid="!!errors.length" />
                          <NumberFieldIncrement class="border-b left-[unset] right-px -translate-y-full border-l border-input py-0 px-1.25 h-[calc(50%-1px)]">
                            <ChevronUp class="h-4 w-4" />
                          </NumberFieldIncrement>
                          <NumberFieldDecrement class="left-[unset] right-px translate-y-0 border-l border-input py-0 px-1.25 h-[calc(50%-1px)]">
                            <ChevronDown class="h-4 w-4" />
                          </NumberFieldDecrement>
                        </NumberFieldContent>
                      </NumberField>
                      <FieldError v-if="errors.length" :errors="errors" />
                    </Field>
                  </VeeField>
                  <VeeField v-slot="{ field, errors }" name="retries_timeout_unit">
                    <Field :data-invalid="!!errors.length">
                      <FieldLabel class="sr-only" for="retries_timeout_unit">
                        Unit
                      </FieldLabel>
                      <ToggleGroup
                        id="retries_timeout_unit"
                        type="single"
                        variant="outline"
                        :aria-invalid="!!errors.length"
                        :model-value="field.value"
                        :name="field.name"
                        :spacing="0"
                        @update:model-value="field.onChange"
                      >
                        <ToggleGroupItem value="ms">
                          Milliseconds
                        </ToggleGroupItem>
                        <ToggleGroupItem value="s">
                          Seconds
                        </ToggleGroupItem>
                      </ToggleGroup>
                      <FieldError v-if="errors.length" :errors="errors" />
                    </Field>
                  </VeeField>
                </div>
              </template>
            </template>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardContent>
          <FieldGroup>
            <VeeField v-slot="{ field, errors }" name="region_ids">
              <FieldSet class="gap-4">
                <FieldLegend>Regions</FieldLegend>
                <FieldDescription class="line-clamp-1">
                  Select the regions where this check will run.
                </FieldDescription>
                <FieldGroup class="flex flex-row flex-wrap gap-2 [--radius:9999rem]" data-slot="checkbox-group">
                  <template v-if="isLoadingRegions">
                    <Skeleton
                      v-for="(width, idx) in ['180px', '200px', '150px', '220px', '190px', '160px', '210px', '175px', '205px', '185px']"
                      :key="`region-skeleton-${idx}`"
                      class="h-10"
                      :style="{ width }"
                    />
                  </template>

                  <template v-if="regions">
                    <FieldLabel
                      v-for="region in regions"
                      :key="region.id"
                      class="w-fit!"
                      :for="`region-${region.id}`"
                    >
                      <Field
                        class="gap-1.5 overflow-hidden px-3! py-1.5! transition-all duration-100 ease-linear group-has-data-[state=checked]/field-label:px-2!"
                        orientation="horizontal"
                        :data-invalid="!!errors.length"
                      >
                        <Checkbox
                          :id="`region-${region.id}`"
                          class="-ml-6 -translate-x-1 rounded-full transition-all duration-100 ease-linear data-[state=checked]:ml-0 data-[state=checked]:translate-x-0"
                          :aria-invalid="!!errors.length"
                          :model-value="field.value?.includes(region.id) ?? false"
                          @update:model-value="(checked: boolean | 'indeterminate') => {
                            const currentRegions = field.value || []
                            const newValue = checked
                              ? [...currentRegions, region.id]
                              : currentRegions.filter((id: string) => id !== region.id)
                            field.onChange(newValue)
                          }"
                        />
                        <FieldTitle>{{ region.name }}</FieldTitle>
                      </Field>
                    </FieldLabel>
                  </template>

                  <template v-if="regionsFetchError">
                    <FieldError :errors="[regionsFetchError.message]" />
                  </template>
                </FieldGroup>
              </FieldSet>
            </VeeField>
          </FieldGroup>
        </CardContent>
      </Card>
    </form>
  </div>
</template>
