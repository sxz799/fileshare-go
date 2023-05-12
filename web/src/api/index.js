import axios from "../utils/request";

import path from "./path"


const api = {
    upload() {
        return path.baseUrl + path.upload
    },
    download(shareCode) {
        return path.baseUrl + path.download + "/" + shareCode
    },
    config() {
        return axios.get(path.baseUrl + path.config)
    },
    exist(shareCode) {
        return axios.get(path.baseUrl + path.exist + "/" + shareCode)
    },

    del(shareCode) {
        return axios.get(path.baseUrl + path.del + "/" + shareCode)
    },
    list(params) {
        return axios.get(path.baseUrl + path.list, params)
    }
}

export default api


