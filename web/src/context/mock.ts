import { deepAssign } from "./deepAssign";
import type { PageID, ExtractSentinelCtx, LoginSentinelCtx, ErrorSentinelCtx, FormRedirectSentinelCtx, FormPostSentinelCtx, InfoSentinelCtx, UserSentinelCtx, User, Account, Provider, URLs, CSRF, UserEditSentinelCtx } from "./types"

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

const urls: URLs = {
  basePath: "/auth",
  resourcePath: "/auth/assets/"
}

const csrf: CSRF = {
  fieldName: "csrf-token",
  value: "uL8/kNYEMPvnwibG8Ap2fsxsHmf1bN3qGO4hm5i4atG0VaniWSv5REllN2RUotO9yP2I7f1SJAcqAjNcPqJEuw=="
}

const user: User = {
  id: "1",
  email: "toni@tester.de",
  name: "Toni Tester",
  picture: "http://www.gravatar.com/avatar/toni@tester.de?size=150",
  username: "tonitester",
}

const accounts: Account[] = [
  {
    provider: "gitlab",
    email: "toni@tester.de",
    name: "Toni Tester",
    picture: "http://www.gravatar.com/avatar/toni@tester.de?size=150",
    username: "tonitester",
  }
]

const providers: Provider[] = [
  {
    "loginUrl": "/auth/github/login",
    "alias": "github",
    "providerId": "github",
    "displayName": "GitHub",
    "icon": "/auth/assets/github.svg",
    "isLinked": false,
  },
  {
    "loginUrl": "",
    "alias": "gitlab",
    "providerId": "gitlab",
    "displayName": "GitLab",
    "icon": "/auth/assets/gitlab.svg",
    "isLinked": true,
  }
]

const mock: { [K in PageID]: ExtractSentinelCtx<K> } = {
  "login.tmpl": {
    pageId: "login.tmpl",
    user: null,
    message: null,
    urls,
    csrf,
    messages: null,
    providers,
  } satisfies LoginSentinelCtx,
  "user.tmpl": {
    pageId: "user.tmpl",
    message: null,
    messages: null,
    csrf,
    user,
    accounts,
    providers,
    urls,
  } satisfies UserSentinelCtx,
  "user-edit.tmpl": {
    pageId: "user-edit.tmpl",
    message: null,
    messages: null,
    csrf,
    user,
    accounts,
    providers,
    urls,
  } satisfies UserEditSentinelCtx,
  "info.tmpl": {
    pageId: "info.tmpl",
    user: null,
    message: {
      type: "info",
      summary: "Everything is ok",
    },
    messages: null,
    urls,
  } satisfies InfoSentinelCtx,
  "error.tmpl": {
    pageId: "error.tmpl",
    user: null,
    message: {
      type: "error",
      summary: "An Error occured",
    },
    messages: null,
    urls,
  } satisfies ErrorSentinelCtx,
  "form-redirect.tmpl": {
    user,
    pageId: "form-redirect.tmpl",
    message: null,
    messages: null,
    urls,
    redirectUrl: "localhost:3000"
  } satisfies FormRedirectSentinelCtx,
  "form-post.tmpl": {
    user,
    pageId: "form-post.tmpl",
    message: null,
    messages: null,
    urls,
    redirectUrl: "localhost:3000"
  } satisfies FormPostSentinelCtx,
}

