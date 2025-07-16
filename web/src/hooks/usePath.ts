import { use } from "react"
import { RouterContext } from "./RouterProvider"

export const useRouter = () => {
  const context = use(RouterContext)

  if (!context) throw new Error("not within a router context")

  return context
}
