/* eslint-disable unused-imports/no-unused-imports */
import axios, { AxiosRequestConfig, InternalAxiosRequestConfig } from 'axios'

// Augment the Axios module to add doNotShowLoading to InternalAxiosRequestConfig, AxiosRequestConfig
// This is used to prevent the loading indicator from being shown for certain requests
declare module 'axios' {
  interface InternalAxiosRequestConfig {
    doNotShowLoading?: boolean
  }

  interface AxiosRequestConfig {
    doNotShowLoading?: boolean
  }
}
