import { Navigate, Route, Routes } from "react-router"
import HomePage from "./pages/HomePage"
import SignUpPage from "./pages/SignUpPage"
import LoginPage from "./pages/LoginPage"
import NotificationsPage from "./pages/NotificationsPage"
import CallPage from "./pages/CallPage"
import ChatPage from "./pages/ChatPage"
import OnboardingPage from "./pages/OnboardingPage"
import { Toaster } from "react-hot-toast"
import useAuthUser from "./hooks/useAuthUser"
import PageLoader from "./components/PageLoader"

function App() {
  const { isLoading, authUser } = useAuthUser();

  const isAuthenticated = Boolean(authUser)
  const isOnboarded = authUser?.isOnboarded

  if (isLoading) return <PageLoader />

  return(
    <div className="h-screen" data-theme="synthwave">
      <Routes>
        <Route 
          path="/"
          element={isAuthenticated && isOnboarded ? (
            <HomePage />
          ) : (
            <Navigate to={!isAuthenticated ? "/login" : "/onboarding"} />
          )}
        />
        <Route 
          path="/signup" 
          element={!isAuthenticated ? <SignUpPage /> : <Navigate to={isOnboarded ? "/" : "/onboarding"} />}
        />
        <Route 
          path="/login" 
          element={!isAuthenticated ? <LoginPage /> : <Navigate to={isOnboarded ? "/" : "/onboarding"} />}
        />
        <Route 
          path="/notifications" 
          element={authUser ? <NotificationsPage /> : <Navigate to="/login" />} />
        <Route 
          path="/call" 
          element={authUser ? <CallPage /> : <Navigate to="/login" />} />
        <Route 
          path="/chat" 
          element={authUser ? <ChatPage /> : <Navigate to="/login" />} />
        <Route 
          path="/onboarding" 
          element={isAuthenticated ? (
            !isOnboarded ? (
              <OnboardingPage />
            ) : (
              <Navigate to="/" />
            )
          ) : (
            <Navigate to="/login" />
          )}
        />
      </Routes>

      <Toaster />
    </div>
  )
}

export default App
