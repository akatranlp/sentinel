import type { Provider } from "@/context/types";

import { Gitlab } from "@/components/icons/Gitlab";
import { Github } from "@/components/icons/Github";
import { Gitea } from "@/components/icons/Gitea";


export const ProviderIcon = (props: { provider: Provider }) => {
  const { provider: p } = props
  switch (p.alias) {
    case "gitlab":
      return <Gitlab className="size-6" />
    case "github":
      return <Github className="size-6" />
    case "gitea":
      return <Gitea className="size-6" />
    default:
      return <img className="size-6" src={p.icon} />
  }
}
