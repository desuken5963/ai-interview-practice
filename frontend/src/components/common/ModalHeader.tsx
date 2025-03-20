'use client';

import { XMarkIcon } from '@heroicons/react/24/outline';

type ModalHeaderProps = {
  title: string;
  onClose: () => void;
};

export function ModalHeader({ title, onClose }: ModalHeaderProps) {
  return (
    <div className="flex items-start justify-between mb-4">
      <h3 className="text-lg sm:text-xl font-semibold text-gray-900">
        {title}
      </h3>
      <button
        type="button"
        onClick={onClose}
        className="text-gray-400 hover:text-gray-500"
        aria-label="閉じる"
      >
        <XMarkIcon className="h-6 w-6" />
      </button>
    </div>
  );
} 