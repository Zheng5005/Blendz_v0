import { create } from "zustand";

interface ThemeState {
  theme: string;
  setTheme: (theme: string) => void
}

export const useThemeStore = create<ThemeState>()((set) => ({
  theme: localStorage.getItem("blendz-theme") || "synthwave",
  setTheme: (theme: string) => {
    localStorage.setItem("blendz-theme", theme);
    set({ theme });
  },
}));
