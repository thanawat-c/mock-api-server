import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Endpoint, EndpointPayload, EndpointDetail } from '../types/endpoint'
import { ElMessage } from 'element-plus'

export const useEndpointStore = defineStore('endpoint', () => {
  const endpoints = ref<Endpoint[]>([])
  const currentEndpoint = ref<EndpointDetail | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  const fetchEndpoints = async () => {
    loading.value = true
    error.value = null
    try {
      const response = await $fetch<{ data: Endpoint[] }>('/api/v1/endpoints')
      endpoints.value = response.data || []
    } catch (err: any) {
      error.value = err.message
      ElMessage.error('Failed to load endpoints')
    } finally {
      loading.value = false
    }
  }

  const fetchEndpoint = async (id: string) => {
    loading.value = true
    error.value = null
    try {
      const response = await $fetch<{ data: EndpointDetail }>(`/api/v1/endpoints/${id}`)
      currentEndpoint.value = response.data
    } catch (err: any) {
      error.value = err.message
      ElMessage.error('Failed to load endpoint details')
    } finally {
      loading.value = false
    }
  }

  const createEndpoint = async (payload: EndpointPayload) => {
    loading.value = true
    error.value = null
    try {
      const response = await $fetch<{ data: Endpoint }>('/api/v1/endpoints', {
        method: 'POST',
        body: payload
      })
      endpoints.value.push(response.data)
      ElMessage.success('Endpoint created successfully')
      return response.data
    } catch (err: any) {
      error.value = err.message
      const errorMsg = err.data?.error || 'Failed to create endpoint'
      ElMessage.error(errorMsg)
      throw err
    } finally {
      loading.value = false
    }
  }

  const updateEndpoint = async (id: string, payload: EndpointPayload) => {
    loading.value = true
    error.value = null
    try {
      await $fetch(`/api/v1/endpoints/${id}`, {
        method: 'PUT',
        body: payload
      })
      
      // Since backend doesn't return the updated data on PUT, we re-fetch it
      await fetchEndpoint(id)
      
      ElMessage.success('Endpoint updated and records re-generated successfully')
      return currentEndpoint.value
    } catch (err: any) {
      error.value = err.message
      ElMessage.error('Failed to update endpoint')
      throw err
    } finally {
      loading.value = false
    }
  }

  const deleteEndpoint = async (id: string) => {
    loading.value = true
    error.value = null
    try {
      await $fetch(`/api/v1/endpoints/${id}`, {
        method: 'DELETE'
      })
      endpoints.value = endpoints.value.filter(e => e.id !== id)
      if (currentEndpoint.value?.id === id) {
        currentEndpoint.value = null
      }
      ElMessage.success('Endpoint deleted successfully')
    } catch (err: any) {
      error.value = err.message
      ElMessage.error('Failed to delete endpoint')
      throw err
    } finally {
      loading.value = false
    }
  }

  const addRecord = async (endpointId: string) => {
    loading.value = true
    try {
      await $fetch(`/api/v1/endpoints/${endpointId}/records`, { method: 'POST' })
      await fetchEndpoint(endpointId)
      ElMessage.success('Record added successfully')
    } catch (err: any) {
      ElMessage.error(err.data?.error || 'Failed to add record')
    } finally {
      loading.value = false
    }
  }

  const updateRecord = async (recordId: string | number, payload: { data: Record<string, any> }) => {
    loading.value = true
    try {
      await $fetch(`/api/v1/records/${recordId}`, {
        method: 'PUT',
        body: payload
      })
      if (currentEndpoint.value) {
        await fetchEndpoint(currentEndpoint.value.id)
      }
      ElMessage.success('Record updated successfully')
    } catch (err: any) {
      ElMessage.error(err.data?.error || 'Failed to update record')
      throw err
    } finally {
      loading.value = false
    }
  }

  const deleteRecord = async (recordId: string | number) => {
    loading.value = true
    try {
      await $fetch(`/api/v1/records/${recordId}`, { method: 'DELETE' })
      if (currentEndpoint.value) {
        await fetchEndpoint(currentEndpoint.value.id)
      }
      ElMessage.success('Record deleted successfully')
    } catch (err: any) {
      ElMessage.error(err.data?.error || 'Failed to delete record')
    } finally {
      loading.value = false
    }
  }

  const resetRecords = async (endpointId: string | number) => {
    loading.value = true
    try {
      await $fetch(`/api/v1/endpoints/${endpointId}/reset`, { method: 'POST' })
      if (currentEndpoint.value) {
        await fetchEndpoint(currentEndpoint.value.id)
      }
      ElMessage.success('Records reset to default successfully')
    } catch (err: any) {
      ElMessage.error(err.data?.error || 'Failed to reset records')
    } finally {
      loading.value = false
    }
  }

  const importEndpoints = async (importedList: EndpointPayload[]) => {
    loading.value = true
    let successCount = 0
    let errors: string[] = []

    for (const item of importedList) {
      try {
        await $fetch('/api/v1/endpoints', {
          method: 'POST',
          body: item
        })
        successCount++
      } catch (err: any) {
        const errorMsg = err.data?.error || `Failed to import ${item.path}`
        errors.push(errorMsg)
      }
    }

    await fetchEndpoints()
    loading.value = false

    if (successCount > 0) {
      ElMessage.success(`Successfully imported ${successCount} endpoint(s)!`)
    }
    if (errors.length > 0) {
      ElMessage.warning(`${errors.length} endpoint(s) skipped or failed: ${errors[0]}`)
    }
    return { successCount, errors }
  }

  return {
    endpoints,
    currentEndpoint,
    loading,
    error,
    fetchEndpoints,
    fetchEndpoint,
    createEndpoint,
    updateEndpoint,
    deleteEndpoint,
    addRecord,
    updateRecord,
    deleteRecord,
    resetRecords,
    importEndpoints
  }
})
