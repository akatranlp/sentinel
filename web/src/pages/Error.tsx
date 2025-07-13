import type { ExtractSentinelCtx } from "@/context/types";
import type { PageProps } from "./PageProps";

export default function Error(props: PageProps<ExtractSentinelCtx<"error.tmpl">>) {
  const { sentinelCtx, Template } = props
  return <Template
    sentinelCtx={sentinelCtx}
  >
    <div />
  </Template>
}
