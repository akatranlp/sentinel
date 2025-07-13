import type { ExtractSentinelCtx } from "@/context/types";
import type { PageProps } from "./PageProps";

export default function FormRedirect(props: PageProps<ExtractSentinelCtx<"form-redirect.tmpl">>) {
  const { sentinelCtx, Template } = props
  return <Template
    sentinelCtx={sentinelCtx}
  >
    <div />
  </Template>
}
