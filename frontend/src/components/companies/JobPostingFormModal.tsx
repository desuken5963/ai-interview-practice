'use client';

import { useState, useEffect } from 'react';
import { XMarkIcon, PlusIcon, TrashIcon } from '@heroicons/react/24/outline';
import { JobPosting, JobPostingInput } from '@/lib/api/types';

type JobPostingFormModalProps = {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (data: JobPostingInput) => Promise<void>;
  initialData?: JobPosting;
  companyName: string;
};

export default function JobPostingFormModal({
  isOpen,
  onClose,
  onSubmit,
  initialData,
  companyName,
}: JobPostingFormModalProps) {
  const [formData, setFormData] = useState<JobPostingInput>({
    title: initialData?.title || '',
    description: initialData?.description || '',
    requirements: initialData?.requirements || '',
    custom_fields: initialData?.custom_fields || [],
  });

  const [errors, setErrors] = useState<{
    title?: string;
    description?: string;
    requirements?: string;
    custom_fields?: string;
    submit?: string;
    [key: `custom_fields.${number}.field_name`]: string;
    [key: `custom_fields.${number}.content`]: string;
  }>({});

  const [isSubmitting, setIsSubmitting] = useState(false);

  // モーダルが開かれた時にフォームをリセット
  useEffect(() => {
    if (isOpen) {
      if (initialData) {
        setFormData({
          title: initialData.title,
          description: initialData.description,
          requirements: initialData.requirements,
          custom_fields: initialData.custom_fields,
        });
      } else {
        setFormData({
          title: '',
          description: '',
          requirements: '',
          custom_fields: [{ field_name: '', content: '' }],
        });
      }
      setErrors({});
    }
  }, [isOpen, initialData]);

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
    field: 'field_name' | 'content',
    value: string
  ) => {
    setFormData((prev) => ({
      ...prev,
      custom_fields: prev.custom_fields.map((item, i) =>
        i === index ? { ...item, [field]: value } : item
      ),
    }));
    // エラーをクリア
    setErrors((prev) => ({
      ...prev,
      [`custom_fields.${index}.${field}`]: undefined,
    }));
  };

  // カスタムフィールドを追加
  const handleAddCustomField = () => {
    setFormData((prev) => ({
      ...prev,
      custom_fields: [...prev.custom_fields, { field_name: '', content: '' }],
    }));
  };

  // カスタムフィールドを削除
  const handleRemoveCustomField = (index: number) => {
    setFormData((prev) => ({
      ...prev,
      custom_fields: prev.custom_fields.filter((_, i) => i !== index),
    }));
  };

  // フォームのバリデーション
  const validateForm = () => {
    const newErrors: typeof errors = {};

    if (!formData.title.trim()) {
      newErrors.title = '求人タイトルを入力してください';
    }

    formData.custom_fields.forEach((field, index) => {
      if (!field.field_name.trim()) {
        newErrors[`custom_fields.${index}.field_name`] = '項目名を入力してください';
      }
      if (!field.content.trim()) {
        newErrors[`custom_fields.${index}.content`] = '内容を入力してください';
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

    try {
      setIsSubmitting(true);
      await onSubmit(formData);
      onClose();
    } catch (error) {
      console.error('Error submitting form:', error);
      setErrors((prev) => ({
        ...prev,
        submit: '保存に失敗しました',
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
                  {initialData ? '求人を編集' : '求人を追加'} - {companyName}
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

                {/* 応募要件 */}
                <div>
                  <label htmlFor="requirements" className="block text-sm font-medium text-gray-700 mb-1">
                    応募要件
                  </label>
                  <textarea
                    id="requirements"
                    name="requirements"
                    rows={4}
                    value={formData.requirements || ''}
                    onChange={handleInputChange}
                    className={`block w-full px-4 py-2 rounded-md border ${
                      errors.requirements ? 'border-red-300' : 'border-gray-300'
                    } text-gray-900 placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500`}
                    placeholder="応募に必要な条件や資格などを入力"
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
                    {formData.custom_fields.map((field, index) => (
                      <div key={index} className="p-4 border border-gray-200 rounded-md bg-gray-50">
                        <div className="flex justify-between items-center mb-2">
                          <label className="text-sm font-medium text-gray-700">項目 {index + 1}</label>
                          {formData.custom_fields.length > 1 && (
                            <button
                              type="button"
                              onClick={() => handleRemoveCustomField(index)}
                              className="p-1 text-gray-400 hover:text-red-600 rounded-full hover:bg-gray-100"
                            >
                              <TrashIcon className="w-5 h-5" />
                            </button>
                          )}
                        </div>
                        <div className="space-y-3">
                          <div>
                            <input
                              type="text"
                              value={field.field_name}
                              onChange={(e) =>
                                handleCustomFieldChange(index, 'field_name', e.target.value)
                              }
                              placeholder="項目名（例：勤務地、給与）"
                              className={`block w-full px-4 py-2 rounded-md border ${
                                errors[`custom_fields.${index}.field_name`] ? 'border-red-300' : 'border-gray-300'
                              } text-gray-900 placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500`}
                            />
                            {errors[`custom_fields.${index}.field_name`] && (
                              <p className="mt-1 text-sm text-red-600">
                                {errors[`custom_fields.${index}.field_name`]}
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
                                errors[`custom_fields.${index}.content`] ? 'border-red-300' : 'border-gray-300'
                              } text-gray-900 placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500`}
                            />
                            {errors[`custom_fields.${index}.content`] && (
                              <p className="mt-1 text-sm text-red-600">
                                {errors[`custom_fields.${index}.content`]}
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
                  initialData ? '更新する' : '登録する'
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