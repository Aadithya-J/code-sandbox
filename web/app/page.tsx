import { SignUpDialog } from "@/components/auth/register";
import { LoginDialog } from "../components/auth/login";

export default function Home() {
  return (
    <main className="flex min-h-screen flex-col items-center justify-between p-24">
      <LoginDialog />
      <SignUpDialog />
    </main>
  );
}