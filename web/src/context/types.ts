export const pageID = ["login.tmpl", "error.tmpl", "info.tmpl", "form-redirect.tmpl", "form-post.tmpl"] as const;
export type PageID = typeof pageID[number]

export type CommonSentinelCtx = {
  pageId: PageID;
  message: Message | null;
  messages: Message[] | null;
  urls: URLs;
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
}

export type CSRF = {
  fieldName: string;
  value: string;
}

export type URLs = {
  basePath: string;
  resourcePath: string;
}

export type SentinelCtx = 
  LoginSentinelCtx |
  ErrorSentinelCtx |
  InfoSentinelCtx |
  FormPostSentinelCtx |
  FormRedirectSentinelCtx

export type ExtractSentinelCtx<T extends PageID> = Extract<SentinelCtx, { pageId: T }>

export type Prettify<T> = {
  [K in keyof T]: T[K]
}

