import fs from "fs/promises";

async function main() {
  let data = await fs.readFile("index.html", { encoding: "utf8" })
  const regex = new RegExp(/<link rel="icon" type="(.*)" href="(.*)" \/>/)
  const match = data.match(regex)
  let faviconLine = match?.[0]
  if (!faviconLine) throw new Error("no favicon found in the index.html")

  faviconLine = faviconLine.replace(/href="(.*)"/, 'href="{{ .URLs.ResourcePath }}/dist$1"')

  const files = await fs.readdir("dist/assets")

  const indexJSRegex = new RegExp(/index-.*\.js/)
  const indexJS = files.find((f) => indexJSRegex.test(f))
  const indexCSSRegex = new RegExp(/index-.*\.css/)
  const indexCSS = files.find((f) => indexCSSRegex.test(f))

  if (!indexJS || !indexCSS) throw new Error("Haven't found an indexJS or indexCSS! Maybe you haven't built the frontend yet?")

  const commonTemplatePath = "common.tmpl.html"
  let commonTemplate = await fs.readFile("src/templates/" + commonTemplatePath, { encoding: "utf8" })
  commonTemplate = commonTemplate.replace("$$favicon$$", faviconLine)
  commonTemplate = commonTemplate.replace("$$index-js$$", indexJS)
  commonTemplate = commonTemplate.replace("$$index-css$$", indexCSS)

  for await (const entry of fs.glob("src/templates/*.tmpl.html")) {
    if (entry.endsWith(commonTemplatePath)) continue
    await fs.copyFile(entry, entry.replace("src/templates/", "dist_templates/"))
  }

  await fs.mkdir("dist_templates", { recursive: true })
  await fs.writeFile("dist_templates/" + commonTemplatePath, commonTemplate, { encoding: "utf8" })
}

await main()

