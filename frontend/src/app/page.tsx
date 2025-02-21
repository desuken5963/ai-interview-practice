import Link from 'next/link';
import { ArrowRightIcon, ChatBubbleLeftRightIcon, MicrophoneIcon, DocumentChartBarIcon } from '@heroicons/react/24/outline';

export default function Home() {
  return (
    <div className="min-h-screen bg-gradient-to-b from-white to-blue-50">
      {/* ヒーローセクション */}
      <section className="relative py-20 px-4">
        <div className="container mx-auto text-center">
          <h1 className="text-4xl md:text-5xl font-bold text-gray-900 mb-6">
            AIで実践的な面接練習を、<br />
            <span className="text-blue-600">いつでもどこでも</span>
          </h1>
          <p className="text-xl text-gray-600 mb-8 max-w-2xl mx-auto">
            企業・求人情報に基づいた模擬面接で、<br />
            本番さながらの練習と即時フィードバックを提供します
          </p>
          <Link
            href="/companies"
            className="inline-flex items-center px-6 py-3 text-lg font-medium text-white bg-blue-600 rounded-lg hover:bg-blue-700 transition-colors"
          >
            無料で始める
            <ArrowRightIcon className="w-5 h-5 ml-2" />
          </Link>
        </div>
      </section>

      {/* 特徴セクション */}
      <section className="py-16 px-4 bg-white">
        <div className="container mx-auto">
          <h2 className="text-3xl font-bold text-center text-gray-900 mb-12">
            メンレンの特徴
          </h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            {/* テキストチャット面接 */}
            <div className="bg-white p-6 rounded-lg shadow-lg">
              <div className="text-blue-600 mb-4">
                <ChatBubbleLeftRightIcon className="w-12 h-12" />
              </div>
              <h3 className="text-xl font-bold text-gray-900 mb-3">
                テキストチャット面接
              </h3>
              <p className="text-gray-600">
                場所を選ばず、スマートフォンやPCから手軽に面接練習が可能。
                テキストベースで丁寧に回答を組み立てられます。
              </p>
            </div>

            {/* 音声面接 */}
            <div className="bg-white p-6 rounded-lg shadow-lg">
              <div className="text-blue-600 mb-4">
                <MicrophoneIcon className="w-12 h-12" />
              </div>
              <h3 className="text-xl font-bold text-gray-900 mb-3">
                音声面接
              </h3>
              <p className="text-gray-600">
                より実践的な面接練習が可能。音声での回答を
                AIがリアルタイムで解析し、フィードバックを提供します。
              </p>
            </div>

            {/* AIフィードバック */}
            <div className="bg-white p-6 rounded-lg shadow-lg">
              <div className="text-blue-600 mb-4">
                <DocumentChartBarIcon className="w-12 h-12" />
              </div>
              <h3 className="text-xl font-bold text-gray-900 mb-3">
                AIフィードバック
              </h3>
              <p className="text-gray-600">
                回答の論理性、説得力、企業文化との適合性など、
                多角的な評価とアドバイスを即時に提供します。
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* 使い方セクション */}
      <section className="py-16 px-4">
        <div className="container mx-auto">
          <h2 className="text-3xl font-bold text-center text-gray-900 mb-12">
            簡単3ステップで始められます
          </h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            {/* ステップ1 */}
            <div className="text-center">
              <div className="text-4xl font-bold text-blue-600 mb-4">1</div>
              <h3 className="text-xl font-bold text-gray-900 mb-3">
                企業・求人情報を登録
              </h3>
              <p className="text-gray-600">
                志望企業の情報を登録するだけ。
                登録せずに一般的な面接練習も可能です。
              </p>
            </div>

            {/* ステップ2 */}
            <div className="text-center">
              <div className="text-4xl font-bold text-blue-600 mb-4">2</div>
              <h3 className="text-xl font-bold text-gray-900 mb-3">
                面接形式を選択
              </h3>
              <p className="text-gray-600">
                テキストチャットか音声面接を選択。
                面接官の役職や面接フェーズも設定できます。
              </p>
            </div>

            {/* ステップ3 */}
            <div className="text-center">
              <div className="text-4xl font-bold text-blue-600 mb-4">3</div>
              <h3 className="text-xl font-bold text-gray-900 mb-3">
                AIと面接練習開始
              </h3>
              <p className="text-gray-600">
                AIが面接官として質問を投げかけ、
                あなたの回答に対して即座にフィードバックを提供します。
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* CTAセクション */}
      <section className="py-20 px-4 bg-blue-600">
        <div className="container mx-auto text-center">
          <h2 className="text-3xl font-bold text-white mb-6">
            今すぐ面接練習を始めましょう
          </h2>
          <p className="text-xl text-blue-100 mb-8 max-w-2xl mx-auto">
            企業研究から面接対策まで、<br />
            あなたの就職・転職活動をAIがサポートします
          </p>
          <Link
            href="/companies"
            className="inline-flex items-center px-8 py-4 text-lg font-medium text-blue-600 bg-white rounded-lg hover:bg-gray-100 transition-colors"
          >
            無料で始める
            <ArrowRightIcon className="w-5 h-5 ml-2" />
          </Link>
        </div>
      </section>

      {/* フッター */}
      <footer className="bg-gray-900 text-white py-12">
        <div className="container mx-auto px-4">
          <div className="text-center">
            <h2 className="text-2xl font-bold mb-4">メンレン</h2>
            <p className="text-gray-400 mb-4">AI面接練習サービス</p>
            <div className="flex justify-center space-x-6">
              <Link href="/companies" className="hover:text-blue-400">
                企業/求人管理
              </Link>
              <Link href="/history" className="hover:text-blue-400">
                面接練習履歴
              </Link>
              <Link href="/settings" className="hover:text-blue-400">
                設定
              </Link>
            </div>
          </div>
        </div>
      </footer>
    </div>
  );
}
