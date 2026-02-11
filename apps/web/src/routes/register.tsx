import { useState } from "react";
import { Link, useNavigate } from "@tanstack/react-router";
import { useAuth } from "@/contexts/AuthContext";
import Button from "@/components/common/Button";
import Input from "@/components/common/Input";
import Card from "@/components/common/Card";

export default function RegisterPage() {
  const navigate = useNavigate();
  const { register, isLoading, error } = useAuth();

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [fullName, setFullName] = useState("");
  const [phone, setPhone] = useState("");
  const [passwordError, setPasswordError] = useState<string | null>(null);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setPasswordError(null);

    if (password !== confirmPassword) {
      setPasswordError("Passwords do not match");
      return;
    }

    if (password.length < 8) {
      setPasswordError("Password must be at least 8 characters");
      return;
    }

    const success = await register({
      email,
      password,
      full_name: fullName,
      phone,
    });

    if (success) {
      navigate({ to: "/dashboard" });
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        <div className="text-center">
          <h1 className="text-4xl font-bold">
            <span className="text-primary-600">Mini</span>
            <span className="text-gray-600 dark:text-gray-400 font-light">
              Bank
            </span>
          </h1>
          <p className="mt-2 text-gray-600 dark:text-gray-400">
            Create your account
          </p>
        </div>

        <Card>
          <form onSubmit={handleSubmit} className="space-y-6">
            <Input
              value={fullName}
              onChange={setFullName}
              label="Full Name"
              placeholder="John Doe"
              required
            />

            <Input
              value={email}
              onChange={setEmail}
              label="Email"
              type="email"
              placeholder="you@example.com"
              required
            />

            <Input
              value={phone}
              onChange={setPhone}
              label="Phone (optional)"
              type="tel"
              placeholder="+1 (555) 000-0000"
            />

            <Input
              value={password}
              onChange={setPassword}
              label="Password"
              type="password"
              placeholder="At least 8 characters"
              required
              error={passwordError || undefined}
            />

            <Input
              value={confirmPassword}
              onChange={setConfirmPassword}
              label="Confirm Password"
              type="password"
              placeholder="Repeat your password"
              required
            />

            {error && (
              <div className="bg-red-50 dark:bg-red-900/30 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-400 px-4 py-3 rounded-lg text-sm">
                {error}
              </div>
            )}

            <Button type="submit" loading={isLoading} className="w-full">
              Create Account
            </Button>

            <p className="text-center text-sm text-gray-600 dark:text-gray-400">
              Already have an account?{" "}
              <Link
                to="/login"
                className="text-primary-600 hover:text-primary-700 font-medium"
              >
                Sign in
              </Link>
            </p>
          </form>
        </Card>
      </div>
    </div>
  );
}
