import type { ExtractSentinelCtx } from "@/context/types";
import type { PageProps } from "./PageProps";
import { GalleryVerticalEnd } from "lucide-react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { ProviderIcon } from "@/components/provider-icon";

export default function Login(props: PageProps<ExtractSentinelCtx<"login.tmpl">>) {
  const { sentinelCtx, Template } = props

  const { urls, providers, csrf } = sentinelCtx

  return <Template
    sentinelCtx={sentinelCtx}
  >
    <div className="flex w-full max-w-sm flex-col gap-6">
      <a href={urls.basePath} className="flex items-center gap-2 self-center font-medium">
        <span className="flex h-6 w-6 items-center justify-center rounded-md bg-primary text-primary-foreground">
          <GalleryVerticalEnd className="size-4" />
        </span>
        <span>GitClassrooms</span>
      </a>
      <div className="flex flex-col gap-6">
        <Card>
          <CardHeader className="text-center">
            <CardTitle className="text-xl">Welcome Back</CardTitle>
            <CardDescription>Login with on of the following Providers</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="grid gap-6">
              <div className="flex flex-col gap-4">
                {providers?.map((p) =>
                  <form
                    method="POST"
                    action={p.loginUrl}
                    key={p.providerId}
                    className="w-full"
                  >
                    <input type="hidden" name={csrf.fieldName} value={csrf.value} />
                    <Button variant="outline" className="w-full">
                      <span className="col-start-2">
                        <ProviderIcon provider={p} />
                      </span>
                      <p>Login with {p.displayName}</p>
                    </Button>
                  </form>
                )}
              </div>
            </div>
          </CardContent>
        </Card>
        <div className="text-balance text-center text-xs text-muted-foreground [&_a]:underline [&_a]:underline-offset-4 [&_a]:hover:text-primary  ">
          By clicking continue, you agree to our <a href="#">Terms of Service</a> and <a href="#">Privacy Policy</a>.
        </div>
      </div>
    </div>
  </Template>
}

