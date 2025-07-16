export const pageID = ["login.tmpl", "error.tmpl", "info.tmpl", "form-redirect.tmpl", "form-post.tmpl", "user.tmpl", "user-edit.tmpl", "logout.tmpl"] as const;
export type PageID = typeof pageID[number]

export type CommonSentinelCtx = {
  pageId: PageID;
  message: Message | null;
  messages: Message[] | null;
  urls: URLs;
  user: User | null;
}

export type LoginSentinelCtx = CommonSentinelCtx & {
  pageId: "login.tmpl"
  providers: Provider[] | null;
  csrf: CSRF;
}

export type FormRedirectSentinelCtx = CommonSentinelCtx & {
  pageId: "form-redirect.tmpl"
  redirectUrl: string;
}

export type FormPostSentinelCtx = CommonSentinelCtx & {
  pageId: "form-post.tmpl"
  redirectUrl: string;
}

export type InfoSentinelCtx = CommonSentinelCtx & {
  pageId: "info.tmpl"
  message: Message;
}

export type ErrorSentinelCtx = CommonSentinelCtx & {
  pageId: "error.tmpl"
  message: Message;
}

export type UserSentinelCtx = CommonSentinelCtx & {
  pageId: "user.tmpl"
  user: User;
  accounts: Account[] | null;
  providers: Provider[] | null;
  csrf: CSRF;
}

export type UserEditSentinelCtx = CommonSentinelCtx & {
  pageId: "user-edit.tmpl"
  user: User;
  accounts: Account[] | null;
  providers: Provider[] | null;
  csrf: CSRF;
}

export type LogoutSentinelCtx = CommonSentinelCtx & {
  pageId: "logout.tmpl"
  csrf: CSRF;
  redirect: string;
  sessionId: string;
}

export const messageType = ["info", "success", "error", "warning"] as const;
export type MessageType = typeof messageType[number]

export type Message = {
  type: MessageType;
  summary: string;
}

export type Provider = {
  loginUrl: string;
  alias: string;
  providerId: string;
  displayName: string;
  icon: string;
  isLinked: boolean;
}

export type CSRF = {
  fieldName: string;
  value: string;
}

export type URLs = {
  basePath: string;
  resourcePath: string;
}

export type User = {
  id: string;
  name: string;
  username: string;
  picture: string;
  email: string;
}

export type Account = {
  provider: string;
  email: string;
  name: string;
  username: string;
  picture: string;
}

export type SentinelCtx = 
  ErrorSentinelCtx |
  FormPostSentinelCtx |
  FormRedirectSentinelCtx |
  InfoSentinelCtx |
  LoginSentinelCtx |
  LogoutSentinelCtx |
  UserEditSentinelCtx |
  UserSentinelCtx

export type ExtractSentinelCtx<T extends PageID> = Extract<SentinelCtx, { pageId: T }>

export type Prettify<T> = {
  [K in keyof T]: T[K]
}

