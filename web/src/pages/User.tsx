import type { CSRF, ExtractSentinelCtx, Provider } from "@/context/types";
import type { PageProps } from "./PageProps";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { ProviderIcon } from "@/components/provider-icon";
import { Input } from "@/components/ui/input";
import { DialogDescription, DialogFooter, DialogHeader, DialogTitle, Dialog, DialogContent, DialogTrigger } from "@/components/ui/dialog";
import { useState } from "react";

const LinkForm = (props: { provider: Provider, csrf: CSRF }) => {
  const [open, setOpen] = useState(false)

  const { provider: p, csrf } = props

  const actionName = p.isLinked ? "Unlink" : "Link"

  return <Dialog open={open} onOpenChange={(open) => setOpen(open)}>
    <DialogTrigger asChild>
      <Button onClick={() => setOpen(true)} variant="outline" className="w-full">
        <span className="col-start-2">
          {/*"h-[28px] w-[28px]"*/}
          <ProviderIcon provider={p} />
        </span>
        <p className="col-start-3">{actionName} with {p.displayName}</p>
      </Button>
    </DialogTrigger>
    <DialogContent>
      <form method="POST" key={p.providerId} action={p.loginUrl} className="w-full">
        <DialogHeader>
          <DialogTitle>
            Are you absolutely sure?
          </DialogTitle>
          <DialogDescription>
            This action cannot be undone.
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <input type="hidden" name={csrf.fieldName} value={csrf.value} />
          <Button variant="outline" type="button" onClick={() => setOpen(false)} className="mt-2 sm:mt-0">Cancel</Button>
          <Button type="submit" variant="default">
            {actionName}
          </Button>
        </DialogFooter>
      </form>
    </DialogContent>
  </Dialog>
}

export default function User(props: PageProps<ExtractSentinelCtx<"user.tmpl">>) {
  const { sentinelCtx, Template } = props

  const { urls, user, providers, csrf } = sentinelCtx

  return <Template
    sentinelCtx={sentinelCtx}
  >
    <div className="flex w-full max-w-lg flex-col gap-6">
      <div className="flex flex-col gap-6">
        <Card>
          <CardHeader className="text-center">
            <CardTitle className="text-xl">User Page</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex flex-col gap-6">
              <div className="grid gap-2 grid-cols-1 md:grid-cols-2">
                <img className="size-32 row-span-3 rounded-full self-center justify-self-center" alt="Profile-picture" src={user.picture} />
                <Input disabled value={user.name} />
                <Input disabled value={user.username} />
                <Input type="email" disabled value={user.email} />
                <Button variant="outline" className="col-start-2" asChild>
                  <a href={urls.basePath + "/user/edit"}>Edit</a>
                </Button>
              </div>
              <div className="flex flex-col gap-4">
                {providers?.map((p) =>
                  <LinkForm {...{ provider: p, csrf }} />
                )}
              </div>
            </div >
          </CardContent >
        </Card >
      </div >
    </div >
  </Template >
}

