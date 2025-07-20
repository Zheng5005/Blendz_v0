import { axiosInstance } from "./axios";

interface SignUpData {
  fullName: string;
  email: string;
  password: string;
}

export const signup = async (signupData: SignUpData) => {
  const res = await axiosInstance.post("/auth/signup", signupData)
  return res.data;
}

export const getAuthUser = async () => {
  try {
    const res = await axiosInstance.get("/auth/me");
    return res.data;
  } catch (error) {
    console.log("Error in getAuthUser:", error);
    return null;
  }
};
