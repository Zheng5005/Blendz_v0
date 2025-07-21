import { axiosInstance } from "./axios";

interface SignUpData {
  fullName: string;
  email: string;
  password: string;
}

interface LoginData {
  email: string;
  password: string;
}

interface UserData {
  fullName: string;
  bio: string;
  nativeLanguage: string;
  learningLanguage: string;
  location: string;
  profilePic: string;
}

export const signup = async (signupData: SignUpData) => {
  try {
    const res = await axiosInstance.post("/auth/signup", signupData)
    return res.data;
  } catch (error) {
    console.log("Error in signup api:", error);
    return null;
  }
}

export const login = async (loginData: LoginData) => {
  try {
    const response = await axiosInstance.post("/auth/login", loginData);
    return response.data;
  } catch (error) {
    console.log("Error in login api:", error);
    return null;
  }
};

export const logout = async () => {
  try {
    const response = await axiosInstance.post("/auth/logout");
    return response.data;
  } catch (error) {
    console.log("Error in logout api:", error);
    return null;
  }
};

export const getAuthUser = async () => {
  try {
    const res = await axiosInstance.get("/users/me");
    return res.data;
  } catch (error) {
    console.log("Error in getAuthUser:", error);
    return null;
  }
};

export const completeOnboarding = async (userData: UserData) => {
  try {
    const response = await axiosInstance.post("/auth/onboarding", userData);
    return response.data;   
  } catch (error) {
    console.log("Error in completeOnboarding api:", error);
    return null;
  }
};
