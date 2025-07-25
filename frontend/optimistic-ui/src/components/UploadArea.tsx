import { useRef, useState } from 'react';

interface UploadAreaProps {
  onFilesSelected: (files: FileList) => void;
}

export default function UploadArea({ onFilesSelected }: UploadAreaProps) {
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [isDragOver, setIsDragOver] = useState(false);

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragOver(true);
  };

  const handleDragLeave = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragOver(false);
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragOver(false);
    
    const files = e.dataTransfer.files;
    if (files.length > 0) {
      onFilesSelected(files);
    }
  };

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files;
    if (files && files.length > 0) {
      onFilesSelected(files);
    }
    // 같은 파일을 다시 선택할 수 있도록 값 초기화
    e.target.value = '';
  };

  const handleClick = () => {
    fileInputRef.current?.click();
  };

  return (
    <div
      className={`
        relative border-2 border-dashed rounded-lg p-8 text-center cursor-pointer
        transition-colors duration-200 ease-in-out
        ${isDragOver 
          ? 'border-blue-500 bg-blue-50' 
          : 'border-gray-300 hover:border-gray-400 hover:bg-gray-50'
        }
      `}
      onDragOver={handleDragOver}
      onDragLeave={handleDragLeave}
      onDrop={handleDrop}
      onClick={handleClick}
    >
      <input
        ref={fileInputRef}
        type="file"
        multiple
        accept="image/*"
        onChange={handleFileSelect}
        className="hidden"
      />
      
      <div className="space-y-4">
        <div className="text-6xl">📁</div>
        
        <div>
          <h3 className="text-lg font-medium text-gray-900 mb-2">
            이미지 파일을 업로드하세요
          </h3>
          <p className="text-sm text-gray-600">
            파일을 드래그 앤 드롭하거나 클릭하여 선택하세요
          </p>
        </div>

        <div className="flex items-center justify-center">
          <button
            type="button"
            className="px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 transition-colors"
            onClick={(e) => {
              e.stopPropagation();
              handleClick();
            }}
          >
            파일 선택
          </button>
        </div>

        <p className="text-xs text-gray-500">
          JPG, PNG, GIF 등의 이미지 파일만 업로드 가능합니다
        </p>
      </div>

      {isDragOver && (
        <div className="absolute inset-0 bg-blue-500 bg-opacity-10 rounded-lg flex items-center justify-center">
          <div className="text-blue-600 font-medium">파일을 놓아주세요!</div>
        </div>
      )}
    </div>
  );
}