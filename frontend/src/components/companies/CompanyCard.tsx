'use client';

import { useState } from 'react';
import { PlayIcon, PencilIcon, TrashIcon } from '@heroicons/react/24/outline';
import dynamic from 'next/dynamic';
import { Company, CompanyInput, CustomField } from '@/lib/api/types';
import { companyAPI } from '@/lib/api/client';

const CompanyFormModal = dynamic(() => import('./CompanyFormModal'));
const JobPostingListModal = dynamic(() => import('./JobPostingListModal'));

type IconButtonProps = {
  icon: React.ReactNode;
  onClick: () => void;
  className?: string;
  label?: string;
  children?: React.ReactNode;
};

function IconButton({ icon, onClick, className = '', label, children }: IconButtonProps) {
  return (
    <button
      className={`inline-flex items-center ${className}`}
      onClick={onClick}
      aria-label={label}
    >
      {icon}
      {children}
    </button>
  );
}

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

  // TODO: 面接練習機能の実装
  const handleStartInterview = () => {
    console.log('面接練習機能は未実装です');
  };

  return (
    <>
      <div className="bg-white rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow">
        {error && (
          <ErrorMessage
            message={error}
            onClose={() => setError(null)}
          />
        )}

        <div className="flex justify-between items-start mb-4">
          <h2 className="text-xl font-semibold text-gray-900">{company.name}</h2>
          <div className="flex gap-2">
            <IconButton
              icon={<PencilIcon className="w-5 h-5" />}
              onClick={() => setIsCompanyFormModalOpen(true)}
              className="p-1 text-gray-500 hover:text-gray-700 transition-colors"
              label="企業情報を編集"
            />
            <IconButton
              icon={<TrashIcon className="w-5 h-5" />}
              onClick={handleDeleteCompany}
              className="p-1 text-gray-500 hover:text-red-600 transition-colors"
              label="企業を削除"
            />
          </div>
        </div>

        <CompanyDescription description={company.businessDescription} />
        <CustomFieldsList fields={company.customFields} />

        <div className="flex items-center justify-between mt-4">
          <IconButton
            icon={null}
            onClick={() => setIsJobPostingListModalOpen(true)}
            className="inline-flex items-center text-sm text-gray-500 hover:text-blue-600 transition-colors gap-1"
            label="求人一覧を表示"
          >
            <span>求人一覧</span>
            <span className="font-semibold">({company.jobPostings?.length ?? 0})</span>
          </IconButton>
          <IconButton
            icon={<PlayIcon className="w-4 h-4 mr-2" />}
            onClick={handleStartInterview}
            className="inline-flex items-center px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 transition-colors"
            label="面接練習を開始"
          >
            面接練習
          </IconButton>
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