import type { SortPrefix } from '@/lib/sort'
import type { Command, CommandStatus } from '@/types/command'
import type { Paging } from '@/types/paging'
import type { AxiosRequestConfig } from 'axios'
import http from '@/lib/http'

export const COMMAND_SORT_VALUES = ['type', 'status', 'source', 'created_at', 'completed_at'] as const
export type CommandSort = (typeof COMMAND_SORT_VALUES)[number]

export interface ListCommandsParams {
  page?: number
  pageSize?: number
  sorts?: SortPrefix<CommandSort>[]
  statuses?: CommandStatus[]
}

const commandsAPI = {
  listCommands: (params: ListCommandsParams, axiosOpts?: AxiosRequestConfig): Promise<Paging<Command>> => {
    return http.get('/commands', {
      params: {
        page: params.page,
        pageSize: params.pageSize,
        sorts: params.sorts?.length !== 0 ? params.sorts?.join(',') : undefined,
        statuses: params.statuses?.length !== 0 ? params.statuses?.join(',') : undefined,
      },
      ...axiosOpts,
    })
  },
}
export default commandsAPI
