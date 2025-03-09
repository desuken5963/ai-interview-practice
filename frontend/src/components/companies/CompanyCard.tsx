'use client';

import { useState, useEffect } from 'react';
import { PlayIcon, PencilIcon, TrashIcon } from '@heroicons/react/24/outline';
import JobListModal from './JobListModal';
import JobFormModal from './JobFormModal';
import CompanyFormModal from './CompanyFormModal';
import { Company, Job, CompanyInput, JobInput } from '@/lib/api/types';
import { jobAPI } from '@/lib/api/client';

type CompanyCardProps = {
  company: Company;
  onEdit?: () => void;
  onDelete?: () => void;
};

export default function CompanyCard({ company, onEdit, onDelete }: CompanyCardProps) {
  const [isJobListModalOpen, setIsJobListModalOpen] = useState(false);
  const [isJobFormModalOpen, setIsJobFormModalOpen] = useState(false);
  const [isCompanyFormModalOpen, setIsCompanyFormModalOpen] = useState(false);
  const [selectedJob, setSelectedJob] = useState<Job | undefined>();
  const [jobs, setJobs] = useState<Job[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // 求人一覧を開く際のハンドラー
  const handleJobListOpen = async () => {
    try {
      setLoading(true);
      // APIから求人一覧を取得
      const data = await jobAPI.getJobs(company.id);
      setJobs(data);
      setIsJobListModalOpen(true);
      
      /* 
      // 仮のモックデータを使用（実際のAPIができたら削除）
      const mockJobs: Job[] = [
        {
          id: 1,
          company_id: company.id,
          title: 'フロントエンドエンジニア',
          description: 'モダンなWebアプリケーション開発のためのフロントエンドエンジニアを募集しています。React、TypeScript、Next.jsなどの技術スタックを使用した開発経験がある方を歓迎します。',
          requirements: null,
          custom_fields: [
            { field_name: '雇用形態', content: '正社員' },
            { field_name: '給与', content: '年収450万円〜800万円' },
            { field_name: '勤務地', content: '東京都渋谷区' },
            { field_name: '必要なスキル', content: 'React, TypeScript, Next.js' }
          ],
          created_at: '2024-03-15T09:00:00Z',
          updated_at: '2024-03-15T09:00:00Z'
        },
        {
          id: 2,
          company_id: company.id,
          title: 'バックエンドエンジニア',
          description: 'スケーラブルなバックエンドシステムの設計・開発を担当していただきます。マイクロサービスアーキテクチャの知識と実践経験がある方を求めています。',
          requirements: null,
          custom_fields: [
            { field_name: '雇用形態', content: '正社員' },
            { field_name: '給与', content: '年収500万円〜900万円' },
            { field_name: '勤務地', content: '東京都港区' },
            { field_name: '必要なスキル', content: 'Go, gRPC, Kubernetes' }
          ],
          created_at: '2024-03-14T10:30:00Z',
          updated_at: '2024-03-14T10:30:00Z'
        }
      ];
      setJobs(mockJobs);
      */
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
    try {
      setLoading(true);
      // 求人削除APIを呼び出す
      await jobAPI.deleteJob(jobId);
      setJobs(prev => prev.filter(job => job.id !== jobId));
    } catch (error) {
      console.error('Error deleting job:', error);
      setError('求人の削除に失敗しました');
    } finally {
      setLoading(false);
    }
  };

  // 求人保存ハンドラー
  const handleSubmitJob = async (data: JobInput) => {
    try {
      setLoading(true);
      if (selectedJob) {
        // 編集の場合
        const updatedJob = await jobAPI.updateJob(selectedJob.id, data);
        setJobs(prev => prev.map(job => 
          job.id === selectedJob.id ? updatedJob : job
        ));
      } else {
        // 新規登録の場合
        const newJob = await jobAPI.createJob(company.id, data);
        setJobs(prev => [...prev, newJob]);
      }
      setIsJobFormModalOpen(false);
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
    try {
      // TODO: APIを呼び出して企業情報を保存
      console.log('Submit company data:', data);
      if (onEdit) {
        onEdit();
      }
    } catch (error) {
      console.error('Error submitting company:', error);
      throw error;
    }
  };

  return (
    <>
      <div className="bg-white rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow">
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
          >
            <span>求人一覧</span>
            <span className="font-semibold">({company.job_count || 0})</span>
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