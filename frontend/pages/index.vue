<template>
  <div>
    <div class="mb-8 flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <h1 class="text-3xl font-bold text-gray-900 tracking-tight">API Endpoints</h1>
        <p class="mt-2 text-gray-500">Manage your mock API endpoints and view their generated data.</p>
      </div>

      <div class="flex items-center gap-3">
        <!-- Hidden file input for import -->
        <input 
          ref="fileInputRef" 
          type="file" 
          accept=".json,application/json" 
          class="hidden" 
          @change="handleFileChange" 
        />

        <el-button 
          v-if="store.endpoints.length > 0"
          plain 
          class="!rounded-full shadow-sm" 
          @click="exportEndpoints"
        >
          <template #icon><Download class="w-4 h-4 mr-1" /></template>
          Export JSON
        </el-button>

        <el-button 
          plain 
          class="!rounded-full shadow-sm" 
          @click="triggerImport"
        >
          <template #icon><Upload class="w-4 h-4 mr-1" /></template>
          Import JSON
        </el-button>
      </div>
    </div>

    <!-- State: Loading -->
    <div v-if="loading" class="flex flex-col items-center justify-center py-16">
      <Loader2 class="w-10 h-10 text-blue-600 animate-spin mb-4" />
      <p class="text-gray-500">Loading endpoints...</p>
    </div>

    <!-- State: Empty -->
    <div v-else-if="!loading && store.endpoints.length === 0" class="bg-white border border-gray-200 rounded-2xl p-12 text-center shadow-sm">
      <div class="mx-auto w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mb-4">
        <Database class="w-8 h-8 text-gray-400" />
      </div>
      <h3 class="text-lg font-medium text-gray-900 mb-1">No endpoints yet</h3>
      <p class="text-gray-500 mb-6 max-w-sm mx-auto">Get started by creating a new mock API endpoint or import existing endpoints from a JSON backup.</p>
      <div class="flex items-center justify-center gap-3">
        <el-button type="primary" size="large" class="!rounded-full shadow-md shadow-blue-500/20" @click="navigateTo('/endpoints/create')">
          <template #icon><Plus /></template>
          Create First Endpoint
        </el-button>
        <el-button size="large" class="!rounded-full shadow-sm" @click="triggerImport">
          <template #icon><Upload class="w-4 h-4 mr-1" /></template>
          Import JSON
        </el-button>
      </div>
    </div>

    <!-- State: Data Grid -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div v-for="endpoint in store.endpoints" :key="endpoint.id" 
        class="bg-white rounded-2xl border border-gray-200 p-6 shadow-sm hover:shadow-md transition-all duration-200 group flex flex-col h-full relative overflow-hidden">
        
        <!-- Decoration line -->
        <div class="absolute top-0 left-0 w-full h-1 bg-gradient-to-r" :class="getMethodColor(endpoint.method).bg"></div>

        <div class="flex justify-between items-start mb-4">
          <div class="flex items-center gap-3 overflow-hidden">
            <span :class="['px-2.5 py-1 text-xs font-bold rounded-md', getMethodColor(endpoint.method).text, getMethodColor(endpoint.method).badgeBg]">
              {{ endpoint.method }}
            </span>
            <h3 class="text-lg font-semibold text-gray-900 truncate" :title="endpoint.path">
              {{ endpoint.path }}
            </h3>
          </div>
          
          <el-dropdown trigger="click">
            <el-button text circle class="!p-2 -mr-2 text-gray-400 hover:text-gray-600">
              <MoreVertical class="w-5 h-5" />
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="navigateTo(`/endpoints/${endpoint.id}`)">
                  <Eye class="w-4 h-4 mr-2" /> View Details
                </el-dropdown-item>
                <el-dropdown-item @click="copyUrl(endpoint.path)">
                  <Copy class="w-4 h-4 mr-2" /> Copy URL
                </el-dropdown-item>
                <el-dropdown-item divided type="danger" @click="confirmDelete(endpoint.id)">
                  <Trash2 class="w-4 h-4 mr-2 text-red-500" /> <span class="text-red-500">Delete</span>
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>

        <div class="flex-1">
          <div class="text-sm text-gray-500 flex items-center gap-2 mb-2">
            <Database class="w-4 h-4 text-gray-400" />
            {{ endpoint.count }} Mock Records
          </div>
          <div class="text-sm text-gray-500 flex items-center gap-2 mb-4">
            <Braces class="w-4 h-4 text-gray-400" />
            {{ Object.keys(endpoint.schema || {}).length }} Fields
          </div>
        </div>
        
        <div class="mt-4 pt-4 border-t border-gray-100 flex items-center justify-between">
          <div class="text-xs text-gray-400">
            Created {{ formatDate(endpoint.createdAt) }}
          </div>
          <el-button size="small" class="!rounded-full opacity-0 group-hover:opacity-100 transition-opacity" @click="navigateTo(`/endpoints/${endpoint.id}`)">
            Explore <ArrowRight class="w-3 h-3 ml-1" />
          </el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useEndpointStore } from '~/stores/endpointStore'
import { Plus, Database, MoreVertical, Copy, Trash2, Eye, Braces, ArrowRight, Download, Upload, Loader2 } from 'lucide-vue-next'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useClipboard } from '@vueuse/core'

const store = useEndpointStore()
const loading = computed(() => store.loading)
const fileInputRef = ref<HTMLInputElement | null>(null)

const triggerImport = () => {
  fileInputRef.value?.click()
}

const handleFileChange = (e: Event) => {
  const target = e.target as HTMLInputElement
  const file = target.files?.[0]
  if (!file) return

  const reader = new FileReader()
  reader.onload = async (event) => {
    try {
      const content = event.target?.result as string
      const parsed = JSON.parse(content)
      
      let listToImport = []
      if (Array.isArray(parsed)) {
        listToImport = parsed
      } else if (parsed && Array.isArray(parsed.endpoints)) {
        listToImport = parsed.endpoints
      } else if (parsed && typeof parsed === 'object') {
        listToImport = [parsed]
      }

      if (listToImport.length === 0) {
        ElMessage.warning('No valid endpoint configuration found in file')
        return
      }

      await store.importEndpoints(listToImport)
    } catch (err: any) {
      ElMessage.error('Invalid JSON file format')
    } finally {
      if (fileInputRef.value) fileInputRef.value.value = ''
    }
  }
  reader.readAsText(file)
}

const exportEndpoints = () => {
  if (store.endpoints.length === 0) {
    ElMessage.warning('No endpoints available to export')
    return
  }

  // Format the endpoints for export
  const exportData = store.endpoints.map(ep => ({
    path: ep.path,
    method: ep.method,
    delay: ep.delay || 0,
    headers: ep.headers || {},
    count: ep.count || 5,
    schema: ep.schema || []
  }))

  const jsonStr = JSON.stringify(exportData, null, 2)
  const blob = new Blob([jsonStr], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  
  const link = document.createElement('a')
  link.href = url
  link.download = `mock-api-endpoints-${new Date().toISOString().split('T')[0]}.json`
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)

  ElMessage.success('Exported endpoints successfully!')
}

onMounted(() => {
  store.fetchEndpoints()
})

const getMethodColor = (method: string) => {
  const methods: Record<string, { text: string, bg: string, badgeBg: string }> = {
    'GET': { text: 'text-blue-700', bg: 'from-blue-400 to-blue-600', badgeBg: 'bg-blue-100' },
    'POST': { text: 'text-emerald-700', bg: 'from-emerald-400 to-emerald-600', badgeBg: 'bg-emerald-100' },
    'PUT': { text: 'text-amber-700', bg: 'from-amber-400 to-amber-600', badgeBg: 'bg-amber-100' },
    'DELETE': { text: 'text-red-700', bg: 'from-red-400 to-red-600', badgeBg: 'bg-red-100' },
  }
  return methods[method?.toUpperCase()] || { text: 'text-gray-700', bg: 'from-gray-400 to-gray-600', badgeBg: 'bg-gray-100' }
}

const { copy, copied } = useClipboard()

const copyUrl = async (path: string) => {
  // Construct the full URL
  const url = `${window.location.origin}/mock${path.startsWith('/') ? path : '/' + path}`
  try {
    await copy(url)
    if (copied.value) {
      ElMessage.success('URL copied to clipboard')
    }
  } catch (e) {
    ElMessage.error('Failed to copy URL')
  }
}

const confirmDelete = (id: string) => {
  ElMessageBox.confirm(
    'Are you sure you want to delete this endpoint? All mock data will be lost.',
    'Warning',
    {
      confirmButtonText: 'Delete',
      cancelButtonText: 'Cancel',
      type: 'warning',
      confirmButtonClass: 'el-button--danger'
    }
  ).then(async () => {
    try {
      await store.deleteEndpoint(id)
    } catch (e) {
      // Error is handled in store
    }
  }).catch(() => {})
}

const formatDate = (dateStr?: string) => {
  if (!dateStr) return 'Recently'
  const date = new Date(dateStr)
  return new Intl.DateTimeFormat('en-US', { month: 'short', day: 'numeric', year: 'numeric' }).format(date)
}
</script>
