import { Link, useLocation } from "@tanstack/react-router";
import { useAuth } from "@/contexts/AuthContext";
import ThemeToggle from "@/components/common/ThemeToggle";
import GopherLogo from "@/components/common/GopherLogo";

export default function Navbar() {
  const { user, logout } = useAuth();
  const location = useLocation();

  function handleLogout() {
    logout();
    window.location.href = "/login";
  }

  const navLinks = [
    { to: "/dashboard" as const, label: "Dashboard" },
    { to: "/transfer" as const, label: "Transfer" },
    { to: "/statement" as const, label: "Statement" },
  ];

  return (
    <nav className="bg-white dark:bg-gray-800 shadow-sm border-b border-gray-200 dark:border-gray-700">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between h-16">
          <div className="flex items-center">
            <Link to="/dashboard" className="flex items-center space-x-2">
              <GopherLogo className="h-8 w-8" />
              <span className="text-2xl font-bold text-primary-600">Mini</span>
              <span className="text-2xl font-light text-gray-600 dark:text-gray-400">
                Bank
              </span>
            </Link>
            <div className="hidden sm:ml-10 sm:flex sm:space-x-4">
              {navLinks.map((link) => (
                <Link
                  key={link.to}
                  to={link.to}
                  className={`px-3 py-2 rounded-md text-sm font-medium transition-colors ${
                    location.pathname === link.to
                      ? "text-primary-600 bg-primary-50 dark:bg-primary-900/30"
                      : "text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-100 hover:bg-gray-100 dark:hover:bg-gray-700"
                  }`}
                >
                  {link.label}
                </Link>
              ))}
            </div>
          </div>
          <div className="flex items-center space-x-4">
            <span className="text-sm text-gray-600 dark:text-gray-400">
              {user?.full_name}
            </span>
            <ThemeToggle />
            <button
              onClick={handleLogout}
              className="px-3 py-2 rounded-md text-sm font-medium text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-100 hover:bg-gray-100 dark:hover:bg-gray-700"
            >
              Logout
            </button>
          </div>
        </div>
      </div>
    </nav>
  );
}
