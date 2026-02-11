import { useState } from "react";
import { Link, useNavigate } from "@tanstack/react-router";
import { useAuth } from "@/contexts/AuthContext";
import Button from "@/components/common/Button";
import Input from "@/components/common/Input";
import Card from "@/components/common/Card";
import GopherLogo from "@/components/common/GopherLogo";

export default function LoginPage() {
  const navigate = useNavigate();
  const { login, isLoading, error } = useAuth();

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();

    const success = await login({ email, password });
    if (success) {
      navigate({ to: "/dashboard" });
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        <div className="text-center">
          <h1 className="text-4xl font-bold flex items-center justify-center gap-2">
            <GopherLogo className="h-10 w-10" />
            <span className="text-primary-600">Mini</span>
            <span className="text-gray-600 dark:text-gray-400 font-light">
              Bank
            </span>
          </h1>
          <p className="mt-2 text-gray-600 dark:text-gray-400">
            Sign in to your account
          </p>
        </div>

        <Card>
          <form onSubmit={handleSubmit} className="space-y-6">
            <Input
              value={email}
              onChange={setEmail}
              label="Email"
              type="email"
              placeholder="you@example.com"
              required
            />

            <Input
              value={password}
              onChange={setPassword}
              label="Password"
              type="password"
              placeholder="Enter your password"
              required
            />

            {error && (
              <div className="bg-red-50 dark:bg-red-900/30 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-400 px-4 py-3 rounded-lg text-sm">
                {error}
              </div>
            )}

            <Button type="submit" loading={isLoading} className="w-full">
              Sign In
            </Button>

            <p className="text-center text-sm text-gray-600 dark:text-gray-400">
              Don't have an account?{" "}
              <Link
                to="/register"
                className="text-primary-600 hover:text-primary-700 font-medium"
              >
                Sign up
              </Link>
            </p>
          </form>
        </Card>
      </div>
    </div>
  );
}
