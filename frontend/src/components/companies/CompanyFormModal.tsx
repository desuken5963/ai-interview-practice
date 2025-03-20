'use client';

import { useState, useEffect } from 'react';
import { XMarkIcon, PlusIcon, TrashIcon } from '@heroicons/react/24/outline';
import { Company, CompanyInput } from '@/lib/api/types';
import { companyAPI } from '@/lib/api/client';

type CompanyFormModalProps = {
  isOpen: boolean;
  onClose: () => void;
  onSuccess: () => void;
  company?: Company;
};

export default function CompanyFormModal({
  isOpen,
  onClose,
  onSuccess,
  company,
}: CompanyFormModalProps) {
  const [formData, setFormData] = useState<CompanyInput>({
    name: '',
    businessDescription: '',
    customFields: [],
  });

  const [errors, setErrors] = useState<{
    name?: string;
    businessDescription?: string;
    customFields?: string;
    submit?: string;
    [key: `customFields.${number}.fieldName`]: string;
    [key: `customFields.${number}.content`]: string;
  }>({});

  const [isSubmitting, setIsSubmitting] = useState(false);

  // モーダルが開かれた時にフォームをリセット
  useEffect(() => {
    if (isOpen) {
      if (company) {
        setFormData({
          name: company.name,
          businessDescription: company.businessDescription,
          customFields: company.customFields.map(field => ({
            fieldName: field.fieldName,
            content: field.content,
          })),
        });
      } else {
        setFormData({
          name: '',
          businessDescription: '',
          customFields: [{ fieldName: '', content: '' }],
        });
      }
      setErrors({});
    }
  }, [isOpen, company]);

  // バリデーション関数
  const validateForm = () => {
    const newErrors: typeof errors = {};

    if (!formData.name.trim()) {
      newErrors.name = '企業名は必須です';
    } else if (formData.name.length > 100) {
      newErrors.name = '企業名は100文字以内で入力してください';
    }

    if (formData.businessDescription && formData.businessDescription.length > 1000) {
      newErrors.businessDescription = '事業内容は1000文字以内で入力してください';
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

  // カスタムフィールドの追加
  const handleAddCustomField = () => {
    setFormData(prev => ({
      ...prev,
      customFields: [
        ...prev.customFields,
        { fieldName: '', content: '' }
      ]
    }));
  };

  // カスタムフィールドの削除
  const handleRemoveCustomField = (index: number) => {
    setFormData(prev => ({
      ...prev,
      customFields: prev.customFields.filter((_, i) => i !== index)
    }));
  };

  // フォーム送信
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!validateForm()) {
      return;
    }

    setIsSubmitting(true);
    try {
      if (company) {
        // 編集の場合
        await companyAPI.updateCompany(company.id, formData);
      } else {
        // 新規登録の場合
        await companyAPI.createCompany(formData);
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
                  {company ? '企業情報を編集' : '企業を登録'}
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
                {/* 企業名 */}
                <div>
                  <label htmlFor="name" className="block text-sm font-medium text-gray-700 mb-1">
                    企業名 <span className="text-red-500">*</span>
                  </label>
                  <input
                    type="text"
                    id="name"
                    value={formData.name}
                    onChange={(e) => setFormData(prev => ({ ...prev, name: e.target.value }))}
                    className={`block w-full px-4 py-2 rounded-md border ${
                      errors.name ? 'border-red-300' : 'border-gray-300'
                    } text-gray-900 placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500`}
                    placeholder="企業名を入力"
                  />
                  {errors.name && (
                    <p className="mt-1 text-sm text-red-600">{errors.name}</p>
                  )}
                </div>

                {/* 事業内容 */}
                <div>
                  <label htmlFor="businessDescription" className="block text-sm font-medium text-gray-700 mb-1">
                    事業内容
                  </label>
                  <textarea
                    id="businessDescription"
                    value={formData.businessDescription || ''}
                    onChange={(e) => setFormData(prev => ({ ...prev, businessDescription: e.target.value }))}
                    rows={4}
                    className={`block w-full px-4 py-2 rounded-md border ${
                      errors.businessDescription ? 'border-red-300' : 'border-gray-300'
                    } text-gray-900 placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500`}
                    placeholder="事業内容を入力"
                  />
                  {errors.businessDescription && (
                    <p className="mt-1 text-sm text-red-600">{errors.businessDescription}</p>
                  )}
                </div>

                {/* カスタム項目 */}
                <div>
                  <div className="flex justify-between items-center mb-2">
                    <label className="block text-sm font-medium text-gray-700">
                      追加情報
                    </label>
                    <button
                      type="button"
                      onClick={handleAddCustomField}
                      className="inline-flex items-center px-3 py-1.5 text-sm bg-blue-50 text-blue-600 rounded-md hover:bg-blue-100 transition-colors"
                    >
                      <PlusIcon className="w-4 h-4 mr-1" />
                      項目を追加
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
                              onChange={(e) => {
                                const newFields = [...formData.customFields];
                                newFields[index].fieldName = e.target.value;
                                setFormData(prev => ({ ...prev, customFields: newFields }));
                              }}
                              placeholder="項目名（例：業界、従業員数）"
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
                              onChange={(e) => {
                                const newFields = [...formData.customFields];
                                newFields[index].content = e.target.value;
                                setFormData(prev => ({ ...prev, customFields: newFields }));
                              }}
                              placeholder="内容を入力"
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

              <div className="mt-6 flex justify-end gap-3">
                <button
                  type="button"
                  onClick={onClose}
                  className="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                >
                  キャンセル
                </button>
                <button
                  type="submit"
                  disabled={isSubmitting}
                  className="inline-flex justify-center px-4 py-2 text-sm font-medium text-white bg-blue-600 border border-transparent rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
                >
                  {isSubmitting ? '保存中...' : '保存'}
                </button>
              </div>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
} 