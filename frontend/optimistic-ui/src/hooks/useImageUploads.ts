import { useState, useCallback } from 'react';
import { ImageUpload, UploadResponse } from '@/types/upload';

export const useImageUploads = () => {
  const [uploads, setUploads] = useState<ImageUpload[]>([]);

  const addUpload = useCallback((file: File) => {
    const newUpload: ImageUpload = {
      id: `${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
      file,
      status: 'uploading',
      uploadStartTime: new Date(),
    };

    setUploads(prev => [newUpload, ...prev]);
    return newUpload.id;
  }, []);

  const updateUploadStatus = useCallback((id: string, updates: Partial<ImageUpload>) => {
    setUploads(prev => prev.map(upload => 
      upload.id === id ? { ...upload, ...updates } : upload
    ));
  }, []);

  const uploadFile = useCallback(async (file: File) => {
    const uploadId = addUpload(file);

    try {
      const formData = new FormData();
      formData.append('file', file);

      // 정확한 시간 측정을 위해 performance.now() 사용
      const startTime = performance.now();
      console.log(`[${file.name}] 업로드 시작: ${file.size} bytes`);

      const response = await fetch('/api/upload', {
        method: 'POST',
        body: formData,
      });

      const endTime = performance.now();
      const duration = endTime - startTime;
      console.log(`[${file.name}] 업로드 완료: ${duration.toFixed(2)}ms`);

      const result: UploadResponse = await response.json();

      if (result.success) {
        // API 응답 도착 시점에 업로드 완료 시간 기록하고 바로 완료 상태로 변경
        updateUploadStatus(uploadId, {
          status: 'completed' as const,
          uploadCompleteTime: new Date(),
          imageUrl: result.imageUrl,
          // 썸네일은 백그라운드에서 처리되므로 일단 시작 시간만 기록
          thumbnailStartTime: new Date(),
        });
      } else {
        updateUploadStatus(uploadId, {
          status: 'error' as const,
          error: result.error || '업로드에 실패했습니다.',
        });
      }
    } catch (error) {
      console.error('업로드 에러:', error);
      updateUploadStatus(uploadId, {
        status: 'error' as const,
        error: '네트워크 오류가 발생했습니다.',
      });
    }
  }, [addUpload, updateUploadStatus]);

  const uploadMultipleFiles = useCallback(async (files: FileList) => {
    const fileArray = Array.from(files);
    const imageFiles = fileArray.filter(file => file.type.startsWith('image/'));
    
    // 동시에 여러 파일 업로드
    await Promise.all(imageFiles.map(file => uploadFile(file)));
  }, [uploadFile]);

  const removeUpload = useCallback((id: string) => {
    setUploads(prev => prev.filter(upload => upload.id !== id));
  }, []);

  return {
    uploads,
    uploadFile,
    uploadMultipleFiles,
    removeUpload,
  };
};