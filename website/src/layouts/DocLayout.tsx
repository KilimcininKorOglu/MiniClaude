import { useState, useCallback } from 'react'
import { Outlet, Link } from 'react-router-dom'
import Sidebar from '../components/Sidebar'
import SearchModal from '../components/SearchModal'
import { useLanguage } from '../contexts/LanguageContext'

export default function DocLayout() {
  const [searchOpen, setSearchOpen] = useState(false)
  const openSearch = useCallback(() => setSearchOpen(true), [])
  const closeSearch = useCallback(() => setSearchOpen(false), [])
  const { t } = useLanguage()

  return (
    <div className="min-h-screen bg-dark text-light font-mono">
      {/* Header */}
      <header className="sticky top-0 z-40 flex h-14 items-center gap-4 border-b border-white/10 bg-dark/95 px-4 backdrop-blur-sm">
        <Link to="/" className="shrink-0 font-display text-lg tracking-wide text-accent">
          MINICLAUDE
        </Link>
        <span className="hidden text-xs text-white/30 sm:inline">/</span>
        <Link to="/guide/quick-start" className="hidden text-xs text-white/50 transition-colors hover:text-white sm:inline">
          {t('nav.guide')}
        </Link>
        <span className="hidden text-xs text-white/30 sm:inline">/</span>
        <Link to="/features/commands" className="hidden text-xs text-white/50 transition-colors hover:text-white sm:inline">
          {t('nav.features')}
        </Link>

        <div className="flex-1" />

        {/* Search trigger */}
        <button
          onClick={openSearch}
          className="flex items-center gap-2 border border-white/20 px-3 py-1.5 text-xs text-white/40 transition-colors hover:border-white/40"
        >
          <svg className="h-3.5 w-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          <span className="hidden sm:inline">{t('search.placeholder')}</span>
          <kbd className="hidden text-xs text-white/20 lg:inline">Ctrl+K</kbd>
        </button>
      </header>

      <div className="flex">
        <Sidebar />
        <main className="min-w-0 flex-1">
          <div className="mx-auto max-w-4xl px-4 py-8 md:px-8">
            <Outlet />
          </div>
        </main>
      </div>

      <SearchModal open={searchOpen} onClose={closeSearch} />
    </div>
  )
}
