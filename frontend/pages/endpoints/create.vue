<template>
  <div class="max-w-4xl mx-auto">
    <div class="mb-8 flex items-center gap-4">
      <el-button circle @click="navigateTo('/')" class="!shadow-sm">
        <ArrowLeft class="w-4 h-4" />
      </el-button>
      <div>
        <h1 class="text-3xl font-bold text-gray-900 tracking-tight">Create Mock API</h1>
        <p class="mt-1 text-gray-500">Define the schema and route for your new endpoint.</p>
      </div>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
      <!-- Main Form -->
      <div class="md:col-span-2 space-y-6">
        <div class="bg-white rounded-2xl border border-gray-200 p-6 shadow-sm">
          <h3 class="text-lg font-semibold text-gray-900 mb-6 flex items-center gap-2">
            <Settings class="w-5 h-5 text-gray-500" />
            Basic Configuration
          </h3>
          
          <el-form label-position="top" @submit.prevent="submit">
            <el-form-item label="Resource Path" required>
              <template #label>
                <span class="inline-flex items-center gap-1 font-medium text-gray-700">
                  Resource Path
                  <el-tooltip content="Must start with a forward slash (/)" placement="top">
                    <Info class="w-4 h-4 text-gray-400" />
                  </el-tooltip>
                </span>
              </template>
              <el-input 
                v-model="form.path" 
                placeholder="/users" 
                size="large"
                class="!text-lg font-mono"
              >
                <template #prepend>
                  <span class="font-semibold text-gray-500">/mock</span>
                </template>
              </el-input>
            </el-form-item>

            <div class="grid grid-cols-2 gap-6 mt-6">
              <el-form-item label="HTTP Method" required>
                <template #label><span class="font-medium text-gray-700">HTTP Method</span></template>
                <el-select v-model="form.method" size="large" class="w-full">
                  <el-option label="GET" value="GET" />
                  <el-option label="POST" value="POST" />
                  <el-option label="PUT" value="PUT" />
                  <el-option label="DELETE" value="DELETE" />
                </el-select>
              </el-form-item>

              <el-form-item v-if="form.method === 'GET'" label="Mock Record Count" required>
                <template #label><span class="font-medium text-gray-700">Initial Mock Records ({{ form.count }})</span></template>
                <el-slider v-model="form.count" :min="1" :max="50" show-input />
              </el-form-item>
              <div v-else class="flex flex-col justify-center text-xs text-gray-500 bg-gray-50 p-3 rounded-lg border border-gray-100">
                <span class="font-semibold text-gray-700 mb-1">Waiting for incoming data</span>
                For {{ form.method }}, the mock endpoint starts with 0 records. It will save new records when clients send requests.
              </div>
            </div>

            <div class="mt-4">
              <el-form-item label="Response Delay (ms)">
                <template #label>
                  <span class="font-medium text-gray-700 flex items-center gap-1">
                    Response Delay ({{ form.delay || 0 }} ms)
                    <el-tooltip content="Simulate slow network latency before responding" placement="top">
                      <Info class="w-4 h-4 text-gray-400" />
                    </el-tooltip>
                  </span>
                </template>
                <el-slider v-model="form.delay" :min="0" :max="3000" :step="100" show-input />
              </el-form-item>
            </div>

            <!-- Custom Response Headers -->
            <div class="mt-6 pt-5 border-t border-gray-100">
              <div class="flex items-center justify-between mb-3">
                <span class="font-medium text-gray-700 flex items-center gap-1.5 text-sm">
                  Custom Response Headers
                  <el-tooltip content="Headers returned along with the mock response (e.g. X-Total-Count, Cache-Control)" placement="top">
                    <Info class="w-4 h-4 text-gray-400" />
                  </el-tooltip>
                </span>
                <el-button size="small" type="primary" plain class="!rounded-lg text-xs" @click="addHeader">
                  + Add Header
                </el-button>
              </div>

              <div v-if="headerList.length === 0" class="text-xs text-gray-400 italic bg-gray-50 p-3 rounded-lg border border-dashed border-gray-200">
                No custom headers defined. Click "+ Add Header" to simulate response headers.
              </div>

              <div v-else class="space-y-2">
                <div v-for="(item, idx) in headerList" :key="idx" class="flex items-center gap-2">
                  <el-input v-model="item.key" placeholder="Header Name (e.g. X-Total-Count)" size="default" class="font-mono text-xs flex-1" />
                  <el-input v-model="item.value" placeholder="Header Value (e.g. 100)" size="default" class="font-mono text-xs flex-1" />
                  <el-button type="danger" plain size="small" circle class="shrink-0" @click="removeHeader(idx)">
                    ✕
                  </el-button>
                </div>
              </div>
            </div>
          </el-form>
        </div>

        <SchemaBuilder v-model="form.schema" />

        <div class="flex justify-end gap-4 pt-4 pb-12">
          <el-button size="large" @click="navigateTo('/')" class="!px-6 !rounded-full">Cancel</el-button>
          <el-button type="primary" size="large" :loading="loading" @click="submit" class="!px-8 !rounded-full shadow-md shadow-blue-500/20">
            Generate Mock API
            <template #icon><Sparkles class="w-4 h-4 ml-1" /></template>
          </el-button>
        </div>
      </div>

      <!-- Info Sidebar -->
      <div class="space-y-6">
        <div class="bg-blue-50 rounded-2xl p-6 border border-blue-100">
          <h4 class="font-semibold text-blue-900 flex items-center gap-2 mb-3">
            <Lightbulb class="w-5 h-5 text-blue-600" />
            Pro Tip
          </h4>
          <p class="text-sm text-blue-800 leading-relaxed mb-4">
            The generated mock endpoint will automatically return an array of objects based on your defined schema.
          </p>
          <p class="text-sm text-blue-800 leading-relaxed">
            You can test this endpoint immediately after creation using any HTTP client like Postman or fetch in your app.
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useEndpointStore } from '~/stores/endpointStore'
import type { EndpointPayload } from '~/types/endpoint'
import SchemaBuilder from '~/components/endpoint/SchemaBuilder.vue'
import { ArrowLeft, Settings, Info, Sparkles, Lightbulb } from 'lucide-vue-next'
import { ElMessage } from 'element-plus'

const store = useEndpointStore()
const loading = ref(false)

const form = ref<EndpointPayload>({
  path: '/users',
  method: 'GET',
  count: 5,
  schema: [
    { name: 'id', type: 'faker.uuid' },
    { name: 'fullName', type: 'faker.name' },
    { name: 'email', type: 'faker.email' },
    { name: 'avatar', type: 'faker.avatar' }
  ]
})

const headerList = ref<{ key: string; value: string }[]>([])

const addHeader = () => {
  headerList.value.push({ key: '', value: '' })
}

const removeHeader = (index: number) => {
  headerList.value.splice(index, 1)
}

const submit = async () => {
  if (!form.value.path || !form.value.path.startsWith('/')) {
    ElMessage.warning('Path must start with a /')
    return
  }
  
  if (form.value.schema.length === 0) {
    ElMessage.warning('Please define at least one field in the schema')
    return
  }

  const customHeaders: Record<string, string> = {}
  headerList.value.forEach(h => {
    if (h.key.trim()) {
      customHeaders[h.key.trim()] = h.value
    }
  })
  form.value.headers = Object.keys(customHeaders).length > 0 ? customHeaders : undefined

  loading.value = true
  try {
    const created = await store.createEndpoint(form.value)
    if (created && created.id) {
      navigateTo(`/endpoints/${created.id}`)
    } else {
      navigateTo('/')
    }
  } catch (error) {
  } finally {
    loading.value = false
  }
}
</script>
