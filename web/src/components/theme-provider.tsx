import { createContext, useCallback, useContext, useEffect, useState } from "react"

type Mode = "dark" | "light" | "system"

export const themes = [
  { name: "default", text: "Default" },
  { name: "red", text: "Red" },
  { name: "rose", text: "Rose" },
  { name: "orange", text: "Orange" },
  { name: "green", text: "Green" },
  { name: "blue", text: "Blue" },
  { name: "yellow", text: "Yellow" },
  { name: "violet", text: "Violet" },
  { name: "stone", text: "Stone" },
  { name: "zinc", text: "Zinc" },
  { name: "neutral", text: "Neutral" },
  { name: "gray", text: "Gray" },
  { name: "slate", text: "Slate" },
] as const

type Theme = typeof themes[number]["name"]

type ThemeProviderProps = {
  children: React.ReactNode
  defaultMode?: Mode
  defaultTheme?: Theme
  modeStorageKey?: string
  themeStorageKey?: string
}

type ThemeProviderState = {
  mode: Mode
  setMode: (mode: Mode) => void
  theme: Theme
  setTheme: (theme: Theme) => void
}

const initialState: ThemeProviderState = {
  mode: "system",
  setMode: () => null,
  theme: "neutral",
  setTheme: () => null,
}

const ThemeProviderContext = createContext<ThemeProviderState>(initialState)

export function ThemeProvider({
  children,
  defaultMode = "system",
  defaultTheme = "neutral",
  modeStorageKey = "vite-ui-mode",
  themeStorageKey = "vite-ui-theme",
  ...props
}: ThemeProviderProps) {
  const [mode, setMode] = useState<Mode>(
    () => (localStorage.getItem(modeStorageKey) as Mode) || defaultMode
  )
  const [theme, setTheme] = useState<Theme>(
    () => (localStorage.getItem(themeStorageKey) as Theme) || defaultTheme
  )

  useEffect(() => {
    const root = window.document.documentElement

    root.classList.remove("light", "dark")

    if (mode === "system") {
      const systemTheme = window.matchMedia("(prefers-color-scheme: dark)")
        .matches
        ? "dark"
        : "light"

      root.classList.add(systemTheme)
      return
    }

    root.classList.add(mode)
  }, [mode])

  useEffect(() => {
    const root = window.document.documentElement

    themes.forEach(({ name }) => {
      root.classList.remove(name)
    })

    root.classList.add(theme)
  }, [theme])

  const setModeImpl = useCallback((mode: Mode) => {
    localStorage.setItem(modeStorageKey, mode)
    setMode(mode)
  }, [])

  const setThemeImpl = useCallback((theme: Theme) => {
    localStorage.setItem(themeStorageKey, theme)
    setTheme(theme)
  }, [])

  return (
    <ThemeProviderContext.Provider {...props} value={{
      theme,
      mode,
      setTheme: setThemeImpl,
      setMode: setModeImpl,
    }}>
      {children}
    </ThemeProviderContext.Provider>
  )
}

export const useTheme = () => {
  const context = useContext(ThemeProviderContext)

  if (context === undefined)
    throw new Error("useTheme must be used within a ThemeProvider")

  return context
}
