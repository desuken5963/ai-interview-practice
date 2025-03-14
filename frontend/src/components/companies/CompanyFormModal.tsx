'use client';

import { useState, useEffect } from 'react';
import { XMarkIcon, PlusIcon, TrashIcon } from '@heroicons/react/24/outline';
import { CompanyInput } from '@/lib/api/types';

type CompanyFormModalProps = {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (data: CompanyInput) => Promise<void>;
  initialData?: CompanyInput;
};

export default function CompanyFormModal({
  isOpen,
  onClose,
  onSubmit,
  initialData,
}: CompanyFormModalProps) {
  const [formData, setFormData] = useState<CompanyInput>({
    name: initialData?.name || '',
    business_description: initialData?.business_description || '',
    custom_fields: initialData?.custom_fields || [],
  });

  const [errors, setErrors] = useState<{
    name?: string;
    business_description?: string;
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
        setFormData(initialData);
      } else {
        setFormData({
          name: '',
          business_description: '',
          custom_fields: [{ field_name: '', content: '' }],
        });
      }
      setErrors({});
    }
  }, [isOpen, initialData]);

  // バリデーション関数
  const validateForm = () => {
    const newErrors: typeof errors = {};

    if (!formData.name.trim()) {
      newErrors.name = '企業名は必須です';
    } else if (formData.name.length > 100) {
      newErrors.name = '企業名は100文字以内で入力してください';
    }

    if (formData.business_description && formData.business_description.length > 1000) {
      newErrors.business_description = '事業内容は1000文字以内で入力してください';
    }

    formData.custom_fields.forEach((field, index) => {
      if (!field.field_name.trim()) {
        newErrors[`custom_fields.${index}.field_name`] = '項目名は必須です';
      } else if (field.field_name.length > 50) {
        newErrors[`custom_fields.${index}.field_name`] = '項目名は50文字以内で入力してください';
      }

      if (!field.content.trim()) {
        newErrors[`custom_fields.${index}.content`] = '内容は必須です';
      } else if (field.content.length > 500) {
        newErrors[`custom_fields.${index}.content`] = '内容は500文字以内で入力してください';
      }
    });

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  // カスタムフィールドの追加
  const addCustomField = () => {
    setFormData(prev => ({
      ...prev,
      custom_fields: [
        ...prev.custom_fields,
        { field_name: '', content: '' }
      ]
    }));
  };

  // カスタムフィールドの削除
  const removeCustomField = (index: number) => {
    setFormData(prev => ({
      ...prev,
      custom_fields: prev.custom_fields.filter((_, i) => i !== index)
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
      await onSubmit(formData);
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
                  {initialData ? '企業情報を編集' : '企業を登録'}
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
                  <label htmlFor="business_description" className="block text-sm font-medium text-gray-700 mb-1">
                    事業内容
                  </label>
                  <textarea
                    id="business_description"
                    value={formData.business_description || ''}
                    onChange={(e) => setFormData(prev => ({ ...prev, business_description: e.target.value }))}
                    rows={4}
                    className={`block w-full px-4 py-2 rounded-md border ${
                      errors.business_description ? 'border-red-300' : 'border-gray-300'
                    } text-gray-900 placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500`}
                    placeholder="事業内容を入力"
                  />
                  {errors.business_description && (
                    <p className="mt-1 text-sm text-red-600">{errors.business_description}</p>
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
                      onClick={addCustomField}
                      className="inline-flex items-center px-3 py-1.5 text-sm bg-blue-50 text-blue-600 rounded-md hover:bg-blue-100 transition-colors"
                    >
                      <PlusIcon className="w-4 h-4 mr-1" />
                      項目を追加
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
                              onClick={() => removeCustomField(index)}
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
                              onChange={(e) => {
                                const newFields = [...formData.custom_fields];
                                newFields[index].field_name = e.target.value;
                                setFormData(prev => ({ ...prev, custom_fields: newFields }));
                              }}
                              placeholder="項目名（例：業界、従業員数）"
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
                              onChange={(e) => {
                                const newFields = [...formData.custom_fields];
                                newFields[index].content = e.target.value;
                                setFormData(prev => ({ ...prev, custom_fields: newFields }));
                              }}
                              placeholder="内容（例：IT業界、100名）"
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