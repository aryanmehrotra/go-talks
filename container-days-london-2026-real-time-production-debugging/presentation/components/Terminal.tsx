import React, { useState, useEffect, useRef } from 'react';
import { COLORS } from '../constants';

interface TerminalProps {
  children: React.ReactNode;
  title?: string;
  className?: string;
  executable?: boolean;
  output?: string[];
}

export const Terminal: React.FC<TerminalProps> = ({ children, title = "main.go", className = "", executable, output = [] }) => {
  const [isRunning, setIsRunning] = useState(false);
  const [displayedOutput, setDisplayedOutput] = useState<string[]>([]);
  const outputEndRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    setIsRunning(false);
    setDisplayedOutput([]);
  }, [children]);

  useEffect(() => {
    if (outputEndRef.current) {
      outputEndRef.current.scrollIntoView({ behavior: 'auto' });
    }
  }, [displayedOutput]);

  const getTime = () => {
    const now = new Date();
    return now.toTimeString().split(' ')[0];
  };

  const handleRun = async (e: React.MouseEvent) => {
    e.preventDefault();
    if (isRunning) return;
    
    setIsRunning(true);
    setDisplayedOutput([]);
    
    const lines = output.length > 0 ? output : [
      "INFO  [" + getTime() + "] compiling " + title + "...",
      "INFO  [" + getTime() + "] starting server on :8000",
      "INFO  [" + getTime() + "] GoFr initialized"
    ];

    for (let i = 0; i < lines.length; i++) {
      await new Promise(resolve => setTimeout(resolve, 40 + Math.random() * 40));
      const timestampedLine = lines[i].includes('[') ? lines[i] : `INFO  [${getTime()}] ${lines[i]}`;
      setDisplayedOutput(prev => [...prev, timestampedLine]);
    }
    
    setIsRunning(false);
  };

  const renderLogLine = (line: string, idx: number) => {
    try {
      if (!line.includes('[') || !line.includes(']')) {
        return <div key={idx} className="text-gray-300">{line}</div>;
      }

      const type = line.split(' [')[0] || "INFO";
      const timePart = line.split('[')[1]?.split(']')[0] || "00:00:00";
      const message = line.split('] ')[1] || line;

      const isErr = type.includes('ERR') || type.includes('CRIT') || type.includes('ALERT');
      const isSafe = type.includes('SAFE') || type.includes('OK') || type.includes('Stable');

      return (
        <div key={idx} className="flex space-x-4 whitespace-nowrap leading-relaxed">
          <span className={`flex-shrink-0 font-bold ${isErr ? 'text-red-500' : isSafe ? 'text-green-500' : 'text-cyan-500'}`}>
            {type}
          </span>
          <span className="text-gray-700 flex-shrink-0 font-medium opacity-30">
            [{timePart}]
          </span>
          <span className={`tracking-tight ${isErr ? 'text-red-400' : isSafe ? 'text-green-400' : 'text-gray-300'}`}>
            {message}
          </span>
        </div>
      );
    } catch (err) {
      return <div key={idx} className="text-gray-400 font-mono text-[11px]">{line}</div>;
    }
  };

  return (
    <div 
      className={`bg-[#080b11] rounded-3xl overflow-hidden shadow-[0_40px_120px_rgba(0,0,0,0.8)] border border-[#1e222c] flex flex-col w-full transition-all duration-500 hover:border-cyan-500/40 hover:shadow-[0_40px_140px_rgba(125,211,252,0.15)] hover:-translate-y-1.5 active:scale-[0.995] ${className}`}
    >
      <div className="bg-[#0d1117] h-16 px-6 flex items-center justify-between border-b border-[#1e222c] flex-shrink-0">
        <div className="flex items-center space-x-4">
          <div className="flex space-x-2 mr-4">
            <div className="w-3 h-3 rounded-full bg-[#ff5f56]"></div>
            <div className="w-3 h-3 rounded-full bg-[#ffbd2e]"></div>
            <div className="w-3 h-3 rounded-full bg-[#27c93f]"></div>
          </div>
          <div className="text-[10px] text-gray-500 font-mono font-bold tracking-[0.3em] uppercase opacity-60">
            {title}
          </div>
        </div>
        
        {executable && (
          <button 
            onClick={handleRun}
            disabled={isRunning}
            style={{ 
              backgroundColor: isRunning ? 'transparent' : `${COLORS.cyan}22`,
              color: COLORS.cyan,
              borderColor: `${COLORS.cyan}44`
            }}
            className="group flex items-center space-x-3 text-[10px] font-black uppercase tracking-[0.2em] px-4 py-2 rounded-lg border hover:brightness-125 transition-colors disabled:opacity-50"
          >
            <div className={`w-2 h-2 rounded-full ${isRunning ? 'animate-pulse bg-yellow-500' : 'bg-cyan-400'}`}></div>
            <span>{isRunning ? 'RUNNING...' : 'RUN'}</span>
          </button>
        )}
      </div>

      <div className="flex-1 overflow-hidden bg-[#080b11] relative flex flex-col">
        <div className="flex-1 px-8 py-8 font-mono text-[16px] leading-[1.6] overflow-auto text-[#c9d1d9] border-b border-white/5 scrollbar-hide">
          <pre className="m-0 whitespace-pre min-w-max">{children}</pre>
        </div>

        <div className="h-[180px] bg-[#020408] font-mono text-[13px] overflow-y-auto flex flex-col p-6 scrollbar-hide">
          <div className="flex items-center justify-between mb-4 flex-shrink-0 h-4">
            <div className="text-[9px] font-black uppercase tracking-[0.4em] text-gray-800">
              Terminal Output
            </div>
          </div>
          
          <div className="flex-1 space-y-1 mt-2">
            {displayedOutput.length === 0 && !isRunning ? (
              <div className="flex items-center space-x-3 text-gray-800 opacity-40">
                <span className="text-cyan-900 font-bold">$</span>
                <span className="italic">Ready for simulation...</span>
              </div>
            ) : (
              displayedOutput.map((line, idx) => renderLogLine(line, idx))
            )}
            {isRunning && (
              <div className="flex items-center space-x-3 text-white/20">
                <span className="text-cyan-500/30 font-bold">$</span>
                <span className="w-2 h-4 bg-cyan-500/20 animate-pulse"></span>
              </div>
            )}
            <div ref={outputEndRef} className="h-2" />
          </div>
        </div>
      </div>
    </div>
  );
};