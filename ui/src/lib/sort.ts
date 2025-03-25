import type { SortingState } from '@tanstack/vue-table'
import { camelCase, snakeCase } from 'lodash-es'

export type SortPrefix<T extends string> = T | `-${T}`

export function convertSortingToParams<T extends string>(
  sorting: SortingState,
): SortPrefix<T>[] {
  return sorting.map((sort) => {
    const prefix = sort.desc ? '-' : ''
    const id = snakeCase(sort.id)
    return `${prefix}${id}` as SortPrefix<T>
  })
}

export function convertParamsToSorting<T extends string>(
  sortPrefixes: SortPrefix<T>[],
): SortingState {
  return sortPrefixes.map((sortPrefix) => {
    const isDescending = sortPrefix.startsWith('-')
    const id = isDescending ? sortPrefix.slice(1) : sortPrefix
    const camelCaseId = camelCase(id)
    return {
      id: camelCaseId,
      desc: isDescending,
    }
  })
}
