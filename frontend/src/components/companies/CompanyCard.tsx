'use client';

import { useState, useEffect } from 'react';
import { PlayIcon, PencilIcon, TrashIcon } from '@heroicons/react/24/outline';
import dynamic from 'next/dynamic';
import { Company, Job, CompanyInput, JobInput } from '@/lib/api/types';
import { jobAPI } from '@/lib/api/client';

// クライアントサイドのみでレンダリングするためにdynamic importを使用
const JobListModal = dynamic(() => import('./JobListModal'), {
  ssr: false,
});

const JobFormModal = dynamic(() => import('./JobFormModal'), {
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
};

export default function CompanyCard({ company, onEdit, onDelete, onRefresh }: CompanyCardProps) {
  const [isJobListModalOpen, setIsJobListModalOpen] = useState(false);
  const [isJobFormModalOpen, setIsJobFormModalOpen] = useState(false);
  const [isCompanyFormModalOpen, setIsCompanyFormModalOpen] = useState(false);
  const [selectedJob, setSelectedJob] = useState<Job | undefined>();
  const [jobs, setJobs] = useState<Job[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [mounted, setMounted] = useState(false);

  // クライアントサイドでのみレンダリングされるようにする
  useEffect(() => {
    setMounted(true);
  }, []);

  // 求人データを取得する関数
  const fetchJobs = async () => {
    if (!mounted) return;
    
    try {
      setLoading(true);
      // APIから求人一覧を取得
      const response = await jobAPI.getJobs(company.id);
      setJobs(response.jobs);
    } catch (error) {
      console.error('Error fetching jobs:', error);
      setError('求人情報の取得に失敗しました');
    } finally {
      setLoading(false);
    }
  };

  // 求人一覧を開く際のハンドラー
  const handleJobListOpen = async () => {
    if (!mounted) return;
    
    try {
      setLoading(true);
      await fetchJobs();
      setIsJobListModalOpen(true);
    } catch (error) {
      console.error('Error fetching jobs:', error);
      setError('求人情報の取得に失敗しました');
    } finally {
      setLoading(false);
    }
  };

  // 求人追加ハンドラー
  const handleAddJob = () => {
    setSelectedJob(undefined);
    setIsJobFormModalOpen(true);
    setIsJobListModalOpen(false);
  };

  // 求人編集ハンドラー
  const handleEditJob = (job: Job) => {
    setSelectedJob(job);
    setIsJobFormModalOpen(true);
    setIsJobListModalOpen(false);
  };

  // 求人削除ハンドラー
  const handleDeleteJob = async (jobId: number) => {
    if (!mounted) return;
    
    try {
      setLoading(true);
      // 求人削除APIを呼び出す
      await jobAPI.deleteJob(company.id, jobId);
      
      // 最新の求人データを取得
      await fetchJobs();
      
      // 企業情報も更新（求人数などが変わるため）
      if (onRefresh) {
        await onRefresh(company.id);
      }
    } catch (error) {
      console.error('Error deleting job:', error);
      setError('求人の削除に失敗しました');
    } finally {
      setLoading(false);
    }
  };

  // 求人保存ハンドラー
  const handleSubmitJob = async (data: JobInput) => {
    if (!mounted) return;
    
    try {
      setLoading(true);
      if (selectedJob) {
        // 編集の場合
        await jobAPI.updateJob(company.id, selectedJob.id, data);
      } else {
        // 新規登録の場合
        await jobAPI.createJob(company.id, data);
      }
      
      // モーダルを閉じる
      setIsJobFormModalOpen(false);
      
      // 最新の求人データを取得
      await fetchJobs();
      
      // 企業情報も更新（求人数などが変わるため）
      if (onRefresh) {
        await onRefresh(company.id);
      }
    } catch (error) {
      console.error('Error submitting job:', error);
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
            onClick={handleJobListOpen}
            disabled={loading}
          >
            {loading ? (
              <span className="inline-block w-4 h-4 border-2 border-gray-300 border-t-blue-600 rounded-full animate-spin mr-2"></span>
            ) : (
              <>
                <span>求人一覧</span>
                <span className="font-semibold">({company.job_count || 0})</span>
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

      <JobListModal
        isOpen={isJobListModalOpen}
        onClose={() => setIsJobListModalOpen(false)}
        companyName={company.name}
        jobs={jobs}
        onAddJob={handleAddJob}
        onEditJob={handleEditJob}
        onDeleteJob={handleDeleteJob}
      />

      <JobFormModal
        isOpen={isJobFormModalOpen}
        onClose={() => setIsJobFormModalOpen(false)}
        onSubmit={handleSubmitJob}
        initialData={selectedJob}
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