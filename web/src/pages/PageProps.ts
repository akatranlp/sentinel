import type { TemplateProps } from "@/TemplateProps";
import type { ReactElement } from "react";

import type { LazyExoticComponent, ComponentType } from "react";

export type LazyOrNot<Component extends ComponentType<any>> =
  | LazyExoticComponent<Component>
  | Component;

export type PageProps<NarrowedSentinelCtx> = {
  Template: LazyOrNot<(props: TemplateProps<any>) => ReactElement | null>;
  sentinelCtx: NarrowedSentinelCtx;
};
