'use client';

import { useState, useEffect } from 'react';
import { PlayIcon, PencilIcon, TrashIcon } from '@heroicons/react/24/outline';
import dynamic from 'next/dynamic';
import { Company, JobPosting, CompanyInput, JobPostingInput } from '@/lib/api/types';
import { jobPostingAPI } from '@/lib/api/client';

// クライアントサイドのみでレンダリングするためにdynamic importを使用
const JobPostingListModal = dynamic(() => import('./JobPostingListModal'), {
  ssr: false,
});

const JobPostingFormModal = dynamic(() => import('./JobPostingFormModal'), {
  ssr: false,
});

const CompanyFormModal = dynamic(() => import('./CompanyFormModal'), {
  ssr: false,
});

type CompanyCardProps = {
  company: Company;
  onEdit?: (companyId: number, data: CompanyInput) => void;
  onDelete?: () => void;
  onRefresh?: (companyId: number) => Promise<void>;
  onJobPostingListOpen?: () => void;
  jobPostings?: JobPosting[];
};

export default function CompanyCard({ 
  company, 
  onEdit, 
  onDelete, 
  onRefresh, 
  onJobPostingListOpen,
  jobPostings = []
}: CompanyCardProps) {
  const [isJobPostingListModalOpen, setIsJobPostingListModalOpen] = useState(false);
  const [isJobPostingFormModalOpen, setIsJobPostingFormModalOpen] = useState(false);
  const [isCompanyFormModalOpen, setIsCompanyFormModalOpen] = useState(false);
  const [selectedJobPosting, setSelectedJobPosting] = useState<JobPosting | undefined>();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [mounted, setMounted] = useState(false);

  // クライアントサイドでのみレンダリングされるようにする
  useEffect(() => {
    setMounted(true);
  }, []);

  // 求人一覧を開く際のハンドラー
  const handleJobPostingListOpen = async () => {
    if (!mounted) return;
    
    try {
      setLoading(true);
      if (onJobPostingListOpen) {
        await onJobPostingListOpen();
      }
      setIsJobPostingListModalOpen(true);
    } catch (error) {
      console.error('Error opening job posting list:', error);
      setError('求人情報の取得に失敗しました');
    } finally {
      setLoading(false);
    }
  };

  // 求人追加ハンドラー
  const handleAddJobPosting = () => {
    setSelectedJobPosting(undefined);
    setIsJobPostingFormModalOpen(true);
    setIsJobPostingListModalOpen(false);
  };

  // 求人編集ハンドラー
  const handleEditJobPosting = (jobPosting: JobPosting) => {
    setSelectedJobPosting(jobPosting);
    setIsJobPostingFormModalOpen(true);
    setIsJobPostingListModalOpen(false);
  };

  // 求人削除ハンドラー
  const handleDeleteJobPosting = async (jobPostingId: number) => {
    if (!mounted) return;
    
    try {
      setLoading(true);
      // 求人削除APIを呼び出す
      await jobPostingAPI.deleteJobPosting(jobPostingId);
      
      // 企業情報を更新
      if (onRefresh) {
        await onRefresh(company.id);
      }
    } catch (error) {
      console.error('Error deleting job posting:', error);
      setError('求人の削除に失敗しました');
    } finally {
      setLoading(false);
    }
  };

  // 求人保存ハンドラー
  const handleSubmitJobPosting = async (data: JobPostingInput) => {
    if (!mounted) return;
    
    try {
      setLoading(true);
      if (selectedJobPosting) {
        // 編集の場合
        await jobPostingAPI.updateJobPosting(selectedJobPosting.id, data);
      } else {
        // 新規登録の場合
        await jobPostingAPI.createJobPosting(company.id, data);
      }
      
      // モーダルを閉じる
      setIsJobPostingFormModalOpen(false);
      
      // 企業情報を更新
      if (onRefresh) {
        await onRefresh(company.id);
      }
    } catch (error) {
      console.error('Error submitting job posting:', error);
      setError('求人情報の保存に失敗しました');
    } finally {
      setLoading(false);
    }
  };

  // 企業情報編集ハンドラー
  const handleEditCompany = () => {
    setIsCompanyFormModalOpen(true);
  };

  // 企業情報保存ハンドラー
  const handleSubmitCompany = async (data: CompanyInput) => {
    if (!mounted) return;
    
    try {
      setLoading(true);
      if (onEdit) {
        // 親コンポーネントの更新処理を呼び出す
        await onEdit(company.id, data);
      }
      setIsCompanyFormModalOpen(false);
    } catch (error) {
      console.error('Error submitting company:', error);
      setError('企業情報の更新に失敗しました');
    } finally {
      setLoading(false);
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

        {company.business_description && (
          <p className="text-gray-600 mb-4 line-clamp-3">
            {company.business_description}
          </p>
        )}

        {company.custom_fields.length > 0 && (
          <div className="mb-4">
            <dl className="grid grid-cols-2 gap-2">
              {company.custom_fields.map((field, index) => (
                <div key={index} className="col-span-1">
                  <dt className="text-sm font-medium text-gray-500">
                    {field.field_name}
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
            onClick={handleJobPostingListOpen}
            disabled={loading}
          >
            {loading ? (
              <span className="inline-block w-4 h-4 border-2 border-gray-300 border-t-blue-600 rounded-full animate-spin mr-2"></span>
            ) : (
              <>
                <span>求人一覧</span>
                <span className="font-semibold">({company.jobPostings?.length || jobPostings.length || 0})</span>
              </>
            )}
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

      <JobPostingListModal
        isOpen={isJobPostingListModalOpen}
        onClose={() => setIsJobPostingListModalOpen(false)}
        companyName={company.name}
        jobPostings={jobPostings}
        onAddJobPosting={handleAddJobPosting}
        onEditJobPosting={handleEditJobPosting}
        onDeleteJobPosting={handleDeleteJobPosting}
      />

      <JobPostingFormModal
        isOpen={isJobPostingFormModalOpen}
        onClose={() => setIsJobPostingFormModalOpen(false)}
        onSubmit={handleSubmitJobPosting}
        initialData={selectedJobPosting}
        companyName={company.name}
      />

      <CompanyFormModal
        isOpen={isCompanyFormModalOpen}
        onClose={() => setIsCompanyFormModalOpen(false)}
        onSubmit={handleSubmitCompany}
        initialData={{
          name: company.name,
          business_description: company.business_description,
          custom_fields: company.custom_fields,
        }}
      />
    </>
  );
} 