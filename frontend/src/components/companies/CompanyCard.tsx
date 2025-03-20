'use client';

import { useState } from 'react';
import { PlayIcon, PencilIcon, TrashIcon } from '@heroicons/react/24/outline';
import dynamic from 'next/dynamic';
import { Company, CompanyInput, CustomField } from '@/lib/api/types';

// クライアントサイドのみでレンダリングするためにdynamic importを使用
const CompanyFormModal = dynamic(() => import('./CompanyFormModal'), {
  ssr: false,
});

type CompanyCardProps = {
  company: Company;
  onEdit?: (companyId: number, data: CompanyInput) => void;
  onDelete?: () => void;
  onJobListClick?: (company: Company) => void;
};

export default function CompanyCard({ 
  company, 
  onEdit, 
  onDelete,
  onJobListClick,
}: CompanyCardProps) {
  const [isCompanyFormModalOpen, setIsCompanyFormModalOpen] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // 企業情報編集ハンドラー
  const handleEditCompany = () => {
    setIsCompanyFormModalOpen(true);
  };

  // 企業情報保存ハンドラー
  const handleSubmitCompany = async (data: CompanyInput) => {
    try {
      if (onEdit) {
        await onEdit(company.id, data);
      }
      setIsCompanyFormModalOpen(false);
    } catch (error) {
      console.error('Error submitting company:', error);
      setError('企業情報の更新に失敗しました');
    }
  };

  return (
    <>
      <div className="bg-white rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow">
        {/* エラーメッセージ */}
        {error && (
          <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
            {error}
            <button 
              className="float-right font-bold"
              onClick={() => setError(null)}
            >
              ×
            </button>
          </div>
        )}

        <div className="flex justify-between items-start mb-4">
          <h2 className="text-xl font-semibold text-gray-900">{company.name}</h2>
          <div className="flex gap-2">
            {onEdit && (
              <button
                className="p-1 text-gray-500 hover:text-gray-700 transition-colors"
                onClick={handleEditCompany}
              >
                <PencilIcon className="w-5 h-5" />
              </button>
            )}
            {onDelete && (
              <button
                className="p-1 text-gray-500 hover:text-red-600 transition-colors"
                onClick={onDelete}
              >
                <TrashIcon className="w-5 h-5" />
              </button>
            )}
          </div>
        </div>

        {company.businessDescription && (
          <p className="text-gray-600 mb-4 line-clamp-3">
            {company.businessDescription}
          </p>
        )}

        {company.customFields.length > 0 && (
          <div className="mb-4">
            <dl className="grid grid-cols-2 gap-2">
              {company.customFields.map((field: CustomField, index: number) => (
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
        )}

        <div className="flex items-center justify-between mt-4">
          <button
            className="inline-flex items-center text-sm text-gray-500 hover:text-blue-600 transition-colors gap-1"
            onClick={() => onJobListClick?.(company)}
          >
            <span>求人一覧</span>
            <span className="font-semibold">({company.jobPostings?.length || 0})</span>
          </button>
          <button
            className="inline-flex items-center px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 transition-colors"
            onClick={() => {/* 面接練習を開始 */}}
          >
            <PlayIcon className="w-4 h-4 mr-2" />
            面接練習
          </button>
        </div>
      </div>

      <CompanyFormModal
        isOpen={isCompanyFormModalOpen}
        onClose={() => setIsCompanyFormModalOpen(false)}
        onSubmit={handleSubmitCompany}
        initialData={{
          name: company.name,
          businessDescription: company.businessDescription,
          customFields: company.customFields,
        }}
      />
    </>
  );
} 