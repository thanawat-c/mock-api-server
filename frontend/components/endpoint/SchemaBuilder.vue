<template>
  <div class="bg-white rounded-2xl border border-gray-200 p-6 shadow-sm">
    <div class="flex items-center justify-between mb-4">
      <h3 class="text-lg font-semibold text-gray-900 flex items-center gap-2">
        <Braces class="w-5 h-5 text-blue-500" />
        Data Schema
      </h3>
      <el-button type="primary" plain size="small" class="!rounded-full" @click="addField">
        <template #icon><Plus /></template>
        Add Field
      </el-button>
    </div>

    <!-- Field Rows -->
    <draggable 
      v-if="fields.length > 0" 
      v-model="fields"
      item-key="_id"
      handle=".drag-handle"
      class="space-y-3 mb-6"
      :animation="200"
    >
      <template #item="{ element: field, index }">
        <div class="flex items-center gap-3 animate-fade-in group">
          <div class="drag-handle cursor-move text-gray-300 hover:text-gray-500 transition-colors p-1 -ml-2 rounded hover:bg-gray-100">
            <GripVertical class="w-5 h-5" />
          </div>
          <div class="flex-1">
            <el-input v-model="field.name" placeholder="Field name (e.g., fullName)" />
          </div>
          <div class="flex-1">
            <el-select v-model="field.type" placeholder="Select faker type" class="w-full">
              <el-option-group v-for="group in fakerOptions" :key="group.label" :label="group.label">
                <el-option
                  v-for="item in group.options"
                  :key="item.value"
                  :label="item.label"
                  :value="item.value"
                >
                  <span class="float-left">{{ item.label }}</span>
                  <span class="float-right text-gray-400 text-xs font-mono">{{ item.value }}</span>
                </el-option>
              </el-option-group>
            </el-select>
          </div>
          <el-button type="danger" text circle class="!p-2 hover:bg-red-50 opacity-0 group-hover:opacity-100 transition-opacity" @click="removeField(index)">
            <Trash2 class="w-4 h-4 text-red-500" />
          </el-button>
        </div>
      </template>
    </draggable>

    <div v-else class="text-center py-8 bg-gray-50 rounded-xl border border-dashed border-gray-300 mb-6">
      <p class="text-gray-500 text-sm">No fields defined yet. Add a field to build your JSON schema.</p>
    </div>

    <!-- Live Preview -->
    <div>
      <h4 class="text-sm font-medium text-gray-700 mb-2 flex items-center gap-1.5">
        <Code2 class="w-4 h-4 text-gray-400" />
        Live Schema Preview
      </h4>
      <div class="bg-gray-900 rounded-xl p-4 overflow-x-auto">
        <pre class="text-sm text-green-400 font-mono m-0 whitespace-pre-wrap">{{ schemaPreview }}</pre>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { Plus, Trash2, Braces, Code2, GripVertical } from 'lucide-vue-next'
import draggable from 'vuedraggable'
import type { SchemaField } from '~/types/endpoint'

const props = defineProps<{
  modelValue: SchemaField[]
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: SchemaField[]): void
}>()

type LocalField = { _id: string, name: string, type: string }
const fields = ref<LocalField[]>([])

watch(() => props.modelValue, (newVal) => {
  if (Array.isArray(newVal) && fields.value.length === 0) {
    fields.value = newVal.map((field) => ({ 
      _id: Math.random().toString(36).substring(7),
      name: field.name, 
      type: field.type 
    }))
  }
}, { immediate: true })

watch(fields, (newFields) => {
  const schema: SchemaField[] = newFields
    .filter(f => f.name && f.type)
    .map(f => ({ name: f.name, type: f.type }))
  
  emit('update:modelValue', schema)
}, { deep: true })

const addField = () => {
  fields.value.push({ _id: Math.random().toString(36).substring(7), name: '', type: 'faker.name' })
}

const removeField = (index: number) => {
  fields.value.splice(index, 1)
}

const schemaPreview = computed(() => {
  const schema: SchemaField[] = fields.value.map(f => ({
    name: f.name || 'undefined',
    type: f.type || 'undefined'
  }))
  return JSON.stringify(schema, null, 2)
})

const fakerOptions = [
  {
    label: 'Person',
    options: [
      { label: 'Full Name', value: 'faker.name' },
      { label: 'First Name', value: 'faker.firstName' },
      { label: 'Last Name', value: 'faker.lastName' },
      { label: 'Email', value: 'faker.email' },
      { label: 'Phone', value: 'faker.phone' },
    ]
  },
  {
    label: 'Media',
    options: [
      { label: 'Image URL', value: 'faker.image' },
      { label: 'Avatar URL', value: 'faker.avatar' },
    ]
  },
  {
    label: 'Data Types',
    options: [
      { label: 'Random ID', value: 'faker.uuid' },
      { label: 'Integer', value: 'faker.int' },
      { label: 'Boolean', value: 'faker.boolean' },
      { label: 'Date', value: 'faker.date' },
      { label: 'Word', value: 'faker.word' },
      { label: 'Sentence', value: 'faker.sentence' },
      { label: 'Paragraph', value: 'faker.paragraph' },
    ]
  },
  {
    label: 'Location',
    options: [
      { label: 'City', value: 'faker.city' },
      { label: 'Country', value: 'faker.country' },
      { label: 'Address', value: 'faker.address' },
      { label: 'ZipCode', value: 'faker.zip' },
    ]
  },
  {
    label: 'Commerce',
    options: [
      { label: 'Product Name', value: 'faker.product' },
      { label: 'Price', value: 'faker.price' },
      { label: 'Company', value: 'faker.company' },
    ]
  }
]
</script>

<style scoped>
.animate-fade-in {
  animation: fadeIn 0.3s ease-in-out;
}
@keyframes fadeIn {
  from { opacity: 0; transform: translateY(-5px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
