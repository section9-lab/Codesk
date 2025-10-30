import plugin from 'tailwindcss/plugin'

/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        // 主色调定义
        background: 'oklch(var(--color-background) / <alpha-value>)',
        foreground: 'oklch(var(--color-foreground) / <alpha-value>)',
        card: 'oklch(var(--color-card) / <alpha-value>)',
        'card-foreground': 'oklch(var(--color-card-foreground) / <alpha-value>)',
        popover: 'oklch(var(--color-popover) / <alpha-value>)',
        'popover-foreground': 'oklch(var(--color-popover-foreground) / <alpha-value>)',
        primary: 'oklch(var(--color-primary) / <alpha-value>)',
        'primary-foreground': 'oklch(var(--color-primary-foreground) / <alpha-value>)',
        secondary: 'oklch(var(--color-secondary) / <alpha-value>)',
        'secondary-foreground': 'oklch(var(--color-secondary-foreground) / <alpha-value>)',
        muted: 'oklch(var(--color-muted) / <alpha-value>)',
        'muted-foreground': 'oklch(var(--color-muted-foreground) / <alpha-value>)',
        accent: 'oklch(var(--color-accent) / <alpha-value>)',
        'accent-foreground': 'oklch(var(--color-accent-foreground) / <alpha-value>)',
        destructive: 'oklch(var(--color-destructive) / <alpha-value>)',
        'destructive-foreground': 'oklch(var(--color-destructive-foreground) / <alpha-value>)',
        border: 'oklch(var(--color-border) / <alpha-value>)',
        input: 'oklch(var(--color-input) / <alpha-value>)',
        ring: 'oklch(var(--color-ring) / <alpha-value>)',
        
        // 额外的绿色调
        'green-500': 'oklch(var(--color-green-500) / <alpha-value>)',
        'green-600': 'oklch(var(--color-green-600) / <alpha-value>)',
      },
      borderRadius: {
        sm: 'var(--radius-sm)',
        DEFAULT: 'var(--radius-base)',
        md: 'var(--radius-md)',
        lg: 'var(--radius-lg)',
        xl: 'var(--radius-xl)',
      },
      fontFamily: {
        sans: [
          'var(--font-sans)',
          '-apple-system',
          'BlinkMacSystemFont',
          '"Segoe UI"',
          '"Roboto"',
          '"Oxygen"',
          '"Ubuntu"',
          '"Cantarell"',
          '"Fira Sans"',
          '"Droid Sans"',
          '"Helvetica Neue"',
          'sans-serif'
        ],
        mono: [
          'ui-monospace',
          'SFMono-Regular',
          '"SF Mono"',
          'Consolas',
          '"Liberation Mono"',
          'Menlo',
          'monospace'
        ],
      },
      fontSize: {
        xs: 'var(--text-xs)',
        sm: 'var(--text-sm)',
        base: 'var(--text-base)',
        lg: 'var(--text-lg)',
        xl: 'var(--text-xl)',
        '2xl': 'var(--text-2xl)',
        '3xl': 'var(--text-3xl)',
        '4xl': 'var(--text-4xl)',
        '5xl': 'var(--text-5xl)',
      },
      lineHeight: {
        none: 'var(--leading-none)',
        tight: 'var(--leading-tight)',
        snug: 'var(--leading-snug)',
        normal: 'var(--leading-normal)',
        relaxed: 'var(--leading-relaxed)',
        loose: 'var(--leading-loose)',
      },
      letterSpacing: {
        tighter: 'var(--tracking-tighter)',
        tight: 'var(--tracking-tight)',
        normal: 'var(--tracking-normal)',
        wide: 'var(--tracking-wide)',
        wider: 'var(--tracking-wider)',
        widest: 'var(--tracking-widest)',
      },
      fontWeight: {
        thin: 'var(--font-weight-thin)',
        extralight: 'var(--font-weight-extralight)',
        light: 'var(--font-weight-light)',
        normal: 'var(--font-weight-normal)',
        medium: 'var(--font-weight-medium)',
        semibold: 'var(--font-weight-semibold)',
        bold: 'var(--font-weight-bold)',
        extrabold: 'var(--font-weight-extrabold)',
        black: 'var(--font-weight-black)',
      },
      animation: {
        'shutter-flash': 'shutterFlash 0.5s ease-in-out',
        'move-to-input': 'moveToInput 0.8s ease-in-out forwards',
        'scanlines': 'scanlines 8s linear infinite',
      },
      keyframes: {
        shutterFlash: {
          '0%': { opacity: '0' },
          '50%': { opacity: '1' },
          '100%': { opacity: '0' },
        },
        moveToInput: {
          '0%': { 
            transform: 'scale(1) translateY(0)',
            opacity: '1',
          },
          '50%': { 
            transform: 'scale(0.3) translateY(50%)',
            opacity: '0.8',
          },
          '100%': { 
            transform: 'scale(0.1) translateY(100vh)',
            opacity: '0',
          },
        },
        scanlines: {
          '0%': { transform: 'translateY(-100%)' },
          '100%': { transform: 'translateY(100%)' },
        },
      },
    },
  },
  plugins: [
    plugin(function({ addUtilities, addComponents }) {
      // 添加自定义工具类
      addUtilities({
        '.scrollbar-hide': {
          'scrollbar-width': 'none',
          '-ms-overflow-style': 'none',
          '&::-webkit-scrollbar': {
            display: 'none',
          },
        },
        '.scrollbar-thin': {
          'scrollbar-width': 'thin',
          'scrollbar-color': 'var(--color-border) transparent',
          '&::-webkit-scrollbar': {
            width: '6px',
            height: '6px',
          },
          '&::-webkit-scrollbar-track': {
            background: 'transparent',
          },
          '&::-webkit-scrollbar-thumb': {
            'background-color': 'var(--color-border)',
            'border-radius': '3px',
          },
          '&::-webkit-scrollbar-thumb:hover': {
            'background-color': 'var(--color-muted-foreground)',
          },
        },
        '.line-clamp-2': {
          overflow: 'hidden',
          display: '-webkit-box',
          '-webkit-box-orient': 'vertical',
          '-webkit-line-clamp': '2',
        },
        '.animate-in': {
          'animation-name': 'enter',
          'animation-duration': '150ms',
          'animation-fill-mode': 'both',
        },
        '.animate-out': {
          'animation-name': 'exit',
          'animation-duration': '150ms',
          'animation-fill-mode': 'both',
        },
        '.shimmer-hover': {
          position: 'relative',
          overflow: 'hidden',
          '&::before': {
            content: '""',
            position: 'absolute',
            top: '0',
            left: '-100%',
            width: '100%',
            height: '100%',
            background: 'linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.05), transparent)',
            transition: 'left 0.5s',
          },
          '&:hover::before': {
            left: '100%',
            animation: 'shimmer 0.5s',
          },
        },
        '.trailing-border': {
          position: 'relative',
          background: 'var(--color-card)',
          'z-index': '0',
          overflow: 'visible',
          '&::after': {
            content: '""',
            position: 'absolute',
            inset: '-2px',
            padding: '2px',
            'border-radius': 'inherit',
            background: 'conic-gradient(from var(--angle, 0deg), transparent 0%, transparent 85%, #d97757 90%, #ff9a7a 92.5%, #d97757 95%, transparent 100%)',
            '-webkit-mask': 'linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0)',
            '-webkit-mask-composite': 'xor',
            'mask-composite': 'exclude',
            opacity: '0',
            transition: 'opacity 0.3s ease',
            'z-index': '-1',
          },
          '&:hover::after': {
            opacity: '1',
            animation: 'trail-rotate 2s linear infinite',
          },
          '& > *': {
            position: 'relative',
            'z-index': '1',
          },
        },
      })

      // 添加自定义组件
      addComponents({
        '.text-display-1': {
          'font-size': 'var(--text-5xl)',
          'font-weight': 'var(--font-weight-bold)',
          'line-height': 'var(--leading-tight)',
          'letter-spacing': 'var(--tracking-tight)',
        },
        '.text-display-2': {
          'font-size': 'var(--text-4xl)',
          'font-weight': 'var(--font-weight-bold)',
          'line-height': 'var(--leading-tight)',
          'letter-spacing': 'var(--tracking-tight)',
        },
        '.text-heading-1': {
          'font-size': 'var(--text-3xl)',
          'font-weight': 'var(--font-weight-semibold)',
          'line-height': 'var(--leading-tight)',
          'letter-spacing': 'var(--tracking-tight)',
        },
        '.text-heading-2': {
          'font-size': 'var(--text-2xl)',
          'font-weight': 'var(--font-weight-semibold)',
          'line-height': 'var(--leading-snug)',
        },
        '.text-heading-3': {
          'font-size': 'var(--text-xl)',
          'font-weight': 'var(--font-weight-semibold)',
          'line-height': 'var(--leading-snug)',
        },
        '.text-heading-4': {
          'font-size': 'var(--text-lg)',
          'font-weight': 'var(--font-weight-medium)',
          'line-height': 'var(--leading-normal)',
        },
        '.text-body-large': {
          'font-size': 'var(--text-lg)',
          'font-weight': 'var(--font-weight-normal)',
          'line-height': 'var(--leading-relaxed)',
        },
        '.text-body': {
          'font-size': 'var(--text-base)',
          'font-weight': 'var(--font-weight-normal)',
          'line-height': 'var(--leading-normal)',
        },
        '.text-body-small': {
          'font-size': 'var(--text-sm)',
          'font-weight': 'var(--font-weight-normal)',
          'line-height': 'var(--leading-normal)',
        },
        '.text-caption': {
          'font-size': 'var(--text-xs)',
          'font-weight': 'var(--font-weight-normal)',
          'line-height': 'var(--leading-normal)',
        },
        '.text-label': {
          'font-size': 'var(--text-sm)',
          'font-weight': 'var(--font-weight-medium)',
          'line-height': 'var(--leading-tight)',
          'letter-spacing': 'var(--tracking-wide)',
        },
        '.text-button': {
          'font-size': 'var(--text-sm)',
          'font-weight': 'var(--font-weight-medium)',
          'line-height': 'var(--leading-tight)',
          'letter-spacing': 'var(--tracking-wide)',
        },
        '.text-overline': {
          'font-size': 'var(--text-xs)',
          'font-weight': 'var(--font-weight-semibold)',
          'line-height': 'var(--leading-tight)',
          'letter-spacing': 'var(--tracking-wider)',
          'text-transform': 'uppercase',
        },
        '.rotating-symbol': {
          display: 'inline-block',
          'vertical-align': 'middle',
          'line-height': '1',
          animation: 'fade-in 0.2s ease-out',
          'font-weight': 'normal',
          'font-size': '1.5rem',
          position: 'relative',
          top: '-2px',
          '&::before': {
            content: '"◐"',
            animation: 'rotate-symbol 1.6s steps(4, end) infinite',
            display: 'inline-block',
            'font-size': 'inherit',
            'line-height': '1',
            'vertical-align': 'baseline',
            'transform-origin': 'center',
          },
        },
      })

      // 添加关键帧动画
      addUtilities({
        '@keyframes rotate-symbol': {
          '0%': { content: '"◐"', transform: 'scale(1)' },
          '25%': { content: '"◓"', transform: 'scale(1)' },
          '50%': { content: '"◑"', transform: 'scale(1)' },
          '75%': { content: '"◒"', transform: 'scale(1)' },
          '100%': { content: '"◐"', transform: 'scale(1)' },
        },
        '@keyframes fade-in': {
          from: {
            opacity: '0',
            transform: 'scale(0.8)',
          },
          to: {
            opacity: '1',
            transform: 'scale(1)',
          },
        },
        '@keyframes shimmer': {
          '0%': {
            'background-position': '-200% 0',
          },
          '100%': {
            'background-position': '200% 0',
          },
        },
        '@keyframes trail-rotate': {
          to: {
            '--angle': '360deg',
          },
        },
      })
    })
  ],
}