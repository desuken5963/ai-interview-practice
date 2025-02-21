'use client';

import { useState } from 'react';
import { XMarkIcon, PlusIcon, PencilIcon, TrashIcon, PlayIcon } from '@heroicons/react/24/outline';

type Job = {
  id: string;
  title: string;
  description: string | null;
  custom_fields: {
    field_name: string;
    content: string;
  }[];
  created_at: string;
  updated_at: string;
};

type JobListModalProps = {
  isOpen: boolean;
  onClose: () => void;
  companyName: string;
  jobs: Job[];
  onAddJob: () => void;
  onEditJob: (job: Job) => void;
  onDeleteJob: (jobId: string) => void;
};

export default function JobListModal({
  isOpen,
  onClose,
  companyName,
  jobs,
  onAddJob,
  onEditJob,
  onDeleteJob,
}: JobListModalProps) {
  const [currentPage, setCurrentPage] = useState(1);
  const itemsPerPage = 5;
  const totalPages = Math.ceil(jobs.length / itemsPerPage);

  // 現在のページの求人を取得
  const currentJobs = jobs.slice(
    (currentPage - 1) * itemsPerPage,
    currentPage * itemsPerPage
  );

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
    <div className="fixed inset-0 z-50 overflow-y-auto">
      <div className="flex min-h-full items-center justify-center p-2 sm:p-4 text-center">
        <div className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" onClick={onClose} />

        <div className="relative transform overflow-hidden rounded-lg bg-white text-left shadow-xl transition-all w-full max-w-4xl mx-2 sm:mx-4">
          <div className="px-4 sm:px-6 py-4 border-b border-gray-200 flex justify-between items-center">
            <h3 className="text-lg sm:text-xl font-semibold text-gray-900 line-clamp-1">
              {companyName}の求人一覧
            </h3>
            <button
              type="button"
              className="text-gray-400 hover:text-gray-500"
              onClick={onClose}
            >
              <XMarkIcon className="h-5 w-5 sm:h-6 sm:w-6" />
            </button>
          </div>

          <div className="p-4 sm:p-6">
            {/* 新規登録ボタン */}
            <div className="mb-4 sm:mb-6">
              <button
                onClick={onAddJob}
                className="inline-flex items-center px-3 sm:px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors text-sm sm:text-base"
              >
                <PlusIcon className="w-4 h-4 sm:w-5 sm:h-5 mr-2" />
                求人を登録
              </button>
            </div>

            {/* 求人リスト */}
            <div className="space-y-4">
              {currentJobs.map((job) => (
                <div
                  key={job.id}
                  className="border border-gray-200 rounded-lg p-3 sm:p-4 hover:shadow-md transition-shadow"
                >
                  <div className="flex justify-between items-start mb-2 gap-4">
                    <h4 className="text-base sm:text-lg font-medium text-gray-900 line-clamp-2">{job.title}</h4>
                    <div className="flex gap-1 sm:gap-2 flex-shrink-0">
                      <button
                        onClick={() => onEditJob(job)}
                        className="p-1 text-gray-500 hover:text-gray-700 transition-colors"
                      >
                        <PencilIcon className="w-4 h-4 sm:w-5 sm:h-5" />
                      </button>
                      <button
                        onClick={() => onDeleteJob(job.id)}
                        className="p-1 text-gray-500 hover:text-red-600 transition-colors"
                      >
                        <TrashIcon className="w-4 h-4 sm:w-5 sm:h-5" />
                      </button>
                    </div>
                  </div>

                  {job.description && (
                    <p className="text-sm sm:text-base text-gray-600 mb-3 line-clamp-2">{job.description}</p>
                  )}

                  {job.custom_fields.length > 0 && (
                    <div className="mb-3">
                      <dl className="grid grid-cols-1 sm:grid-cols-2 gap-2">
                        {job.custom_fields.map((field, index) => (
                          <div key={index} className="col-span-1">
                            <dt className="text-xs sm:text-sm font-medium text-gray-500">
                              {field.field_name}
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
                      登録日: {formatDate(job.created_at)}
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

              {jobs.length === 0 && (
                <div className="text-center py-8">
                  <p className="text-sm sm:text-base text-gray-500">登録されている求人はありません</p>
                </div>
              )}
            </div>

            {/* ページネーション */}
            {totalPages > 1 && (
              <div className="flex justify-center mt-4 sm:mt-6">
                <nav className="flex items-center gap-2">
                  <button
                    onClick={() => setCurrentPage(prev => Math.max(prev - 1, 1))}
                    disabled={currentPage === 1}
                    className="px-2 sm:px-3 py-1 text-sm rounded border border-gray-300 disabled:opacity-50"
                  >
                    前へ
                  </button>
                  <span className="px-3 sm:px-4 py-1 text-sm">
                    {currentPage} / {totalPages}
                  </span>
                  <button
                    onClick={() => setCurrentPage(prev => Math.min(prev + 1, totalPages))}
                    disabled={currentPage === totalPages}
                    className="px-2 sm:px-3 py-1 text-sm rounded border border-gray-300 disabled:opacity-50"
                  >
                    次へ
                  </button>
                </nav>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
} 