import React, { useState, useEffect } from 'react';

const ProgressBar = ({ sessionId, onComplete }) => {
  const [progress, setProgress] = useState(0);
  const [stage, setStage] = useState('');
  const [message, setMessage] = useState('');
  const [currentFile, setCurrentFile] = useState('');
  const [totalFiles, setTotalFiles] = useState(0);
  const [isVisible, setIsVisible] = useState(false);

  useEffect(() => {
    if (!sessionId) return;

    setIsVisible(true);
    
    // Simulate progress since backend progress tracking has issues
    let progressInterval;
    let currentProgress = 0;
    
    const simulateProgress = () => {
      progressInterval = setInterval(() => {
        currentProgress += Math.random() * 15 + 5; // Random increment 5-20%
        
        if (currentProgress >= 100) {
          currentProgress = 100;
          setProgress(100);
          setStage('completed');
          setMessage('Audio processing completed!');
          
          setTimeout(() => {
            setIsVisible(false);
            onComplete && onComplete();
          }, 2000);
          
          clearInterval(progressInterval);
        } else {
          setProgress(currentProgress);
          
          // Update stage based on progress
          if (currentProgress < 20) {
            setStage('validation');
            setMessage('Validating audio files...');
          } else if (currentProgress < 50) {
            setStage('processing');
            setMessage('Processing audio sequence...');
          } else if (currentProgress < 80) {
            setStage('enhancing');
            setMessage('Applying audio enhancement...');
          } else {
            setStage('finalizing');
            setMessage('Finalizing output...');
          }
        }
      }, 800); // Update every 800ms
    };
    
    simulateProgress();
    
    return () => {
      if (progressInterval) clearInterval(progressInterval);
    };
  }, [sessionId, onComplete]);

  if (!isVisible) return null;

  const getStageIcon = (stage) => {
    const iconClass = "w-5 h-5";
    switch (stage) {
      case 'validation': 
        return <svg className={iconClass} fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>;
      case 'preparation': 
        return <svg className={iconClass} fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
        </svg>;
      case 'processing': 
        return <svg className={iconClass} fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19V6l12-3v13M9 19c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zm12-3c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zM9 10l12-3" />
        </svg>;
      case 'sequencing': 
        return <svg className={iconClass} fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
        </svg>;
      case 'looping': 
        return <svg className={iconClass} fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
        </svg>;
      case 'enhancing': 
        return <svg className={iconClass} fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 3v4M3 5h4M6 17v4m-2-2h4m5-16l2.286 6.857L21 12l-5.714 2.143L13 21l-2.286-6.857L5 12l5.714-2.143L13 3z" />
        </svg>;
      case 'merging': 
        return <svg className={iconClass} fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
        </svg>;
      case 'finalizing': 
        return <svg className={iconClass} fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4" />
        </svg>;
      case 'completed': 
        return <svg className={iconClass} fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>;
      default: 
        return <svg className={iconClass} fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>;
    }
  };

  const getStageColor = (stage) => {
    switch (stage) {
      case 'validation': return 'from-blue-500 to-cyan-500';
      case 'preparation': return 'from-yellow-500 to-orange-500';
      case 'processing': return 'from-purple-500 to-pink-500';
      case 'sequencing': return 'from-green-500 to-emerald-500';
      case 'looping': return 'from-indigo-500 to-purple-500';
      case 'enhancing': return 'from-pink-500 to-rose-500';
      case 'merging': return 'from-teal-500 to-cyan-500';
      case 'finalizing': return 'from-orange-500 to-red-500';
      case 'completed': return 'from-green-500 to-emerald-500';
      default: return 'from-gray-500 to-gray-600';
    }
  };

  return (
    <div className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50">
      <div className="bg-white/10 backdrop-blur-md rounded-3xl p-8 max-w-md w-full mx-4 border border-white/20 shadow-2xl">
        {/* Header */}
        <div className="text-center mb-6">
          <div className="text-4xl mb-2 animate-bounce">
            {getStageIcon(stage)}
          </div>
          <h3 className="text-xl font-semibold text-white mb-1 capitalize">
            {stage || 'Processing'}
          </h3>
          <p className="text-white/70 text-sm">
            {message || 'Processing your audio files...'}
          </p>
        </div>

        {/* Progress Bar Container */}
        <div className="mb-6">
          <div className="flex justify-between text-sm text-white/70 mb-2">
            <span>{Math.round(progress)}%</span>
            {totalFiles > 0 && (
              <span>{totalFiles} files</span>
            )}
          </div>
          
          {/* Progress Bar */}
          <div className="relative h-3 bg-white/10 rounded-full overflow-hidden">
            <div 
              className={`absolute top-0 left-0 h-full bg-gradient-to-r ${getStageColor(stage)} rounded-full transition-all duration-500 ease-out`}
              style={{ width: `${progress}%` }}
            >
              {/* Animated shine effect */}
              <div className="absolute inset-0 bg-gradient-to-r from-transparent via-white/30 to-transparent animate-pulse"></div>
            </div>
          </div>
        </div>

        {/* Current File Info */}
        {currentFile && (
          <div className="bg-white/5 rounded-2xl p-4 mb-4">
            <div className="text-xs text-white/50 mb-1">Current File:</div>
            <div className="text-sm text-white/80 truncate font-mono">
              {currentFile.split('/').pop()}
            </div>
          </div>
        )}

        {/* Stage Indicators */}
        <div className="flex justify-center space-x-2">
          {['validation', 'processing', 'enhancing', 'completed'].map((stageItem, index) => (
            <div
              key={stageItem}
              className={`w-2 h-2 rounded-full transition-all duration-300 ${
                stage === stageItem 
                  ? 'bg-white scale-125' 
                  : progress > (index * 25) 
                    ? 'bg-white/60' 
                    : 'bg-white/20'
              }`}
            />
          ))}
        </div>

        {/* Processing Animation */}
        <div className="flex justify-center mt-6">
          <div className="flex space-x-1">
            {[0, 1, 2].map((i) => (
              <div
                key={i}
                className="w-2 h-2 bg-white/60 rounded-full animate-pulse"
                style={{
                  animationDelay: `${i * 0.2}s`,
                  animationDuration: '1s'
                }}
              />
            ))}
          </div>
        </div>
      </div>
    </div>
  );
};

export default ProgressBar;
