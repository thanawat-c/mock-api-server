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

  // Proxy API requests to Go backend (defaults to localhost:8080 or BACKEND_URL env)
  routeRules: {
    '/api/**': { proxy: `${process.env.BACKEND_URL || 'http://localhost:8080'}/api/**` },
    '/mock/**': { proxy: `${process.env.BACKEND_URL || 'http://localhost:8080'}/mock/**` }
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
