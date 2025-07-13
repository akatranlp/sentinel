import type { ExtractSentinelCtx } from "@/context/types";
import type { PageProps } from "./PageProps";

export default function FormPost(props: PageProps<ExtractSentinelCtx<"form-post.tmpl">>) {
  const { sentinelCtx, Template } = props
  return <Template
    sentinelCtx={sentinelCtx}
  >
    <div />
  </Template>
}
