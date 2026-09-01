/**
 * axios统一封装
 *
 * @author xiaoRui
 *
 */
import { Message } from 'element-plus'
import axios from 'axios'
import router from "@/router/router"
import storage from "./storage"

// 创建axios对象，添加全局配置
const service = axios.create({
    baseURL: (process.env.VUE_APP_API_BASE_URL || '').replace(/\/$/, ''), // 通过环境变量指定后端地址
    timeout: 15000, // 增加到15秒
    headers: {
        'Content-Type': 'application/json'
    }
})

// 是否正在跳转登录页的标志，避免重复跳转和多次提示
let isRedirectingToLogin = false

// 处理token过期的统一函数
function handleTokenExpired(msg = 'Token已过期，正在跳转登录页...') {
    if (isRedirectingToLogin) {
        return
    }

    isRedirectingToLogin = true

    // 1. 清除本地存储
    storage.clearAll()

    // 2. 显示提示
    Message({ message: msg, type: 'warning', duration: 2000 })

    // 3. 跳转到登录页
    router.push('/login')

    // 重置标志
    setTimeout(() => {
        isRedirectingToLogin = false
    }, 1000)
}

// 请求拦截
service.interceptors.request.use((req) => {
    // 如果正在跳转登录页，取消所有新的请求
    if (isRedirectingToLogin) {
        return Promise.reject(new Error('正在跳转登录页'))
    }

    console.log('API Base URL:', process.env.VUE_APP_API_BASE_URL)

    const headers = req.headers
    const token = storage.getItem("token") || {}
    if(!headers.Authorization) {
        // 如果token是对象，获取其access_token属性；如果是字符串直接使用
        const tokenValue = typeof token === 'object' ? token.access_token || token.token : token
        headers.Authorization = 'Bearer ' + tokenValue
    }

    // 确保URL以/api/v1开头
    if (req.url && !req.url.startsWith('/api/v1/')) {
        req.url = '/api/v1' + (req.url.startsWith('/') ? req.url : '/' + req.url)
    }

    return req
})

// 响应拦截
service.interceptors.response.use((res) => {
    // 检查响应数据类型，只对对象类型进行解构
    if (res.data && typeof res.data === 'object') {
        const {code, message} = res.data

        // 401或406: token过期或无效
        if (code === 401 || code === 406) {
            handleTokenExpired(message || 'Token已过期，正在跳转登录页...')
            // 返回一个永远pending的Promise，阻止业务代码继续执行
            return new Promise(() => {})
        }
    }
    return res
}, (error) => {
    // 检查error对象是否存在
    if (!error) {
        console.error('响应拦截器收到undefined错误')
        return Promise.reject(new Error('未知错误'))
    }

    // HTTP错误处理
    if (error.response) {
        const status = error.response.status
        const data = error.response.data

        // HTTP 401也是认证失败
        if (status === 401) {
            handleTokenExpired(data?.message || '登录已过期，正在跳转登录页...')
            return new Promise(() => {})
        }

        // 其他HTTP错误
        if (!isRedirectingToLogin) {
            const errorMsg = data?.message || `请求失败(${status})`
            Message.error(errorMsg)
        }
    } else if (error.request && !isRedirectingToLogin) {
        Message.error('网络连接失败，请检查网络设置')
    } else if (!isRedirectingToLogin && error.message) {
        Message.error(error.message)
    }

    return Promise.reject(error)
})

// 请求核心函数
function request(options) {
    options.method = options.method || 'get'

    if (options.method.toLowerCase() === 'get') {
        options.params = options.data || options.params
    }

    if (options.headers) {
        delete options.headers['x-upload-lock']
        delete options.headers['X-Upload-Lock']
    }

    return service(options)
}

export default request