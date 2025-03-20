'use client';

import { useState, useEffect } from 'react';
import { PlusIcon } from '@heroicons/react/24/outline';
import dynamic from 'next/dynamic';
import { Company } from '@/lib/api/types';
import { companyAPI } from '@/lib/api/client';

// クライアントサイドのみでレンダリングするためにdynamic importを使用
const CompanyCard = dynamic(() => import('@/components/companies/CompanyCard'));
const CompanyFormModal = dynamic(() => import('@/components/companies/CompanyFormModal'));

export default function CompaniesPage() {
  const [companies, setCompanies] = useState<Company[]>([]);
  const [isCompanyFormModalOpen, setIsCompanyFormModalOpen] = useState(false);

  const fetchCompanies = async () => {
    try {
      const response = await companyAPI.getCompanies();
      setCompanies(response.companies);
    } catch (error) {
      console.error('Error fetching companies:', error);
    }
  };

  useEffect(() => {
    fetchCompanies();
  }, []);

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold text-gray-900">企業一覧</h1>
        <button
          onClick={() => setIsCompanyFormModalOpen(true)}
          className="inline-flex items-center px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
        >
          <PlusIcon className="w-5 h-5 mr-2" />
          企業を追加
        </button>
      </div>

      {companies.length === 0 ? (
        <div className="text-center py-8">
          <p className="text-gray-500">登録されている企業はありません</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {companies.map((company) => (
            <CompanyCard
              key={company.id}
              company={company}
              onRefresh={fetchCompanies}
            />
          ))}
        </div>
      )}

      <CompanyFormModal
        isOpen={isCompanyFormModalOpen}
        onClose={() => setIsCompanyFormModalOpen(false)}
        onSuccess={() => {
          setIsCompanyFormModalOpen(false);
          fetchCompanies();
        }}
      />
    </div>
  );
} 