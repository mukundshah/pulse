import type { useBreadcrumbItems } from '#imports'
import type { DropdownMenuItemProps } from 'reka-ui'
import type { ButtonProps } from '@/components/ui/button'

import { toValue } from 'vue'

type BreadcrumbOverride = NonNullable<Parameters<typeof useBreadcrumbItems>[0]>['overrides']

export interface LayoutAction {
  label: string
  to?: string
  icon?: string
  onClick?: (event: MouseEvent) => void
  // children?: LayoutAction[] | { label: string, children: LayoutAction[] }
  children?: { label: string, children: (Omit<LayoutAction, 'children' | 'props'> & { props?: DropdownMenuItemProps })[] } | (Omit<LayoutAction, 'children' | 'props'> & { props?: DropdownMenuItemProps })[]
  props?: Omit<ButtonProps, 'onClick' | 'asChild'>
}

export interface LayoutContext {
  breadcrumbOverrides: BreadcrumbOverride
  actions: LayoutAction[]
}

// TODO: This is a hack to get the breadcrumb overrides to work. We need to find a better way to do this. or is this the best way?
export const useLayoutContext = (context: Partial<LayoutContext> = {}) => {
  const layoutContext = useState<LayoutContext>('layout-context', () => ({
    breadcrumbOverrides: [],
    actions: [],
    ...context,
  }))

  layoutContext.value = {
    ...layoutContext.value,
    ...context,
  }

  return {
    breadcrumbOverrides: computed(() => toValue(layoutContext.value.breadcrumbOverrides)),
    actions: computed(() => toValue(layoutContext.value.actions)),
  }
}
