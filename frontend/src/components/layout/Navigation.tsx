'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';

const navigation = [
  { name: '企業/求人管理', href: '/companies' },
  { name: '面接練習履歴', href: '/history' },
  { name: '設定', href: '/settings' },
];

export default function Navigation() {
  const pathname = usePathname();

  return (
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
  );
} 