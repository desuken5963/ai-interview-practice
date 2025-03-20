// APIクライアントの基本設定
import {
  Company,
  CompanyInput,
  CompanyListResponse,
  APICompany,
  APICompanyListResponse,
  CustomField,
  APICustomField,
  JobPosting,
  APIJobPosting
} from './types';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

// スネークケースからキャメルケースへの変換関数
function convertCustomField(field: APICustomField): CustomField {
  return {
    id: field.id,
    fieldName: field.field_name,
    content: field.content,
    createdAt: field.created_at,
    updatedAt: field.updated_at,
  };
}

function convertJobPosting(posting: APIJobPosting): JobPosting {
  return {
    id: posting.id,
    companyId: posting.company_id,
    title: posting.title,
    description: posting.description,
    customFields: posting.custom_fields.map(convertCustomField),
    createdAt: posting.created_at,
    updatedAt: posting.updated_at,
  };
}

function convertCompany(company: APICompany): Company {
  return {
    id: company.id,
    name: company.name,
    businessDescription: company.business_description,
    customFields: company.custom_fields.map(convertCustomField),
    jobPostings: company.job_postings.map(convertJobPosting),
    createdAt: company.created_at,
    updatedAt: company.updated_at,
  };
}

// キャメルケースからスネークケースへの変換関数
function convertToAPICustomField(field: Partial<CustomField>): Partial<APICustomField> {
  return {
    field_name: field.fieldName,
    content: field.content,
  };
}

function convertToAPICompanyInput(input: CompanyInput): any {
  return {
    name: input.name,
    business_description: input.businessDescription,
    custom_fields: input.customFields.map(convertToAPICustomField),
  };
}

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
  getCompanies: async () => {
    const response = await fetchAPI<APICompanyListResponse>('/api/v1/companies');
    return {
      companies: response.companies.map(convertCompany),
      total: response.total,
      page: response.page,
      limit: response.limit,
    } as CompanyListResponse;
  },
  
  // 企業を作成
  createCompany: async (data: CompanyInput) => {
    const apiData = convertToAPICompanyInput(data);
    const response = await fetchAPI<APICompany>('/api/v1/companies', {
      method: 'POST',
      body: JSON.stringify(apiData),
    });
    return convertCompany(response);
  },
  
  // 企業を更新
  updateCompany: async (id: number, data: CompanyInput) => {
    const apiData = convertToAPICompanyInput(data);
    const response = await fetchAPI<APICompany>(`/api/v1/companies/${id}`, {
      method: 'PUT',
      body: JSON.stringify(apiData),
    });
    return convertCompany(response);
  },
  
  // 企業を削除
  deleteCompany: (id: number) => fetchAPI<void>(`/api/v1/companies/${id}`, {
    method: 'DELETE',
  }),
};

// 求人情報のAPI関数
export const jobPostingAPI = {
  // 求人一覧を取得（企業IDが指定されている場合は企業に紐づく求人一覧を取得）
  getJobPostings: (companyId?: number) => {
    const endpoint = companyId 
      ? `/api/v1/companies/${companyId}/job-postings` 
      : '/api/v1/job-postings';
    return fetchAPI<JobPostingResponse>(endpoint);
  },
  
  // 求人詳細を取得
  getJobPosting: (id: number) => fetchAPI<JobPosting>(`/api/v1/job-postings/${id}`),
  
  // 求人を作成
  createJobPosting: (companyId: number, data: JobPostingInput) => fetchAPI<JobPosting>(`/api/v1/companies/${companyId}/job-postings`, {
    method: 'POST',
    body: JSON.stringify(data),
  }),
  
  // 求人を更新
  updateJobPosting: (id: number, data: JobPostingInput) => fetchAPI<JobPosting>(`/api/v1/job-postings/${id}`, {
    method: 'PUT',
    body: JSON.stringify(data),
  }),
  
  // 求人を削除
  deleteJobPosting: (id: number) => fetchAPI<void>(`/api/v1/job-postings/${id}`, {
    method: 'DELETE',
  }),
}; 