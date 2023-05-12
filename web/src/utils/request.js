import axios from "axios";
const instance = axios.create({
    timeout: 5000
})

instance.interceptors.request.use(
    vueConfig => {
        return vueConfig
    },
    error => {
        return Promise.reject(error)
    }
)

instance.interceptors.response.use(
    response => {
        return response.status === 200 ? Promise.resolve(response) : Promise.reject(response)
    },
    error => {
        const { response } = error;
        return response.status
    }
)

export default instance
