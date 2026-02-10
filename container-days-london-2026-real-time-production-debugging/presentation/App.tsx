import React, { useState, useEffect, useCallback, useRef } from 'react';
import { Terminal } from './components/Terminal';
import { GoFrLogo } from './components/GoFrLogo';
import { Gopher } from './components/Gopher';
import { ImageFrame } from './components/ImageFrame';
import { GrafanaGraph } from './components/GrafanaGraph';
import { COLORS } from './constants';
import { SlideData } from './types';

const BackgroundGlow = ({ color = COLORS.cyan, size = "90%", opacity = "15" }: { color?: string, size?: string, opacity?: string }) => (
  <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 -z-10 pointer-events-none rounded-full" 
       style={{ 
         width: size, 
         height: size, 
         background: `radial-gradient(circle, ${color}${opacity} 0%, transparent 70%)`, 
         filter: 'blur(160px)' 
       }}></div>
);

const SocialRow = ({ icon, handle }: { icon: React.ReactNode, handle: string }) => (
  <div className="flex items-center space-x-6 text-gray-400 w-full group/row">
    <div className="w-12 h-12 flex-shrink-0 flex items-center justify-center bg-white/5 border border-white/10 rounded-xl text-cyan-400 group-hover/row:border-cyan-500/50 group-hover/row:bg-cyan-500/10 transition-all">
      {icon}
    </div>
    <span className="font-mono text-2xl tracking-tight text-gray-300 group-hover/row:text-white transition-colors">{handle}</span>
  </div>
);

const SpeakerSocials = ({ github, linkedin, x }: { github: string, linkedin: string, x: string }) => (
  <div className="flex flex-col space-y-4 mt-12 items-start w-fit mx-auto">
    <SocialRow 
      handle={`github.com/${github}`}
      icon={<svg className="w-6 h-6" fill="currentColor" viewBox="0 0 24 24"><path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-