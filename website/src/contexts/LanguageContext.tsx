import { createContext, useContext, ReactNode } from 'react'

type Language = 'en'

interface LanguageContextType {
  language: Language
  lang: Language
  setLanguage: (lang: Language) => void
  t: (key: string) => string
}

const translations = {
  'hero.tagline': 'AI CODING ASSISTANT · PURE LOCAL',
  'hero.subtitle1': 'Local AI Coding Assistant',
  'hero.subtitle2': 'Smart Chat · Code Generation · Project Understanding',
  'hero.subtitle3': 'Streamlined from Claude Code',
  'hero.cta.start': 'GET STARTED →',
  'hero.cta.learn': 'LEARN MORE',
  'hero.cta.docs': 'DOCS →',
  'hero.stats.chat': 'AI CHAT',
  'hero.stats.code': 'CODE GEN',
  'hero.stats.files': 'FILE OPS',
  'hero.stats.git': 'GIT INTEGRATION',

  'features.title': 'CORE FEATURES',
  'features.heading1': 'WHAT MINICLAUDE',
  'features.heading2': 'CAN DO',
  'features.ai.title': 'AI CHAT',
  'features.ai.metric': 'SMART',
  'features.ai.unit': 'INTERACTION',
  'features.ai.desc': 'Natural-language interaction, context understanding, and multi-model switching (Sonnet/Opus/Haiku).',
  'features.code.title': 'CODE GEN',
  'features.code.metric': 'AUTO',
  'features.code.unit': 'CODING',
  'features.code.desc': 'Create, modify, and refactor code with project-aware analysis and intelligent completion.',
  'features.files.title': 'FILE OPS',
  'features.files.metric': 'COMPLETE',
  'features.files.unit': 'TOOLSET',
  'features.files.desc': 'Read, write, and edit files with syntax highlighting, file search, and ripgrep-powered content lookup.',
  'features.dev.title': 'DEV INTEGRATION',
  'features.dev.metric': 'FULL',
  'features.dev.unit': 'SUPPORT',
  'features.dev.desc': 'Shell command execution, Git integration, LSP language services, and MCP protocol support.',
  'features.tech.runtime': 'RUNTIME',
  'features.tech.language': 'LANGUAGE',
  'features.tech.protocol': 'PROTOCOL',

  'comparison.title': 'COMPARISON',
  'comparison.heading1': 'VS CLAUDE',
  'comparison.heading2': 'CODE',
  'comparison.table.feature': 'FEATURE',
  'comparison.features.core': 'Core Development Features',
  'comparison.features.cloud': 'Cloud Service Integration',
  'comparison.features.telemetry': 'Telemetry Data Reporting',
  'comparison.features.sync': 'Settings Sync',
  'comparison.features.collab': 'Team Collaboration',
  'comparison.features.local': '100% Local',
  'comparison.features.lightweight': 'Lightweight Deployment',
  'comparison.features.simple': 'Simplified Commands',
  'comparison.stats.lines': 'LINES REMOVED',
  'comparison.stats.files': 'FILES DELETED',
  'comparison.stats.commands': 'COMMANDS CUT',

  'install.title': 'INSTALLATION',
  'install.heading1': 'GET',
  'install.heading2': 'STARTED',
  'install.platform.windows': 'WINDOWS',
  'install.platform.macos': 'MACOS',
  'install.platform.linux': 'LINUX',
  'install.copy': 'COPY',
  'install.copied': 'COPIED',
  'install.req.runtime': 'RUNTIME',
  'install.req.system': 'SYSTEM',
  'install.req.apikey': 'API KEY',

  'nav.guide': 'Guide',
  'nav.features': 'Features',
  'nav.reference': 'Reference',
  'search.placeholder': 'Search docs...',
  'search.empty': 'Type to search documentation',
  'search.noResults': 'No results found',
  'search.tip': 'Ctrl+K',

  'footer.desc': 'Lightweight local AI coding assistant\nStreamlined from Claude Code',
  'footer.links': 'LINKS',
  'footer.based': 'BASED ON',
  'footer.copyright': '© 2026 MINICLAUDE · NOT AFFILIATED WITH ANTHROPIC',
  'footer.disclaimer': 'This project is not an official Anthropic project and is not authorized or endorsed by Anthropic.\nUse it at your own risk for learning and research purposes only.',
} as const

const LanguageContext = createContext<LanguageContextType | undefined>(undefined)

export function LanguageProvider({ children }: { children: ReactNode }) {
  const language: Language = 'en'

  const setLanguage = (_lang: Language) => {
    return
  }

  const t = (key: string): string => {
    return translations[key as keyof typeof translations] || key
  }

  return (
    <LanguageContext.Provider value={{ language, lang: language, setLanguage, t }}>
      {children}
    </LanguageContext.Provider>
  )
}

export function useLanguage() {
  const context = useContext(LanguageContext)
  if (!context) {
    throw new Error('useLanguage must be used within LanguageProvider')
  }
  return context
}
