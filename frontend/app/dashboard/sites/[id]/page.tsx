'use client';

import { useState, useEffect } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { apiClient } from '@/lib/api-client';

export default function SiteDetailPage() {
  const params = useParams();
  const router = useRouter();
  const siteId = params.id as string;
  
  const [site, setSite] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [licenseData, setLicenseData] = useState<any>(null);
  const [fingerprint, setFingerprint] = useState<any>(null);

  useEffect(() => {
    loadSiteDetails();
  }, [siteId]);

  const loadSiteDetails = async () => {
    try {
      const response = await apiClient.getSite(siteId);
      const siteData = response.data.license;
      setSite(siteData);
      
      // Parse license data and fingerprint
      if (siteData.license_data) {
        try {
          const parsed = JSON.parse(siteData.license_data);
          setLicenseData(parsed);
        } catch (e) {
          console.error('Failed to parse license_data:', e);
        }
      }
      
      if (siteData.fingerprint) {
        try {
          const parsed = JSON.parse(siteData.fingerprint);
          setFingerprint(parsed);
        } catch (e) {
          setFingerprint({});
        }
      }
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to load site details');
    } finally {
      setLoading(false);
    }
  };

  const handleDownload = async () => {
    try {
      const response = await apiClient.getSite(siteId);
      const licenseDataStr = response.data.license?.license_data;
      
      if (!licenseDataStr) {
        alert('No license data available');
        return;
      }

      // Create downloadable JSON file
      const blob = new Blob([licenseDataStr], { type: 'application/json' });
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `site_${siteId}.lic`;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      window.URL.revokeObjectURL(url);
    } catch (err) {
      alert('Failed to download license file');
    }
  };

  const handleRevoke = async () => {
    if (!confirm(`Are you sure you want to revoke site ${siteId}?`)) return;
    
    try {
      await apiClient.deleteSite(siteId);
      router.push('/dashboard/sites');
    } catch (err: any) {
      alert(err.response?.data?.error || 'Failed to revoke site');
    }
  };

  if (loading) return <div className="flex items-center justify-center h-96">Loading...</div>;
  if (error) return <div className="text-center py-8 text-red-600">{error}</div>;
  if (!site) return <div className="text-center py-8">Site not found</div>;

  return (
    <div className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
      <div className="px-4 py-6 sm:px-0">
        {/* Header */}
        <div className="flex justify-between items-center mb-6">
          <div>
            <button
              onClick={() => router.push('/dashboard/sites')}
              className="text-indigo-600 hover:text-indigo-900 mb-2"
            >
              ← Back to Sites
            </button>
            <h1 className="text-3xl font-bold text-gray-900">Site Details</h1>
            <p className="text-gray-600 mt-1">{siteId}</p>
          </div>
          <div className="flex gap-2">
            <button
              onClick={handleDownload}
              className="bg-indigo-600 text-white px-4 py-2 rounded-md hover:bg-indigo-700"
            >
              Download License
            </button>
            {site.status !== 'revoked' && (
              <button
                onClick={handleRevoke}
                className="bg-red-600 text-white px-4 py-2 rounded-md hover:bg-red-700"
              >
                Revoke Site
              </button>
            )}
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
          {/* Basic Info Card */}
          <div className="bg-white shadow rounded-lg p-6">
            <h2 className="text-xl font-semibold mb-4">Basic Information</h2>
            <dl className="grid grid-cols-1 gap-4">
              <div>
                <dt className="text-sm font-medium text-gray-500">Site ID</dt>
                <dd className="mt-1 text-sm text-gray-900">{site.site_id}</dd>
              </div>
              <div>
                <dt className="text-sm font-medium text-gray-500">Status</dt>
                <dd className="mt-1">
                  <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
                    site.status === 'active' ? 'bg-green-100 text-green-800' : 
                    site.status === 'revoked' ? 'bg-red-100 text-red-800' : 
                    'bg-gray-100 text-gray-800'
                  }`}>
                    {site.status}
                  </span>
                </dd>
              </div>
              <div>
                <dt className="text-sm font-medium text-gray-500">Issued At</dt>
                <dd className="mt-1 text-sm text-gray-900">{new Date(site.issued_at).toLocaleString()}</dd>
              </div>
              <div>
                <dt className="text-sm font-medium text-gray-500">Last Seen</dt>
                <dd className="mt-1 text-sm text-gray-900">
                  {site.last_seen ? new Date(site.last_seen).toLocaleString() : 'Never'}
                </dd>
              </div>
            </dl>
          </div>

          {/* Fingerprint Info Card */}
          <div className="bg-white shadow rounded-lg p-6">
            <h2 className="text-xl font-semibold mb-4">Fingerprint</h2>
            {fingerprint && Object.keys(fingerprint).length > 0 ? (
              <dl className="grid grid-cols-1 gap-4">
                {Object.entries(fingerprint).map(([key, value]) => (
                  <div key={key}>
                    <dt className="text-sm font-medium text-gray-500 capitalize">{key.replace('_', ' ')}</dt>
                    <dd className="mt-1 text-sm text-gray-900">{String(value)}</dd>
                  </div>
                ))}
              </dl>
            ) : (
              <p className="text-sm text-gray-500">No fingerprint information available</p>
            )}
          </div>
        </div>

        {/* License Data Card */}
        <div className="bg-white shadow rounded-lg p-6 mb-6">
          <h2 className="text-xl font-semibold mb-4">License Data</h2>
          {licenseData ? (
            <div className="bg-gray-50 rounded-lg p-4 overflow-auto">
              <pre className="text-sm text-gray-900 whitespace-pre-wrap">
                {JSON.stringify(licenseData, null, 2)}
              </pre>
            </div>
          ) : (
            <p className="text-sm text-gray-500">No license data available</p>
          )}
        </div>

        {/* Signature Card */}
        <div className="bg-white shadow rounded-lg p-6">
          <h2 className="text-xl font-semibold mb-4">Signature</h2>
          <div className="bg-gray-50 rounded-lg p-4">
            <p className="text-sm text-gray-900 font-mono break-all">{site.signature}</p>
            {site.signature && !site.signature.startsWith('TODO:') ? (
              <p className="text-sm text-green-600 mt-2">✓ Valid ECDSA signature</p>
            ) : (
              <p className="text-sm text-yellow-600 mt-2">⚠ Placeholder signature</p>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

