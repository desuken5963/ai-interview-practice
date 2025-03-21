'use client';

import { useState, useEffect } from 'react';
import { XMarkIcon, PlusIcon, TrashIcon } from '@heroicons/react/24/outline';
import { JobPosting, JobPostingInput } from '@/lib/api/types';
import { jobPostingAPI } from '@/lib/api/client';

type JobPostingFormModalProps = {
  isOpen: boolean;
  onClose: () => void;
  onSuccess: () => void;
  jobPosting?: JobPosting;
  companyId: number;
  companyName: string;
};

export default function JobPostingFormModal({
  isOpen,
  onClose,
  onSuccess,
  jobPosting,
  companyId,
  companyName,
}: JobPostingFormModalProps) {
  const [formData, setFormData] = useState<JobPostingInput>({
    title: '',
    description: '',
    customFields: [],
  });

  const [errors, setErrors] = useState<{
    title?: string;
    description?: string;
    customFields?: string;
    submit?: string;
    [key: `customFields.${number}.fieldName`]: string;
    [key: `customFields.${number}.content`]: string;
  }>({});

  const [isSubmitting, setIsSubmitting] = useState(false);

  // モーダルが開かれた時にフォームをリセット
  useEffect(() => {
    if (isOpen) {
      if (jobPosting) {
        setFormData({
          title: jobPosting.title,
          description: jobPosting.description || '',
          customFields: jobPosting.customFields.map(field => ({
            fieldName: field.fieldName,
            content: field.content,
          })),
        });
      } else {
        setFormData({
          title: '',
          description: '',
          customFields: [{ fieldName: '', content: '' }],
        });
      }
      setErrors({});
    }
  }, [isOpen, jobPosting]);

  // フォームの入力値を更新
  const handleInputChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
  ) => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }));
    // エラーをクリア
    setErrors((prev) => ({
      ...prev,
      [name]: undefined,
    }));
  };

  // カスタムフィールドの入力値を更新
  const handleCustomFieldChange = (
    index: number,
    field: 'fieldName' | 'content',
    value: string
  ) => {
    setFormData((prev) => ({
      ...prev,
      customFields: prev.customFields.map((item: { fieldName: string; content: string }, i: number) =>
        i === index ? { ...item, [field]: value } : item
      ),
    }));
    // エラーをクリア
    setErrors((prev) => ({
      ...prev,
      [`customFields.${index}.${field}`]: undefined,
    }));
  };

  // カスタムフィールドを追加
  const handleAddCustomField = () => {
    setFormData((prev) => ({
      ...prev,
      customFields: [...prev.customFields, { fieldName: '', content: '' }],
    }));
  };

  // カスタムフィールドを削除
  const handleRemoveCustomField = (index: number) => {
    setFormData((prev) => ({
      ...prev,
      customFields: prev.customFields.filter((_, i) => i !== index),
    }));
  };

  // フォームのバリデーション
  const validateForm = () => {
    const newErrors: typeof errors = {};

    if (!formData.title.trim()) {
      newErrors.title = '求人タイトルは必須です';
    } else if (formData.title.length > 100) {
      newErrors.title = '求人タイトルは100文字以内で入力してください';
    }

    if (formData.description && formData.description.length > 1000) {
      newErrors.description = '求人詳細は1000文字以内で入力してください';
    }

    formData.customFields.forEach((field, index) => {
      if (!field.fieldName.trim()) {
        newErrors[`customFields.${index}.fieldName`] = '項目名は必須です';
      } else if (field.fieldName.length > 50) {
        newErrors[`customFields.${index}.fieldName`] = '項目名は50文字以内で入力してください';
      }

      if (!field.content.trim()) {
        newErrors[`customFields.${index}.content`] = '内容は必須です';
      } else if (field.content.length > 500) {
        newErrors[`customFields.${index}.content`] = '内容は500文字以内で入力してください';
      }
    });

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  // フォームの送信
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!validateForm()) {
      return;
    }

    setIsSubmitting(true);
    try {
      if (jobPosting) {
        // 編集の場合
        await jobPostingAPI.updateJobPosting(jobPosting.id, {
          ...formData,
          company_id: companyId,
        });
      } else {
        // 新規登録の場合
        await jobPostingAPI.createJobPosting(companyId, formData);
      }
      onSuccess();
      onClose();
    } catch (error) {
      console.error('Error submitting form:', error);
      setErrors(prev => ({
        ...prev,
        submit: '保存に失敗しました。もう一度お試しください。'
      }));
    } finally {
      setIsSubmitting(false);
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 overflow-y-auto">
      <div className="flex min-h-full items-center justify-center p-2 sm:p-4 text-center">
        <div className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" onClick={onClose} />

        <div className="relative transform overflow-hidden rounded-lg bg-white text-left shadow-xl transition-all w-full max-w-2xl">
          <form onSubmit={handleSubmit}>
            <div className="bg-white px-4 pb-6 pt-5 sm:p-6">
              <div className="flex items-start justify-between mb-4">
                <h3 className="text-lg sm:text-xl font-semibold text-gray-900">
                  {jobPosting ? '求人を編集' : '求人を追加'} - {companyName}
                </h3>
                <button
                  type="button"
                  onClick={onClose}
                  className="text-gray-400 hover:text-gray-500"
                >
                  <XMarkIcon className="h-6 w-6" />
                </button>
              </div>

              {errors.submit && (
                <div className="mb-4 p-3 bg-red-100 text-red-700 rounded-md">
                  {errors.submit}
                </div>
              )}

              <div className="space-y-4">
                {/* 求人タイトル */}
                <div>
                  <label htmlFor="title" className="block text-sm font-medium text-gray-700 mb-1">
                    求人タイトル <span className="text-red-500">*</span>
                  </label>
                  <input
                    type="text"
                    id="title"
                    name="title"
                    value={formData.title}
                    onChange={handleInputChange}
                    className={`block w-full px-4 py-2 rounded-md border ${
                      errors.title ? 'border-red-300' : 'border-gray-300'
                    } text-gray-900 placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500`}
                    placeholder="求人タイトルを入力"
                  />
                  {errors.title && (
                    <p className="mt-1 text-sm text-red-600">{errors.title}</p>
                  )}
                </div>

                {/* 求人詳細 */}
                <div>
                  <label htmlFor="description" className="block text-sm font-medium text-gray-700 mb-1">
                    求人詳細
                  </label>
                  <textarea
                    id="description"
                    name="description"
                    rows={4}
                    value={formData.description || ''}
                    onChange={handleInputChange}
                    className={`block w-full px-4 py-2 rounded-md border ${
                      errors.description ? 'border-red-300' : 'border-gray-300'
                    } text-gray-900 placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500`}
                    placeholder="求人の詳細情報を入力"
                  />
                </div>

                {/* カスタムフィールド */}
                <div>
                  <div className="flex justify-between items-center mb-2">
                    <label className="block text-sm font-medium text-gray-700">
                      カスタムフィールド
                    </label>
                    <button
                      type="button"
                      onClick={handleAddCustomField}
                      className="inline-flex items-center px-3 py-1.5 text-sm bg-blue-50 text-blue-600 rounded-md hover:bg-blue-100 transition-colors"
                    >
                      <PlusIcon className="w-4 h-4 mr-1" />
                      追加
                    </button>
                  </div>

                  <div className="space-y-4">
                    {formData.customFields.map((field, index) => (
                      <div key={index} className="p-4 border border-gray-200 rounded-md bg-gray-50">
                        <div className="flex justify-between items-center mb-2">
                          <label className="text-sm font-medium text-gray-700">項目 {index + 1}</label>
                          <button
                            type="button"
                            onClick={() => handleRemoveCustomField(index)}
                            className="p-1 text-gray-400 hover:text-red-600 rounded-full hover:bg-gray-100"
                          >
                            <TrashIcon className="w-5 h-5" />
                          </button>
                        </div>
                        <div className="space-y-3">
                          <div>
                            <input
                              type="text"
                              value={field.fieldName}
                              onChange={(e) =>
                                handleCustomFieldChange(index, 'fieldName', e.target.value)
                              }
                              placeholder="項目名（例：勤務地、給与）"
                              className={`block w-full px-4 py-2 rounded-md border ${
                                errors[`customFields.${index}.fieldName`] ? 'border-red-300' : 'border-gray-300'
                              } text-gray-900 placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500`}
                            />
                            {errors[`customFields.${index}.fieldName`] && (
                              <p className="mt-1 text-sm text-red-600">
                                {errors[`customFields.${index}.fieldName`]}
                              </p>
                            )}
                          </div>
                          <div>
                            <input
                              type="text"
                              value={field.content}
                              onChange={(e) =>
                                handleCustomFieldChange(index, 'content', e.target.value)
                              }
                              placeholder="内容（例：東京都、年収500万円〜）"
                              className={`block w-full px-4 py-2 rounded-md border ${
                                errors[`customFields.${index}.content`] ? 'border-red-300' : 'border-gray-300'
                              } text-gray-900 placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500`}
                            />
                            {errors[`customFields.${index}.content`] && (
                              <p className="mt-1 text-sm text-red-600">
                                {errors[`customFields.${index}.content`]}
                              </p>
                            )}
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              </div>
            </div>

            <div className="bg-gray-50 px-4 py-4 sm:px-6 sm:flex sm:flex-row-reverse border-t border-gray-200">
              <button
                type="submit"
                disabled={isSubmitting}
                className="w-full sm:w-auto inline-flex justify-center items-center rounded-md bg-blue-600 px-5 py-2.5 text-sm font-medium text-white shadow-sm hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed sm:ml-3"
              >
                {isSubmitting ? (
                  <>
                    <span className="inline-block w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin mr-2"></span>
                    保存中...
                  </>
                ) : (
                  jobPosting ? '更新する' : '登録する'
                )}
              </button>
              <button
                type="button"
                onClick={onClose}
                disabled={isSubmitting}
                className="mt-3 sm:mt-0 w-full sm:w-auto inline-flex justify-center items-center rounded-md bg-white px-5 py-2.5 text-sm font-medium text-gray-700 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-gray-300 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                キャンセル
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
} 