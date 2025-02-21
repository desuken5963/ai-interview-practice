'use client';

import { useState } from 'react';
import { PlayIcon, PencilIcon, TrashIcon } from '@heroicons/react/24/outline';
import JobListModal from './JobListModal';
import JobFormModal from './JobFormModal';
import CompanyFormModal from './CompanyFormModal';

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

type CompanyCardProps = {
  company: {
    id: string;
    name: string;
    business_description: string | null;
    custom_fields: {
      field_name: string;
      content: string;
    }[];
    job_count: number;
    created_at: string;
    updated_at: string;
  };
  onEdit?: () => void;
  onDelete?: () => void;
};

export default function CompanyCard({ company, onEdit, onDelete }: CompanyCardProps) {
  const [isJobListModalOpen, setIsJobListModalOpen] = useState(false);
  const [isJobFormModalOpen, setIsJobFormModalOpen] = useState(false);
  const [isCompanyFormModalOpen, setIsCompanyFormModalOpen] = useState(false);
  const [selectedJob, setSelectedJob] = useState<Job | undefined>();
  const [jobs, setJobs] = useState<Job[]>([]); // 実際のAPIができたら削除

  // 求人一覧を開く際のハンドラー
  const handleJobListOpen = async () => {
    try {
      // TODO: APIから求人一覧を取得
      // 仮のモックデータを使用（実際のAPIができたら削除）
      const mockJobs: Job[] = [
        {
          id: '1',
          title: 'フロントエンドエンジニア',
          description: 'モダンなWebアプリケーション開発のためのフロントエンドエンジニアを募集しています。React、TypeScript、Next.jsなどの技術スタックを使用した開発経験がある方を歓迎します。',
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
          id: '2',
          title: 'バックエンドエンジニア',
          description: 'スケーラブルなバックエンドシステムの設計・開発を担当していただきます。マイクロサービスアーキテクチャの知識と実践経験がある方を求めています。',
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
      setIsJobListModalOpen(true);
    } catch (error) {
      console.error('Error fetching jobs:', error);
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
  const handleDeleteJob = async (jobId: string) => {
    try {
      // TODO: 求人削除APIを呼び出す
      console.log('Delete job:', jobId);
      // 仮実装：モックデータから削除
      setJobs(prev => prev.filter(job => job.id !== jobId));
    } catch (error) {
      console.error('Error deleting job:', error);
    }
  };

  // 求人保存ハンドラー
  const handleSubmitJob = async (data: {
    title: string;
    description: string | null;
    custom_fields: { field_name: string; content: string; }[];
  }) => {
    try {
      // TODO: APIを呼び出して求人情報を保存
      console.log('Submit job data:', data);
      
      if (selectedJob) {
        // 編集の場合
        setJobs(prev => prev.map(job => 
          job.id === selectedJob.id
            ? { ...job, ...data, updated_at: new Date().toISOString() }
            : job
        ));
      } else {
        // 新規登録の場合
        const newJob: Job = {
          id: String(Date.now()), // 一時的なID
          ...data,
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
        };
        setJobs(prev => [...prev, newJob]);
      }
    } catch (error) {
      console.error('Error submitting job:', error);
      throw error;
    }
  };

  // 企業情報編集ハンドラー
  const handleEditCompany = () => {
    setIsCompanyFormModalOpen(true);
  };

  // 企業情報保存ハンドラー
  const handleSubmitCompany = async (data: {
    name: string;
    business_description: string | null;
    custom_fields: { field_name: string; content: string; }[];
  }) => {
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
            <span className="font-semibold">({company.job_count})</span>
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