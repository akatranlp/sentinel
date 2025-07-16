import type { ExtractSentinelCtx } from "@/context/types";
import type { PageProps } from "./PageProps";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";

export default function Logout(props: PageProps<ExtractSentinelCtx<"logout.tmpl">>) {
  const { sentinelCtx, Template } = props

  const { urls, user, csrf, sessionId, redirect } = sentinelCtx

  return <Template sentinelCtx={sentinelCtx}>
    <div className="flex w-full max-w-lg flex-col gap-6">
      <div className="flex flex-col gap-6">
        <Card>
          <CardHeader className="text-center">
            <CardTitle className="text-xl">Logout</CardTitle>
            <CardDescription>Willst du dich wirklich abmelden?</CardDescription>
          </CardHeader>
          <CardContent>
            <form method="POST" className="flex flex-col gap-4">
              <input type="hidden" name={csrf.fieldName} value={csrf.value} />
              <input type="hidden" name="sid" value={sessionId} />
              <input type="hidden" name="redirect" value={redirect} />
              <img
                className="size-32 rounded-full self-center"
                alt="Profile-picture"
                src={user?.picture}
              />
              <div className="text-center">
                <p className="font-semibold">{user?.name}</p>
                <p className="text-sm text-muted-foreground">{user?.email}</p>
              </div>
              <div className="flex gap-4 justify-center mt-4">
                <Button variant="outline" asChild>
                  <a href={urls.basePath}>No Stay Here</a>
                </Button>
                <Button type="submit">
                  Yes, Log Me Out
                </Button>
              </div>
            </form>
          </CardContent>
        </Card>
      </div>
    </div>
  </Template>
}
