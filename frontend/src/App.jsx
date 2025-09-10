import React, { useState } from 'react'
import Landing from './components/Landing'
import Workspace from './components/Workspace'

function App() {
  const [currentView, setCurrentView] = useState('landing')

  return (
    <div className="min-h-screen bg-black text-white">
      {currentView === 'landing' && (
        <Landing onStart={() => setCurrentView('workspace')} />
      )}
      {currentView === 'workspace' && (
        <Workspace onBack={() => setCurrentView('landing')} />
      )}
    </div>
  )
}

export default App
