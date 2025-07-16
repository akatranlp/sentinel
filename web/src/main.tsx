import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import { SentinelPage } from './SentinelPage.tsx'

// The following block can be uncommented to test a specific page with `pnpm dev`
// Don't forget to comment back or your bundle size will increase
// import { getSentinelContextMock } from "./context/mock.ts";
//
// if (import.meta.env.DEV) {
//   window.sentinelCtx = getSentinelContextMock({
//     pageId: "user-edit.tmpl",
//     overrides: {
//       // message: {
//       //   summary: "Fabse is cool!"
//       // }
//     }
//   });
// }

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    {!window.sentinelCtx ? (
      <h1>No Sentinel Context</h1>
    ) : (
      < SentinelPage sentinelCtx={window.sentinelCtx} />
    )}
  </StrictMode>,
)
