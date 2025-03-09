// APIクライアントの基本設定
import { Company, CompanyInput, CompanyListResponse, Job, JobInput } from './types';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

// 基本的なフェッチ関数
async function fetchAPI<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const url = `${API_BASE_URL}${endpoint}`;
  
  const defaultHeaders = {
    'Content-Type': 'application/json',
  };

  const response = await fetch(url, {
    ...options,
    headers: {
      ...defaultHeaders,
      ...options.headers,
    },
  });

  if (!response.ok) {
    const error = await response.json().catch(() => ({}));
    throw new Error(error.message || `API error: ${response.status}`);
  }

  return response.json();
}

// 企業情報のAPI関数
export const companyAPI = {
  // 企業一覧を取得
  getCompanies: () => fetchAPI<CompanyListResponse>('/api/v1/companies'),
  
  // 企業詳細を取得
  getCompany: (id: number) => fetchAPI<Company>(`/api/v1/companies/${id}`),
  
  // 企業を作成
  createCompany: (data: CompanyInput) => fetchAPI<Company>('/api/v1/companies', {
    method: 'POST',
    body: JSON.stringify(data),
  }),
  
  // 企業を更新
  updateCompany: (id: number, data: CompanyInput) => fetchAPI<Company>(`/api/v1/companies/${id}`, {
    method: 'PUT',
    body: JSON.stringify(data),
  }),
  
  // 企業を削除
  deleteCompany: (id: number) => fetchAPI<void>(`/api/v1/companies/${id}`, {
    method: 'DELETE',
  }),
};

// 求人情報のAPI関数
export const jobAPI = {
  // 求人一覧を取得
  getJobs: (companyId?: number) => {
    const endpoint = companyId ? `/api/v1/companies/${companyId}/jobs` : '/api/v1/jobs';
    return fetchAPI<Job[]>(endpoint);
  },
  
  // 求人詳細を取得
  getJob: (id: number) => fetchAPI<Job>(`/api/v1/jobs/${id}`),
  
  // 求人を作成
  createJob: (companyId: number, data: JobInput) => fetchAPI<Job>(`/api/v1/companies/${companyId}/jobs`, {
    method: 'POST',
    body: JSON.stringify(data),
  }),
  
  // 求人を更新
  updateJob: (id: number, data: JobInput) => fetchAPI<Job>(`/api/v1/jobs/${id}`, {
    method: 'PUT',
    body: JSON.stringify(data),
  }),
  
  // 求人を削除
  deleteJob: (id: number) => fetchAPI<void>(`/api/v1/jobs/${id}`, {
    method: 'DELETE',
  }),
}; 