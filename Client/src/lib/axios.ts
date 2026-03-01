import type { AxiosInstance } from "axios";
import axios from "axios";

const API_URL = import.meta.env.VITE_API_URL;

export const axiosInstance: AxiosInstance = axios.create({
  withCredentials: true, // send cookies with the request
  baseURL: API_URL,
})
