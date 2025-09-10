import React, { useState, useRef } from 'react';
import axios from 'axios';
import ProgressBar from './ProgressBar';

function Workspace({ onBack }) {
  const [audioFiles, setAudioFiles] = useState([])
  const [loops, setLoops] = useState(2)
  const [crossfade, setCrossfade] = useState(2)
  const [isProcessing, setIsProcessing] = useState(false)
  const [resultUrl, setResultUrl] = useState(null)
  const [error, setError] = useState(null)
  const [dragOver, setDragOver] = useState(false)
  const [sessionId, setSessionId] = useState(null)
  const [enhance, setEnhance] = useState(true)
  const [dolbyStereo, setDolbyStereo] = useState(false)
  const [format, setFormat] = useState('mp3')
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

  const handleMixAudio = async () => {
    if (audioFiles.length === 0) {
      setError('Please select audio files first')
      return
    }

    setIsProcessing(true)
    setError(null)
    setResultUrl(null)
    
    // Generate session ID for progress tracking
    const newSessionId = Date.now().toString()
    setSessionId(newSessionId)

    const formData = new FormData()
    audioFiles.forEach((file, index) => {
      formData.append('audio_files', file)
    })
    formData.append('loops', loops.toString())
    formData.append('crossfade', crossfade.toString())
    formData.append('enhance', enhance.toString())
    formData.append('dolby_stereo', dolbyStereo.toString())
    formData.append('format', format)
    formData.append('session_id', newSessionId)

    try {
      const response = await axios.post('http://localhost:8081/api/mix', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
        responseType: 'blob',
      })

      const contentType = format === 'wav' ? 'audio/wav' : 'audio/mpeg'
      const blob = new Blob([response.data], { type: contentType })
      const url = URL.createObjectURL(blob)
      setResultUrl(url)
    } catch (err) {
      setError('Failed to process audio files. Please try again.')
      console.error('Error:', err)
    } finally {
      setIsProcessing(false)
      setSessionId(null)
    }
  }

  const handleProgressComplete = () => {
    setSessionId(null)
  }

  const handleDownload = () => {
    if (resultUrl && audioFiles.length > 0) {
      const firstName = audioFiles[0].name.replace(/\.[^/.]+$/, "") // Remove extension
      const lastName = audioFiles[audioFiles.length - 1].name.replace(/\.[^/.]+$/, "") // Remove extension
      
      let filename
      if (audioFiles.length === 1) {
        filename = `${firstName}_bitzy.id.${format}`
      } else {
        filename = `${firstName}_${lastName}_bitzy.id.${format}`
      }
      
      const a = document.createElement('a')
      a.href = resultUrl
      a.download = filename
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
    <div className="min-h-screen bg-gradient-to-br from-black via-gray-900 to-black text-white p-2 sm:p-4">
      <div className="glass-card p-3 sm:p-4 w-full max-w-6xl mx-auto">
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between mb-3 sm:mb-4">
          <div>
            <h2 className="text-xl sm:text-2xl font-bold text-white mb-1">Audio Workspace</h2>
            <div className="text-xs text-gray-500">
              Powered by <span className="text-blue-400 font-semibold">BITZY.ID</span>
            </div>
          </div>
          <button
            onClick={onBack}
            className="text-gray-400 hover:text-white transition-colors text-sm self-start sm:self-auto mt-2 sm:mt-0"
          >
            ← Back to Home
          </button>
        </div>

        <div className="space-y-3 sm:space-y-4">
          {/* File Upload Section */}
          <div className="glass-card p-3 sm:p-4">
            <label className="block mb-3 sm:mb-4 text-base sm:text-lg font-semibold flex items-center justify-center space-x-2">
              <div className="w-2 h-2 bg-green-400 rounded-full animate-pulse shadow-lg shadow-green-400/50"></div>
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
              <div className="p-4 sm:p-6 text-center space-y-3 sm:space-y-4">
                <div className="relative">
                  <svg className="w-16 h-16 sm:w-20 sm:h-20 mx-auto text-white/80" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M9 19V6l12-3v13M9 19c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zm12-3c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zM9 10l12-3" />
                  </svg>
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
              <div className="mt-3 sm:mt-4 space-y-2 sm:space-y-3">
                <h3 className="text-sm sm:text-base font-semibold text-white/90 flex items-center space-x-2">
                  <span className="w-2 h-2 bg-blue-400 rounded-full"></span>
                  <span>Selected Files ({audioFiles.length})</span>
                </h3>
                <div className="space-y-1.5 sm:space-y-2 max-h-32 sm:max-h-40 overflow-y-auto scrollbar-thin scrollbar-thumb-white/20 scrollbar-track-transparent">
                  {audioFiles.map((file, index) => (
                    <div key={index} className="flex items-center justify-between p-2 sm:p-3 bg-white/5 rounded-lg border border-white/10 hover:bg-white/10 transition-colors group">
                      <div className="flex items-center space-x-2 sm:space-x-3 flex-1 min-w-0">
                        <div className="w-6 h-6 sm:w-8 sm:h-8 bg-gradient-to-br from-purple-500 to-pink-500 rounded-md flex items-center justify-center flex-shrink-0">
                          <span className="text-white text-xs font-bold">{index + 1}</span>
                        </div>
                        <div className="flex-1 min-w-0">
                          <p className="text-white font-medium truncate text-xs sm:text-sm">{file.name}</p>
                          <p className="text-white/60 text-xs">{(file.size / 1024 / 1024).toFixed(2)} MB</p>
                        </div>
                      </div>
                      <button
                        onClick={() => removeFile(index)}
                        className="ml-2 p-1 text-red-400 hover:text-red-300 hover:bg-red-400/10 rounded-md transition-colors opacity-0 group-hover:opacity-100 flex-shrink-0"
                      >
                        <svg className="w-3 h-3 sm:w-4 sm:h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                        </svg>
                      </button>
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>

          {/* Controls Section */}
          <div className="glass-card p-3 sm:p-4">
            <div className="grid grid-cols-2 lg:grid-cols-5 gap-3 sm:gap-4">
              <div className="space-y-2 sm:space-y-3">
                <label className="block text-xs sm:text-sm font-semibold text-white/90 flex items-center space-x-1.5">
                  <span className="w-1.5 h-1.5 bg-purple-400 rounded-full"></span>
                  <span>Loops</span>
                </label>
                <input
                  type="number"
                  min="1"
                  max="10"
                  value={loops}
                  onChange={(e) => setLoops(parseInt(e.target.value))}
                  className="w-full px-2 sm:px-3 py-1.5 sm:py-2 bg-white/10 border border-white/20 rounded-lg text-white placeholder-white/50 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent transition-all text-xs sm:text-sm"
                />
              </div>
              <div className="space-y-2 sm:space-y-3">
                <label className="block text-xs sm:text-sm font-semibold text-white/90 flex items-center space-x-1.5">
                  <span className="w-1.5 h-1.5 bg-blue-400 rounded-full"></span>
                  <span>Crossfade (sec)</span>
                </label>
                <input
                  type="number"
                  min="0"
                  max="10"
                  step="0.1"
                  value={crossfade}
                  onChange={(e) => setCrossfade(parseFloat(e.target.value))}
                  className="w-full px-2 sm:px-3 py-1.5 sm:py-2 bg-white/10 border border-white/20 rounded-lg text-white placeholder-white/50 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all text-xs sm:text-sm"
                />
              </div>
              <div className="space-y-2 sm:space-y-3">
                <label className="block text-xs sm:text-sm font-semibold text-white/90 flex items-center space-x-1.5">
                  <span className="w-1.5 h-1.5 bg-green-400 rounded-full"></span>
                  <span>Enhancement</span>
                </label>
                <button
                  onClick={() => setEnhance(!enhance)}
                  className={`w-full px-2 sm:px-3 py-1.5 sm:py-2 rounded-lg font-medium transition-all text-xs sm:text-sm ${
                    enhance
                      ? 'bg-green-500 text-white shadow-lg shadow-green-500/25'
                      : 'bg-white/10 text-white/70 border border-white/20 hover:bg-white/20'
                  }`}
                >
                  {enhance ? '✓ Enabled' : 'Disabled'}
                </button>
              </div>
              <div className="space-y-2 sm:space-y-3">
                <label className="block text-xs sm:text-sm font-semibold text-white/90 flex items-center space-x-1.5">
                  <span className="w-1.5 h-1.5 bg-pink-400 rounded-full"></span>
                  <span>Dolby Stereo</span>
                </label>
                <button
                  onClick={() => setDolbyStereo(!dolbyStereo)}
                  className={`w-full px-2 sm:px-3 py-1.5 sm:py-2 rounded-lg font-medium transition-all text-xs sm:text-sm ${
                    dolbyStereo
                      ? 'bg-pink-500 text-white shadow-lg shadow-pink-500/25'
                      : 'bg-white/10 text-white/70 border border-white/20 hover:bg-white/20'
                  }`}
                >
                  {dolbyStereo ? '✓ Enabled' : 'Disabled'}
                </button>
              </div>
              <div className="space-y-2 sm:space-y-3">
                <label className="block text-xs sm:text-sm font-semibold text-white/90 flex items-center space-x-1.5">
                  <span className="w-1.5 h-1.5 bg-orange-400 rounded-full"></span>
                  <span>Output Format</span>
                </label>
                <select
                  value={format}
                  onChange={(e) => setFormat(e.target.value)}
                  className="w-full px-2 sm:px-3 py-1.5 sm:py-2 bg-white/10 border border-white/20 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-orange-500 focus:border-transparent transition-all text-xs sm:text-sm"
                >
                  <option value="mp3" className="bg-gray-800">MP3 (320k)</option>
                  <option value="wav" className="bg-gray-800">WAV (24-bit)</option>
                </select>
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

          {/* Process Button */}
          <div className="text-center">
            <button
              onClick={handleMixAudio}
              disabled={audioFiles.length === 0 || isProcessing}
              className={`px-6 sm:px-8 py-2.5 sm:py-3 rounded-xl font-bold text-sm sm:text-base transition-all duration-300 ${
                audioFiles.length === 0 || isProcessing
                  ? 'bg-gray-600 text-gray-400 cursor-not-allowed'
                  : 'bg-gradient-to-r from-purple-600 to-pink-600 text-white hover:from-purple-700 hover:to-pink-700 shadow-lg hover:shadow-xl transform hover:scale-105'
              }`}
            >
              {isProcessing ? (
                <div className="flex items-center space-x-2">
                  <div className="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
                  <span>Processing...</span>
                </div>
              ) : (
                <div className="flex items-center space-x-2">
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19V6l12-3v13M9 19c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zm12-3c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zM9 10l12-3" />
                </svg>
                <span>Mix {audioFiles.length} Audio Files</span>
              </div>
              )}
            </button>
          </div>

          {/* Result Section */}
          {resultUrl && (
            <div className="glass-card p-4 sm:p-6 space-y-4 sm:space-y-6 animate-slide-up">
              <div className="flex items-center space-x-2">
                <svg className="w-5 h-5 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
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
                <div className="flex items-center space-x-2">
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                  </svg>
                  <span>Download Mix</span>
                </div>
              </button>
            </div>
          )}
        </div>

        {/* Progress Bar */}
        {sessionId && (
          <ProgressBar 
            sessionId={sessionId} 
            onComplete={handleProgressComplete}
          />
        )}
      </div>
    </div>
  )
}

export default Workspace
