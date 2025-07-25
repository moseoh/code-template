import { ImageUpload } from '@/types/upload';
import Image from 'next/image';

interface ImageCardProps {
  upload: ImageUpload;
  onRemove?: (id: string) => void;
}

export default function ImageCard({ upload, onRemove }: ImageCardProps) {
  const formatTime = (date: Date) => {
    return date.toLocaleTimeString('ko-KR', { 
      hour: '2-digit', 
      minute: '2-digit', 
      second: '2-digit' 
    });
  };

  const getElapsedTime = (start: Date, end?: Date) => {
    const endTime = end || new Date();
    const diff = endTime.getTime() - start.getTime();
    return `${(diff / 1000).toFixed(1)}초`;
  };

  const getStatusIcon = () => {
    switch (upload.status) {
      case 'uploading':
        return (
          <div className="animate-spin w-5 h-5 border-2 border-blue-500 border-t-transparent rounded-full"></div>
        );
      case 'completed':
        return (
          <div className="w-5 h-5 bg-green-500 rounded-full flex items-center justify-center">
            <span className="text-white text-xs">✓</span>
          </div>
        );
      case 'error':
        return (
          <div className="w-5 h-5 bg-red-500 rounded-full flex items-center justify-center">
            <span className="text-white text-xs">✕</span>
          </div>
        );
    }
  };

  const getStatusText = () => {
    switch (upload.status) {
      case 'uploading':
        return '업로드 중...';
      case 'completed':
        return '업로드 완료';
      case 'error':
        return '오류';
    }
  };

  return (
    <div className="bg-white rounded-lg shadow-md overflow-hidden border border-gray-200">
      <div className="relative aspect-video bg-gray-100">
        {upload.imageUrl ? (
          <Image
            src={upload.imageUrl}
            alt={upload.file.name}
            fill
            sizes="(max-width: 768px) 100vw, (max-width: 1200px) 50vw, 25vw"
            className={`object-cover ${upload.status === 'uploading' ? 'opacity-50' : ''}`}
          />
        ) : (
          <div className="w-full h-full flex items-center justify-center">
            <div className="text-center">
              <div className="text-4xl mb-2">📷</div>
              <div className="text-sm text-gray-500">미리보기 생성 중...</div>
            </div>
          </div>
        )}
        
        {/* 상태 오버레이 */}
        {upload.status !== 'completed' && (
          <div className="absolute inset-0 bg-black bg-opacity-30 flex items-center justify-center">
            {getStatusIcon()}
          </div>
        )}
        
        {/* 삭제 버튼 */}
        {onRemove && (
          <button
            onClick={() => onRemove(upload.id)}
            className="absolute top-2 right-2 w-6 h-6 bg-red-500 text-white rounded-full text-xs hover:bg-red-600"
          >
            ✕
          </button>
        )}
      </div>
      
      <div className="p-4">
        <div className="flex items-center gap-2 mb-3">
          {getStatusIcon()}
          <span className="text-sm font-medium">{getStatusText()}</span>
        </div>
        
        <div className="text-xs text-gray-600 space-y-1">
          <div className="font-medium truncate">{upload.file.name}</div>
          <div>크기: {(upload.file.size / 1024 / 1024).toFixed(2)} MB</div>
          
          <div className="border-t pt-2 mt-2 space-y-1">
            <div>업로드 시작: {formatTime(upload.uploadStartTime)}</div>
            
            {upload.uploadCompleteTime && (
              <div>업로드 완료: {formatTime(upload.uploadCompleteTime)} 
                <span className="text-blue-600 ml-1">
                  ({getElapsedTime(upload.uploadStartTime, upload.uploadCompleteTime)})
                </span>
              </div>
            )}
            
            {upload.thumbnailStartTime && (
              <div>썸네일 처리: 백그라운드에서 진행 중</div>
            )}
            
            {upload.status === 'completed' && (
              <div className="text-green-600 font-medium">
                업로드 소요시간: {getElapsedTime(upload.uploadStartTime, upload.uploadCompleteTime)}
              </div>
            )}
            
            {upload.error && (
              <div className="text-red-500 text-xs mt-2">{upload.error}</div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}