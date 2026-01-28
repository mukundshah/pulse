<script setup lang="ts">
import type { IconCategory } from '@/constants/icons'
import { Input } from '@/components/ui/input'
import { AVAILABLE_ICONS } from '@/constants/icons'

defineProps<{
  modelValue: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const open = ref(false)
const searchQuery = ref('')

// Filter icons based on search query
const filteredCategories = computed(() => {
  if (!searchQuery.value.trim()) {
    return AVAILABLE_ICONS
  }

  const query = searchQuery.value.toLowerCase().trim()
  return AVAILABLE_ICONS
    .map((category: IconCategory) => ({
      ...category,
      icons: category.icons.filter(icon =>
        icon.toLowerCase().includes(query)
        || category.category.toLowerCase().includes(query),
      ),
    }))
    .filter((category: IconCategory) => category.icons.length > 0)
})

function selectIcon(icon: string) {
  emit('update:modelValue', icon)
  open.value = false
  searchQuery.value = ''
}
</script>

<template>
  <Popover v-model:open="open">
    <PopoverTrigger as-child>
      <button
        class="flex size-10 shrink-0 cursor-pointer items-center justify-center rounded-full bg-[#ECF4E9] text-lg text-[#1E4841] transition-colors hover:bg-[#D9E8D4]"
        type="button"
      >
        <Icon :name="modelValue" />
      </button>
    </PopoverTrigger>
    <PopoverContent align="start" class="w-82 p-0" side="bottom">
      <div class="space-y-3 overflow-clip py-3">
        <div class="space-y-2 px-3">
          <h4 class="sr-only text-sm font-semibold text-[#1E4841]">
            Select Icon
          </h4>
          <div class="relative">
            <Input
              v-model="searchQuery"
              class="w-full pr-8"
              placeholder="Search icons..."
              type="search"
            />
            <Icon
              class="pointer-events-none absolute top-1/2 right-2.5 size-4 -translate-y-1/2 text-[#6B7271]"
              name="lucide:search"
            />
          </div>
        </div>
        <div class="h-96 space-y-4 overflow-y-auto px-3 pr-0">
          <div
            v-for="category in filteredCategories"
            :key="category.category"
            class="space-y-2"
          >
            <h5 class="text-xs font-medium tracking-wide text-[#6B7271] uppercase">
              {{ category.category }}
            </h5>
            <div class="flex flex-wrap gap-2">
              <button
                v-for="icon in category.icons"
                :key="icon"
                class="flex size-10 items-center justify-center rounded-lg bg-[#F5F5F5] text-lg text-[#1E4841] transition-colors hover:bg-[#ECF4E9]"
                type="button"
                :class="{ 'bg-[#ECF4E9] ring-2 ring-[#1E4841]': modelValue === icon }"
                @click="selectIcon(icon)"
              >
                <Icon :name="icon" />
              </button>
            </div>
          </div>
          <div
            v-if="filteredCategories.length === 0"
            class="py-8 text-center text-sm text-[#6B7271]"
          >
            No icons found
          </div>
        </div>
      </div>
    </PopoverContent>
  </Popover>
</template>
