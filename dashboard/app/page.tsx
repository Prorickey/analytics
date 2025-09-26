import { auth, signIn, signOut } from "@/auth";
import { Dashboard } from "@/components/dashboard";

export default async function Home() {
  const session = await auth()
  if(!session) return signIn()
  return <Dashboard />
}
