import configAPI from '@/api/config'
import { useMutation, useQuery } from '@tanstack/vue-query'

export const LOG_CONFIG_QUERY_KEY = 'logConfig'
export const HARDWARE_CONFIG_QUERY_KEY = 'hardwareConfig'
export const CLOUD_CONFIG_QUERY_KEY = 'cloudConfig'
export const GRPC_CONFIG_QUERY_KEY = 'grpcConfig'
export const HTTP_CONFIG_QUERY_KEY = 'httpConfig'
export const CARGO_CONFIG_QUERY_KEY = 'cargoConfig'
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

export function useGRPCConfigQuery() {
  return useQuery({
    queryKey: [GRPC_CONFIG_QUERY_KEY],
    queryFn: configAPI.getGrpcConfig,
  })
}

export function useGRPCConfigMutation() {
  return useMutation({
    mutationFn: configAPI.updateGrpcConfig,
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

export function useCargoConfigQuery() {
  return useQuery({
    queryKey: [CARGO_CONFIG_QUERY_KEY],
    queryFn: configAPI.getCargoConfig,
  })
}

export function useCargoConfigMutation() {
  return useMutation({
    mutationFn: configAPI.updateCargoConfig,
  })
}
