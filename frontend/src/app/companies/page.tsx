'use client';

import { useEffect, useState } from 'react';
import { PlusIcon } from '@heroicons/react/24/outline';
import { useRouter } from 'next/navigation';
import dynamic from 'next/dynamic';
import { Company, CompanyInput } from '@/lib/api/types';
import { companyAPI } from '@/lib/api/client';

// クライアントサイドのみでレンダリングするためにdynamic importを使用
const CompanyCard = dynamic(() => import('@/components/companies/CompanyCard'));
const CompanyFormModal = dynamic(() => import('@/components/companies/CompanyFormModal'));
const JobPostingListModal = dynamic(() => import('@/components/companies/JobPostingListModal'));

// Company型からCompanyInput型への変換関数
const convertToCompanyInput = (company: Company | null): CompanyInput | undefined => {
  if (!company) return undefined;
  return {
    name: company.name,
    businessDescription: company.businessDescription,
    customFields: company.customFields.map(field => ({
      fieldName: field.fieldName,
      content: field.content,
    })),
  };
};

export default function CompaniesPage() {
  const router = useRouter();
  const [companies, setCompanies] = useState<Company[]>([]);
  const [isCompanyFormModalOpen, setIsCompanyFormModalOpen] = useState(false);
  const [isJobPostingListModalOpen, setIsJobPostingListModalOpen] = useState(false);
  const [selectedCompany, setSelectedCompany] = useState<Company | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    fetchCompanies();
  }, []);

  const fetchCompanies = async () => {
    try {
      const response = await companyAPI.getCompanies();
      setCompanies(response.companies);
    } catch (error) {
      console.error('Failed to fetch companies:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleAddCompany = () => {
    setSelectedCompany(null);
    setIsCompanyFormModalOpen(true);
  };

  const handleEditCompany = (company: Company) => {
    setSelectedCompany(company);
    setIsCompanyFormModalOpen(true);
  };

  const handleDeleteCompany = async (company: Company) => {
    if (!window.confirm('この企業を削除してもよろしいですか？この操作は取り消せません。')) {
      return;
    }

    try {
      await companyAPI.deleteCompany(company.id);
      await fetchCompanies();
    } catch (error) {
      console.error('Error deleting company:', error);
    }
  };

  const handleJobListClick = (company: Company) => {
    setSelectedCompany(company);
    setIsJobPostingListModalOpen(true);
  };

  const handleCompanyFormSubmit = async () => {
    await fetchCompanies();
    setIsCompanyFormModalOpen(false);
  };

  if (isLoading) {
    return (
      <div className="flex justify-center items-center min-h-screen">
        <div className="animate-spin rounded-full h-8 w-8 border-t-2 border-b-2 border-blue-500"></div>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold text-gray-900">企業一覧</h1>
        <button
          onClick={handleAddCompany}
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
              onEdit={() => handleEditCompany(company)}
              onDelete={() => handleDeleteCompany(company)}
              onJobListClick={() => handleJobListClick(company)}
            />
          ))}
        </div>
      )}

      <CompanyFormModal
        isOpen={isCompanyFormModalOpen}
        onClose={() => setIsCompanyFormModalOpen(false)}
        onSubmit={handleCompanyFormSubmit}
        initialData={convertToCompanyInput(selectedCompany)}
      />

      <JobPostingListModal
        isOpen={isJobPostingListModalOpen}
        onClose={() => setIsJobPostingListModalOpen(false)}
        companyName={selectedCompany?.name || ''}
        jobPostings={selectedCompany?.jobPostings || []}
        onAddJobPosting={() => {/* 必要に応じて実装 */}}
        onEditJobPosting={() => {/* 必要に応じて実装 */}}
        onDeleteJobPosting={() => {/* 必要に応じて実装 */}}
      />
    </div>
  );
} 