import type { SortPrefix } from '@/lib/sort'
import type { Command } from '@/types/command'
import type { Paging } from '@/types/paging'
import type { AxiosRequestConfig } from 'axios'
import http from '@/lib/http'

export const COMMAND_SORT_VALUES = ['type', 'status', 'source', 'created_at', 'completed_at'] as const
export type CommandSort = (typeof COMMAND_SORT_VALUES)[number]

export interface ListCommandsParams {
  page?: number
  pageSize?: number
  sorts?: SortPrefix<CommandSort>[]
}

const command = {
  getCommandInProgress: (opt?: AxiosRequestConfig): Promise<Command> => {
    return http.get('/commands/in-progress', opt)
  },
  listCommands: (params: ListCommandsParams): Promise<Paging<Command>> => {
    return http.get('/commands', {
      params: {
        page: params.page,
        pageSize: params.pageSize,
        sorts: params.sorts?.length !== 0 ? params.sorts?.join(',') : undefined,
      },
    })
  },
}
export default command
