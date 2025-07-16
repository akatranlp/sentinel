import { lazy, Suspense, type ReactNode } from "react";
import Template from "./Template";
import type { SentinelCtx } from "./context/types";
import { assert, type Equals } from "tsafe/assert";

declare global {
  interface Window {
    sentinelCtx?: SentinelCtx;
  }
}

export const Login = lazy(() => import("./pages/Login"));
export const Info = lazy(() => import("./pages/Info"));
export const Error = lazy(() => import("./pages/Error"));
export const FormRedirect = lazy(() => import("./pages/FormRedirect"));
export const FormPost = lazy(() => import("./pages/FormPost"));
export const User = lazy(() => import("./pages/User"));
export const UserEdit = lazy(() => import("./pages/UserEdit"));
// export const Terms = lazy(() => import("./pages/Terms"));
// export const Code = lazy(() => import("./pages/Code"));
// export const DeleteAccountConfirm = lazy(() => import("./pages/DeleteAccountConfirm"));
// export const LogoutConfirm = lazy(() => import("./pages/LogoutConfirm"));
// export const LoginIdpLinkConfirm = lazy(() => import("./pages/LoginIdpLinkConfirm"));
// export const LoginIdpLinkEmail = lazy(() => import("./pages/LoginIdpLinkEmail"));
// export const LoginPageExpired = lazy(() => import("./pages/LoginPageExpired"));

export function SentinelPage(props: { sentinelCtx: SentinelCtx; fallback?: ReactNode }) {
  const { sentinelCtx, fallback } = props
  return <Suspense fallback={fallback}>
    {(() => {
      switch (sentinelCtx.pageId) {
        case "login.tmpl":
          return <Login {...{ sentinelCtx }} Template={Template} />;
        case "info.tmpl":
          return <Info {...{ sentinelCtx }} Template={Template} />;
        case "error.tmpl":
          return <Error {...{ sentinelCtx }} Template={Template} />;
        case "form-redirect.tmpl":
          return <FormRedirect {...{ sentinelCtx }} Template={Template} />;
        case "form-post.tmpl":
          return <FormPost {...{ sentinelCtx }} Template={Template} />;
        case "user.tmpl":
          return <User {...{ sentinelCtx }} Template={Template} />;
        case "user-edit.tmpl":
          return <UserEdit {...{ sentinelCtx }} Template={Template} />;
          {/* case "terms.tmpl": */ }
          {/*   return <Terms {...{ sentinelCtx }} Template={Template} />; */ }
          {/* case "login-idp-link-confirm.tmpl": */ }
          {/*   return <LoginIdpLinkConfirm {...{ sentinelCtx }} Template={Template} />; */ }
          {/* case "login-idp-link-email.tmpl": */ }
          {/*   return <LoginIdpLinkEmail {...{ sentinelCtx }} Template={Template} />; */ }
          {/* case "login-page-expired.tmpl": */ }
          {/*   return <LoginPageExpired {...{ sentinelCtx }} Template={Template} />; */ }
          {/* case "logout-confirm.tmpl": */ }
          {/*   return <LogoutConfirm {...{ sentinelCtx }} Template={Template} />; */ }
          {/* case "code.tmpl": */ }
          {/*   return <Code {...{ sentinelCtx }} Template={Template} />; */ }
          {/* case "delete-account-confirm.tmpl": */ }
          {/*   return <DeleteAccountConfirm {...{ sentinelCtx }} Template={Template} />; */ }
      }
      assert<Equals<typeof sentinelCtx, never>>(false);
    })()}
  </Suspense>
}

