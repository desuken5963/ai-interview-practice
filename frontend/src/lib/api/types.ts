// APIレスポンスの型（スネークケース）
export type APICustomField = {
  id: number;
  field_name: string;
  content: string;
  created_at: string;
  updated_at: string;
};

export type APIJobPosting = {
  id: number;
  company_id: number;
  title: string;
  description: string | null;
  custom_fields: APICustomField[];
  created_at: string;
  updated_at: string;
};

export type APICompany = {
  id: number;
  name: string;
  business_description: string | null;
  custom_fields: APICustomField[];
  job_postings: APIJobPosting[];
  created_at: string;
  updated_at: string;
};

export type APICompanyListResponse = {
  companies: APICompany[];
  total: number;
  page: number;
  limit: number;
};

// フロントエンド用の型（キャメルケース）
export type CustomField = {
  id: number;
  fieldName: string;
  content: string;
  createdAt: string;
  updatedAt: string;
};

export type JobPosting = {
  id: number;
  companyId: number;
  title: string;
  description: string | null;
  customFields: CustomField[];
  createdAt: string;
  updatedAt: string;
};

export type Company = {
  id: number;
  name: string;
  businessDescription: string | null;
  customFields: CustomField[];
  jobPostings: JobPosting[];
  createdAt: string;
  updatedAt: string;
};

export type CompanyListResponse = {
  companies: Company[];
  total: number;
  page: number;
  limit: number;
};

// 入力用の型
export type CompanyInput = {
  name: string;
  businessDescription: string | null;
  customFields: {
    fieldName: string;
    content: string;
  }[];
};

// 求人一覧のレスポンス型
export type JobPostingResponse = {
  jobPostings: JobPosting[];
  total: number;
  page: number;
  limit: number;
};

// 求人作成・更新用のデータ型
export type JobPostingInput = {
  title: string;
  description: string | null;
  requirements: string | null;
  custom_fields: CustomField[];
}; 