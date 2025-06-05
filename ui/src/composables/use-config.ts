import { useMutation, useQuery } from '@tanstack/vue-query'
import configAPI from '@/api/config'

export const LOG_CONFIG_QUERY_KEY = 'logConfig'
export const HARDWARE_CONFIG_QUERY_KEY = 'hardwareConfig'
export const CLOUD_CONFIG_QUERY_KEY = 'cloudConfig'
export const HTTP_CONFIG_QUERY_KEY = 'httpConfig'
export const WIFI_CONFIG_QUERY_KEY = 'wifiConfig'
export const COMMAND_CONFIG_QUERY_KEY = 'commandConfig'
export function useLogConfigQuery() {
  return useQuery({
    queryKey: [LOG_CONFIG_QUERY_KEY],
    queryFn: configAPI.getLogConfig,
  })
}

export function useLogConfigMutation() {
  return useMutation({
    mutationFn: configAPI.updateLogConfig,
  })
}

export function useHardwareConfigQuery() {
  return useQuery({
    queryKey: [HARDWARE_CONFIG_QUERY_KEY],
    queryFn: configAPI.getHardwareConfig,
  })
}

export function useHardwareConfigMutation() {
  return useMutation({
    mutationFn: configAPI.updateHardwareConfig,
  })
}

export function useCloudConfigQuery() {
  return useQuery({
    queryKey: [CLOUD_CONFIG_QUERY_KEY],
    queryFn: configAPI.getCloudConfig,
  })
}

export function useCloudConfigMutation() {
  return useMutation({
    mutationFn: configAPI.updateCloudConfig,
  })
}

export function useHTTPConfigQuery() {
  return useQuery({
    queryKey: [HTTP_CONFIG_QUERY_KEY],
    queryFn: configAPI.getHttpConfig,
  })
}

export function useHTTPConfigMutation() {
  return useMutation({
    mutationFn: configAPI.updateHttpConfig,
  })
}

export function useWifiConfigQuery() {
  return useQuery({
    queryKey: [WIFI_CONFIG_QUERY_KEY],
    queryFn: configAPI.getWifiConfig,
  })
}

export function useWifiConfigMutation() {
  return useMutation({
    mutationFn: configAPI.updateWifiConfig,
  })
}

export function useCommandConfigQuery() {
  return useQuery({
    queryKey: [COMMAND_CONFIG_QUERY_KEY],
    queryFn: configAPI.getCommandConfig,
  })
}

export function useCommandConfigMutation() {
  return useMutation({
    mutationFn: configAPI.updateCommandConfig,
  })
}
