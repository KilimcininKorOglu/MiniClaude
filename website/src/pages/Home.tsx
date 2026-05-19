import Hero from '../components/Hero'
import Features from '../components/Features'
import Comparison from '../components/Comparison'
import Installation from '../components/Installation'
import Footer from '../components/Footer'
import { Link } from 'react-router-dom'
import { useLanguage } from '../contexts/LanguageContext'

export default function HomePage() {
  const { t } = useLanguage()

  return (
    <div className="relative min-h-screen">
      {/* Header with docs link */}
      <div className="absolute top-4 right-4 z-20 flex items-center gap-3">
        <Link
          to="/guide/quick-start"
          className="border border-white/20 px-3 py-1.5 font-mono text-xs text-white/50 transition-colors hover:border-white/40 hover:text-white"
        >
          {t('nav.guide')}
        </Link>
      </div>
      <Hero />
      <Features />
      <Comparison />
      <Installation />
      <Footer />
    </div>
  )
}
