// APIクライアントの基本設定
import { Company, CompanyInput, CompanyListResponse, Job, JobInput, JobResponse } from './types';

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

  // 204 No Content または Content-Length が 0 の場合は空のオブジェクトを返す
  if (response.status === 204 || response.headers.get('Content-Length') === '0') {
    return {} as T;
  }

  try {
    return await response.json();
  } catch (error) {
    console.warn('Failed to parse JSON response:', error);
    return {} as T;
  }
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

  // 企業と紐づく求人一覧を一括取得
  getCompanyWithJobs: (id: number) => fetchAPI<Company>(`/api/v1/companies/${id}/with-jobs`),
};

// 求人情報のAPI関数
export const jobAPI = {
  // 求人一覧を取得（企業IDが指定されている場合は企業に紐づく求人一覧を取得）
  getJobs: (companyId?: number) => {
    const endpoint = companyId 
      ? `/api/v1/companies/${companyId}/job-postings` 
      : '/api/v1/job-postings';
    return fetchAPI<JobResponse>(endpoint);
  },
  
  // 求人詳細を取得
  getJob: (id: number) => fetchAPI<Job>(`/api/v1/job-postings/${id}`),
  
  // 求人を作成
  createJob: (companyId: number, data: JobInput) => fetchAPI<Job>(`/api/v1/companies/${companyId}/job-postings`, {
    method: 'POST',
    body: JSON.stringify(data),
  }),
  
  // 求人を更新
  updateJob: (id: number, data: JobInput) => fetchAPI<Job>(`/api/v1/job-postings/${id}`, {
    method: 'PUT',
    body: JSON.stringify(data),
  }),
  
  // 求人を削除
  deleteJob: (id: number) => fetchAPI<void>(`/api/v1/job-postings/${id}`, {
    method: 'DELETE',
  }),

  // 企業と紐づく求人一覧を一括取得
  getCompanyWithJobs: (companyId: number) => fetchAPI<Company>(`/api/v1/companies/${companyId}/with-job-postings`),
}; 