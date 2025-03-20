'use client';

import { XMarkIcon } from '@heroicons/react/24/outline';

type ErrorMessageProps = {
  message: string;
  className?: string;
  onClose?: () => void;
};

export function ErrorMessage({ 
  message, 
  className = '',
  onClose
}: ErrorMessageProps) {
  return (
    <div className={`p-3 bg-red-100 text-red-700 rounded-md flex justify-between items-start ${className}`}>
      <div>{message}</div>
      {onClose && (
        <button
          type="button"
          onClick={onClose}
          className="text-red-700 hover:text-red-900"
          aria-label="エラーメッセージを閉じる"
        >
          <XMarkIcon className="h-5 w-5" />
        </button>
      )}
    </div>
  );
} 