import type { useBreadcrumbItems } from '#imports'
import { toValue } from 'vue'

type BreadcrumbOverride = NonNullable<Parameters<typeof useBreadcrumbItems>[0]>['overrides']

export interface LayoutContext {
  breadcrumbOverrides: BreadcrumbOverride
}

// TODO: This is a hack to get the breadcrumb overrides to work. We need to find a better way to do this. or is this the best way?
export const useLayoutContext = (context: Partial<LayoutContext> = {}) => {
  const layoutContext = useState<LayoutContext>('layout-context', () => ({
    breadcrumbOverrides: [],
    ...context,
  }))

  layoutContext.value = {
    ...layoutContext.value,
    ...context,
  }

  return {
    breadcrumbOverrides: computed(() => toValue(layoutContext.value.breadcrumbOverrides)),
  }
}
