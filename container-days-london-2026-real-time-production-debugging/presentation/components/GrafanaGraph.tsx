import React from 'react';
import { COLORS } from '../constants';

interface GraphSeries {
  label: string;
  color: string;
  data: number[]; 
}

interface GrafanaGraphProps {
  title: string;
  yLabel?: string;
  series: GraphSeries[];
  height?: number;
  showLegend?: boolean;
  maxValue?: number; 
}

export const GrafanaGraph: React.FC<GrafanaGraphProps> = ({ 
  title, 
  yLabel, 
  series, 
  height = 240, 
  showLegend = true,
  maxValue
}) => {
  const points = series[0]?.data.length || 20;
  const width = 800;
  const paddingX = 80;
  const paddingTop = 40;
  const paddingBottom = 50;

  // Determine the max value for scaling. 
  const dataMax = Math.max(...series.flatMap(s => s.data), 10);
  const scaleMax = maxValue || (dataMax > 100 ? Math.ceil(dataMax / 100) * 100 : 100);
  
  const getY = (val: number) => {
    // Clamp values to prevent "going off the chart"
    const clampedVal = Math.min(scaleMax, Math.max(0, val));
    const availableHeight = height - paddingTop - paddingBottom;
    // We add 2px of inner padding to ensure the stroke doesn't touch the top edge
    return (height - paddingBottom) - (clampedVal * (availableHeight - 4) / scaleMax) - 2;
  };

  const getX = (index: number) => {
    return paddingX + (index * (width - paddingX * 2) / (points - 1));
  };
  
  const getPath = (data: number[]) => {
    return data.map((val, i) => {
      const x = getX(i);
      const y = getY(val);
      return `${i === 0 ? 'M' : 'L'} ${x} ${y}`;
    }).join(' ');
  };

  const getAreaPath = (data: number[]) => {
    const p = getPath(data);
    const lastX = getX(data.length - 1);
    return `${p} L ${lastX} ${height - paddingBottom} L ${paddingX} ${height - paddingBottom} Z`;
  };

  return (
    <div className="bg-[#11141d] border border-white/10 rounded-xl p-6 shadow-inner w-full flex flex-col group transition-all hover:border-white/20">
      <div className="flex justify-between items-center mb-4">
        <h4 className="text-gray-400 font-bold uppercase tracking-wider text-lg">{title}</h4>
        <div className="flex space-x-1.5 opacity-30">
          <div className="w-1.5 h-1.5 rounded-full bg-gray-500"></div>
          <div className="w-1.5 h-1.5 rounded-full bg-gray-500"></div>
        </div>
      </div>
      
      <div className="relative flex-1" style={{ height: `${height}px` }}>
        {/* Y Axis Label */}
        {yLabel && (
          <div className="absolute left-1 top-1/2 -translate-y-1/2 -rotate-90 text-[9px] text-gray-600 font-bold tracking-[0.2em] whitespace-nowrap origin-center">
            {yLabel}
          </div>
        )}
        
        <svg viewBox={`0 0 ${width} ${height}`} className="w-full h-full overflow-hidden">
          {/* Grid Lines (Horizontal) */}
          {[0, 0.25, 0.5, 0.75, 1].map((ratio) => {
            const v = ratio * scaleMax;
            const y = getY(v);
            return (
              <g key={ratio}>
                <line x1={paddingX} y1={y} x2={width - paddingX} y2={y} stroke="rgba(255,255,255,0.06)" strokeWidth="1" />
                <text x={paddingX - 12} y={y + 4} textAnchor="end" className="fill-gray-600 text-[10px] font-mono">
                  {v === scaleMax && scaleMax >= 500 ? '600 MB' : v === 0 ? '0' : v.toFixed(0)}
                </text>
              </g>
            );
          })}
          
          {/* Vertical Grid Lines */}
          {[...Array(6)].map((_, i) => {
            const x = getX(i * (points-1) / 5);
            return <line key={i} x1={x} y1={paddingTop} x2={x} y2={height - paddingBottom} stroke="rgba(255,255,255,0.04)" strokeWidth="1" />;
          })}

          {/* Series Rendering */}
          {series.map((s, idx) => (
            <g key={`series-${idx}`}>
              <path d={getAreaPath(s.data)} fill={`${s.color}10`} />
              <path d={getPath(s.data)} stroke={s.color} strokeWidth="2.5" fill="none" strokeLinecap="round" strokeLinejoin="round" />
            </g>
          ))}
          
          {/* Time markers on X Axis */}
          <text x={paddingX} y={height - paddingBottom + 25} className="fill-gray-700 text-[11px] font-mono">19:00</text>
          <text x={width/2} y={height - paddingBottom + 25} textAnchor="middle" className="fill-gray-700 text-[11px] font-mono">19:05</text>
          <text x={width - paddingX} y={height - paddingBottom + 25} textAnchor="end" className="fill-gray-700 text-[11px] font-mono">19:10</text>
        </svg>
      </div>

      {showLegend && (
        <div className="mt-4 flex flex-wrap gap-x-6 gap-y-2 border-t border-white/5 pt-3">
          {series.map((s, i) => (
            <div key={i} className="flex items-center space-x-2">
              <div className="w-3 h-1 rounded-full" style={{ backgroundColor: s.color }}></div>
              <span className="text-[10px] text-gray-500 font-mono tracking-tighter uppercase">{s.label}</span>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};