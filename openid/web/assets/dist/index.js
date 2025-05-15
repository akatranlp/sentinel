const modeButton = document.getElementById("change-mode")

const modeIcons = {
  "dark": `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-moon-icon lucide-moon"><path d="M12 3a6 6 0 0 0 9 9 9 9 0 1 1-9-9Z"/></svg>`,
  "system": `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-sun-moon-icon lucide-sun-moon"><path d="M12 8a2.83 2.83 0 0 0 4 4 4 4 0 1 1-4-4"/><path d="M12 2v2"/><path d="M12 20v2"/><path d="m4.9 4.9 1.4 1.4"/><path d="m17.7 17.7 1.4 1.4"/><path d="M2 12h2"/><path d="M20 12h2"/><path d="m6.3 17.7-1.4 1.4"/><path d="m19.1 4.9-1.4 1.4"/></svg>`,
  "light": `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-sun-icon lucide-sun"><circle cx="12" cy="12" r="4"/><path d="M12 2v2"/><path d="M12 20v2"/><path d="m4.93 4.93 1.41 1.41"/><path d="m17.66 17.66 1.41 1.41"/><path d="M2 12h2"/><path d="M20 12h2"/><path d="m6.34 17.66-1.41 1.41"/><path d="m19.07 4.93-1.41 1.41"/></svg>`,
}

let currentMode = localStorage.getItem("mode") || "system"
let currentTheme = localStorage.getItem("theme") || ""

function setMode(mode) {
  switch (mode) {
    case "dark":
      currentMode = "dark"
      break
    case "system":
      currentMode = "system"
      break
    default:
      currentMode = "light"
      break
  }
  localStorage.setItem("mode", currentMode)
  updateMode()
}

const themes = {
  "red": "red",
  "rose": "rose",
  "orange": "orange",
  "green": "green",
  "blue": "blue",
  "yellow": "yellow",
  "violet": "violet",
  "stone": "stone",
  "zinc": "zinc",
  "neutral": "neutral",
  "gray": "gray",
  "slate": "slate",
}

function setTheme(theme) {
  if (!Object.keys(themes).includes(theme)) {
    currentTheme = ""
  } else {
    currentTheme = themes[theme]
  }
  localStorage.setItem("theme", currentTheme)
  updateTheme()
}

const updateMode = () => {
  modeButton.innerHTML = modeIcons[currentMode]
  const root = window.document.documentElement
  root.classList.remove("light", "dark")

  if (currentMode === "system") {
    const systemmode = window.matchMedia("(prefers-color-scheme: dark)")
      .matches
      ? "dark"
      : "light"
    root.classList.add(systemmode)
    return
  }

  if (currentMode === "dark") {
    root.classList.add("dark")
  } else {
    root.classList.add("light")
  }
}

const updateTheme = () => {
  // modeButton.innerHTML = modeIcons[currentMode]
  const root = window.document.documentElement

  Object.values(themes).forEach(theme => {
    root.classList.remove(theme)
  })

  if (currentTheme !== "") {
    root.classList.add(currentTheme)
  }
}

updateMode()
updateTheme()

window.matchMedia("(prefers-color-scheme: dark)").addEventListener("change", () => {
  if (!currentMode === "system") return
  updateMode()
})

