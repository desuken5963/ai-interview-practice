'use client';

import { CustomField } from '@/lib/api/types';

type CustomFieldsListProps = {
  fields: CustomField[];
  variant?: 'default' | 'compact';
  className?: string;
};

export function CustomFieldsList({ 
  fields, 
  variant = 'default',
  className = ''
}: CustomFieldsListProps) {
  if (fields.length === 0) return null;

  const gridClass = variant === 'compact' 
    ? 'grid-cols-1 sm:grid-cols-2' 
    : 'grid-cols-2';
  const textClass = variant === 'compact'
    ? 'text-xs sm:text-sm'
    : 'text-sm';

  return (
    <div className={`mb-4 ${className}`}>
      <dl className={`grid gap-2 ${gridClass}`}>
        {fields.map((field, index) => (
          <div key={index} className="col-span-1">
            <dt className={`font-medium text-gray-500 ${textClass}`}>
              {field.fieldName}
            </dt>
            <dd className={`text-gray-900 ${textClass}`}>
              {field.content}
            </dd>
          </div>
        ))}
      </dl>
    </div>
  );
} 