import type { AxiosInstance } from "axios";
import axios from "axios";

export const axiosInstance: AxiosInstance = axios.create({
  withCredentials: true, // send cookies with the request
  baseURL: "http://localhost:8080/api",
})
