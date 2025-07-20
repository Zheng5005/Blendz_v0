import type { AxiosInstance } from "axios";
import axios from "axios";

export const axiosInstance: AxiosInstance = axios.create({
  baseURL: "http://localhost:8080/api",
  withCredentials: true // send cookies with the request
})
