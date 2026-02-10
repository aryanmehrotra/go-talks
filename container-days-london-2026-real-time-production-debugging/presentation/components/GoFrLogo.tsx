
import React from 'react';
import { COLORS } from '../constants';

export const GoFrLogo: React.FC<{ className?: string }> = ({ className = "" }) => (
  <div className={`flex items-center ${className}`}>
    <div 
      className="text-[2rem] tracking-tight flex items-baseline" 
      style={{ fontFamily: "'Inter', sans-serif", lineHeight: '2.5rem' }}
    >
      <span 
        className="font-bold italic" 
        style={{ color: COLORS.cyan }}
      >
        Go
      </span>
      <span className="inline-block w-3"></span>
      <span className="font-bold text-white">
        Fr
      </span>
    </div>
  </div>
);
