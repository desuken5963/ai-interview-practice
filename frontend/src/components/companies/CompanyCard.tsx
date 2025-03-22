'use client';

import { useState } from 'react';
import { PlayIcon, PencilIcon, TrashIcon } from '@heroicons/react/24/outline';
import dynamic from 'next/dynamic';
import { Company, CustomField } from '@/lib/api/types';
import { companyAPI } from '@/lib/api/client';

const CompanyFormModal = dynamic(() => import('./CompanyFormModal'));
const JobPostingListModal = dynamic(() => import('./JobPostingListModal'));

type ErrorMessageProps = {
  message: string;
  onClose: () => void;
};

function ErrorMessage({ message, onClose }: ErrorMessageProps) {
  return (
    <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
      {message}
      <button 
        className="float-right font-bold"
        onClick={onClose}
      >
        ×
      </button>
    </div>
  );
}

function CompanyDescription({ description }: { description: string | null }) {
  if (!description) return null;
  return (
    <p className="text-gray-600 mb-4 line-clamp-3">
      {description}
    </p>
  );
}

function CustomFieldsList({ fields }: { fields: CustomField[] }) {
  if (fields.length === 0) return null;
  return (
    <div className="mb-4">
      <dl className="grid grid-cols-2 gap-2">
        {fields.map((field, index) => (
          <div key={index} className="col-span-1">
            <dt className="text-sm font-medium text-gray-500">
              {field.fieldName}
            </dt>
            <dd className="text-sm text-gray-900">
              {field.content}
            </dd>
          </div>
        ))}
      </dl>
    </div>
  );
}

type CompanyCardProps = {
  company: Company;
  onRefresh: () => void;
};

export default function CompanyCard({ 
  company,
  onRefresh,
}: CompanyCardProps) {
  const [isCompanyFormModalOpen, setIsCompanyFormModalOpen] = useState(false);
  const [isJobPostingListModalOpen, setIsJobPostingListModalOpen] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleModalSuccess = (closeModal: () => void) => {
    closeModal();
    onRefresh();
  };

  const handleDeleteCompany = async () => {
    if (!window.confirm('この企業を削除してもよろしいですか？この操作は取り消せません。')) {
      return;
    }
    try {
      await companyAPI.deleteCompany(company.id);
      onRefresh();
    } catch (error) {
      console.error('Error deleting company:', error);
      setError('企業の削除に失敗しました');
    }
  };

  return (
    <>
      <div className="bg-white rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow flex flex-col min-h-[300px]">
        {error && (
          <ErrorMessage
            message={error}
            onClose={() => setError(null)}
          />
        )}

        <div className="flex-1">
          <div className="flex justify-between items-start mb-4">
            <h2 className="text-xl font-semibold text-gray-900">{company.name}</h2>
            <div className="flex gap-2">
              <button
                className="p-1 text-gray-500 hover:text-gray-700 transition-colors"
                onClick={() => setIsCompanyFormModalOpen(true)}
              >
                <PencilIcon className="w-5 h-5" />
              </button>
              <button
                className="p-1 text-gray-500 hover:text-red-600 transition-colors"
                onClick={handleDeleteCompany}
              >
                <TrashIcon className="w-5 h-5" />
              </button>
            </div>
          </div>

          <CompanyDescription description={company.businessDescription} />
          <CustomFieldsList fields={company.customFields} />
        </div>

        <div className="flex items-center justify-between mt-4 pt-4 border-t border-gray-200">
          <button
            className="inline-flex items-center text-sm text-gray-500 hover:text-blue-600 transition-colors gap-1"
            onClick={() => setIsJobPostingListModalOpen(true)}
          >
            <span>求人一覧</span>
            <span className="font-semibold">({company.jobPostings?.length ?? 0})</span>
          </button>
          <button
            className="inline-flex items-center px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 transition-colors"
            onClick={() => {/* TODO: 面接練習機能の実装 */}}
          >
            <PlayIcon className="w-4 h-4 mr-2" />
            面接練習
          </button>
        </div>
      </div>

      <CompanyFormModal
        isOpen={isCompanyFormModalOpen}
        onClose={() => setIsCompanyFormModalOpen(false)}
        onSuccess={() => handleModalSuccess(() => setIsCompanyFormModalOpen(false))}
        company={company}
      />

      <JobPostingListModal
        isOpen={isJobPostingListModalOpen}
        onClose={() => setIsJobPostingListModalOpen(false)}
        company={company}
        onSuccess={() => handleModalSuccess(() => setIsJobPostingListModalOpen(false))}
      />
    </>
  );
} 