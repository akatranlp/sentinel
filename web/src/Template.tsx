import { useEffect } from "react";
import type { SentinelCtx } from "./context/types";
import type { TemplateProps } from "./TemplateProps";
import { GalleryVerticalEnd, LogOut } from "lucide-react";
// import { Button } from "./components/ui/button";
import { ThemeProvider } from "./components/theme-provider";
import { ThemeToggle } from "./components/theme-toggle";
import { ModeToggle } from "./components/mode-toggle";
import { Alert, AlertDescription, AlertTitle } from "./components/ui/alert";
import { Terminal } from "lucide-react";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "./components/ui/dropdown-menu";
import { Avatar, AvatarFallback, AvatarImage } from "./components/ui/avatar";
export default function Template(props: TemplateProps<SentinelCtx>) {
  const {
    documentTitle,
    sentinelCtx,
    children
  } = props;

  const { message, pageId, urls, user } = sentinelCtx;

  // TODO: add appURL, user
  // const appURL = undefined
  // const user = undefined

  useEffect(() => {
    document.title = documentTitle ?? pageId;
  }, []);

  {/* return <div> */ }
  {/*   <div> */ }
  {/*     {message && message.summary} */ }
  {/*   </div> */ }
  {/*   {children} */ }
  {/* </div> */ }

  return (
    <ThemeProvider>
      <div className="fixed z-50 w-full bg-sidebar text-sidebar-foreground top-0 left-0 h-16">
        <div className="h-full flex justify-between items-center m-auto max-w-5xl">
          <div className="flex gap-2">
            {/* <Button variant="link" asChild> */}
            <a href={urls.basePath} className="flex items-center gap-2 self-center font-medium">
              <span className="flex h-6 w-6 items-center justify-center rounded-md bg-primary text-primary-foreground">
                <GalleryVerticalEnd className="size-4" />
              </span>
              <span>GitClassrooms</span>
            </a>
            {/* </Button> */}
            {/* if appURL != "" { */}
            {/*  <a href={ templ.SafeURL(appURL) } className="flex items-center gap-2 self-center font-medium"> */}
            {/*   <div className="flex h-6 w-6 items-center justify-center rounded-md bg-primary text-primary-foreground"> */}
            {/*    <svg xmlns="http://www.w3.org/2000/svg" className="size-4" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-gallery-vertical-end-icon lucide-gallery-vertical-end"><path d="M7 2h10"></path><path d="M5 6h14"></path><rect width="18" height="12" x="3" y="10" rx="2"></rect></svg> */}
            {/*   </div> */}
            {/*   GitClassrooms App */}
            {/*  </a> */}
            {/* } */}
          </div>
          <div className="flex gap-4">
            <ThemeToggle />
            <ModeToggle />
            {user && (
              <DropdownMenu>
                <DropdownMenuTrigger>
                  <Avatar>
                    <AvatarImage src={user.picture} alt={"@" + user.username} />
                    <AvatarFallback>{user.name.split(" ").map(v => v[0])}</AvatarFallback>
                  </Avatar>
                </DropdownMenuTrigger>
                <DropdownMenuContent>
                  <DropdownMenuItem className="px-4 text-sm font-bold">
                    {"@" + user.username}
                  </DropdownMenuItem>
                  <DropdownMenuItem className="px-4" asChild>
                    <a
                      href={urls.basePath + "/logout"}
                    >
                      <LogOut /> Logout
                    </a>
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            )}
          </div>
        </div>
      </div>
      <div className="fixed empty:hidden bg-sidebar top-16 left-[50%] translate-x-[-50%] z-40 p-2 rounded-b-xl w-full max-w-3xl flex flex-col gap-2">
        {message &&
          <Alert variant={message.type === "error" ? "destructive" : "default"}>
            <Terminal />
            <AlertTitle>{message.summary}</AlertTitle>
            <AlertDescription>{message.summary}</AlertDescription>
          </Alert>
        }
      </div>
      <div className="min-h-svh pt-22 md:pt-26 flex flex-col items-center gap-6 bg-muted px-6 md:px-10">
        {children}
      </div>
    </ThemeProvider>
  )
}
