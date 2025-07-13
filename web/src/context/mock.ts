import { deepAssign } from "./deepAssign";
import { type PageID, type ExtractSentinelCtx, type LoginSentinelCtx, type ErrorSentinelCtx, type FormRedirectSentinelCtx, type FormPostSentinelCtx, type InfoSentinelCtx } from "./types"

export type DeepPartial<T> = {
  [P in keyof T]?: DeepPartial<T[P]>;
};

export const getSentinelContextMock = <T extends PageID>(args: { pageId: T, overrides?: DeepPartial<ExtractSentinelCtx<T>> }): ExtractSentinelCtx<T> => {
  const { pageId, overrides } = args

  const sentinelCtx = structuredClone(mock[pageId])

  if (overrides)
    deepAssign({
      target: sentinelCtx,
      source: overrides
    })

  return sentinelCtx
}

const mock: { [K in PageID]: ExtractSentinelCtx<K> } = {
  "login.tmpl": {
    pageId: "login.tmpl",
    message: null,
    csrf: {
      fieldName: "csrf-token",
      value: "uL8/kNYEMPvnwibG8Ap2fsxsHmf1bN3qGO4hm5i4atG0VaniWSv5REllN2RUotO9yP2I7f1SJAcqAjNcPqJEuw=="
    },
    messages: null,
    providers: [
      {
        "loginUrl": "/auth/github/login",
        "alias": "github",
        "providerId": "github",
        "displayName": "GitHub",
        "icon": "/auth/assets/github.svg"
      },
      {
        "loginUrl": "",
        "alias": "gitlab",
        "providerId": "gitlab",
        "displayName": "GitLab",
        "icon": "/auth/assets/gitlab.svg"
      }
    ],
    urls: {
      basePath: "/auth",
      resourcePath: "/auth/assets/"
    }
  } satisfies LoginSentinelCtx,
  "info.tmpl": {
    pageId: "info.tmpl",
    message: {
      type: "info",
      summary: "Everything is ok",
    },
    messages: null,
    urls: {
      basePath: "/auth",
      resourcePath: ""
    }
  } satisfies InfoSentinelCtx,
  "error.tmpl": {
    pageId: "error.tmpl",
    message: {
      type: "error",
      summary: "An Error occured",
    },
    messages: null,
    urls: {
      basePath: "/auth",
      resourcePath: ""
    }
  } satisfies ErrorSentinelCtx,
  "form-redirect.tmpl": {
    pageId: "form-redirect.tmpl",
    message: null,
    messages: null,
    urls: {
      basePath: "/auth",
      resourcePath: ""
    },
    redirectUrl: "localhost:3000"
  } satisfies FormRedirectSentinelCtx,
  "form-post.tmpl": {
    pageId: "form-post.tmpl",
    message: null,
    messages: null,
    urls: {
      basePath: "/auth",
      resourcePath: ""
    },
    redirectUrl: "localhost:3000"
  } satisfies FormPostSentinelCtx,
}

