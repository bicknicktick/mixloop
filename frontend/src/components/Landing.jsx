import React from 'react'

function Landing({ onStart }) {
  return (
    <div className="flex flex-col items-center justify-center min-h-screen animate-fade-in px-4 sm:px-6 lg:px-8">
      <div className="text-center space-y-8 sm:space-y-12 max-w-lg mx-auto">
        <div className="space-y-4">
          <h1 className="text-6xl md:text-8xl font-black tracking-tight text-white">
            MixLoop
          </h1>
          <p className="text-xl md:text-2xl text-gray-300 font-light max-w-2xl mx-auto">
            Create seamless audio sequences with crossfade transitions
          </p>
          <div className="text-sm text-gray-500 font-medium">
            ";)" <span className="text-blue-400 font-semibold">BITZY.ID</span>
          </div>
        </div>
        
        <div className="pt-2 sm:pt-4">
          <button
            onClick={onStart}
            className="modern-button text-base sm:text-lg px-12 sm:px-16 py-4 sm:py-5 w-full sm:w-auto"
          >
            Start Mixing
          </button>
        </div>
        
        <div className="flex flex-col sm:flex-row items-center justify-center space-y-4 sm:space-y-0 sm:space-x-6 lg:space-x-8 pt-6 sm:pt-8 text-xs sm:text-sm text-gray-500">
          <div className="flex items-center space-x-2">
            <svg className="w-4 h-4 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
            </svg>
            <span>Upload Audio</span>
          </div>
          <div className="hidden sm:block w-px h-4 bg-gray-700"></div>
          <div className="flex items-center space-x-2">
            <svg className="w-4 h-4 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
            <span>Set Loops</span>
          </div>
          <div className="hidden sm:block w-px h-4 bg-gray-700"></div>
          <div className="flex items-center space-x-2">
            <svg className="w-4 h-4 text-purple-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
            <span>Export Mix</span>
          </div>
        </div>
      </div>
    </div>
  )
}

export default Landing
