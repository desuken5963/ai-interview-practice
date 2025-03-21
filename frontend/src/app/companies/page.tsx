'use client';

import { useState, useEffect } from 'react';
import { PlusIcon } from '@heroicons/react/24/outline';
import { ChevronLeftIcon, ChevronRightIcon } from '@heroicons/react/20/solid';
import dynamic from 'next/dynamic';
import { Company } from '@/lib/api/types';
import { companyAPI } from '@/lib/api/client';

// クライアントサイドのみでレンダリングするためにdynamic importを使用
const CompanyCard = dynamic(() => import('@/components/companies/CompanyCard'));
const CompanyFormModal = dynamic(() => import('@/components/companies/CompanyFormModal'));

export default function CompaniesPage() {
  const [companies, setCompanies] = useState<Company[]>([]);
  const [isCompanyFormModalOpen, setIsCompanyFormModalOpen] = useState(false);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [totalItems, setTotalItems] = useState(0);
  const itemsPerPage = 6; // 1ページあたり6件に変更

  const fetchCompanies = async (page: number = 1) => {
    try {
      const response = await companyAPI.getCompanies(page, itemsPerPage);
      setCompanies(response.companies);
      setTotalItems(response.total);
      setTotalPages(Math.ceil(response.total / itemsPerPage));
    } catch (error) {
      console.error('Error fetching companies:', error);
    }
  };

  useEffect(() => {
    fetchCompanies(currentPage);
  }, [currentPage]);

  const handlePageChange = (newPage: number) => {
    if (newPage >= 1 && newPage <= totalPages) {
      setCurrentPage(newPage);
    }
  };

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold text-gray-900">企業一覧</h1>
        <button
          onClick={() => setIsCompanyFormModalOpen(true)}
          className="inline-flex items-center px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
        >
          <PlusIcon className="w-5 h-5 mr-2" />
          企業を追加
        </button>
      </div>

      {companies.length === 0 ? (
        <div className="text-center py-8">
          <p className="text-gray-500">登録されている企業はありません</p>
        </div>
      ) : (
        <>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {companies.map((company) => (
              <CompanyCard
                key={company.id}
                company={company}
                onRefresh={() => fetchCompanies(currentPage)}
              />
            ))}
          </div>

          {totalPages > 1 && (
            <div className="flex justify-center items-center space-x-2 mt-8">
              <button
                onClick={() => handlePageChange(currentPage - 1)}
                disabled={currentPage === 1}
                className="p-2 rounded-md border border-gray-300 disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50"
              >
                <ChevronLeftIcon className="w-5 h-5 text-gray-900" />
              </button>
              <div className="flex items-center space-x-1">
                {[...Array(totalPages)].map((_, index) => (
                  <button
                    key={index + 1}
                    onClick={() => handlePageChange(index + 1)}
                    className={`px-3 py-1 rounded-md text-gray-900 ${
                      currentPage === index + 1
                        ? 'bg-blue-600 text-white'
                        : 'border border-gray-300 hover:bg-gray-50'
                    }`}
                  >
                    {index + 1}
                  </button>
                ))}
              </div>
              <button
                onClick={() => handlePageChange(currentPage + 1)}
                disabled={currentPage === totalPages}
                className="p-2 rounded-md border border-gray-300 disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50"
              >
                <ChevronRightIcon className="w-5 h-5 text-gray-900" />
              </button>
            </div>
          )}
        </>
      )}

      <CompanyFormModal
        isOpen={isCompanyFormModalOpen}
        onClose={() => setIsCompanyFormModalOpen(false)}
        onSuccess={() => {
          setIsCompanyFormModalOpen(false);
          fetchCompanies(currentPage);
        }}
      />
    </div>
  );
} 