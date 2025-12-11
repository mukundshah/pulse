<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { z } from 'zod'

import { HTTP_METHODS, IP_VERSIONS, STATUS_CODES } from '~/constants/http'

const route = useRoute()

const { projectId, type } = route.params

const schema = z.object({
  name: z.string().min(1),
  tag_ids: z.uuidv7().array(),
  region_ids: z.uuidv7().array(),
  method: z.enum(HTTP_METHODS),
  expected_status: z.number().int().min(100).max(599),
  timeout_ms: z.number().int().min(1000).max(30000),
  interval_seconds: z.number().int().min(10).max(86400),
  alert_threshold: z.number().int().min(1).max(100),
  is_enabled: z.boolean(),
  is_muted: z.boolean(),
  should_fail: z.boolean(),
  headers: z.object({
    key: z.string(),
    value: z.string(),
  }).array(),
  url: z.url(),
  ip_version: z.enum(IP_VERSIONS),
  ssl_verification: z.boolean(),
  follow_redirects: z.boolean(),
  playwright_script: z.string().optional(),
  assertions: z.object({
    source: z.enum(['status_code', 'body', 'response_time']),
    property: z.string(),
    comparison: z.enum(['equals', 'not_equals', 'contains', 'not_contains', 'is_empty', 'is_not_empty']),
    target: z.string(),
  }).array(),
})

const { handleSubmit, isSubmitting, values } = useForm({
  validationSchema: toTypedSchema(schema),
  initialValues: {
    ip_version: 'ipv4',
    method: 'GET',
    ssl_verification: true,
    follow_redirects: true,
    should_fail: false,
    headers: [{ key: '', value: '' }],
    assertions: [{ source: 'status_code', property: '', comparison: 'equals', target: '' }],
  },
})

const onSubmit = handleSubmit(async (data) => {
  console.log(data)
})
</script>

<template>
  <div>
    <form @submit="onSubmit">
      <div class="space-y-6">
        <Card>
          <CardHeader>
            <CardTitle>New {{ type }} Check</CardTitle>
          </CardHeader>
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

        <Card>
          <CardHeader>
            URL Configuration
          </CardHeader>
          <CardContent>
            <div class="flex flex-col gap-4">
              <div class="flex flex-row gap-4 justify-between">
                <div>
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
                </div>
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
                        Should Fail
                      </FormLabel>
                    </FormItem>
                  </FormField>
                </div>
              </div>
              <div class="flex flex-col gap-4">
                <FormField v-slot="{ componentField }" name="url">
                  <FormItem>
                    <FormLabel class="sr-only">
                      URL
                    </FormLabel>
                    <FormControl>
                      <InputGroup>
                        <InputGroupInput v-bind="componentField" />
                        <InputGroupAddon>
                          <FormField v-slot="{ componentField: methodComponentField }" name="method">
                            <FormItem>
                              <FormLabel class="sr-only">
                                HTTP Method
                              </FormLabel>
                            </FormItem>
                            <FormControl>
                              <Select v-bind="methodComponentField">
                                <SelectTrigger size="sm">
                                  <SelectValue placeholder="Select HTTP Method" />
                                </SelectTrigger>
                                <SelectContent>
                                  <SelectItem v-for="method in HTTP_METHODS" :key="method" :value="method">
                                    {{ method }}
                                  </SelectItem>
                                </SelectContent>
                              </Select>
                            </FormControl>
                          </FormField>
                        </InputGroupAddon>
                      </InputGroup>
                    </FormControl>
                  </FormItem>
                </FormField>

                <!-- <FormField v-slot="{ componentField }" name="status_code">
                  <FormItem>
                    <FormLabel>Status Code</FormLabel>
                    <FormControl>
                      <Select v-bind="componentField">
                        <SelectTrigger>
                          <SelectValue placeholder="Select Status Code" />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem v-for="code in STATUS_CODES" :key="code" :value="code">
                            {{ code }}
                          </SelectItem>
                        </SelectContent>
                      </Select>
                    </FormControl>
                  </FormItem>
                </FormField> -->
              </div>
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
                  <!-- <div class="grid grid-cols-[1fr_1fr_auto] gap-4">
                    <div>
                      <Label>Key</Label>
                    </div>
                    <div>
                      <Label>Value</Label>
                    </div>
                    <div></div>
                  </div> -->
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
              <div v-if="values.method && ['POST', 'PUT', 'PATCH', 'DELETE'].includes(values.method)" class="flex flex-col gap-4">
                <FormField v-slot="{ componentField }" name="body">
                  <FormItem>
                    <FormLabel>Body</FormLabel>
                    <FormControl>
                      <Input v-bind="componentField" placeholder="Body" />
                    </FormControl>
                  </FormItem>
                </FormField>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            Assertions
          </CardHeader>
          <CardContent>
            <FormFieldArray v-slot="{ fields, push, remove }" name="assertions">
              <div v-for="(field, idx) in fields" :key="field.key" class="grid grid-cols-[auto_1fr_auto_1fr_auto] gap-4">
                <FormField v-slot="{ componentField }" :name="`assertions[${idx}].source`">
                  <FormItem>
                    <FormLabel class="sr-only">
                      Source
                    </FormLabel>
                    <FormControl>
                      <Select v-bind="componentField">
                        <SelectTrigger size="sm">
                          <SelectValue placeholder="Select Source" />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="status_code">
                            Status Code
                          </SelectItem>
                          <SelectItem value="body">
                            Body
                          </SelectItem>
                          <SelectItem value="response_time">
                            Response Time
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
                      <Input v-bind="componentField" placeholder="Property" size="sm" />
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
                        <SelectTrigger size="sm">
                          <SelectValue placeholder="Select Comparison" />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="equals">
                            Equals
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
                      <Input v-bind="componentField" placeholder="Target" size="sm" />
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
            </FormFieldArray>
          </CardContent>
        </Card>
      </div>
    </form>
  </div>
</template>
