import React from 'react';
import { COLORS } from '../constants';

interface ImageFrameProps {
  src?: string;
  alt?: string;
  label?: string;
  className?: string;
  placeholderText?: string;
}

export const ImageFrame: React.FC<ImageFrameProps> = ({ src, alt, label, className = "", placeholderText = "Screenshot Placeholder" }) => {
  return (
    <div className={`relative group ${className}`}>
      {/* Decorative Outer Glow */}
      <div className="absolute -inset-1 bg-gradient-to-r from-[#7dd3fc33] to-[#00ADD833] rounded-2xl blur opacity-25 group-hover:opacity-50 transition duration-1000 group-hover:duration-200"></div>
      
      <div className="relative bg-[#0E172A] border border-white/10 rounded-xl overflow-hidden shadow-2xl">
        {/* Terminal-like header for the frame */}
        <div className="bg-[#161b22] px-4 py-2 flex items-center justify-between border-b border-white/5">
          <div className="flex space-x-1.5">
            <div className="w-2 h-2 rounded-full bg-red-500/20"></div>
            <div className="w-2 h-2 rounded-full bg-yellow-500/20"></div>
            <div className="w-2 h-2 rounded-full bg-green-500/20"></div>
          </div>
          {label && <span className="text-[10px] font-bold uppercase tracking-widest text-gray-500">{label}</span>}
        </div>

        {/* Image or Placeholder Content */}
        <div className="aspect-video w-full flex items-center justify-center bg-[#0d1117]">
          {src ? (
            <img src={src} alt={alt || label} className="w-full h-full object-cover" />
          ) : (
            <div className="flex flex-col items-center space-y-4">
              <svg className="w-12 h-12 text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="1.5" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
              </svg>
              <span className="text-gray-600 font-mono text-sm tracking-tighter uppercase">{placeholderText}</span>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};