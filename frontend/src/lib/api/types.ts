// 企業情報の型定義
export type Company = {
  id: number;
  name: string;
  business_description: string | null;
  custom_fields: CustomField[];
  job_count?: number;
  created_at: string;
  updated_at: string;
};

// カスタムフィールドの型定義
export type CustomField = {
  field_name: string;
  content: string;
};

// 企業作成・更新用のデータ型
export type CompanyInput = {
  name: string;
  business_description: string | null;
  custom_fields: CustomField[];
};

// 企業一覧のレスポンス型
export type CompanyListResponse = {
  companies: Company[];
  total: number;
  page: number;
  limit: number;
};

// 求人情報の型定義
export type Job = {
  id: number;
  company_id: number;
  title: string;
  description: string | null;
  requirements: string | null;
  custom_fields: CustomField[];
  created_at: string;
  updated_at: string;
};

// 求人作成・更新用のデータ型
export type JobInput = {
  title: string;
  description: string | null;
  requirements: string | null;
  custom_fields: CustomField[];
}; 