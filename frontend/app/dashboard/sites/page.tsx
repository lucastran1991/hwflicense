'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { apiClient } from '@/lib/api-client';

export default function SitesPage() {
  const [sites, setSites] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [showCreate, setShowCreate] = useState(false);
  const [newSiteId, setNewSiteId] = useState('');
  const [fingerprint, setFingerprint] = useState({ address: '', dns_suffix: '', deployment_tag: '' });
  const router = useRouter();

  useEffect(() => {
    loadSites();
  }, []);

  const loadSites = async () => {
    try {
      const response = await apiClient.listSites();
      setSites(response.data.sites || []);
    } catch (error) {
      console.error('Failed to load sites:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      // Filter out empty fingerprint fields
      const fingerprintData: any = {};
      if (fingerprint.address) fingerprintData.address = fingerprint.address;
      if (fingerprint.dns_suffix) fingerprintData.dns_suffix = fingerprint.dns_suffix;
      if (fingerprint.deployment_tag) fingerprintData.deployment_tag = fingerprint.deployment_tag;
      
      await apiClient.createSite(newSiteId, Object.keys(fingerprintData).length > 0 ? fingerprintData : undefined);
      setShowCreate(false);
      setNewSiteId('');
      setFingerprint({ address: '', dns_suffix: '', deployment_tag: '' });
      loadSites();
    } catch (error: any) {
      alert(error.response?.data?.error || 'Failed to create site');
    }
  };

  const handleDownload = async (siteId: string) => {
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

  const handleDelete = async (siteId: string) => {
    if (!confirm(`Delete site ${siteId}?`)) return;
    try {
      await apiClient.deleteSite(siteId);
      loadSites();
    } catch (error) {
      alert('Failed to delete site');
    }
  };

  if (loading) return <div className="text-center py-8">Loading...</div>;

  return (
    <div className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
      <div className="px-4 py-6 sm:px-0">
        <div className="flex justify-between items-center mb-6">
          <h1 className="text-3xl font-bold text-gray-900">Site Licenses</h1>
          <button
            onClick={() => setShowCreate(true)}
            className="bg-indigo-600 text-white px-4 py-2 rounded-md hover:bg-indigo-700"
          >
            Create Site License
          </button>
        </div>

        {showCreate && (
          <div className="bg-white shadow rounded-lg p-6 mb-6">
            <h2 className="text-xl font-semibold mb-4">Create New Site License</h2>
            <form onSubmit={handleCreate}>
              <div className="mb-4">
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Site ID
                </label>
                <input
                  type="text"
                  value={newSiteId}
                  onChange={(e) => setNewSiteId(e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md"
                  required
                />
              </div>
              
              <div className="mb-4">
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Fingerprint (Optional)
                </label>
                <div className="grid grid-cols-3 gap-4">
                  <div>
                    <label className="block text-xs text-gray-500 mb-1">Address</label>
                    <input
                      type="text"
                      value={fingerprint.address}
                      onChange={(e) => setFingerprint({...fingerprint, address: e.target.value})}
                      placeholder="192.168.1.1"
                      className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm"
                    />
                  </div>
                  <div>
                    <label className="block text-xs text-gray-500 mb-1">DNS Suffix</label>
                    <input
                      type="text"
                      value={fingerprint.dns_suffix}
                      onChange={(e) => setFingerprint({...fingerprint, dns_suffix: e.target.value})}
                      placeholder="hwf.local"
                      className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm"
                    />
                  </div>
                  <div>
                    <label className="block text-xs text-gray-500 mb-1">Deployment Tag</label>
                    <input
                      type="text"
                      value={fingerprint.deployment_tag}
                      onChange={(e) => setFingerprint({...fingerprint, deployment_tag: e.target.value})}
                      placeholder="production"
                      className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm"
                    />
                  </div>
                </div>
              </div>
              
              <div className="flex gap-2">
                <button
                  type="submit"
                  className="bg-indigo-600 text-white px-4 py-2 rounded-md hover:bg-indigo-700"
                >
                  Create
                </button>
                <button
                  type="button"
                  onClick={() => setShowCreate(false)}
                  className="bg-gray-200 text-gray-800 px-4 py-2 rounded-md hover:bg-gray-300"
                >
                  Cancel
                </button>
              </div>
            </form>
          </div>
        )}

        <div className="bg-white shadow rounded-lg overflow-hidden">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Site ID
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Status
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Issued At
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Last Seen
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {sites.length === 0 ? (
                <tr>
                  <td colSpan={5} className="px-6 py-4 text-center text-gray-500">
                    No sites found
                  </td>
                </tr>
              ) : (
                sites.map((site) => (
                  <tr key={site.site_id}>
                    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                      {site.site_id}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
                        site.status === 'active' ? 'bg-green-100 text-green-800' : 
                        site.status === 'revoked' ? 'bg-red-100 text-red-800' : 
                        'bg-gray-100 text-gray-800'
                      }`}>
                        {site.status}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {new Date(site.issued_at).toLocaleString()}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {site.last_seen ? new Date(site.last_seen).toLocaleString() : 'Never'}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium space-x-2">
                      <button
                        onClick={() => router.push(`/dashboard/sites/${site.site_id}`)}
                        className="text-indigo-600 hover:text-indigo-900 mr-2"
                      >
                        View
                      </button>
                      <button
                        onClick={() => handleDownload(site.site_id)}
                        className="text-green-600 hover:text-green-900 mr-2"
                      >
                        Download
                      </button>
                      {site.status !== 'revoked' && (
                        <button
                          onClick={() => handleDelete(site.site_id)}
                          className="text-red-600 hover:text-red-900"
                        >
                          Revoke
                        </button>
                      )}
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}

