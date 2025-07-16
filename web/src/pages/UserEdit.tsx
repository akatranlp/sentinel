import type { Account, ExtractSentinelCtx, User } from "@/context/types";
import type { PageProps } from "./PageProps";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { useMemo, useState } from "react";
import { SelectTrigger, Select, SelectContent, SelectItem, SelectValue } from "@/components/ui/select";

const RenderSelect = (props: { fieldName: keyof Omit<User, "id">, user: Omit<User, "id">, accounts: Account[] }) => {
  const { fieldName: name, user, accounts } = props

  const [provider, setProvider] = useState(() => {
    const acc = accounts.find((a) => a[name] === user[name])
    if (!acc) return accounts[0].provider
    return acc.provider
  })

  const selected = useMemo(() => {
    const acc = accounts.find((a) => a.provider === provider)
    if (!acc) return accounts[0][name]
    return acc[name]
  }, [provider])

  return <Select name={name} onValueChange={(v) => setProvider(v)} defaultValue={provider}>
    <SelectTrigger className="w-full">
      <SelectValue placeholder={`Select one of your ${name}s to be your active ones`} defaultValue={selected} />
    </SelectTrigger>
    <SelectContent>
      {accounts?.map(a =>
        <SelectItem key={a.provider} value={a.provider}>{a[name]}</SelectItem>
      )}
    </SelectContent>
  </Select>
}

export default function UserEdit(props: PageProps<ExtractSentinelCtx<"user-edit.tmpl">>) {
  const { sentinelCtx, Template } = props

  const { urls, user, accounts: accs, csrf } = sentinelCtx
  const accounts = accs!

  const [pictureProvider, setPictureProvider] = useState(() => {
    const acc = accounts.find((a) => a.picture === user.picture)
    if (!acc) return accounts[0].provider
    return acc.provider
  })

  const picture = useMemo(() => {
    const acc = accounts.find((a) => a.provider === pictureProvider)
    if (!acc) return accounts[0].picture
    return acc.picture
  }, [pictureProvider])

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
              <form method="POST" className="grid gap-2 grid-cols-1 md:grid-cols-2">
                <input type="hidden" name={csrf.fieldName} value={csrf.value} />
                <img className="size-32 row-span-3 rounded-full self-center justify-self-center" alt="Profile-picture" src={picture} />
                <RenderSelect {...{ fieldName: "name", user, accounts: accounts }} />
                <RenderSelect {...{ fieldName: "username", user, accounts: accounts }} />
                <RenderSelect {...{ fieldName: "email", user, accounts: accounts }} />
                <Select name="picture" onValueChange={(v) => setPictureProvider(v)} defaultValue={pictureProvider}>
                  <SelectTrigger className="w-full">
                    <SelectValue placeholder={`Select one of your pictures to be your active ones`} />
                  </SelectTrigger>
                  <SelectContent>
                    {accounts?.map(a =>
                      <SelectItem key={a.provider} value={a.provider}>{a.provider}</SelectItem>
                    )}
                  </SelectContent>
                </Select>
                <div className="flex justify-between">
                  <Button variant="outline" asChild>
                    <a href={urls.basePath + "/user"}>Cancel</a>
                  </Button>
                  <Button type="submit">Edit</Button>
                </div>
              </form>
            </div>
          </CardContent >
        </Card >
      </div >
    </div >
  </Template >
}

