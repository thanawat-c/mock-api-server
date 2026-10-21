<template>
  <div v-if="loading" class="flex flex-col items-center justify-center py-20">
    <Loader2 class="w-10 h-10 text-blue-600 animate-spin mb-4" />
    <p class="text-gray-500">Loading endpoint details...</p>
  </div>
  
  <div v-else-if="endpoint" class="max-w-6xl mx-auto pb-16">
    <div class="mb-6 flex items-center justify-between">
      <div class="flex items-center gap-4">
        <el-button circle @click="navigateTo('/')" class="!shadow-sm">
          <ArrowLeft class="w-4 h-4" />
        </el-button>
        <div>
          <div class="flex items-center gap-3">
            <h1 class="text-3xl font-bold text-gray-900 tracking-tight">{{ endpoint.path }}</h1>
            <span :class="['px-3 py-1 text-sm font-bold rounded-md', getMethodColor(endpoint.method).text, getMethodColor(endpoint.method).badgeBg]">
              {{ endpoint.method }}
            </span>
          </div>
          <p class="mt-1 text-gray-500 text-sm">Created on {{ formatDate(endpoint.createdAt) }}</p>
        </div>
      </div>
      
      <div class="flex items-center gap-3">
        <el-button @click="openTestUrl" class="!rounded-full shadow-sm">
          <ExternalLink class="w-4 h-4 mr-2" />
          Test Fetch
        </el-button>
        <el-button @click="exportSingleEndpoint" class="!rounded-full shadow-sm">
          <Download class="w-4 h-4 mr-2" />
          Export JSON
        </el-button>
        <el-button type="primary" @click="copyUrl" class="!rounded-full shadow-md shadow-blue-500/20">
          <Copy class="w-4 h-4 mr-2" />
          Copy URL
        </el-button>
      </div>
    </div>

    <!-- URL Bar -->
    <div class="bg-gray-900 rounded-xl p-4 mb-8 flex items-center justify-between overflow-hidden shadow-lg shadow-gray-900/10">
      <div class="flex items-center gap-3 truncate text-green-400 font-mono text-sm">
        <Terminal class="w-4 h-4 text-gray-500 shrink-0" />
        <span class="truncate">{{ getFullUrl() }}</span>
      </div>
      <el-button size="small" type="primary" plain class="shrink-0 ml-4 !border-green-500/30 !text-green-400 hover:!bg-green-500/10" @click="copyUrl">
        Copy
      </el-button>
    </div>

    <!-- Tabs -->
    <div class="bg-white rounded-2xl border border-gray-200 shadow-sm overflow-hidden">
      <el-tabs v-model="activeTab" class="px-6 pt-4 border-b border-gray-100">
        <el-tab-pane name="data">
          <template #label>
            <span class="flex items-center gap-2">
              <Database class="w-4 h-4" /> Mock Data
            </span>
          </template>
        </el-tab-pane>
        <el-tab-pane name="docs">
          <template #label>
            <span class="flex items-center gap-2">
              <BookOpen class="w-4 h-4" /> API Docs
            </span>
          </template>
        </el-tab-pane>
        <el-tab-pane name="test">
          <template #label>
            <span class="flex items-center gap-2">
              <Play class="w-4 h-4" /> Test Console
            </span>
          </template>
        </el-tab-pane>
        <el-tab-pane name="schema">
          <template #label>
            <span class="flex items-center gap-2">
              <Settings class="w-4 h-4" /> Edit Schema
            </span>
          </template>
        </el-tab-pane>
      </el-tabs>

      <div class="p-6">
        <!-- Mock Data View -->
        <div v-if="activeTab === 'data'">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-lg font-semibold text-gray-900">Generated Records <span class="text-sm font-normal text-gray-500 ml-2">({{ endpoint.records?.length || 0 }} items)</span></h3>
            <div class="flex items-center gap-3">
              <el-button type="primary" size="small" @click="handleAddRecord" class="!rounded-full shadow-sm">
                <template #icon><Plus class="w-4 h-4" /></template>
                Add Record
              </el-button>
              <el-button size="small" type="warning" plain @click="handleResetRecords" class="!rounded-full shadow-sm">
                <template #icon><RotateCcw class="w-4 h-4" /></template>
                Reset Default Data
              </el-button>
              <div class="flex items-center bg-gray-100 rounded-lg p-1">
                <button 
                  @click="viewMode = 'table'" 
                  :class="['px-3 py-1.5 text-sm rounded-md font-medium transition-colors', viewMode === 'table' ? 'bg-white text-gray-900 shadow-sm' : 'text-gray-500 hover:text-gray-700']"
                >
                  Table
                </button>
                <button 
                  @click="viewMode = 'json'" 
                  :class="['px-3 py-1.5 text-sm rounded-md font-medium transition-colors', viewMode === 'json' ? 'bg-white text-gray-900 shadow-sm' : 'text-gray-500 hover:text-gray-700']"
                >
                  JSON
                </button>
              </div>

              <!-- Visible Copy JSON Button in Header -->
              <el-button 
                v-if="viewMode === 'json'" 
                size="small" 
                @click="copyText(JSON.stringify(endpoint.records, null, 2), 'Records JSON copied!')" 
                class="!rounded-lg"
              >
                <Copy class="w-3.5 h-3.5 mr-1" /> Copy JSON
              </el-button>
            </div>
          </div>

          <!-- Table View -->
          <div v-if="viewMode === 'table'" class="border border-gray-200 rounded-xl overflow-hidden">
            <el-table :data="endpoint.records || []" style="width: 100%" max-height="600" stripe>
              <!-- Numeric Record ID Column -->
              <el-table-column label="#" width="70" align="center">
                <template #default="scope">
                  <span class="inline-flex items-center px-2 py-0.5 rounded text-xs font-bold font-mono bg-gray-100 text-gray-700 border border-gray-200">
                    {{ scope.row.id }}
                  </span>
                </template>
              </el-table-column>

              <el-table-column v-for="field in endpoint.schema || []" :key="field.name" :prop="field.name" :label="field.name" min-width="150" show-overflow-tooltip>
                <template #default="scope">
                  <div v-if="scope.row.data && typeof scope.row.data[field.name] === 'string' && scope.row.data[field.name].startsWith('http')" class="flex items-center gap-2">
                    <img v-if="scope.row.data[field.name].match(/\.(jpeg|jpg|gif|png)$/) || field.type.includes('avatar') || field.type.includes('image')" 
                         :src="scope.row.data[field.name]" 
                         class="w-8 h-8 rounded-full object-cover border border-gray-200" />
                    <a :href="scope.row.data[field.name]" target="_blank" class="text-blue-600 hover:underline truncate">{{ scope.row.data[field.name] }}</a>
                  </div>
                  <span v-else>{{ scope.row.data?.[field.name] }}</span>
                </template>
              </el-table-column>
              <el-table-column label="Actions" width="120" fixed="right" align="center">
                <template #default="scope">
                  <div class="flex items-center justify-center gap-2">
                    <el-button type="primary" text circle class="!p-2" @click="openEditRecordDialog(scope.row)">
                      <Pencil class="w-4 h-4" />
                    </el-button>
                    <el-button type="danger" text circle class="!p-2 hover:bg-red-50" @click="handleDeleteRecord(scope.row)">
                      <Trash2 class="w-4 h-4 text-red-500" />
                    </el-button>
                  </div>
                </template>
              </el-table-column>
            </el-table>
          </div>

          <!-- JSON View -->
          <div v-if="viewMode === 'json'" class="bg-gray-900 rounded-xl p-4 overflow-x-auto relative group">
            <button 
              @click="copyText(JSON.stringify(endpoint.records, null, 2), 'Records JSON copied!')"
              class="absolute top-3 right-3 opacity-0 group-hover:opacity-100 transition-opacity bg-gray-800 hover:bg-gray-700 text-gray-300 text-xs px-2.5 py-1 rounded-md border border-gray-600 flex items-center gap-1 font-mono"
            >
              <Copy class="w-3 h-3" /> Copy JSON
            </button>
            <pre class="text-sm text-green-400 font-mono m-0 whitespace-pre-wrap">{{ JSON.stringify(endpoint.records, null, 2) }}</pre>
          </div>
        </div>

        <!-- API Docs Tab (Clear Guide for Users/Developers) -->
        <div v-if="activeTab === 'docs'" class="p-6">
          <div class="max-w-4xl space-y-8">
            <div>
              <h3 class="text-xl font-bold text-gray-900 tracking-tight">API Documentation & Usage Guide</h3>
              <p class="text-sm text-gray-500 mt-1">
                Share or use these endpoints directly in your web, mobile, or backend applications.
              </p>
            </div>

            <!-- List of Endpoints with Method, Path, and Actions -->
            <div class="space-y-4">
              <div 
                v-for="doc in apiEndpointsDocs" 
                :key="doc.method + doc.samplePath" 
                class="border border-gray-200 rounded-2xl p-5 hover:border-gray-300 transition-all bg-white shadow-2xs space-y-3"
              >
                <div class="flex items-center justify-between flex-wrap gap-3">
                  <div class="flex items-center gap-3">
                    <span :class="['px-3 py-1 text-xs font-bold rounded-md font-mono', getMethodColor(doc.method).text, getMethodColor(doc.method).badgeBg]">
                      {{ doc.method }}
                    </span>
                    <span class="font-mono text-gray-900 font-semibold text-sm">
                      {{ doc.fullUrl }}
                    </span>
                  </div>

                  <div class="flex items-center gap-2">
                    <el-button size="small" @click="copyText(doc.fullUrl, 'URL copied!')" class="!rounded-lg">
                      <Copy class="w-3.5 h-3.5 mr-1" /> Copy URL
                    </el-button>
                    <el-button size="small" type="primary" plain @click="useInTestConsole(doc.method, doc.samplePath)" class="!rounded-lg">
                      <Play class="w-3.5 h-3.5 mr-1" /> Try in Console
                    </el-button>
                  </div>
                </div>

                <p class="text-xs text-gray-600 pl-1">
                  {{ doc.description }}
                </p>

                <!-- Sample Payload if available -->
                <div v-if="doc.body" class="bg-gray-50 rounded-xl p-3 border border-gray-100 mt-2">
                  <div class="flex items-center justify-between text-xs text-gray-500 font-semibold mb-1">
                    <span>Expected Request Body (JSON):</span>
                    <button @click="copyText(doc.body, 'Payload copied!')" class="text-blue-600 hover:underline text-xs flex items-center gap-1">
                      <Copy class="w-3 h-3" /> Copy Body
                    </button>
                  </div>
                  <pre class="text-xs font-mono text-gray-800 m-0 overflow-x-auto whitespace-pre-wrap">{{ doc.body }}</pre>
                </div>
              </div>
            </div>

            <!-- Ready-to-use Code Snippets (JavaScript / cURL) -->
            <div class="bg-gray-900 rounded-2xl p-6 text-white space-y-4">
              <div class="flex items-center justify-between">
                <div class="flex items-center gap-2">
                  <Code2 class="w-5 h-5 text-indigo-400" />
                  <h4 class="font-bold text-sm">Frontend Code Snippet (Fetch API)</h4>
                </div>
                <button 
                  @click="copyText(sampleFetchCode, 'Code snippet copied!')" 
                  class="text-xs text-indigo-400 hover:text-indigo-300 flex items-center gap-1 font-medium bg-gray-800 px-3 py-1.5 rounded-lg border border-gray-700"
                >
                  <Copy class="w-3.5 h-3.5" /> Copy Code
                </button>
              </div>

              <pre class="text-xs font-mono text-green-400 m-0 overflow-x-auto whitespace-pre-wrap leading-relaxed">{{ sampleFetchCode }}</pre>
            </div>
          </div>
        </div>

        <!-- Test Console Tab -->
        <div v-if="activeTab === 'test'" class="p-6">
          <div class="max-w-3xl space-y-6">
            <div>
              <h3 class="text-lg font-semibold text-gray-900 mb-1">Interactive Request Tester</h3>
              <p class="text-sm text-gray-500">Test sending requests directly from this page to verify schema validation and response codes.</p>
            </div>

            <!-- Interactive URL Bar (Postman-style) -->
            <div class="space-y-2">
              <label class="text-sm font-semibold text-gray-700">Request URL & Method</label>
              <div class="flex items-center gap-2">
                <el-select v-model="testMethod" class="!w-32" size="large">
                  <el-option label="GET" value="GET" />
                  <el-option label="POST" value="POST" />
                  <el-option label="PUT" value="PUT" />
                  <el-option label="DELETE" value="DELETE" />
                </el-select>
                <el-input 
                  v-model="testPath" 
                  placeholder="/mock/users or /mock/users/1" 
                  size="large"
                  class="font-mono flex-1"
                >
                  <template #prepend>
                    <span class="text-gray-500 font-semibold">{{ currentOrigin }}</span>
                  </template>
                </el-input>
                <el-button 
                  type="primary" 
                  size="large" 
                  :loading="isSendingTest" 
                  @click="sendTestRequest" 
                  class="!px-6 !rounded-xl shadow-md shadow-blue-500/20"
                >
                  Send
                  <template #icon><Send class="w-4 h-4 ml-1" /></template>
                </el-button>
              </div>
              <p class="text-xs text-gray-500">
                Tip: You can test specific records by adding their ID, e.g. <code class="bg-gray-100 px-1 py-0.5 rounded text-gray-700 font-mono">{{ endpoint.path }}/{{ endpoint.records?.[0]?.id || '1' }}</code>
              </p>
            </div>

            <!-- Expected Schema Guide -->
            <div class="bg-blue-50 border border-blue-200 rounded-xl p-4">
              <div class="text-xs font-bold text-blue-900 uppercase tracking-wider mb-2 flex items-center gap-1.5">
                <CheckCircle2 class="w-4 h-4 text-blue-600" /> Expected Schema Fields (Contract)
              </div>
              <div class="flex flex-wrap gap-2">
                <span v-for="field in endpoint.schema || []" :key="field.name" class="inline-flex items-center px-2.5 py-1 rounded-md text-xs font-medium bg-white text-blue-800 border border-blue-200 shadow-2xs">
                  <span class="font-semibold">{{ field.name }}</span>
                  <span class="text-gray-400 mx-1">:</span>
                  <span class="text-gray-500">{{ field.type }}</span>
                </span>
              </div>
            </div>

            <!-- Request Body Input (for POST / PUT) -->
            <div v-if="testMethod === 'POST' || testMethod === 'PUT'" class="space-y-2">
              <div class="flex items-center justify-between">
                <label class="text-sm font-semibold text-gray-700">Request Body (JSON)</label>
                <div class="flex items-center gap-2">
                  <el-button size="small" text @click="fillSamplePayload">
                    Fill Valid Payload
                  </el-button>
                  <el-button size="small" text type="danger" @click="testPayload = '{}'">
                    Empty Payload
                  </el-button>
                </div>
              </div>
              <el-input 
                v-model="testPayload" 
                type="textarea" 
                :rows="6" 
                placeholder="{ ... }" 
                class="font-mono text-sm"
              />
            </div>

            <!-- Response Section -->
            <div v-if="testResult" class="space-y-3 pt-2">
              <div class="flex items-center justify-between">
                <h4 class="text-sm font-bold text-gray-900 uppercase tracking-wider">Response Result</h4>
                <div class="flex items-center gap-3">
                  <span :class="['px-2.5 py-0.5 rounded text-xs font-bold font-mono', testResult.status >= 200 && testResult.status < 300 ? 'bg-emerald-100 text-emerald-800' : 'bg-red-100 text-red-800']">
                    Status: {{ testResult.status }} {{ testResult.statusText }}
                  </span>
                  <el-button 
                    size="small" 
                    plain 
                    @click="copyText(JSON.stringify(testResult.data, null, 2), 'Response JSON copied!')" 
                    class="!rounded-lg"
                  >
                    <Copy class="w-3.5 h-3.5 mr-1" /> Copy JSON
                  </el-button>
                </div>
              </div>

              <!-- Error Alert when validation fails -->
              <div v-if="testResult.data?.missing_fields?.length" class="bg-red-50 border border-red-200 rounded-xl p-4">
                <div class="flex items-start gap-3">
                  <AlertCircle class="w-5 h-5 text-red-600 shrink-0 mt-0.5" />
                  <div>
                    <h5 class="text-sm font-bold text-red-900">Schema Validation Error!</h5>
                    <p class="text-xs text-red-700 mt-1">
                      The payload you sent is missing required fields defined in your schema:
                    </p>
                    <div class="flex flex-wrap gap-1.5 mt-2">
                      <span v-for="missing in testResult.data.missing_fields" :key="missing" class="px-2 py-0.5 rounded bg-red-100 text-red-800 border border-red-300 font-mono text-xs font-semibold">
                        {{ missing }}
                      </span>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Response Headers Display (if any) -->
              <div v-if="testResult.headers && Object.keys(testResult.headers).length" class="bg-gray-50 border border-gray-200 rounded-xl p-3">
                <div class="text-xs font-bold text-gray-700 uppercase tracking-wider mb-2 flex items-center justify-between">
                  <span>Response Headers:</span>
                  <span class="text-gray-400 font-normal font-mono lowercase">({{ Object.keys(testResult.headers).length }} headers)</span>
                </div>
                <div class="space-y-1 max-h-32 overflow-y-auto">
                  <div v-for="(v, k) in testResult.headers" :key="k" class="text-xs font-mono flex items-center gap-2">
                    <span class="text-blue-600 font-semibold">{{ k }}:</span>
                    <span class="text-gray-800 break-all">{{ v }}</span>
                  </div>
                </div>
              </div>

              <!-- Raw Response JSON Box -->
              <div class="bg-gray-900 rounded-xl p-4 overflow-x-auto relative group">
                <button 
                  @click="copyText(JSON.stringify(testResult.data, null, 2), 'Response JSON copied!')"
                  class="absolute top-3 right-3 opacity-0 group-hover:opacity-100 transition-opacity bg-gray-800 hover:bg-gray-700 text-gray-300 text-xs px-2.5 py-1 rounded-md border border-gray-600 flex items-center gap-1 font-mono"
                >
                  <Copy class="w-3 h-3" /> Copy
                </button>
                <pre class="text-sm font-mono m-0 whitespace-pre-wrap" :class="testResult.status >= 200 && testResult.status < 300 ? 'text-green-400' : 'text-red-400'">{{ JSON.stringify(testResult.data, null, 2) }}</pre>
              </div>
            </div>
          </div>
        </div>

        <!-- Edit Schema View -->
        <div v-if="activeTab === 'schema'" class="p-6">
          <div class="max-w-3xl space-y-6">
            <div class="bg-gray-50 border border-gray-200 rounded-xl p-4">
              <label class="text-sm font-semibold text-gray-700 block mb-2">Response Delay ({{ editForm.delay || 0 }} ms)</label>
              <el-slider v-model="editForm.delay" :min="0" :max="3000" :step="100" show-input />
              <p class="text-xs text-gray-500 mt-2">Simulate network delay for testing frontend loading spinners.</p>
            </div>

            <!-- Custom Response Headers -->
            <div class="bg-gray-50 border border-gray-200 rounded-xl p-4">
              <div class="flex items-center justify-between mb-3">
                <span class="text-sm font-semibold text-gray-700 flex items-center gap-1.5">
                  Custom Response Headers
                  <el-tooltip content="Headers sent back with the mock response" placement="top">
                    <Info class="w-4 h-4 text-gray-400" />
                  </el-tooltip>
                </span>
                <el-button size="small" type="primary" plain class="!rounded-lg text-xs" @click="addEditHeader">
                  + Add Header
                </el-button>
              </div>

              <div v-if="editHeaderList.length === 0" class="text-xs text-gray-400 italic bg-white p-3 rounded-lg border border-dashed border-gray-200">
                No custom headers configured. Click "+ Add Header" to define response headers.
              </div>

              <div v-else class="space-y-2">
                <div v-for="(item, idx) in editHeaderList" :key="idx" class="flex items-center gap-2">
                  <el-input v-model="item.key" placeholder="Header Name (e.g. X-Total-Count)" size="default" class="font-mono text-xs flex-1 bg-white" />
                  <el-input v-model="item.value" placeholder="Header Value (e.g. 100)" size="default" class="font-mono text-xs flex-1 bg-white" />
                  <el-button type="danger" plain size="small" circle class="shrink-0" @click="removeEditHeader(idx)">
                    ✕
                  </el-button>
                </div>
              </div>
            </div>

            <SchemaBuilder v-model="editForm.schema" class="!border-gray-200 !shadow-none mb-6" />
            
            <div class="flex justify-end">
              <el-button type="primary" :loading="isUpdating" @click="saveSchema" class="!px-8 !rounded-full">
                Save & Regenerate Data
                <template #icon><RefreshCw class="w-4 h-4 ml-1" :class="{ 'animate-spin': isUpdating }" /></template>
              </el-button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>

  <div v-else-if="!loading && !endpoint" class="text-center py-20">
    <h3 class="text-xl font-medium text-gray-900 mb-2">Endpoint not found</h3>
    <p class="text-gray-500 mb-6">The mock API you are looking for might have been deleted.</p>
    <el-button @click="navigateTo('/')">Go to Dashboard</el-button>
  </div>

  <!-- Edit Record Dialog -->
  <el-dialog v-model="editRecordDialogVisible" title="Edit Record" width="500px" destroy-on-close>
    <el-form label-position="top" @submit.prevent="submitEditRecord">
      <el-form-item v-for="field in endpoint?.schema || []" :key="field.name" :label="field.name">
        <el-input v-model="editRecordForm[field.name]" :placeholder="field.type" />
      </el-form-item>
    </el-form>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="editRecordDialogVisible = false">Cancel</el-button>
        <el-button type="primary" :loading="isSavingRecord" @click="submitEditRecord">Save Changes</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useEndpointStore } from '~/stores/endpointStore'
import SchemaBuilder from '~/components/endpoint/SchemaBuilder.vue'
import { ArrowLeft, Copy, ExternalLink, Terminal, Database, Settings, RefreshCw, Plus, Pencil, Trash2, Play, Send, CheckCircle2, AlertCircle, RotateCcw, BookOpen, Code2, Download, Loader2 } from 'lucide-vue-next'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useClipboard } from '@vueuse/core'

const route = useRoute()
const store = useEndpointStore()
const id = route.params.id as string

const endpoint = computed(() => store.currentEndpoint)
const loading = computed(() => store.loading)

const activeTab = ref('data')
const viewMode = ref<'table' | 'json'>('table')
const isUpdating = ref(false)

const editForm = ref({
  path: '',
  method: 'GET',
  delay: 0,
  count: 5,
  schema: [] as { name: string; type: string }[]
})

onMounted(async () => {
  if (id) {
    await store.fetchEndpoint(id)
  }
})

const editHeaderList = ref<{ key: string; value: string }[]>([])

const addEditHeader = () => {
  editHeaderList.value.push({ key: '', value: '' })
}

const removeEditHeader = (index: number) => {
  editHeaderList.value.splice(index, 1)
}

watch(endpoint, (newVal) => {
  if (newVal) {
    editForm.value = {
      path: newVal.path,
      method: newVal.method,
      delay: newVal.delay || 0,
      count: newVal.count,
      schema: Array.isArray(newVal.schema) ? [...newVal.schema] : []
    }

    if (newVal.headers && typeof newVal.headers === 'object') {
      editHeaderList.value = Object.entries(newVal.headers).map(([key, value]) => ({
        key,
        value: String(value)
      }))
    } else {
      editHeaderList.value = []
    }
  }
}, { immediate: true, deep: true })

const getMethodColor = (method: string) => {
  const methods: Record<string, { text: string, bg: string, badgeBg: string }> = {
    'GET': { text: 'text-blue-700', bg: 'from-blue-400 to-blue-600', badgeBg: 'bg-blue-100' },
    'POST': { text: 'text-emerald-700', bg: 'from-emerald-400 to-emerald-600', badgeBg: 'bg-emerald-100' },
    'PUT': { text: 'text-amber-700', bg: 'from-amber-400 to-amber-600', badgeBg: 'bg-amber-100' },
    'DELETE': { text: 'text-red-700', bg: 'from-red-400 to-red-600', badgeBg: 'bg-red-100' },
  }
  return methods[method?.toUpperCase()] || { text: 'text-gray-700', bg: 'from-gray-400 to-gray-600', badgeBg: 'bg-gray-100' }
}

const getFullUrl = () => {
  if (!endpoint.value) return ''
  const path = endpoint.value.path
  return `${window.location.origin}/mock${path.startsWith('/') ? path : '/' + path}`
}

const { copy, copied } = useClipboard()

const copyUrl = async () => {
  try {
    await copy(getFullUrl())
    if (copied.value) {
      ElMessage.success('URL copied to clipboard')
    }
  } catch (e) {
    ElMessage.error('Failed to copy URL')
  }
}

const openTestUrl = () => {
  window.open(getFullUrl(), '_blank')
}

const exportSingleEndpoint = () => {
  if (!endpoint.value) return
  const data = {
    path: endpoint.value.path,
    method: endpoint.value.method,
    delay: endpoint.value.delay || 0,
    headers: endpoint.value.headers || {},
    count: endpoint.value.count || 5,
    schema: endpoint.value.schema || []
  }

  const jsonStr = JSON.stringify(data, null, 2)
  const blob = new Blob([jsonStr], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  
  const link = document.createElement('a')
  link.href = url
  const cleanName = endpoint.value.path.replace(/\
  link.download = `mock-endpoint-${cleanName || 'api'}.json`
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)

  ElMessage.success('Exported endpoint JSON successfully!')
}

const saveSchema = async () => {
  isUpdating.value = true
  try {
    const customHeaders: Record<string, string> = {}
    editHeaderList.value.forEach(h => {
      if (h.key.trim()) {
        customHeaders[h.key.trim()] = h.value
      }
    })

    await store.updateEndpoint(id, {
      ...editForm.value,
      headers: Object.keys(customHeaders).length > 0 ? customHeaders : undefined
    })
    activeTab.value = 'data'
  } catch (error) {
  } finally {
    isUpdating.value = false
  }
}

const formatDate = (dateStr?: string) => {
  if (!dateStr) return 'Unknown'
  const date = new Date(dateStr)
  return new Intl.DateTimeFormat('en-US', { 
    month: 'short', 
    day: 'numeric', 
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  }).format(date)
}

const handleAddRecord = async () => {
  await store.addRecord(id)
}

const handleResetRecords = () => {
  ElMessageBox.confirm(
    'This will delete all current records and regenerate 5 new mock records based on the schema. Are you sure?',
    'Reset Records',
    {
      confirmButtonText: 'Reset',
      cancelButtonText: 'Cancel',
      type: 'warning',
    }
  ).then(async () => {
    await store.resetRecords(id)
  }).catch(() => {})
}

const editRecordDialogVisible = ref(false)
const currentEditRecordId = ref<string | number | null>(null)
const editRecordForm = ref<Record<string, any>>({})
const isSavingRecord = ref(false)

const openEditRecordDialog = (row: any) => {
  currentEditRecordId.value = row.id
  editRecordForm.value = JSON.parse(JSON.stringify(row.data || {}))
  editRecordDialogVisible.value = true
}

const submitEditRecord = async () => {
  if (!currentEditRecordId.value) return
  isSavingRecord.value = true
  try {
    await store.updateRecord(currentEditRecordId.value, { data: editRecordForm.value })
    editRecordDialogVisible.value = false
  } catch (error) {
  } finally {
    isSavingRecord.value = false
  }
}

const handleDeleteRecord = (row: any) => {
  ElMessageBox.confirm(
    'Are you sure you want to delete this record?',
    'Warning',
    {
      confirmButtonText: 'Delete',
      cancelButtonText: 'Cancel',
      type: 'warning',
    }
  ).then(async () => {
    await store.deleteRecord(row.id)
  }).catch(() => {
  })
}

const testMethod = ref('GET')
const testPath = ref('')
const currentOrigin = ref('')
const testPayload = ref('{}')
const isSendingTest = ref(false)
const testResult = ref<{ status: number; statusText: string; data: any; headers?: Record<string, string> } | null>(null)

onMounted(() => {
  if (typeof window !== 'undefined') {
    currentOrigin.value = window.location.origin
  }
})

const fillSamplePayload = () => {
  if (!endpoint.value?.schema) return
  const sample: Record<string, any> = {}
  endpoint.value.schema.forEach((f: any) => {
    if (f.type.includes('email')) sample[f.name] = 'test@example.com'
    else if (f.type.includes('int') || f.type.includes('price')) sample[f.name] = 100
    else if (f.type.includes('bool')) sample[f.name] = true
    else if (f.type.includes('uuid')) sample[f.name] = 'sample-uuid-123'
    else if (f.type.includes('avatar') || f.type.includes('image')) sample[f.name] = 'https://i.pravatar.cc/300'
    else sample[f.name] = `Sample ${f.name}`
  })
  testPayload.value = JSON.stringify(sample, null, 2)
}

watch(endpoint, (ep) => {
  if (ep) {
    if (!testPath.value) {
      testPath.value = `/mock${ep.path.startsWith('/') ? ep.path : '/' + ep.path}`
    }
    if (!testMethod.value || testMethod.value === 'GET') {
      testMethod.value = ep.method || 'GET'
    }
    if (ep.schema?.length && testPayload.value === '{}') {
      fillSamplePayload()
    }
  }
}, { immediate: true })

const sendTestRequest = async () => {
  if (!testPath.value) return
  isSendingTest.value = true
  testResult.value = null

  let targetUrl = testPath.value.trim()
  if (!targetUrl.startsWith('http')) {
    if (!targetUrl.startsWith('/')) {
      targetUrl = '/' + targetUrl
    }
    if (!targetUrl.startsWith('/mock')) {
      targetUrl = '/mock' + targetUrl
    }
  }

  const method = testMethod.value || 'GET'

  try {
    let bodyData: any = undefined
    if (method === 'POST' || method === 'PUT') {
      try {
        bodyData = JSON.parse(testPayload.value)
      } catch (e) {
        ElMessage.error('Invalid JSON in Request Body')
        isSendingTest.value = false
        return
      }
    }

    const res = await fetch(targetUrl, {
      method,
      headers: {
        'Content-Type': 'application/json',
      },
      body: bodyData ? JSON.stringify(bodyData) : undefined,
    })

    const data = await res.json().catch(() => null)
    const respHeaders: Record<string, string> = {}
    res.headers.forEach((val, key) => {
      respHeaders[key] = val
    })

    testResult.value = {
      status: res.status,
      statusText: res.statusText || (res.ok ? 'OK' : 'Error'),
      data,
      headers: respHeaders,
    }

    if (res.ok) {
      ElMessage.success(`${method} request completed successfully!`)
      await store.fetchEndpoint(id)
    } else {
      ElMessage.warning(`Server responded with ${res.status}`)
    }
  } catch (err: any) {
    testResult.value = {
      status: 0,
      statusText: 'Network / Connection Error',
      data: { error: err.message },
    }
    ElMessage.error('Failed to send request')
  } finally {
    isSendingTest.value = false
  }
}

const copyText = async (text: string, successMsg = 'Copied to clipboard!') => {
  try {
    await copy(text)
    ElMessage.success(successMsg)
  } catch {
    ElMessage.error('Failed to copy')
  }
}

const useInTestConsole = (method: string, samplePath: string) => {
  testMethod.value = method
  testPath.value = samplePath
  activeTab.value = 'test'
}

const apiEndpointsDocs = computed(() => {
  if (!endpoint.value) return []
  const origin = currentOrigin.value || ''
  const basePath = endpoint.value.path.startsWith('/') ? endpoint.value.path : '/' + endpoint.value.path
  const firstId = endpoint.value.records?.[0]?.id || '1'

  const sampleObj: Record<string, any> = {}
  if (endpoint.value.schema) {
    endpoint.value.schema.forEach((f: any) => {
      if (f.type.includes('email')) sampleObj[f.name] = 'john@example.com'
      else if (f.type.includes('int') || f.type.includes('price')) sampleObj[f.name] = 100
      else if (f.type.includes('bool')) sampleObj[f.name] = true
      else if (f.type.includes('uuid')) sampleObj[f.name] = 'user-uuid-123'
      else if (f.type.includes('avatar') || f.type.includes('image')) sampleObj[f.name] = 'https://i.pravatar.cc/300'
      else sampleObj[f.name] = `Sample ${f.name}`
    })
  }
  const sampleBodyStr = JSON.stringify(sampleObj, null, 2)

  return [
    {
      method: 'GET',
      fullUrl: `${origin}/mock${basePath}`,
      samplePath: `/mock${basePath}`,
      description: 'Fetch all records in this mock resource. Supports query filtering and pagination (?page=1&limit=10, ?search=...)',
      body: ''
    },
    {
      method: 'GET',
      fullUrl: `${origin}/mock${basePath}/:id`,
      samplePath: `/mock${basePath}/${firstId}`,
      description: `Fetch a single record by its specific ID (e.g. /mock${basePath}/${firstId})`,
      body: ''
    },
    {
      method: 'POST',
      fullUrl: `${origin}/mock${basePath}`,
      samplePath: `/mock${basePath}`,
      description: 'Create a new record. Must provide all required fields matching the schema.',
      body: sampleBodyStr
    },
    {
      method: 'PUT',
      fullUrl: `${origin}/mock${basePath}/:id`,
      samplePath: `/mock${basePath}/${firstId}`,
      description: `Update an existing record by providing its ID either in the URL or in the JSON body payload.`,
      body: sampleBodyStr
    },
    {
      method: 'DELETE',
      fullUrl: `${origin}/mock${basePath}/:id`,
      samplePath: `/mock${basePath}/${firstId}`,
      description: `Permanently delete a record by specifying its ID in the URL.`,
      body: ''
    }
  ]
})

const sampleFetchCode = computed(() => {
  if (!endpoint.value) return ''
  const origin = currentOrigin.value || 'http://localhost:3000'
  const basePath = endpoint.value.path.startsWith('/') ? endpoint.value.path : '/' + endpoint.value.path
  return `
fetch('${origin}/mock${basePath}')
  .then(res => res.json())
  .then(data => console.log(data));

fetch('${origin}/mock${basePath}', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ /* your schema fields */ })
});

fetch('${origin}/mock${basePath}/1', {
  method: 'PUT',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ /* updated fields */ })
});

fetch('${origin}/mock${basePath}/1', {
  method: 'DELETE'
});`
})
</script>
