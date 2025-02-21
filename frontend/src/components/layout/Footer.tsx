import Link from 'next/link';

const footerNavigation = {
  company: [
    { name: '運営会社', href: '/company' },
    { name: 'プライバシーポリシー', href: '/privacy' },
    { name: '利用規約', href: '/terms' },
    { name: 'お問い合わせ', href: '/contact' },
  ],
};

export default function Footer() {
  return (
    <footer className="bg-white border-t border-gray-200">
      <div className="container mx-auto px-4 py-8">
        <div className="flex flex-col md:flex-row justify-between items-center space-y-4 md:space-y-0">
          {/* ロゴ部分 */}
          <div className="flex items-center">
            <Link href="/" className="text-lg font-semibold text-blue-600">
              AI Interview
            </Link>
          </div>

          {/* ナビゲーション */}
          <nav className="flex flex-wrap justify-center gap-x-8 gap-y-2">
            {footerNavigation.company.map((item) => (
              <Link
                key={item.name}
                href={item.href}
                className="text-sm text-gray-500 hover:text-gray-900 transition-colors"
              >
                {item.name}
              </Link>
            ))}
          </nav>
        </div>

        {/* コピーライト */}
        <div className="mt-8 text-center text-sm text-gray-400">
          <p>© {new Date().getFullYear()} AI Interview. All rights reserved.</p>
        </div>
      </div>
    </footer>
  );
} 