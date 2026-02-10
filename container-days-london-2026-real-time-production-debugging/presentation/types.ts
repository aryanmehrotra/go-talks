
import React from 'react';

export type SlideLayout = 
  | 'centered' 
  | 'split' 
  | 'grid'
  | 'code' 
  | 'metric' 
  | 'matrix' 
  | 'diagram' 
  | 'flamegraph' 
  | 'table' 
  | 'comparison' 
  | 'checklist' 
  | 'qr';

export interface SlideData {
  id: number;
  title: string;
  subtitle?: string;
  layout: SlideLayout;
  content: (step: number) => React.ReactNode;
  steps?: number;
  speakerNotes?: string;
  footer?: string;
  background?: string;
  speaker?: 'A' | 'B' | 'Both';
}
