// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
    compatibilityDate: '2024-04-03',
    devtools: {enabled: false},
    modules: ["@nuxt/ui", '@nuxt/icon', '@nuxtjs/color-mode', '@vueuse/nuxt', 'dayjs-nuxt'],
    ssr: false,
    nitro: {
        preset: 'static'
    },
    dayjs: {
        locales: ['zh'],
        defaultLocale: 'zh'
    },
    icon: {
        clientBundle: {
            scan: {
                globInclude: ['**/*.{vue,jsx,tsx}', 'node_modules/@nuxt/ui/**/*.js'],
                globExclude: ['.*', 'coverage', 'test', 'tests', 'dist', 'build'],
            },
        },
    },
    tailwindcss: {
        safelist: [
            'grid-cols-1',
            'grid-cols-3',
        ]
    },
    vue: {
        compilerOptions: {
            isCustomElement: (tag:string) => ['meting-js'].includes(tag),
        },
    },
    app: {
        baseURL: '/', // 如果你部署在子路径下，需要改成 /子路径名/
        head: {
            meta: [
                { name: "viewport", content: "width=device-width, initial-scale=1, user-scalable=no" },
                { charset: "utf-8" },
            ],
            link: [
                {href: `/css/APlayer.min.css`, rel: 'stylesheet'},
            ],
            script: [
                {src: `/js/APlayer.min.js`, type: 'text/javascript', async: true, defer: true},
                {src: `/js/Meting.min.js`, type: 'text/javascript', async: true, defer: true},
                {src: `/js/main.js`, type: 'text/javascript', async: true, defer: true},
            ]
        }
    },
    runtimeConfig: {
        public: {
            apiBase: process.env.NUXT_PUBLIC_API_BASE || 'https://moments.alamentation.xyz'
            // apiBase: process.env.NUXT_PUBLIC_API_BASE || 'http://192.168.0.102:8888'
        }
    },
    vite: {
        server: {
            proxy: {
                "/api": {
                    target: "http://192.168.0.102:8888",
                    // changeOrigin: true,
                },
                "/upload": {
                    target: "http://192.168.0.102:8888",
                    // changeOrigin: true,
                },
                "/rss": {
                    target: "http://192.168.0.102:8888",
                    // changeOrigin: true,
                },
                "/swagger": {
                    target: "http://192.168.0.102:8888",
                    // changeOrigin: true,
                },
            },
        },
        build: {
            rollupOptions: {
                output: {
                    hashCharacters: 'base36'
                }
            }
        }
    }
})