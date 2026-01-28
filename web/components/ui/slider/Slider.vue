<script setup lang="ts">
import type { SliderRootEmits, SliderRootProps } from 'reka-ui'
import type { HTMLAttributes } from 'vue'

import { reactiveOmit } from '@vueuse/core'
import { SliderRange, SliderRoot, SliderThumb, SliderTrack, useForwardPropsEmits } from 'reka-ui'

import { cn } from '@/utils/style'

const props = defineProps<SliderRootProps & { class?: HTMLAttributes['class'], ticks?: number[] | Record<number, string>, max?: number }>()
const emits = defineEmits<SliderRootEmits>()

const delegatedProps = reactiveOmit(props, 'class', 'ticks')

const forwarded = useForwardPropsEmits(delegatedProps, emits)
</script>

<template>
  <SliderRoot
    v-slot="{ modelValue }"
    data-slot="slider"
    :class="cn(
      'relative flex w-full touch-none items-center select-none data-disabled:opacity-50 data-[orientation=vertical]:h-full data-[orientation=vertical]:min-h-44 data-[orientation=vertical]:w-auto data-[orientation=vertical]:flex-col',
      props.class,
    )"
    v-bind="forwarded"
  >
    <SliderTrack
      as="div"
      class="relative mb-6 grow overflow-visible rounded-full bg-muted data-[orientation=horizontal]:h-1.5 data-[orientation=horizontal]:w-full data-[orientation=vertical]:h-full data-[orientation=vertical]:w-1.5"
      data-slot="slider-track"
    >
      <SliderRange
        as="div"
        class="absolute bg-primary transition-[width,height] duration-500 data-[orientation=horizontal]:h-full data-[orientation=vertical]:w-full"
        data-slot="slider-range"
      />

      <template v-if="ticks">
        <template v-for="(tick, index) in ticks" :key="index">
          <div
            class="absolute top-0 right-0"
            :style="{ right: `${100 - (typeof tick === 'number' ? tick : index / (max ?? 1)) * 100}%` }"
          >
            <div class="relative">
              <div class="absolute top-0 size-1.5 -translate-x-1/2 rounded-full bg-primary"></div>
              <span class="absolute top-4 -translate-x-1/2 text-xs text-foreground">
                {{ tick }}
              </span>
            </div>
          </div>
        </template>
      </template>
    </SliderTrack>

    <SliderThumb
      v-for="(_, key) in modelValue"
      :key="key"
      as="div"
      class="-top-1.25 block size-4 shrink-0 rounded-full border border-primary bg-white shadow-sm ring-ring/50 transition-[color,box-shadow] hover:ring-4 focus-visible:ring-4 focus-visible:outline-hidden disabled:pointer-events-none disabled:opacity-50"
      data-slot="slider-thumb"
    />
  </SliderRoot>
</template>
