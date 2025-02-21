'use client';

import { useState } from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { Bars3Icon, XMarkIcon } from '@heroicons/react/24/outline';

const navigation = [
  { name: '企業/求人管理', href: '/companies' },
  { name: '面接練習履歴', href: '/history' },
  { name: '設定', href: '/settings' },
];

export default function Navigation() {
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const pathname = usePathname();

  return (
    <div className="flex items-center">
      {/* デスクトップナビゲーション */}
      <nav className="hidden md:flex space-x-8">
        {navigation.map((item) => (
          <Link
            key={item.name}
            href={item.href}
            className={`
              px-3 py-2 text-sm font-medium rounded-md transition-colors
              ${pathname === item.href
                ? 'text-blue-600 bg-blue-50'
                : 'text-gray-600 hover:text-blue-600 hover:bg-gray-50'
              }
            `}
          >
            {item.name}
          </Link>
        ))}
      </nav>

      {/* ハンバーガーメニューボタン */}
      <div className="md:hidden">
        <button
          type="button"
          className="p-2 rounded-md text-gray-600 hover:text-gray-900 hover:bg-gray-100"
          onClick={() => setIsMenuOpen(!isMenuOpen)}
        >
          <span className="sr-only">メニューを開く</span>
          {isMenuOpen ? (
            <XMarkIcon className="h-6 w-6" />
          ) : (
            <Bars3Icon className="h-6 w-6" />
          )}
        </button>

        {/* モバイルメニュー */}
        {isMenuOpen && (
          <div className="absolute top-16 left-0 right-0 bg-white shadow-lg">
            <div className="px-2 pt-2 pb-3 space-y-1">
              {navigation.map((item) => (
                <Link
                  key={item.name}
                  href={item.href}
                  className={`
                    block px-3 py-2 rounded-md text-base font-medium transition-colors
                    ${pathname === item.href
                      ? 'text-blue-600 bg-blue-50'
                      : 'text-gray-600 hover:text-blue-600 hover:bg-gray-50'
                    }
                  `}
                  onClick={() => setIsMenuOpen(false)}
                >
                  {item.name}
                </Link>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
} 