import type { ExtractSentinelCtx } from "@/context/types";
import type { PageProps } from "./PageProps";

export default function Info(props: PageProps<ExtractSentinelCtx<"info.tmpl">>) {
  const { sentinelCtx, Template } = props
  return <Template
    sentinelCtx={sentinelCtx}
  >
    <div />
  </Template>
}
