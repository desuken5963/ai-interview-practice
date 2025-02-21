import { PlayIcon, PencilIcon, TrashIcon } from '@heroicons/react/24/outline';

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
};

export default function CompanyCard({ company }: CompanyCardProps) {
  return (
    <div className="bg-white rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow">
      <div className="flex justify-between items-start mb-4">
        <h2 className="text-xl font-semibold text-gray-900">{company.name}</h2>
        <div className="flex gap-2">
          <button
            className="p-1 text-gray-500 hover:text-gray-700 transition-colors"
            onClick={() => {/* 編集モーダルを開く */}}
          >
            <PencilIcon className="w-5 h-5" />
          </button>
          <button
            className="p-1 text-gray-500 hover:text-red-600 transition-colors"
            onClick={() => {/* 削除確認モーダルを開く */}}
          >
            <TrashIcon className="w-5 h-5" />
          </button>
        </div>
      </div>

      {company.business_description && (
        <p className="text-gray-600 mb-4 line-clamp-3">
          {company.business_description}
        </p>
      )}

      <div className="flex items-center justify-between mt-4">
        <div className="text-sm text-gray-500">
          求人数: <span className="font-semibold">{company.job_count}</span>
        </div>
        <button
          className="inline-flex items-center px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 transition-colors"
          onClick={() => {/* 面接練習を開始 */}}
        >
          <PlayIcon className="w-4 h-4 mr-2" />
          面接練習
        </button>
      </div>
    </div>
  );
} 