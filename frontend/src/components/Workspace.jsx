import React, { useState, useRef } from 'react'
import axios from 'axios'

function Workspace({ onBack }) {
  const [audioFiles, setAudioFiles] = useState([])
  const [loopCount, setLoopCount] = useState(2)
  const [crossfadeDuration, setCrossfadeDuration] = useState(2)
  const [isProcessing, setIsProcessing] = useState(false)
  const [resultUrl, setResultUrl] = useState(null)
  const [error, setError] = useState(null)
  const [dragOver, setDragOver] = useState(false)
  const fileInputRef = useRef(null)
  const audioRef = useRef(null)

  const handleFileSelect = (e) => {
    const files = Array.from(e.target.files)
    setAudioFiles(files)
    setError(null)
    setResultUrl(null)
  }

  const handleDragOver = (e) => {
    e.preventDefault()
    setDragOver(true)
  }

  const handleDragLeave = (e) => {
    e.preventDefault()
    setDragOver(false)
  }

  const handleDrop = (e) => {
    e.preventDefault()
    setDragOver(false)
    const files = Array.from(e.dataTransfer.files).filter(file => 
      file.type.includes('audio')
    )
    if (files.length > 0) {
      setAudioFiles(files)
      setError(null)
      setResultUrl(null)
    }
  }

  const handleMix = async () => {
    if (audioFiles.length === 0) {
      setError('Please select at least one audio file')
      return
    }

    setIsProcessing(true)
    setError(null)
    setResultUrl(null)

    const formData = new FormData()
    audioFiles.forEach((file) => {
      formData.append('audio', file)
    })
    formData.append('loops', loopCount)
    formData.append('crossfade', crossfadeDuration)

    try {
      const response = await axios.post('/api/mix', formData, {
        responseType: 'blob',
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      })

      const url = window.URL.createObjectURL(response.data)
      setResultUrl(url)
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to process audio')
    } finally {
      setIsProcessing(false)
    }
  }

  const handleDownload = () => {
    if (resultUrl) {
      const a = document.createElement('a')
      a.href = resultUrl
      a.download = 'mixloop_output.mp3'
      document.body.appendChild(a)
      a.click()
      document.body.removeChild(a)
    }
  }

  const removeFile = (index) => {
    const newFiles = audioFiles.filter((_, i) => i !== index)
    setAudioFiles(newFiles)
  }

  return (
    <div className="min-h-screen p-4 sm:p-6 lg:p-8 animate-fade-in">
      <div className="glass-card p-6 sm:p-8 w-full max-w-4xl mx-auto">
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between mb-6 sm:mb-8">
          <div>
            <h2 className="text-2xl sm:text-3xl font-bold text-white mb-2">Audio Workspace</h2>
            <div className="text-xs text-gray-500">
              Powered by <span className="text-blue-400 font-semibold">BITZY.ID</span>
            </div>
          </div>
          <button
            onClick={onBack}
            className="text-gray-400 hover:text-white transition-colors text-sm sm:text-base self-start sm:self-auto mt-4 sm:mt-0"
          >
            ← Back to Home
          </button>
        </div>

        <div className="space-y-6 sm:space-y-8">
          {/* File Upload Section */}
          <div className="glass-card p-4 sm:p-6 lg:p-8">
            <label className="block mb-6 sm:mb-8 text-lg sm:text-xl font-semibold flex items-center justify-center space-x-3">
              <div className="w-3 h-3 bg-green-400 rounded-full animate-pulse shadow-lg shadow-green-400/50"></div>
              <span>Audio Files</span>
            </label>
            <input
              ref={fileInputRef}
              type="file"
              multiple
              accept="audio/mp3,audio/wav"
              onChange={handleFileSelect}
              className="hidden"
            />
            <div
              onClick={() => fileInputRef.current?.click()}
              onDragOver={handleDragOver}
              onDragLeave={handleDragLeave}
              onDrop={handleDrop}
              className={`relative overflow-hidden rounded-2xl border-2 border-dashed transition-all duration-500 cursor-pointer
                         ${dragOver 
                           ? 'border-white bg-white/10 scale-[1.02]' 
                           : 'border-white/20 bg-white/5 hover:border-white/40 hover:bg-white/8'
                         }`}
            >
              <div className="p-6 sm:p-8 lg:p-12 text-center space-y-4 sm:space-y-6">
                <div className="relative">
                  <div className="text-4xl sm:text-5xl lg:text-6xl opacity-80">🎵</div>
                  <div className="absolute inset-0 bg-gradient-to-t from-transparent to-white/10 rounded-full blur-xl"></div>
                </div>
                <div className="space-y-2">
                  <div className="text-lg sm:text-xl font-semibold px-2">
                    {audioFiles.length === 0
                      ? 'Tap to select audio files'
                      : `${audioFiles.length} file(s) selected`}
                  </div>
                  <div className="text-xs sm:text-sm text-gray-400 font-medium px-2">
                    MP3 and WAV • Max 50MB per file
                  </div>
                </div>
              </div>
              {dragOver && (
                <div className="absolute inset-0 bg-gradient-to-br from-white/20 to-transparent pointer-events-none"></div>
              )}
            </div>
            
            {audioFiles.length > 0 && (
              <div className="mt-6 sm:mt-8 space-y-3 sm:space-y-4">
                {audioFiles.map((file, index) => (
                  <div key={index} className="glass-card p-3 sm:p-4 flex flex-col sm:flex-row sm:items-center sm:justify-between space-y-3 sm:space-y-0">
                    <div className="flex items-center space-x-3 sm:space-x-4 min-w-0 flex-1">
                      <div className="w-3 h-3 bg-gradient-to-r from-blue-400 to-purple-400 rounded-full animate-pulse shadow-lg shadow-blue-400/30 flex-shrink-0"></div>
                      <div className="flex flex-col min-w-0 flex-1">
                        <span className="font-semibold text-white truncate text-sm sm:text-base">{file.name}</span>
                        <span className="text-xs text-gray-400 font-medium">
                          {(file.size / 1024 / 1024).toFixed(1)} MB • Audio File
                        </span>
                      </div>
                    </div>
                    <button
                      onClick={() => removeFile(index)}
                      className="px-3 sm:px-4 py-2 bg-red-500/20 hover:bg-red-500/30 text-red-400 hover:text-red-300 
                               rounded-lg transition-all duration-300 text-xs sm:text-sm font-medium border border-red-500/30 
                               w-full sm:w-auto flex-shrink-0"
                    >
                      Remove
                    </button>
                  </div>
                ))}
              </div>
            )}
          </div>

          {/* Parameters Section */}
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 sm:gap-6 lg:gap-8">
            <div className="glass-card p-4 sm:p-6 lg:p-8">
              <label className="block mb-4 sm:mb-6 lg:mb-8 text-lg sm:text-xl font-semibold flex items-center justify-center space-x-3">
                <div className="w-3 h-3 bg-blue-400 rounded-full animate-pulse shadow-lg shadow-blue-400/50"></div>
                <span>Loop Count</span>
              </label>
              <div className="relative">
                <input
                  type="number"
                  min="1"
                  max="10"
                  value={loopCount}
                  onChange={(e) => setLoopCount(parseInt(e.target.value) || 1)}
                  className="modern-input w-full text-2xl sm:text-3xl font-bold"
                />
                <div className="absolute -bottom-2 left-1/2 transform -translate-x-1/2 w-8 sm:w-12 h-1 bg-gradient-to-r from-transparent via-blue-400 to-transparent opacity-60"></div>
              </div>
              <div className="mt-4 sm:mt-6 text-center text-xs sm:text-sm text-gray-400 font-medium">
                Times to repeat the audio
              </div>
            </div>

            <div className="glass-card p-4 sm:p-6 lg:p-8">
              <label className="block mb-4 sm:mb-6 lg:mb-8 text-lg sm:text-xl font-semibold flex items-center justify-center space-x-3">
                <div className="w-3 h-3 bg-purple-400 rounded-full animate-pulse shadow-lg shadow-purple-400/50"></div>
                <span>Crossfade Duration</span>
              </label>
              <div className="relative">
                <input
                  type="number"
                  min="0"
                  max="10"
                  step="0.5"
                  value={crossfadeDuration}
                  onChange={(e) => setCrossfadeDuration(parseFloat(e.target.value) || 0)}
                  className="modern-input w-full text-2xl sm:text-3xl font-bold"
                />
                <div className="absolute -bottom-2 left-1/2 transform -translate-x-1/2 w-8 sm:w-12 h-1 bg-gradient-to-r from-transparent via-purple-400 to-transparent opacity-60"></div>
              </div>
              <div className="mt-4 sm:mt-6 text-center text-xs sm:text-sm text-gray-400 font-medium">
                Transition time in seconds
              </div>
            </div>
          </div>

          {/* Error Display */}
          {error && (
            <div className="glass-card p-4 border-red-500/50 bg-red-900/20">
              <div className="flex items-center space-x-2">
                <div className="w-2 h-2 bg-red-500 rounded-full"></div>
                <span className="text-red-400">{error}</span>
              </div>
            </div>
          )}

          {/* Mix Button */}
          <div className="relative">
            <button
              onClick={handleMix}
              disabled={isProcessing || audioFiles.length === 0}
              className={`modern-button w-full py-4 sm:py-5 lg:py-6 text-base sm:text-lg font-semibold ${
                isProcessing || audioFiles.length === 0
                  ? 'opacity-50 cursor-not-allowed'
                  : ''
              }`}
            >
              {isProcessing ? (
                <div className="flex flex-col sm:flex-row items-center justify-center space-y-2 sm:space-y-0 sm:space-x-3">
                  <div className="loading-dots scale-75 sm:scale-100">
                    <div></div>
                    <div></div>
                    <div></div>
                    <div></div>
                  </div>
                  <span className="text-sm sm:text-base">Processing Audio...</span>
                </div>
              ) : (
                'Mix & Export'
              )}
            </button>
          </div>

          {/* Result Section */}
          {resultUrl && (
            <div className="glass-card p-4 sm:p-6 space-y-4 sm:space-y-6 animate-slide-up">
              <div className="flex items-center space-x-2">
                <div className="w-2 h-2 bg-green-500 rounded-full animate-pulse"></div>
                <h3 className="text-base sm:text-lg font-medium">Your Mix is Ready</h3>
              </div>
              
              <audio
                ref={audioRef}
                controls
                src={resultUrl}
                className="w-full rounded-lg h-12 sm:h-auto"
              />

              <button
                onClick={handleDownload}
                className="modern-button w-full py-3 sm:py-4 bg-gradient-to-r from-gray-800 to-gray-700 
                         hover:from-gray-700 hover:to-gray-600 text-white text-sm sm:text-base"
              >
                Download Mix
              </button>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

export default Workspace
