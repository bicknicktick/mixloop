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
            Developed by <span className="text-blue-400 font-semibold">BITZY.ID</span>
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
            <div className="w-2 h-2 bg-green-500 rounded-full animate-pulse"></div>
            <span>Upload Audio</span>
          </div>
          <div className="hidden sm:block w-px h-4 bg-gray-700"></div>
          <div className="flex items-center space-x-2">
            <div className="w-2 h-2 bg-blue-500 rounded-full animate-pulse"></div>
            <span>Set Loops</span>
          </div>
          <div className="hidden sm:block w-px h-4 bg-gray-700"></div>
          <div className="flex items-center space-x-2">
            <div className="w-2 h-2 bg-purple-500 rounded-full animate-pulse"></div>
            <span>Export Mix</span>
          </div>
        </div>
      </div>
    </div>
  )
}

export default Landing
