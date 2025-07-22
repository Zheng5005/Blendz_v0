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

export async function getUserFriends() {
  try {
    const response = await axiosInstance.get("/users/friends");
    return response.data;
  } catch (error) {
    console.log("Error in getUserFriends api:", error);
    return null;
  }
}

export async function getRecommendedUsers() {
  try {
    const response = await axiosInstance.get("/users");
    return response.data;
  } catch (error) {
    console.log("Error in getRecommendedUsers api:", error);
    return null;
  }
}

export async function getOutgoingFriendReqs() {
  try {
    const response = await axiosInstance.get("/users/outgoing-friend-requests");
    return response.data;
  } catch (error) {
    console.log("Error in getOutgoingFriendReqs api:", error);
    return null;
  }
}

export async function sendFriendRequest(userId: string) {
  try {
    const response = await axiosInstance.post(`/users/friend-request/${userId}`);
    return response.data;
  } catch (error) {
    console.log("Error in sendFriendRequest api:", error);
    return null;
  }
}

export async function getFriendRequests() {
  try {
    const response = await axiosInstance.get("/users/friend-requests");
    return response.data;
  } catch (error) {
    console.log("Error in getFriendRequests api:", error);
    return null;
  }
}

export async function acceptFriendRequest(requestId: number) {
  try {
    const response = await axiosInstance.put(`/users/friend-request/${requestId}/accept`);
    return response.data;
  } catch (error) {
    console.log("Error in acceptFriendRequest api:", error);
    return null;
  }
}

export async function getStreamToken() {
  try {
    const response = await axiosInstance.get("/chat/token")
    return response.data
  } catch (error) {
    console.log("Error in getStreamToken api:", error);
    return null;
  }
}
