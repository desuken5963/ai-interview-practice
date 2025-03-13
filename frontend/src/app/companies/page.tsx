'use client';

import { useState, useEffect, useCallback } from 'react';
import { PlusIcon } from '@heroicons/react/24/outline';
import dynamic from 'next/dynamic';
import { companyAPI } from '@/lib/api/client';
import { Company, CompanyInput } from '@/lib/api/types';

// クライアントサイドのみでレンダリングするためにdynamic importを使用
const CompanyCard = dynamic(() => import('@/components/companies/CompanyCard'), {
  ssr: false,
});

const CompanyFormModal = dynamic(() => import('@/components/companies/CompanyFormModal'), {
  ssr: false,
});

export default function CompaniesPage() {
  const [companies, setCompanies] = useState<Company[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [mounted, setMounted] = useState(false);

  // クライアントサイドでのみレンダリングされるようにする
  useEffect(() => {
    setMounted(true);
  }, []);

  // 企業情報の取得
  const fetchCompanies = useCallback(async () => {
    if (!mounted) return;

    try {
      setLoading(true);
      const response = await companyAPI.getCompanies();
      setCompanies(response.companies);
      // ページネーションの設定
      setTotalPages(Math.ceil(response.total / response.limit));
      setCurrentPage(response.page);
    } catch (err) {
      console.error('Failed to fetch companies:', err);
      setError('企業情報の取得に失敗しました。');
    } finally {
      setLoading(false);
    }
  }, [mounted]);

  // マウント時に企業情報を取得
  useEffect(() => {
    if (mounted) {
      fetchCompanies();
    }
  }, [mounted, fetchCompanies]);

  // 企業情報の登録処理
  const handleSubmit = async (data: CompanyInput) => {
    try {
      setLoading(true);
      const newCompany = await companyAPI.createCompany(data);
      setCompanies(prev => [newCompany, ...prev]);
      setIsModalOpen(false);
    } catch (error) {
      console.error('Error submitting company:', error);
      setError('企業情報の登録に失敗しました。');
    } finally {
      setLoading(false);
    }
  };

  // 企業情報の更新処理
  const handleUpdate = async (companyId: number, data: CompanyInput) => {
    try {
      setLoading(true);
      const updatedCompany = await companyAPI.updateCompany(companyId, data);
      
      // 更新データの整合性チェック
      if (!updatedCompany || !updatedCompany.id) {
        throw new Error('更新後のデータが不正です');
      }
      
      // バッチ更新を使用して、ステートの整合性を保証
      setCompanies(prev => {
        const newCompanies = [...prev];
        const index = newCompanies.findIndex(c => c.id === companyId);
        if (index !== -1) {
          newCompanies[index] = {
            ...newCompanies[index],
            ...updatedCompany,
            id: companyId  // IDを明示的に保持
          };
        }
        return newCompanies;
      });
    } catch (error) {
      console.error('Error updating company:', error);
      setError('企業情報の更新に失敗しました。');
    } finally {
      setLoading(false);
    }
  };

  // 企業情報の削除処理
  const handleDelete = async (companyId: number) => {
    try {
      setLoading(true);
      await companyAPI.deleteCompany(companyId);
      // 削除後は配列から該当企業を削除
      setCompanies(prev => prev.filter(company => company.id !== companyId));
    } catch (error) {
      console.error('Error deleting company:', error);
      setError('企業情報の削除に失敗しました。');
    } finally {
      setLoading(false);
    }
  };

  // 特定の企業情報を更新
  const handleRefreshCompany = async (companyId: number) => {
    try {
      setLoading(true);
      const updatedCompany = await companyAPI.getCompany(companyId);
      
      // 更新データの整合性チェック
      if (!updatedCompany || !updatedCompany.id) {
        throw new Error('更新後のデータが不正です');
      }
      
      // バッチ更新を使用して、ステートの整合性を保証
      setCompanies(prev => {
        const newCompanies = [...prev];
        const index = newCompanies.findIndex(c => c.id === companyId);
        if (index !== -1) {
          newCompanies[index] = {
            ...newCompanies[index],
            ...updatedCompany,
            id: companyId  // IDを明示的に保持
          };
        }
        return newCompanies;
      });
    } catch (error) {
      console.error('Error refreshing company:', error);
      setError('企業情報の更新に失敗しました。');
    } finally {
      setLoading(false);
    }
  };

  // クライアントサイドでのみレンダリングする
  if (!mounted) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="flex justify-between items-center mb-8">
          <h1 className="text-2xl font-bold text-gray-900">企業/求人管理</h1>
          <button
            className="inline-flex items-center px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
          >
            <PlusIcon className="w-5 h-5 mr-2" />
            企業を登録
          </button>
        </div>
        <div className="text-center py-12">
          <div className="inline-block animate-spin rounded-full h-8 w-8 border-4 border-blue-500 border-t-transparent"></div>
          <p className="mt-2 text-gray-600">読み込み中...</p>
        </div>
      </div>
    );
  }

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

      {/* エラーメッセージ */}
      {error && (
        <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative mb-4">
          <span className="block sm:inline">{error}</span>
          <span
            className="absolute top-0 bottom-0 right-0 px-4 py-3 cursor-pointer"
            onClick={() => setError(null)}
          >
            <svg
              className="fill-current h-6 w-6 text-red-500"
              role="button"
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 20 20"
            >
              <title>Close</title>
              <path d="M14.348 14.849a1.2 1.2 0 0 1-1.697 0L10 11.819l-2.651 3.029a1.2 1.2 0 1 1-1.697-1.697l2.758-3.15-2.759-3.152a1.2 1.2 0 1 1 1.697-1.697L10 8.183l2.651-3.031a1.2 1.2 0 1 1 1.697 1.697l-2.758 3.152 2.758 3.15a1.2 1.2 0 0 1 0 1.698z" />
            </svg>
          </span>
        </div>
      )}

      {/* ローディング表示 */}
      {loading && (
        <div className="text-center py-12">
          <div className="inline-block animate-spin rounded-full h-8 w-8 border-4 border-blue-500 border-t-transparent"></div>
          <p className="mt-2 text-gray-600">読み込み中...</p>
        </div>
      )}

      {/* 企業一覧 */}
      {!loading && (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {companies.map((company) => {
            // IDの存在確認
            if (!company?.id) {
              console.warn('Company without ID:', company);
              return null;
            }
            
            return (
              <CompanyCard
                key={company.id}
                company={company}
                onEdit={handleUpdate}
                onDelete={() => handleDelete(company.id)}
                onRefresh={handleRefreshCompany}
              />
            );
          })}
        </div>
      )}

      {!loading && companies.length === 0 && (
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

      {/* 企業情報登録/編集モーダル */}
      <CompanyFormModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSubmit={handleSubmit}
      />
    </div>
  );
} 