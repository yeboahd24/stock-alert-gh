import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'

import App from './App.jsx'
import { PostHogProvider } from 'posthog-js/react'

const posthogKey = import.meta.env.VITE_PUBLIC_POSTHOG_KEY
const posthogHost = import.meta.env.VITE_PUBLIC_POSTHOG_HOST

const options = {
  api_host: posthogHost || 'https://eu.i.posthog.com',
}

const AppWithPostHog = () => {
  if (posthogKey) {
    return (
      <PostHogProvider apiKey={posthogKey} options={options}>
        <App />
      </PostHogProvider>
    )
  }
  return <App />
}

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <AppWithPostHog />
  </StrictMode>,
)