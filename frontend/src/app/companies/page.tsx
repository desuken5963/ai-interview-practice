'use client';

import { useState } from 'react';
import { PlusIcon } from '@heroicons/react/24/outline';
import CompanyCard from '@/components/companies/CompanyCard';

type Company = {
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

export default function CompaniesPage() {
  const mockCompanies: Company[] = [
    {
      id: '1',
      name: '株式会社テックイノベーション',
      business_description: 'AIと機械学習を活用した革新的なソリューションを提供する企業です。クラウドサービス、データ分析、自然言語処理など、最先端技術を駆使したサービスを展開しています。',
      custom_fields: [
        { field_name: '業界', content: 'IT・テクノロジー' },
        { field_name: '従業員数', content: '150名' }
      ],
      job_count: 5,
      created_at: '2024-03-15T09:00:00Z',
      updated_at: '2024-03-15T09:00:00Z'
    },
    {
      id: '2',
      name: 'グローバルコンサルティング株式会社',
      business_description: '世界各国の企業に対して、経営戦略、デジタルトランスフォーメーション、組織改革などのコンサルティングサービスを提供しています。',
      custom_fields: [
        { field_name: '業界', content: 'コンサルティング' },
        { field_name: '従業員数', content: '300名' }
      ],
      job_count: 3,
      created_at: '2024-03-14T10:30:00Z',
      updated_at: '2024-03-14T10:30:00Z'
    },
    {
      id: '3',
      name: '未来フィンテック株式会社',
      business_description: 'ブロックチェーン技術を活用した次世代の金融サービスを開発。個人向けおよび法人向けの革新的な決済ソリューションを提供しています。',
      custom_fields: [
        { field_name: '業界', content: 'フィンテック' },
        { field_name: '従業員数', content: '80名' }
      ],
      job_count: 2,
      created_at: '2024-03-13T15:45:00Z',
      updated_at: '2024-03-13T15:45:00Z'
    },
    {
      id: '4',
      name: 'エコテクノロジー株式会社',
      business_description: '再生可能エネルギーとスマートグリッド技術を組み合わせた環境配慮型のエネルギーマネジメントシステムを開発・提供しています。',
      custom_fields: [
        { field_name: '業界', content: 'エネルギー・環境' },
        { field_name: '従業員数', content: '120名' }
      ],
      job_count: 4,
      created_at: '2024-03-12T11:20:00Z',
      updated_at: '2024-03-12T11:20:00Z'
    },
    {
      id: '5',
      name: 'ヘルスケアソリューションズ株式会社',
      business_description: 'IoTとAIを活用した遠隔医療プラットフォームの開発・運営。予防医療から治療後のケアまで、包括的な医療サービスを提供しています。',
      custom_fields: [
        { field_name: '業界', content: 'ヘルスケア' },
        { field_name: '従業員数', content: '200名' }
      ],
      job_count: 6,
      created_at: '2024-03-11T14:15:00Z',
      updated_at: '2024-03-11T14:15:00Z'
    },
    {
      id: '6',
      name: 'デジタルエデュケーション株式会社',
      business_description: 'オンライン教育プラットフォームの開発・運営。個別最適化された学習体験を提供し、生涯学習をサポートしています。',
      custom_fields: [
        { field_name: '業界', content: 'エドテック' },
        { field_name: '従業員数', content: '90名' }
      ],
      job_count: 3,
      created_at: '2024-03-10T16:40:00Z',
      updated_at: '2024-03-10T16:40:00Z'
    },
    {
      id: '7',
      name: 'スマートモビリティ株式会社',
      business_description: '次世代モビリティサービスの開発。自動運転技術とシェアリングエコノミーを組み合わせた革新的な移動手段を提供しています。',
      custom_fields: [
        { field_name: '業界', content: 'モビリティ' },
        { field_name: '従業員数', content: '180名' }
      ],
      job_count: 5,
      created_at: '2024-03-09T13:25:00Z',
      updated_at: '2024-03-09T13:25:00Z'
    },
    {
      id: '8',
      name: 'サイバーセキュリティプロ株式会社',
      business_description: '最新のセキュリティ技術を活用したサイバー攻撃対策ソリューションを提供。企業のデジタルアセットを守ります。',
      custom_fields: [
        { field_name: '業界', content: 'セキュリティ' },
        { field_name: '従業員数', content: '110名' }
      ],
      job_count: 4,
      created_at: '2024-03-08T10:50:00Z',
      updated_at: '2024-03-08T10:50:00Z'
    },
    {
      id: '9',
      name: 'クラウドインフラ株式会社',
      business_description: 'エンタープライズ向けクラウドインフラストラクチャサービスを提供。高可用性と安全性を重視したシステム基盤を構築します。',
      custom_fields: [
        { field_name: '業界', content: 'クラウド・インフラ' },
        { field_name: '従業員数', content: '250名' }
      ],
      job_count: 7,
      created_at: '2024-03-07T09:15:00Z',
      updated_at: '2024-03-07T09:15:00Z'
    }
  ];

  const [companies, setCompanies] = useState<Company[]>(mockCompanies);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);

  // 企業情報登録モーダルの状態
  const [isModalOpen, setIsModalOpen] = useState(false);

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-2xl font-bold text-gray-900">企業/求人管理</h1>
        <button
          onClick={() => setIsModalOpen(true)}
          className="inline-flex items-center px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
        >
          <PlusIcon className="w-5 h-5 mr-2" />
          企業を登録
        </button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {companies.map((company) => (
          <CompanyCard key={company.id} company={company} />
        ))}
      </div>

      {companies.length === 0 && (
        <div className="text-center py-12">
          <p className="text-gray-500">登録されている企業はありません</p>
        </div>
      )}

      {/* ページネーション */}
      {totalPages > 1 && (
        <div className="flex justify-center mt-8">
          <nav className="flex items-center gap-2">
            <button
              onClick={() => setCurrentPage(prev => Math.max(prev - 1, 1))}
              disabled={currentPage === 1}
              className="px-3 py-1 rounded border border-gray-300 disabled:opacity-50"
            >
              前へ
            </button>
            <span className="px-4 py-1">
              {currentPage} / {totalPages}
            </span>
            <button
              onClick={() => setCurrentPage(prev => Math.min(prev + 1, totalPages))}
              disabled={currentPage === totalPages}
              className="px-3 py-1 rounded border border-gray-300 disabled:opacity-50"
            >
              次へ
            </button>
          </nav>
        </div>
      )}
    </div>
  );
} 