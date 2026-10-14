// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2024-04-03',
  srcDir: '.',
  devtools: { enabled: false },

  modules: [
    '@nuxtjs/tailwindcss',
    '@pinia/nuxt',
    '@element-plus/nuxt',
  ],

  // Proxy API requests to Go backend
  routeRules: {
    '/api/**': { proxy: 'http://localhost:8080/api/**' },
    '/mock/**': { proxy: 'http://localhost:8080/mock/**' }
  },

  css: [
    'element-plus/dist/index.css',
    '~/assets/css/main.css'
  ],

  elementPlus: {
    icon: 'ElIcon',
    importStyle: 'css',
  }
})
