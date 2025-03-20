'use client';

import { useState } from 'react';
import { XMarkIcon, PlusIcon, PencilIcon, TrashIcon, PlayIcon } from '@heroicons/react/24/outline';
import { Company, JobPosting } from '@/lib/api/types';
import { jobPostingAPI } from '@/lib/api/client';
import dynamic from 'next/dynamic';

const JobPostingFormModal = dynamic(() => import('./JobPostingFormModal'));

type JobPostingListModalProps = {
  isOpen: boolean;
  onClose: () => void;
  company: Company | null;
  onSuccess: () => void;
};

export default function JobPostingListModal({
  isOpen,
  onClose,
  company,
  onSuccess,
}: JobPostingListModalProps) {
  const [currentPage, setCurrentPage] = useState(1);
  const [selectedJobPosting, setSelectedJobPosting] = useState<JobPosting | null>(null);
  const [isJobPostingFormModalOpen, setIsJobPostingFormModalOpen] = useState(false);
  const itemsPerPage = 5;
  const jobPostings = company?.jobPostings || [];
  const totalPages = Math.ceil((jobPostings?.length || 0) / itemsPerPage);

  // 現在のページの求人を取得
  const currentJobPostings = jobPostings?.slice(
    (currentPage - 1) * itemsPerPage,
    currentPage * itemsPerPage
  ) || [];

  const handleAddJobPosting = () => {
    setSelectedJobPosting(null);
    setIsJobPostingFormModalOpen(true);
  };

  const handleEditJobPosting = (jobPosting: JobPosting) => {
    setSelectedJobPosting(jobPosting);
    setIsJobPostingFormModalOpen(true);
  };

  const handleDeleteJobPosting = async (jobPostingId: number) => {
    if (!window.confirm('この求人を削除してもよろしいですか？この操作は取り消せません。')) {
      return;
    }
    try {
      await jobPostingAPI.deleteJobPosting(jobPostingId);
      onSuccess();
    } catch (error) {
      console.error('Error deleting job posting:', error);
    }
  };

  const handleJobPostingFormSuccess = async () => {
    onSuccess();
    setIsJobPostingFormModalOpen(false);
    setSelectedJobPosting(null);
  };

  // 日付フォーマット関数
  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('ja-JP', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
  };

  if (!isOpen) return null;

  return (
    <>
      <div className="fixed inset-0 z-50 overflow-y-auto">
        <div className="flex min-h-full items-center justify-center p-2 sm:p-4 text-center">
          <div className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" onClick={onClose} />

          <div className="relative transform overflow-hidden rounded-lg bg-white text-left shadow-xl transition-all w-full max-w-2xl">
            <div className="bg-white px-4 pb-4 pt-5 sm:p-6">
              <div className="flex items-start justify-between mb-4">
                <h3 className="text-lg sm:text-xl font-semibold text-gray-900">
                  {company?.name || ''}の求人一覧
                </h3>
                <button
                  onClick={onClose}
                  className="text-gray-400 hover:text-gray-500"
                >
                  <XMarkIcon className="h-6 w-6" />
                </button>
              </div>

              <div className="mb-4">
                <button
                  onClick={handleAddJobPosting}
                  className="inline-flex items-center px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
                >
                  <PlusIcon className="w-5 h-5 mr-2" />
                  求人を追加
                </button>
              </div>

              <div className="space-y-4">
                {currentJobPostings.map((jobPosting) => (
                  <div
                    key={jobPosting.id}
                    className="border border-gray-200 rounded-lg p-3 sm:p-4 hover:shadow-md transition-shadow"
                  >
                    <div className="flex justify-between items-start mb-2 gap-4">
                      <h4 className="text-base sm:text-lg font-medium text-gray-900 line-clamp-2">{jobPosting.title}</h4>
                      <div className="flex gap-1 sm:gap-2 flex-shrink-0">
                        <button
                          onClick={() => handleEditJobPosting(jobPosting)}
                          className="p-1 text-gray-500 hover:text-gray-700 transition-colors"
                        >
                          <PencilIcon className="w-4 h-4 sm:w-5 sm:h-5" />
                        </button>
                        <button
                          onClick={() => handleDeleteJobPosting(jobPosting.id)}
                          className="p-1 text-gray-500 hover:text-red-600 transition-colors"
                        >
                          <TrashIcon className="w-4 h-4 sm:w-5 sm:h-5" />
                        </button>
                      </div>
                    </div>

                    {jobPosting.description && (
                      <p className="text-sm sm:text-base text-gray-600 mb-3 line-clamp-2">{jobPosting.description}</p>
                    )}

                    {jobPosting.customFields.length > 0 && (
                      <div className="mb-3">
                        <dl className="grid grid-cols-1 sm:grid-cols-2 gap-2">
                          {jobPosting.customFields.map((field, index) => (
                            <div key={index} className="col-span-1">
                              <dt className="text-xs sm:text-sm font-medium text-gray-500">
                                {field.fieldName}
                              </dt>
                              <dd className="text-xs sm:text-sm text-gray-900">
                                {field.content}
                              </dd>
                            </div>
                          ))}
                        </dl>
                      </div>
                    )}

                    <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-2 sm:gap-0">
                      <div className="text-xs sm:text-sm text-gray-500">
                        登録日: {formatDate(jobPosting.createdAt)}
                      </div>
                      <button
                        className="w-full sm:w-auto inline-flex items-center justify-center px-3 sm:px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 transition-colors text-sm"
                        onClick={() => {/* 面接練習を開始 */}}
                      >
                        <PlayIcon className="w-4 h-4 mr-2" />
                        面接練習
                      </button>
                    </div>
                  </div>
                ))}

                {jobPostings.length === 0 && (
                  <div className="text-center py-8">
                    <p className="text-sm sm:text-base text-gray-500">登録されている求人はありません</p>
                  </div>
                )}
              </div>

              {totalPages > 1 && (
                <div className="flex justify-center items-center gap-2 mt-4">
                  <button
                    onClick={() => setCurrentPage(prev => Math.max(prev - 1, 1))}
                    disabled={currentPage === 1}
                    className="px-3 py-1 text-sm bg-gray-100 rounded-md disabled:opacity-50"
                  >
                    前へ
                  </button>
                  <span className="text-sm text-gray-600">
                    {currentPage} / {totalPages}
                  </span>
                  <button
                    onClick={() => setCurrentPage(prev => Math.min(prev + 1, totalPages))}
                    disabled={currentPage === totalPages}
                    className="px-3 py-1 text-sm bg-gray-100 rounded-md disabled:opacity-50"
                  >
                    次へ
                  </button>
                </div>
              )}
            </div>
          </div>
        </div>
      </div>

      <JobPostingFormModal
        isOpen={isJobPostingFormModalOpen}
        onClose={() => setIsJobPostingFormModalOpen(false)}
        onSuccess={handleJobPostingFormSuccess}
        jobPosting={selectedJobPosting || undefined}
        companyId={company?.id || 0}
        companyName={company?.name || ''}
      />
    </>
  );
} 