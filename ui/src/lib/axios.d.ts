/* eslint-disable unused-imports/no-unused-imports */
import axios, { AxiosRequestConfig, InternalAxiosRequestConfig } from 'axios'

// Augment the Axios module to add doNotShowLoading to InternalAxiosRequestConfig
declare module 'axios' {
  interface InternalAxiosRequestConfig {
    doNotShowLoading?: boolean
  }

  interface AxiosRequestConfig {
    doNotShowLoading?: boolean
  }
}
