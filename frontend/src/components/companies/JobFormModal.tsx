'use client';

import { useState, useEffect } from 'react';
import { XMarkIcon, PlusIcon, TrashIcon } from '@heroicons/react/24/outline';

type JobFormData = {
  title: string;
  description: string | null;
  custom_fields: {
    field_name: string;
    content: string;
  }[];
};

type JobFormModalProps = {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (data: JobFormData) => Promise<void>;
  initialData?: JobFormData;
  companyName: string;
};

export default function JobFormModal({
  isOpen,
  onClose,
  onSubmit,
  initialData,
  companyName,
}: JobFormModalProps) {
  const [formData, setFormData] = useState<JobFormData>({
    title: initialData?.title || '',
    description: initialData?.description || '',
    custom_fields: initialData?.custom_fields || [],
  });

  const [errors, setErrors] = useState<{
    title?: string;
    description?: string;
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
          title: '',
          description: '',
          custom_fields: [{ field_name: '', content: '' }],
        });
      }
      setErrors({});
    }
  }, [isOpen, initialData]);

  // バリデーション関数
  const validateForm = () => {
    const newErrors: typeof errors = {};

    if (!formData.title.trim()) {
      newErrors.title = '求人タイトルは必須です';
    } else if (formData.title.length > 100) {
      newErrors.title = '求人タイトルは100文字以内で入力してください';
    }

    if (formData.description && formData.description.length > 1000) {
      newErrors.description = '仕事内容は1000文字以内で入力してください';
    }

    formData.custom_fields.forEach((field, index) => {
      if (field.field_name.trim() || field.content.trim()) {
        if (field.field_name.length > 50) {
          newErrors[`custom_fields.${index}.field_name`] = '項目名は50文字以内で入力してください';
        }
        if (field.content.length > 500) {
          newErrors[`custom_fields.${index}.content`] = '内容は500文字以内で入力してください';
        }
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
      <div className="flex min-h-full items-center justify-center p-4 text-center">
        <div className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" onClick={onClose} />

        <div className="relative transform overflow-hidden rounded-lg bg-white text-left shadow-xl transition-all w-full max-w-2xl">
          <div className="px-6 py-4 border-b border-gray-200 flex justify-between items-center">
            <h3 className="text-xl font-semibold text-gray-900">
              {companyName} - {initialData ? '求人情報を編集' : '求人を登録'}
            </h3>
            <button
              type="button"
              className="text-gray-400 hover:text-gray-500"
              onClick={onClose}
            >
              <XMarkIcon className="h-6 w-6" />
            </button>
          </div>

          <form onSubmit={handleSubmit} className="p-6">
            {/* 求人タイトル */}
            <div className="mb-4">
              <label htmlFor="title" className="block text-sm font-medium text-gray-700 mb-1">
                求人タイトル <span className="text-red-500">*</span>
              </label>
              <input
                type="text"
                id="title"
                value={formData.title}
                onChange={(e) => setFormData(prev => ({ ...prev, title: e.target.value }))}
                className={`w-full rounded-md border ${
                  errors.title ? 'border-red-300' : 'border-gray-300'
                } px-3 py-2 text-gray-900 focus:outline-none focus:ring-2 focus:ring-blue-500`}
                placeholder="求人タイトルを入力"
              />
              {errors.title && (
                <p className="mt-1 text-sm text-red-600">{errors.title}</p>
              )}
            </div>

            {/* 仕事内容 */}
            <div className="mb-4">
              <label htmlFor="description" className="block text-sm font-medium text-gray-700 mb-1">
                仕事内容
              </label>
              <textarea
                id="description"
                value={formData.description || ''}
                onChange={(e) => setFormData(prev => ({ ...prev, description: e.target.value }))}
                rows={4}
                className={`w-full rounded-md border ${
                  errors.description ? 'border-red-300' : 'border-gray-300'
                } px-3 py-2 text-gray-900 focus:outline-none focus:ring-2 focus:ring-blue-500`}
                placeholder="仕事内容を入力"
              />
              {errors.description && (
                <p className="mt-1 text-sm text-red-600">{errors.description}</p>
              )}
            </div>

            {/* カスタム項目 */}
            <div className="mb-4">
              <div className="flex justify-between items-center mb-2">
                <label className="block text-sm font-medium text-gray-700">
                  追加情報
                </label>
                <button
                  type="button"
                  onClick={addCustomField}
                  className="inline-flex items-center px-3 py-1 text-sm text-blue-600 hover:text-blue-700"
                >
                  <PlusIcon className="h-4 w-4 mr-1" />
                  項目を追加
                </button>
              </div>

              {formData.custom_fields.map((field, index) => (
                <div key={index} className="mb-3 p-3 border border-gray-200 rounded-md">
                  <div className="flex justify-between items-start mb-2">
                    <div className="flex-grow mr-2">
                      <input
                        type="text"
                        value={field.field_name}
                        onChange={(e) => {
                          const newFields = [...formData.custom_fields];
                          newFields[index].field_name = e.target.value;
                          setFormData(prev => ({ ...prev, custom_fields: newFields }));
                        }}
                        className="w-full rounded-md border border-gray-300 px-3 py-2 text-sm text-gray-900 focus:outline-none focus:ring-2 focus:ring-blue-500"
                        placeholder="項目名"
                      />
                    </div>
                    <button
                      type="button"
                      onClick={() => removeCustomField(index)}
                      className="text-gray-400 hover:text-red-500"
                    >
                      <TrashIcon className="h-5 w-5" />
                    </button>
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
                      className="w-full rounded-md border border-gray-300 px-3 py-2 text-sm text-gray-900 focus:outline-none focus:ring-2 focus:ring-blue-500"
                      placeholder="内容"
                    />
                  </div>
                  {errors[`custom_fields.${index}.field_name`] && (
                    <p className="mt-1 text-sm text-red-600">
                      {errors[`custom_fields.${index}.field_name`]}
                    </p>
                  )}
                  {errors[`custom_fields.${index}.content`] && (
                    <p className="mt-1 text-sm text-red-600">
                      {errors[`custom_fields.${index}.content`]}
                    </p>
                  )}
                </div>
              ))}
            </div>

            {/* エラーメッセージ */}
            {errors.submit && (
              <div className="mb-4 p-3 bg-red-50 text-red-600 rounded-md">
                {errors.submit}
              </div>
            )}

            {/* ボタン */}
            <div className="flex justify-end gap-3">
              <button
                type="button"
                onClick={onClose}
                className="px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 border border-gray-300 rounded-md"
              >
                キャンセル
              </button>
              <button
                type="submit"
                disabled={isSubmitting}
                className="px-4 py-2 text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 rounded-md disabled:opacity-50"
              >
                {isSubmitting ? '保存中...' : '保存'}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
} 