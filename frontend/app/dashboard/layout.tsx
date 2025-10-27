'use client';

import { useAuth } from '@/lib/auth-context';
import { useRouter } from 'next/navigation';
import { useEffect } from 'react';

function LogoutButton() {
  const { logout } = useAuth();
  return (
    <button
      onClick={logout}
      className="text-gray-600 hover:text-gray-900 px-4 py-2"
    >
      Logout
    </button>
  );
}

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const { isAuthenticated, loading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!loading && !isAuthenticated) {
      router.push('/login');
    }
  }, [isAuthenticated, loading, router]);

  if (loading) {
    return <div>Loading...</div>;
  }

  if (!isAuthenticated) {
    return null;
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <nav className="bg-white shadow-sm border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16">
            <div className="flex items-center gap-4">
              <h1 className="text-xl font-semibold text-gray-900">TaskMaster License Hub</h1>
              <div className="flex gap-2">
                <a href="/dashboard" className="text-gray-600 hover:text-gray-900 px-3 py-2 text-sm font-medium">Dashboard</a>
                <a href="/dashboard/sites" className="text-gray-600 hover:text-gray-900 px-3 py-2 text-sm font-medium">Sites</a>
                <a href="/dashboard/manifests" className="text-gray-600 hover:text-gray-900 px-3 py-2 text-sm font-medium">Manifests</a>
              </div>
            </div>
            <div className="flex items-center">
              <LogoutButton />
            </div>
          </div>
        </div>
      </nav>
      <main>{children}</main>
    </div>
  );
}

