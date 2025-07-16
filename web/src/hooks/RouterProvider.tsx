import type { PageID, SentinelCtx } from "@/context/types"
import { createContext, useCallback, useMemo, useState } from "react"

type NavigateFn = (opts: { to: string, toPageId: PageID }) => void

interface RouterContextState {
  path: string
  sentinelCtx: SentinelCtx
  navigate: NavigateFn
}

export const RouterContext = createContext<RouterContextState | null>(null)

export const RouterProvider = ({ children, ...initalSentinelCtx }: React.PropsWithChildren<SentinelCtx>) => {
  const [path, setPath] = useState<string>(window.location.pathname)
  const [pageId, setPageId] = useState<PageID>(initalSentinelCtx.pageId)

  const navigate = useCallback<NavigateFn>(({ to, toPageId }) => {
    setPath(to)
    setPageId(toPageId)
    window.history.pushState({}, "", to)
  }, [])

  const sentinelCtx = useMemo(() => ({ ...initalSentinelCtx, pageId }), [pageId])

  return <RouterContext value={{ path, navigate, sentinelCtx: sentinelCtx as SentinelCtx }}> {children}</RouterContext>
}
