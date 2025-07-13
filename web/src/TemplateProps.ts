import type { ReactNode } from "react";

export type TemplateProps<SentinelCtx> = {
  sentinelCtx: SentinelCtx;
  children: ReactNode;

  // displayInfo?: boolean;
  // displayMessage?: boolean;
  // displayRequiredFields?: boolean;
  // headerNode: ReactNode;
  // infoNode?: ReactNode;
  documentTitle?: string;
  // bodyClassName?: string;
};

